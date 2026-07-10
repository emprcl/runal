//go:build js

// Web (wasm) backend: the sketch runs in the browser's native JS engine and
// calls into this Go canvas engine through //go:wasmexport functions (direct
// wasm exports — no JS interpreter embedded, and no syscall/js call overhead).
//
// All functions operate on a single package-global Canvas (the browser has one
// sketch/canvas). Numbers cross the boundary directly; colors are passed as
// packed 0xRRGGBB uint32; strings (stroke/fill/background text) are written by
// JS into the exported scratch buffer and read here by length.
//
// The rendered frame is produced by rRender, which fills webBlob in linear
// memory (rune + packed fg + packed bg per display glyph) and returns its
// pointer; JS reads it directly from the wasm memory and paints the canvas.
package canvas

import (
	"encoding/binary"
	"unsafe"

	"github.com/charmbracelet/x/ansi"
)

const bytesPerCell = 12

var (
	web     *Canvas
	webBlob []byte
	scratch = make([]byte, 8192)
)

func rgb(v uint32) ansi.Color {
	return ansi.RGBColor{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v)}
}

func packRGB(col ansi.Color) uint32 {
	if col == nil {
		return 0
	}
	r, g, b, _ := col.RGBA()
	return (r>>8)<<16 | (g>>8)<<8 | (b >> 8)
}

func scratchStr(n int32) string {
	if n < 0 || int(n) > len(scratch) {
		return ""
	}
	return string(scratch[:n])
}

// --- lifecycle / config ---

//go:wasmexport rInit
func rInit(cols, rows int32) {
	web = newCanvas(int(cols), int(rows))
}

//go:wasmexport rScratchPtr
func rScratchPtr() int32 { return int32(uintptr(unsafe.Pointer(&scratch[0]))) }

//go:wasmexport rResize
func rResize(cols, rows int32) { web.resize(int(cols), int(rows)) }

//go:wasmexport rSize
func rSize(w, h int32) { web.Size(int(w), int(h)) }

//go:wasmexport rAutoResize
func rAutoResize() int32 {
	if web.autoResize {
		return 1
	}
	return 0
}

//go:wasmexport rClear
func rClear() { web.Clear() }

//go:wasmexport rFps
func rFps(fps int32) { web.fps = int(fps) }

//go:wasmexport rGetFps
func rGetFps() int32 { return int32(web.fps) }

//go:wasmexport rWidth
func rWidth() int32 { return int32(web.Width) }

//go:wasmexport rHeight
func rHeight() int32 { return int32(web.Height) }

//go:wasmexport rFramecount
func rFramecount() int32 { return int32(web.Framecount) }

// --- stroke / fill / background state ---

//go:wasmexport rStroke
func rStroke(textLen int32, fg, bg uint32) {
	web.StrokeText(scratchStr(textLen))
	web.strokeFg = rgb(fg)
	web.strokeBg = rgb(bg)
}

//go:wasmexport rStrokeText
func rStrokeText(textLen int32) { web.StrokeText(scratchStr(textLen)) }

//go:wasmexport rStrokeFg
func rStrokeFg(fg uint32) { web.stroke = true; web.strokeFg = rgb(fg) }

//go:wasmexport rStrokeBg
func rStrokeBg(bg uint32) { web.stroke = true; web.strokeBg = rgb(bg) }

//go:wasmexport rFill
func rFill(textLen int32, fg, bg uint32) {
	web.FillText(scratchStr(textLen))
	web.fillFg = rgb(fg)
	web.fillBg = rgb(bg)
}

//go:wasmexport rFillText
func rFillText(textLen int32) { web.FillText(scratchStr(textLen)) }

//go:wasmexport rFillFg
func rFillFg(fg uint32) { web.fill = true; web.fillFg = rgb(fg) }

//go:wasmexport rFillBg
func rFillBg(bg uint32) { web.fill = true; web.fillBg = rgb(bg) }

//go:wasmexport rBackground
func rBackground(textLen int32, fg, bg uint32) {
	web.BackgroundText(scratchStr(textLen))
	web.backgroundFg = rgb(fg)
	web.backgroundBg = rgb(bg)
}

//go:wasmexport rBackgroundText
func rBackgroundText(textLen int32) { web.BackgroundText(scratchStr(textLen)) }

//go:wasmexport rBackgroundFg
func rBackgroundFg(fg uint32) { web.backgroundFg = rgb(fg) }

//go:wasmexport rBackgroundBg
func rBackgroundBg(bg uint32) { web.backgroundBg = rgb(bg) }

//go:wasmexport rNoStroke
func rNoStroke() { web.NoStroke() }

//go:wasmexport rNoFill
func rNoFill() { web.NoFill() }

// --- transforms ---

//go:wasmexport rPush
func rPush() { web.Push() }

//go:wasmexport rPop
func rPop() { web.Pop() }

//go:wasmexport rTranslate
func rTranslate(x, y int32) { web.Translate(int(x), int(y)) }

//go:wasmexport rRotate
func rRotate(a float64) { web.Rotate(a) }

//go:wasmexport rScale
func rScale(s float64) { web.Scale(s) }

// --- shapes ---

//go:wasmexport rPoint
func rPoint(x, y int32) { web.Point(int(x), int(y)) }

//go:wasmexport rLine
func rLine(x1, y1, x2, y2 int32) { web.Line(int(x1), int(y1), int(x2), int(y2)) }

//go:wasmexport rRect
func rRect(x, y, w, h int32) { web.Rect(int(x), int(y), int(w), int(h)) }

//go:wasmexport rSquare
func rSquare(x, y, s int32) { web.Square(int(x), int(y), int(s)) }

//go:wasmexport rEllipse
func rEllipse(cx, cy, rx, ry int32) { web.Ellipse(int(cx), int(cy), int(rx), int(ry)) }

//go:wasmexport rCircle
func rCircle(cx, cy, r int32) { web.Circle(int(cx), int(cy), int(r)) }

//go:wasmexport rTriangle
func rTriangle(x1, y1, x2, y2, x3, y3 int32) {
	web.Triangle(int(x1), int(y1), int(x2), int(y2), int(x3), int(y3))
}

//go:wasmexport rQuad
func rQuad(x1, y1, x2, y2, x3, y3, x4, y4 int32) {
	web.Quad(int(x1), int(y1), int(x2), int(y2), int(x3), int(y3), int(x4), int(y4))
}

//go:wasmexport rBezier
func rBezier(x1, y1, x2, y2, x3, y3, x4, y4 int32) {
	web.Bezier(int(x1), int(y1), int(x2), int(y2), int(x3), int(y3), int(x4), int(y4))
}

//go:wasmexport rText
func rText(textLen, x, y int32) { web.Text(scratchStr(textLen), int(x), int(y)) }

// --- noise / random (stateful: kept in the Go engine) ---

//go:wasmexport rNoise1D
func rNoise1D(x float64) float64 { return web.Noise1D(x) }

//go:wasmexport rNoise2D
func rNoise2D(x, y float64) float64 { return web.Noise2D(x, y) }

//go:wasmexport rNoiseSeed
func rNoiseSeed(seed int64) { web.NoiseSeed(seed) }

//go:wasmexport rNoiseLoop
func rNoiseLoop(angle, radius float64) float64 { return web.NoiseLoop(angle, radius) }

//go:wasmexport rNoiseLoop1D
func rNoiseLoop1D(angle, radius float64, x int32) float64 {
	return web.NoiseLoop1D(angle, radius, int(x))
}

//go:wasmexport rNoiseLoop2D
func rNoiseLoop2D(angle, radius float64, x, y int32) float64 {
	return web.NoiseLoop2D(angle, radius, int(x), int(y))
}

//go:wasmexport rLoopAngle
func rLoopAngle(duration int32) float64 { return web.LoopAngle(int(duration)) }

//go:wasmexport rRandom
func rRandom(minimum, maximum int32) float64 { return web.Random(int(minimum), int(maximum)) }

//go:wasmexport rRandomSeed
func rRandomSeed(seed int64) { web.RandomSeed(seed) }

// --- colors (parsing/space conversion stays in the Go engine) ---
// Each returns a packed 0xRRGGBB value; JS formats it back to a hex string so
// sketches see the same string values as the terminal API.

//go:wasmexport rColorRGB
func rColorRGB(r, g, b int32) uint32 { return packRGB(color(web.ColorRGB(int(r), int(g), int(b)))) }

//go:wasmexport rColorHSL
func rColorHSL(h, s, l int32) uint32 { return packRGB(color(web.ColorHSL(int(h), int(s), int(l)))) }

//go:wasmexport rColorHSV
func rColorHSV(h, s, v int32) uint32 { return packRGB(color(web.ColorHSV(int(h), int(s), int(v)))) }

// rResolveColor parses a color string (hex, CSS name, or ansi index) from the
// scratch buffer into a packed 0xRRGGBB value, using the engine's color().
//
//go:wasmexport rResolveColor
func rResolveColor(textLen int32) uint32 { return packRGB(color(scratchStr(textLen))) }

// --- cell mode ---

//go:wasmexport rCellModeDouble
func rCellModeDouble() { web.CellModeDouble() }

//go:wasmexport rCellModeCustom
func rCellModeCustom(r int32) { web.CellModeCustom(string(rune(r))) }

//go:wasmexport rCellModeDefault
func rCellModeDefault() { web.CellModeDefault() }

// --- render ---

//go:wasmexport rDispCols
func rDispCols() int32 {
	if web.cellMode.enabled() {
		return int32(web.Width * 2)
	}
	return int32(web.Width)
}

//go:wasmexport rRows
func rRows() int32 { return int32(web.Height) }

// rRender builds the frame blob into linear memory and returns its pointer.
// Each logical cell expands to one or two display glyphs (cell mode), mirroring
// the terminal renderer.
//
//go:wasmexport rRender
func rRender() int32 {
	glyphs := 1
	if web.cellMode.enabled() {
		glyphs = 2
	}
	rows := web.Height
	dispCols := web.Width * glyphs
	need := dispCols * rows * bytesPerCell
	if cap(webBlob) < need {
		webBlob = make([]byte, need)
	} else {
		webBlob = webBlob[:need]
	}

	i := 0
	put := func(r rune, fg, bg ansi.Color) {
		binary.LittleEndian.PutUint32(webBlob[i:], uint32(r))
		binary.LittleEndian.PutUint32(webBlob[i+4:], packRGB(fg))
		binary.LittleEndian.PutUint32(webBlob[i+8:], packRGB(bg))
		i += bytesPerCell
	}

	for y := range web.buffer {
		for x := range web.buffer[y] {
			cll := web.buffer[y][x]
			var first, second rune
			var fg, bg ansi.Color
			if cll.char == 0 {
				first = web.nextBackgroundRune()
				fg, bg = web.backgroundFg, web.backgroundBg
				if web.cellMode == cellModeCustom {
					second = web.cellModeRune
				} else {
					second = first
				}
			} else {
				first = cll.char
				fg, bg = cll.foreground, cll.background
				switch {
				case cll.padChar != 0:
					second = cll.padChar
				case web.cellMode == cellModeCustom:
					second = web.cellModeRune
				default:
					second = cll.char
				}
			}
			put(first, fg, bg)
			if glyphs == 2 {
				put(second, fg, bg)
			}
		}
	}

	web.Framecount++
	web.reset()
	if need == 0 {
		return 0
	}
	return int32(uintptr(unsafe.Pointer(&webBlob[0])))
}
