//go:build js

package canvas

import (
	"context"
	"strings"
	"sync"
	"syscall/js"
)

// Run starts the runal event loop and blocks until the sketch exits.
func Run(ctx context.Context, setup, draw func(c *Canvas), opts ...CallbackOption) {
	Start(ctx, nil, setup, draw, opts...).Wait()
}

// Start runs the runal event loop driven by requestAnimationFrame and DOM
// events, and returns a WaitGroup that completes when the sketch stops.
//
// requestAnimationFrame is used rather than a time.Ticker so that the loop
// always yields to the browser between frames: a Go ticker whose interval is
// shorter than a (slow) frame stays perpetually ready and would starve the JS
// event loop, freezing the page.
func Start(ctx context.Context, done chan struct{}, setup, draw func(c *Canvas), opts ...CallbackOption) *sync.WaitGroup {
	cols, rows := 80, 24
	if display != nil {
		if nc, nr := display.metrics(); nc > 0 && nr > 0 {
			cols, rows = nc, nr
		}
	}
	c := newCanvas(cols, rows)

	eventCallbacks := callbacks{}
	for _, opt := range opts {
		opt(&eventCallbacks)
	}

	// DOM event bridge. Handlers run on the JS event-loop turn and only push
	// into buffered channels; all canvas/goja access happens while draining
	// them in the frame callback below, keeping the JS engine single-threaded.
	keyCh := make(chan KeyEvent, 256)
	moveCh := make(chan MouseEvent, 256)
	clickCh := make(chan MouseEvent, 256)
	releaseCh := make(chan MouseEvent, 256)
	wheelCh := make(chan MouseEvent, 256)
	resizeCh := make(chan struct{}, 1)

	listeners := attachDOMEvents(keyCh, moveCh, clickCh, releaseCh, wheelCh, resizeCh)

	render := func() {
		draw(c)
		c.render()
	}

	perf := js.Global().Get("performance")
	now := func() float64 { return perf.Call("now").Float() }
	interval := 1000.0 / float64(defaultFPS)

	setup(c)
	render()
	lastFrame := now()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	var frameFn js.Func
	stopped := false
	stop := func() {
		if stopped {
			return
		}
		stopped = true
		listeners.release()
		frameFn.Release()
		if done != nil {
			select {
			case done <- struct{}{}:
			default:
			}
		}
		wg.Done()
	}

	// drain* helpers consume all buffered events without blocking.
	drainResize := func() {
		select {
		case <-resizeCh:
			nc, nr := display.metrics()
			c.termWidth, c.termHeight = nc, nr
			if c.autoResize && nc > 0 && nr > 0 {
				c.resize(nc, nr)
			}
		default:
		}
	}
	drainBus := func() (exit bool) {
		for {
			select {
			case event := <-c.bus:
				switch event.name {
				case "fps":
					if event.value > 0 {
						interval = 1000.0 / float64(event.value)
					}
				case "render":
					render()
				case "exit":
					return true
				}
			default:
				return false
			}
		}
	}
	drainInput := func() {
		for {
			select {
			case e := <-keyCh:
				if eventCallbacks.onKey != nil {
					eventCallbacks.onKey(c, e)
				}
			case e := <-moveCh:
				c.setMousePostion(e.X, e.Y)
				if eventCallbacks.onMouseMove != nil {
					eventCallbacks.onMouseMove(c, MouseEvent{X: c.MouseX, Y: c.MouseY})
				}
			case e := <-clickCh:
				c.setMousePostion(e.X, e.Y)
				if eventCallbacks.onMouseClick != nil {
					eventCallbacks.onMouseClick(c, MouseEvent{X: c.MouseX, Y: c.MouseY, Button: e.Button})
				}
			case e := <-releaseCh:
				c.setMousePostion(e.X, e.Y)
				if eventCallbacks.onMouseRelease != nil {
					eventCallbacks.onMouseRelease(c, MouseEvent{X: c.MouseX, Y: c.MouseY, Button: e.Button})
				}
			case e := <-wheelCh:
				c.setMousePostion(e.X, e.Y)
				if eventCallbacks.onMouseWheel != nil {
					eventCallbacks.onMouseWheel(c, MouseEvent{X: c.MouseX, Y: c.MouseY, Button: e.Button})
				}
			default:
				return
			}
		}
	}

	frameFn = js.FuncOf(func(_ js.Value, _ []js.Value) any {
		if ctx.Err() != nil {
			stop()
			return nil
		}
		drainResize()
		drainInput()
		if drainBus() {
			stop()
			return nil
		}
		if c.IsLooping {
			if t := now(); t-lastFrame >= interval {
				render()
				lastFrame = t
			}
		}
		js.Global().Call("requestAnimationFrame", frameFn)
		return nil
	})
	js.Global().Call("requestAnimationFrame", frameFn)

	return wg
}

// domListeners tracks registered js.Func callbacks so they can be released
// when the loop stops.
type domListeners struct {
	el    js.Value
	funcs map[string]js.Func
	ro    js.Value // ResizeObserver
}

func (l domListeners) release() {
	for name, fn := range l.funcs {
		l.el.Call("removeEventListener", name, fn)
		fn.Release()
	}
	if !l.ro.IsUndefined() {
		l.ro.Call("disconnect")
	}
}

func attachDOMEvents(keyCh chan KeyEvent, moveCh, clickCh, releaseCh, wheelCh chan MouseEvent, resizeCh chan struct{}) domListeners {
	el := display.el
	funcs := map[string]js.Func{}

	cellCoords := func(e js.Value) (int, int) {
		x, y := 0, 0
		if display.cellW > 0 {
			x = int(e.Get("offsetX").Float() / display.cellW)
		}
		if display.cellH > 0 {
			y = int(e.Get("offsetY").Float() / display.cellH)
		}
		return x, y
	}

	send := func(ch chan MouseEvent, ev MouseEvent) {
		select {
		case ch <- ev:
		default:
		}
	}

	funcs["keydown"] = js.FuncOf(func(_ js.Value, args []js.Value) any {
		e := args[0]
		select {
		case keyCh <- KeyEvent{Key: mapKey(e.Get("key").String()), Code: e.Get("keyCode").Int()}:
		default:
		}
		return nil
	})
	funcs["mousemove"] = js.FuncOf(func(_ js.Value, args []js.Value) any {
		x, y := cellCoords(args[0])
		send(moveCh, MouseEvent{X: x, Y: y})
		return nil
	})
	funcs["mousedown"] = js.FuncOf(func(_ js.Value, args []js.Value) any {
		x, y := cellCoords(args[0])
		send(clickCh, MouseEvent{X: x, Y: y, Button: mapButton(args[0].Get("button").Int())})
		return nil
	})
	funcs["mouseup"] = js.FuncOf(func(_ js.Value, args []js.Value) any {
		x, y := cellCoords(args[0])
		send(releaseCh, MouseEvent{X: x, Y: y, Button: mapButton(args[0].Get("button").Int())})
		return nil
	})
	funcs["wheel"] = js.FuncOf(func(_ js.Value, args []js.Value) any {
		x, y := cellCoords(args[0])
		dir := "up"
		if args[0].Get("deltaY").Float() > 0 {
			dir = "down"
		}
		send(wheelCh, MouseEvent{X: x, Y: y, Button: dir})
		return nil
	})

	for name, fn := range funcs {
		el.Call("addEventListener", name, fn)
	}
	// The canvas needs focus to receive key events.
	if el.Get("tabIndex").Int() < 0 {
		el.Set("tabIndex", 0)
	}

	ro := js.Value{}
	if obs := js.Global().Get("ResizeObserver"); !obs.IsUndefined() {
		roFn := js.FuncOf(func(_ js.Value, _ []js.Value) any {
			select {
			case resizeCh <- struct{}{}:
			default:
			}
			return nil
		})
		funcs["__resize"] = roFn
		ro = obs.New(roFn)
		ro.Call("observe", el)
	}

	return domListeners{el: el, funcs: funcs, ro: ro}
}

// mapButton converts a DOM MouseEvent.button index to a runal button name.
func mapButton(b int) string {
	switch b {
	case 1:
		return "middle"
	case 2:
		return "right"
	default:
		return "left"
	}
}

// mapKey converts a DOM KeyboardEvent.key to the runal key naming used by the
// terminal backend (charmbracelet/x/input).
func mapKey(k string) string {
	switch k {
	case " ":
		return "space"
	case "ArrowUp":
		return "up"
	case "ArrowDown":
		return "down"
	case "ArrowLeft":
		return "left"
	case "ArrowRight":
		return "right"
	case "Enter":
		return "enter"
	case "Escape":
		return "esc"
	case "Backspace":
		return "backspace"
	case "Tab":
		return "tab"
	case "Delete":
		return "delete"
	}
	if len(k) == 1 {
		return strings.ToLower(k)
	}
	return k
}
