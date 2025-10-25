package runtime

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/log"

	"github.com/dop251/goja"
	"github.com/emprcl/runal"
	"github.com/fsnotify/fsnotify"
)

type console struct {
	canvas *runal.Canvas
}

func (c console) Log(messages ...string) {
	v := make([]any, len(messages))
	for i, m := range messages {
		v[i] = m
	}
	c.canvas.Debug(v...)
}

type runtime struct {
	watcher  *fsnotify.Watcher
	filename string
}

func New(filename string, watcher *fsnotify.Watcher) runtime {
	return runtime{
		watcher:  watcher,
		filename: filename,
	}
}

func (s runtime) Run() {
	done := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	content, err := os.ReadFile(s.filename)
	vm, setup, draw, cb, err := parseJS(string(content))
	var wg *sync.WaitGroup
	if err != nil {
		log.Error(err)
	} else {
		wg = s.runSketch(ctx, done, vm, setup, draw, cb)
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
					vm, setup, draw, cb, err := parseJS(string(content))
					if err != nil {
						log.Error(err)
						continue
					}
					ctx, cancel = context.WithCancel(context.Background())
					wg = s.runSketch(ctx, done, vm, setup, draw, cb)
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

func (s runtime) RunInternal(sketch string) {
	vm, setup, draw, callbacks, err := parseJS(sketch)
	if err != nil {
		log.Error(err)
		return
	}
	s.runSketch(context.Background(), nil, vm, setup, draw, callbacks).Wait()
}

func (s runtime) runSketch(ctx context.Context, done chan struct{}, vm *goja.Runtime, setup, draw goja.Callable, cb callbacks) *sync.WaitGroup {
	panicRecover := func(c *runal.Canvas) {
		if r := recover(); r != nil {
			log.Errorf("%v", r)
			c.DisableRendering()
		}
	}

	var onKeyCallback func(c *runal.Canvas, e runal.KeyEvent)
	if cb.onKey != nil {
		onKeyCallback = func(c *runal.Canvas, e runal.KeyEvent) {
			defer panicRecover(c)
			_, err := cb.onKey(goja.Undefined(), vm.ToValue(c), vm.ToValue(e))
			if err != nil {
				log.Error(err)
				c.DisableRendering()
			}
		}
	}

	var onMouseClickCallback func(c *runal.Canvas, e runal.MouseEvent)
	if cb.onMouseClick != nil {
		onMouseClickCallback = func(c *runal.Canvas, e runal.MouseEvent) {
			defer panicRecover(c)
			_, err := cb.onMouseClick(goja.Undefined(), vm.ToValue(c), vm.ToValue(e))
			if err != nil {
				log.Error(err)
				c.DisableRendering()
			}
		}
	}

	var onMouseReleaseCallback func(c *runal.Canvas, e runal.MouseEvent)
	if cb.onMouseRelease != nil {
		onMouseReleaseCallback = func(c *runal.Canvas, e runal.MouseEvent) {
			defer panicRecover(c)
			_, err := cb.onMouseRelease(goja.Undefined(), vm.ToValue(c), vm.ToValue(e))
			if err != nil {
				log.Error(err)
				c.DisableRendering()
			}
		}
	}

	var onMouseWheelCallback func(c *runal.Canvas, e runal.MouseEvent)
	if cb.onMouseWheel != nil {
		onMouseWheelCallback = func(c *runal.Canvas, e runal.MouseEvent) {
			defer panicRecover(c)
			_, err := cb.onMouseWheel(goja.Undefined(), vm.ToValue(c), vm.ToValue(e))
			if err != nil {
				log.Error(err)
				c.DisableRendering()
			}
		}
	}

	var onMouseMoveCallback func(c *runal.Canvas, e runal.MouseEvent)
	if cb.onMouseMove != nil {
		onMouseMoveCallback = func(c *runal.Canvas, e runal.MouseEvent) {
			defer panicRecover(c)
			_, err := cb.onMouseMove(goja.Undefined(), vm.ToValue(c), vm.ToValue(e))
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
			vm.Set("console", console{
				canvas: c,
			})
			_, err := setup(goja.Undefined(), vm.ToValue(c))
			if err != nil {
				log.Error(err)
				c.DisableRendering()
			}
		},
		func(c *runal.Canvas) {
			defer panicRecover(c)
			_, err := draw(goja.Undefined(), vm.ToValue(c))
			if err != nil {
				log.Error(err)
				c.DisableRendering()
			}
		},
		runal.WithOnKey(onKeyCallback),
		runal.WithOnMouseMove(onMouseMoveCallback),
		runal.WithOnMouseClick(onMouseClickCallback),
		runal.WithOnMouseRelease(onMouseReleaseCallback),
		runal.WithOnMouseWheel(onMouseWheelCallback),
	)
}
