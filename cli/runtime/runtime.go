package runtime

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
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
	vm, setup, draw, onKey, err := parseJS(string(content))
	var wg *sync.WaitGroup
	if err != nil {
		log.Error(err)
	} else {
		wg = s.runSketch(ctx, done, vm, setup, draw, onKey)
	}

	go func() {
		for {
			select {
			case event, ok := <-s.watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					if !strings.HasSuffix(event.Name, s.filename) {
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
					vm, setup, draw, onKey, err := parseJS(string(content))
					if err != nil {
						log.Error(err)
						continue
					}
					ctx, cancel = context.WithCancel(context.Background())
					wg = s.runSketch(ctx, done, vm, setup, draw, onKey)
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

func (s runtime) RunDemo(demo string) {
	done := make(chan os.Signal, 1)
	vm, setup, draw, onKey, err := parseJS(demo)
	if err != nil {
		log.Error(err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	wg := s.runSketch(ctx, done, vm, setup, draw, onKey)
	<-done
	cancel()
	wg.Wait()
}

func (s runtime) runSketch(ctx context.Context, done chan os.Signal, vm *goja.Runtime, setup, draw goja.Callable, onKey goja.Callable) *sync.WaitGroup {
	panicRecover := func(c *runal.Canvas) {
		if r := recover(); r != nil {
			log.Errorf("%v", r)
			c.DisableRendering()
		}
	}

	var onKeyCallback func(c *runal.Canvas, key string)
	if onKey != nil {
		onKeyCallback = func(c *runal.Canvas, key string) {
			defer panicRecover(c)
			_, err := onKey(goja.Undefined(), vm.ToValue(c), vm.ToValue(key))
			if err != nil {
				log.Error(err)
				c.DisableRendering()
			}
		}
	}

	return runal.Start(
		ctx,
		done,
		func(c *runal.Canvas) {
			defer panicRecover(c)
			vm.Set("console", s.console)
			_, err := setup(goja.Undefined(), vm.ToValue(c))
			if err != nil {
				log.Error(err)
				c.DisableRendering()
			}
		},
		func(c *runal.Canvas) {
			defer panicRecover(c)
			vm.Set("c", c)
			_, err := draw(goja.Undefined(), vm.ToValue(c))
			if err != nil {
				log.Error(err)
				c.DisableRendering()
			}
		},
		onKeyCallback,
	)
}
