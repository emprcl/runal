package runtime

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/dop251/goja"
	"github.com/emprcl/runal"
	"github.com/fsnotify/fsnotify"
)

type console struct {
	logger io.Writer
}

func (c console) Log(msg string) {
	fmt.Fprintln(c.logger, msg)
}

type runtime struct {
	watcher  *fsnotify.Watcher
	console  console
	filename string
}

func New(filename string, watcher *fsnotify.Watcher, logger io.Writer) runtime {
	return runtime{
		watcher: watcher,
		console: console{
			logger: logger,
		},
		filename: filename,
	}
}

func (s runtime) Run() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	content, err := os.ReadFile(s.filename)
	vm, setup, draw, err := parseJS(string(content))
	var wg *sync.WaitGroup
	if err != nil {
		log.Error(err)
	} else {
		wg = s.runSketch(ctx, vm, setup, draw)
	}

	go func() {
		for {
			select {
			case event, ok := <-s.watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					if event.Name != s.filename {
						continue
					}
					cancel()
					if wg != nil {
						wg.Wait()
					}
					content, err := os.ReadFile(event.Name)
					if err != nil {
						log.Error(err)
						continue
					}
					vm, setup, draw, err := parseJS(string(content))
					if err != nil {
						log.Error(err)
						continue
					}
					ctx, cancel = context.WithCancel(context.Background())
					wg = s.runSketch(ctx, vm, setup, draw)
				}
			case err, ok := <-s.watcher.Errors:
				if !ok {
					return
				}
				log.Error(err)
			}
		}
	}()

	err = s.watcher.Add(filepath.Dir(s.filename))
	if err != nil {
		log.Fatal(err)
	}

	<-done
	cancel()
	wg.Wait()
}

func (s runtime) runSketch(ctx context.Context, vm *goja.Runtime, setup, draw goja.Callable) *sync.WaitGroup {
	return runal.Start(
		ctx,
		func(c *runal.Canvas) {
			vm.Set("console", s.console)
			vm.Set("c", c)
			_, err := setup(goja.Undefined())
			if err != nil {
				log.Error(err)
				c.DisableRendering()
			}
		},
		func(c *runal.Canvas) {
			vm.Set("c", c)
			_, err := draw(goja.Undefined())
			if err != nil {
				log.Error(err)
				c.DisableRendering()
			}
		},
	)
}
