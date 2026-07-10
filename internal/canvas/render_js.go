//go:build js

package canvas

import (
	"encoding/binary"
	"syscall/js"

	"github.com/charmbracelet/x/ansi"
)

// bytesPerCell is the wire size of a single cell pushed to JS:
// rune code point (uint32) + packed foreground RGB (uint32) + packed
// background RGB (uint32), all little-endian.
const bytesPerCell = 12

// jsDisplay holds the browser rendering target and the reusable buffers
// used to ship a frame to JavaScript in a single copy.
type jsDisplay struct {
	el       js.Value // the <canvas> element
	runal    js.Value // the global __runal helper object
	fontSize int

	buf     []byte   // scratch: cols*rows*bytesPerCell
	jsArray js.Value // Uint8Array view shared with JS
	arrSize int

	cellW, cellH float64 // cell size in CSS pixels (for mouse mapping)
}

var display *jsDisplay

// SetDisplayCanvas registers the DOM <canvas> element to render into and the
// font size (in CSS pixels) used for cell metrics. Must be called before the
// sketch is started.
func SetDisplayCanvas(el js.Value, fontSize int) {
	if fontSize <= 0 {
		fontSize = 16
	}
	display = &jsDisplay{
		el:       el,
		runal:    js.Global().Get("__runal"),
		fontSize: fontSize,
	}
}

// metrics asks the JS side for the number of whole cells that fit the current
// element size, caching the per-cell pixel size for mouse mapping. Returns
// (cols, rows), or (0, 0) if unavailable.
func (d *jsDisplay) metrics() (int, int) {
	if d == nil || d.runal.IsUndefined() {
		return 0, 0
	}
	m := d.runal.Call("metrics", d.el, d.fontSize)
	d.cellW = m.Index(2).Float()
	d.cellH = m.Index(3).Float()
	return m.Index(0).Int(), m.Index(1).Int()
}

func (d *jsDisplay) ensure(size int) {
	if d.arrSize == size && len(d.buf) == size {
		return
	}
	d.buf = make([]byte, size)
	d.jsArray = js.Global().Get("Uint8Array").New(size)
	d.arrSize = size
}

func packRGB(col ansi.Color) uint32 {
	if col == nil {
		return 0
	}
	r, g, b, _ := col.RGBA()
	return (r>>8)<<16 | (g>>8)<<8 | (b >> 8)
}

func (c *Canvas) render() {
	if c.disabled || display == nil || display.runal.IsUndefined() {
		return
	}

	// Each logical cell expands to one or two display glyphs depending on the
	// cell mode, mirroring the terminal renderer (render_term.go): a cell is
	// drawn as a single glyph by default, or as two adjacent glyphs when a
	// cell mode is active (the char doubled, followed by the custom spacing
	// rune, or followed by a per-cell pad char set by Text()).
	glyphs := 1
	if c.cellMode.enabled() {
		glyphs = 2
	}
	rows := c.Height
	dispCols := c.Width * glyphs
	display.ensure(dispCols * rows * bytesPerCell)

	i := 0
	putGlyph := func(r rune, fg, bg ansi.Color) {
		binary.LittleEndian.PutUint32(display.buf[i:], uint32(r))
		binary.LittleEndian.PutUint32(display.buf[i+4:], packRGB(fg))
		binary.LittleEndian.PutUint32(display.buf[i+8:], packRGB(bg))
		i += bytesPerCell
	}

	for y := range c.buffer {
		for x := range c.buffer[y] {
			cll := c.buffer[y][x]
			var first, second rune
			var fg, bg ansi.Color
			if cll.char == 0 {
				first = c.nextBackgroundRune()
				fg, bg = c.backgroundFg, c.backgroundBg
				if c.cellMode == cellModeCustom {
					second = c.cellModeRune
				} else {
					second = first
				}
			} else {
				first = cll.char
				fg, bg = cll.foreground, cll.background
				switch {
				case cll.padChar != 0:
					second = cll.padChar
				case c.cellMode == cellModeCustom:
					second = c.cellModeRune
				default:
					second = cll.char
				}
			}
			putGlyph(first, fg, bg)
			if glyphs == 2 {
				putGlyph(second, fg, bg)
			}
		}
	}

	js.CopyBytesToJS(display.jsArray, display.buf)
	display.runal.Call("draw", display.el, dispCols, rows, display.jsArray, display.fontSize)

	c.Framecount++
	c.reset()
	c.flushDebug()
}

// flushDebug mirrors terminal debug output by forwarding console.log messages
// (collected via Canvas.Debug) to the browser console.
func (c *Canvas) flushDebug() {
	if len(c.debugBuffer) == 0 {
		return
	}
	console := js.Global().Get("console")
	for _, msg := range c.debugBuffer {
		console.Call("log", msg)
	}
	c.debugBuffer = c.debugBuffer[:0]
}
