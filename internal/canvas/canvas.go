package canvas

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"math"
	"math/rand"
	"os"
	"strings"

	perlin "github.com/aquilax/go-perlin"
	"github.com/charmbracelet/x/ansi"
	ansitoimage "github.com/xaviergodart/go-ansi-to-image"
)

const (
	defaultPaddingRune    = ' '
	defaultStrokeText     = "."
	defaultFillText       = " "
	defaultBackgroundText = " "
)

type cellMode uint8

const (
	cellModeDisabled cellMode = iota
	cellModeDouble
	cellModeCustom
)

func (c cellMode) enabled() bool {
	return c != cellModeDisabled
}

type videoFormat uint8

const (
	videoFormatGif videoFormat = iota
	videoFormatMp4
)

// Canvas represents a drawable area where shapes, text, and effects
// can be rendered.
type Canvas struct {
	buffer  frame
	output  io.Writer
	capture *ansitoimage.Converter
	frames  []image.Image
	videoFormat
	noise  *perlin.Perlin
	random *rand.Rand

	debugBuffer []string

	state *state

	strokeFg, strokeBg         ansi.Color
	fillFg, fillBg             ansi.Color
	backgroundFg, backgroundBg ansi.Color
	lastStyle                  style

	strokeText, fillText, backgroundText string
	strokeRunes, fillRunes, backgroundRunes []rune

	saveFilename string

	bus chan event

	Width, Height  int
	MouseX, MouseY int
	Framecount     int
	fps            int

	strokeIndex, backgroundIndex int
	termWidth, termHeight        int
	originX, originY             int
	rotationAngle                float64
	scale                        float64

	cellModeRune rune
	cellMode     cellMode
	stroke       bool
	fill         bool
	isFilling    bool
	IsLooping    bool
	save         bool
	autoResize   bool
	disabled     bool
}

func newCanvas(width, height int) *Canvas {
	return &Canvas{
		Width:          width,
		Height:         height,
		fps:            defaultFPS,
		bus:            make(chan event, 16),
		termWidth:      width,
		termHeight:     height,
		cellModeRune:   defaultPaddingRune,
		cellMode:       cellModeDisabled,
		buffer:         newFrame(width, height),
		output:         os.Stdout,
		capture:        newCapture(width, height),
		noise:          newNoise(),
		random:         newRandom(),
		debugBuffer:    make([]string, 0, maxDebugBufferSize),
		scale:          1,
		strokeFg:       color("#ffffff"),
		strokeBg:       color("#000000"),
		fillFg:         color("#ffffff"),
		fillBg:         color("#000000"),
		backgroundFg:   color("#ffffff"),
		backgroundBg:   color("#000000"),
		strokeText:     defaultStrokeText,
		strokeRunes:    []rune(defaultStrokeText),
		fillText:       defaultFillText,
		fillRunes:      []rune(defaultFillText),
		backgroundText: defaultBackgroundText,
		backgroundRunes: []rune(defaultBackgroundText),
		stroke:         true,
		autoResize:     true,
		IsLooping:      true,
	}
}

func mockCanvas(width, height int) *Canvas {
	c := newCanvas(width, height)
	c.output = new(bytes.Buffer)
	return c
}

func (c *Canvas) render() {
	if c.disabled {
		return
	}
	var output strings.Builder
	output.Grow(c.Width * c.Height * 20)

	for y := range c.buffer {
		lineLen := 0
		for x := range c.buffer[y] {
			if lineLen >= c.termWidth {
				continue
			}
			if c.buffer[y][x].char == 0 {
				lineLen += c.writeBackgroundCell(&output)
			} else {
				lineLen += c.writeRenderCell(&output, c.buffer[y][x])
			}
		}
		output.WriteString(resetStyleStr)
		writePadding(&output, lineLen, c.termWidth, ' ')
		if y < len(c.buffer)-1 {
			output.WriteString("\r\n")
			c.lastStyle = style{}
		}
	}
	// Clear garbage outside the canvas
	if c.Height < c.termHeight {
		for range c.termHeight - c.Height - 1 {
			output.WriteString(clearLineStr)
		}
	}
	if c.save {
		c.exportCanvasToPNG(output.String())
		c.save = false
	}
	if c.frames != nil {
		c.recordFrame(output.String())
	}
	c.Framecount++
	c.reset()
	fmt.Fprint(c.output, output.String())
	c.renderDebug()
}

func (c *Canvas) reset() {
	c.originX = 0
	c.originY = 0
	c.rotationAngle = 0
	c.scale = 1
	c.strokeIndex = 0
	c.backgroundIndex = 0
	c.lastStyle = style{}
}

func (c *Canvas) resize(width, height int) {
	newWidth := c.widthWithPadding(width)
	newHeight := height
	newBuffer := newFrame(newWidth, newHeight)

	minWidth := min(c.Width, newWidth)
	minHeight := min(c.Height, newHeight)

	for y := range minHeight {
		for x := range minWidth {
			newBuffer[y][x] = c.buffer[y][x]
		}
	}

	c.Width = newWidth
	c.Height = newHeight
	c.buffer = newBuffer
	c.captureResize(width, height)
}

func (c *Canvas) char(char rune, x, y int) {
	if !c.stroke && !c.isFilling {
		return
	}
	c.write(cell{
		char:       char,
		foreground: c.strokeFg,
		background: c.strokeBg,
	}, x, y, 1)
}

func (c *Canvas) write(cll cell, x, y int, minBlockSize int) {
	var destX, destY, blockSize int

	if c.rotationAngle == 0 && c.scale == 1 {
		destX = x + c.originX
		destY = y + c.originY
		blockSize = minBlockSize
	} else {
		scaledX := float64(x) * c.scale
		scaledY := float64(y) * c.scale

		sinA, cosA := math.Sincos(c.rotationAngle)
		rotatedX := scaledX*cosA - scaledY*sinA
		rotatedY := scaledX*sinA + scaledY*cosA

		destX = int(math.Round(rotatedX)) + c.originX
		destY = int(math.Round(rotatedY)) + c.originY
		blockSize = max(int(math.Round(c.scale)), minBlockSize)
	}

	for dy := range blockSize {
		for dx := range blockSize {
			sx := destX + dx
			sy := destY + dy
			if c.outOfBounds(sx, sy) {
				continue
			}
			c.buffer[sy][sx] = cll

			if c.isFilling {
				c.forceFill(sx, sy, c.buffer[sy][sx])
			}
		}
	}
}

func (c *Canvas) forceFill(sx, sy int, cll cell) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			px := sx + dx
			py := sy + dy
			if c.outOfBounds(px, py) || (dx == 0 && dy == 0) {
				continue
			}
			if c.buffer[py][px].char == 0 {
				if c.inBoundsAndMatch(px+dx, py+dy, cll.char) && c.inBoundsAndMatch(px-dx, py-dy, cll.char) {
					c.buffer[py][px] = cll
				}
			}
		}
	}
}

func (c *Canvas) inBoundsAndMatch(x, y int, char rune) bool {
	return !c.outOfBounds(x, y) && c.buffer[y][x].char == char
}

// writeRenderCell writes a cell directly to the builder and returns the display width.
func (c *Canvas) writeRenderCell(b *strings.Builder, cll cell) int {
	currentStyle := style{foreground: cll.foreground, background: cll.background}

	if !c.lastStyle.equals(currentStyle) {
		c.lastStyle = currentStyle
		currentStyle.writeTo(b)
	}

	if cll.padChar != 0 {
		b.WriteRune(cll.char)
		b.WriteRune(cll.padChar)
		return 2
	}
	switch c.cellMode {
	case cellModeCustom:
		b.WriteRune(cll.char)
		b.WriteRune(c.cellModeRune)
		return 2
	case cellModeDouble:
		b.WriteRune(cll.char)
		b.WriteRune(cll.char)
		return 2
	default:
		b.WriteRune(cll.char)
		return 1
	}
}

// writeBackgroundCell writes a background cell directly to the builder and returns the display width.
func (c *Canvas) writeBackgroundCell(b *strings.Builder) int {
	currentStyle := style{foreground: c.backgroundFg, background: c.backgroundBg}

	if !c.lastStyle.equals(currentStyle) {
		c.lastStyle = currentStyle
		currentStyle.writeTo(b)
	}

	switch c.cellMode {
	case cellModeCustom:
		b.WriteRune(c.nextBackgroundRune())
		b.WriteRune(c.cellModeRune)
		return 2
	case cellModeDouble:
		next := c.nextBackgroundRune()
		b.WriteRune(next)
		b.WriteRune(next)
		return 2
	default:
		b.WriteRune(c.nextBackgroundRune())
		return 1
	}
}

func (c *Canvas) widthWithPadding(w int) int {
	if c.cellMode.enabled() {
		return w / 2
	}
	return w
}

func (c *Canvas) toggleFill() {
	c.isFilling = !c.isFilling
	c.strokeText, c.fillText = c.fillText, c.strokeText
	c.strokeRunes, c.fillRunes = c.fillRunes, c.strokeRunes
	c.strokeBg, c.fillBg = c.fillBg, c.strokeBg
	c.strokeFg, c.fillFg = c.fillFg, c.strokeFg
}

func (c *Canvas) nextBackgroundRune() rune {
	r := c.backgroundRunes[c.backgroundIndex]
	c.backgroundIndex = (c.backgroundIndex + 1) % len(c.backgroundRunes)
	return r
}

func (c *Canvas) nextStrokeRune() rune {
	r := c.strokeRunes[c.strokeIndex]
	c.strokeIndex = (c.strokeIndex + 1) % len(c.strokeRunes)
	return r
}

func (c *Canvas) outOfBounds(x, y int) bool {
	return x < 0 || y < 0 || x > c.Width-1 || y > c.Height-1
}

func (c *Canvas) setMousePostion(x, y int) {
	c.MouseX = x
	if c.cellMode.enabled() {
		c.MouseX = x / 2
	}
	c.MouseY = y
}
