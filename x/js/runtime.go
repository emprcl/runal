package js

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/charmbracelet/log"

	"github.com/dop251/goja"
	"github.com/emprcl/runal"
	"github.com/fsnotify/fsnotify"
)

// reloadDebounce is how long the sketch file must go quiet before a reload.
const reloadDebounce = 50 * time.Millisecond

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

type sketch struct {
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

// stop cancels the sketch and waits for the terminal to be restored.
// A nil sketch (one that failed to load) stops cleanly.
func (s *sketch) stop() {
	if s == nil {
		return
	}
	s.cancel()
	s.wg.Wait()
}

// Run starts the sketch and reloads it whenever its file changes,
// blocking until the sketch exits.
func (r runtime) Run() {
	done := make(chan struct{}, 1)

	// Watch the directory, not the file: atomic saves replace the file and
	// would drop a watch registered on the original inode.
	dir := filepath.Dir(r.filename)
	if err := r.watcher.Add(dir); err != nil {
		log.Fatalf("can't watch %s: %v", dir, err)
	}

	current := r.start(done)
	defer func() { current.stop() }()

	// A single save arrives as a burst of events; wait for it to go quiet
	// so an in-place truncate isn't read back while momentarily empty.
	var (
		debounce *time.Timer
		settled  <-chan time.Time
	)
	defer func() {
		if debounce != nil {
			debounce.Stop()
		}
	}()

	for {
		select {
		case <-done:
			return

		case event, ok := <-r.watcher.Events:
			if !ok {
				return
			}
			if !r.shouldReload(event) {
				continue
			}
			if debounce == nil {
				debounce = time.NewTimer(reloadDebounce)
				settled = debounce.C
			} else {
				debounce.Reset(reloadDebounce)
			}

		case <-settled:
			debounce, settled = nil, nil
			current.stop()
			current = r.start(done)

		case err, ok := <-r.watcher.Errors:
			if !ok {
				return
			}
			log.Error(err)
		}
	}
}

// RunInternal runs a sketch embedded in the binary, without file watching.
func (r runtime) RunInternal(script string) {
	vm, setup, draw, cb, err := parseJS(script)
	if err != nil {
		log.Error(err)
		return
	}

	wg, err := runSketch(context.Background(), nil, vm, setup, draw, cb)
	if err != nil {
		log.Error(err)
		return
	}
	wg.Wait()
}

// start loads and runs the sketch, returning nil if it fails to load.
func (r runtime) start(done chan struct{}) *sketch {
	content, err := os.ReadFile(r.filename)
	if err != nil {
		log.Errorf("can't read %s: %v", r.filename, err)
		return nil
	}

	vm, setup, draw, cb, err := parseJS(string(content))
	if err != nil {
		log.Error(err)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg, err := runSketch(ctx, done, vm, setup, draw, cb)
	if err != nil {
		log.Errorf("can't start sketch: %v", err)
		cancel()
		return nil
	}
	return &sketch{
		cancel: cancel,
		wg:     wg,
	}
}

func (r runtime) shouldReload(event fsnotify.Event) bool {
	if !event.Has(fsnotify.Write) && !event.Has(fsnotify.Create) && !event.Has(fsnotify.Rename) {
		return false
	}
	return sameFile(event.Name, r.filename)
}

func sameFile(a, b string) bool {
	return resolvePath(a) == resolvePath(b)
}

func resolvePath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return abs
	}
	return resolved
}

func runSketch(ctx context.Context, done chan struct{}, vm *goja.Runtime, setup, draw goja.Callable, cb callbacks) (*sync.WaitGroup, error) {
	// call invokes a javascript function, turning thrown errors and panics
	// into a stopped sketch rather than a crash.
	call := func(c *runal.Canvas, fn goja.Callable, args ...goja.Value) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("%v", r)
				c.DisableRendering()
			}
		}()
		if _, err := fn(goja.Undefined(), args...); err != nil {
			log.Error(err)
			c.DisableRendering()
		}
	}

	// A nil javascript callback must stay a nil go callback so the canvas
	// can skip the event entirely.
	mouseHandler := func(fn goja.Callable) func(*runal.Canvas, runal.MouseEvent) {
		if fn == nil {
			return nil
		}
		return func(c *runal.Canvas, e runal.MouseEvent) {
			call(c, fn, vm.ToValue(c), vm.ToValue(e))
		}
	}

	var onKey func(*runal.Canvas, runal.KeyEvent)
	if cb.onKey != nil {
		onKey = func(c *runal.Canvas, e runal.KeyEvent) {
			call(c, cb.onKey, vm.ToValue(c), vm.ToValue(e))
		}
	}

	return runal.Start(
		ctx,
		done,
		func(c *runal.Canvas) {
			if err := vm.Set("console", console{canvas: c}); err != nil {
				log.Error(err)
			}
			call(c, setup, vm.ToValue(c))
		},
		func(c *runal.Canvas) {
			call(c, draw, vm.ToValue(c))
		},
		runal.WithOnKey(onKey),
		runal.WithOnMouseMove(mouseHandler(cb.onMouseMove)),
		runal.WithOnMouseClick(mouseHandler(cb.onMouseClick)),
		runal.WithOnMouseRelease(mouseHandler(cb.onMouseRelease)),
		runal.WithOnMouseWheel(mouseHandler(cb.onMouseWheel)),
	)
}
