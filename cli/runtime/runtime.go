package runtime

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

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
	fps      int
}

func New(filename string, watcher *fsnotify.Watcher, logger io.Writer, fps int) runtime {
	return runtime{
		watcher: watcher,
		console: console{
			logger: logger,
		},
		filename: filename,
		fps:      fps,
	}
}

func (s runtime) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	content, err := os.ReadFile(s.filename)
	vm, setup, draw, err := parseJS(string(content))
	if err != nil {
		log.Error(err)
	} else {
		s.runSketch(ctx, vm, setup, draw)
	}

	go func() {
		for {
			select {
			case event, ok := <-s.watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					cancel()
					// let finish the last frame rendering
					time.Sleep(time.Duration(2*1000/s.fps) * time.Millisecond)
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
					s.runSketch(ctx, vm, setup, draw)
				}
			case err, ok := <-s.watcher.Errors:
				if !ok {
					return
				}
				log.Error(err)
			}
		}
	}()

	err = s.watcher.Add(s.filename)
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}

func (s runtime) runSketch(ctx context.Context, vm *goja.Runtime, setup, draw goja.Callable) {
	runal.Run(
		ctx,
		func(c *runal.Canvas) {
			vm.Set("console", s.console)
			vm.Set("runal", c)
			_, err := setup(goja.Undefined())
			if err != nil {
				panic(err)
			}
		},
		func(c *runal.Canvas) {
			vm.Set("runal", c)
			_, err := draw(goja.Undefined())
			if err != nil {
				panic(err)
			}
		},
		runal.WithFPS(s.fps),
	)
}
