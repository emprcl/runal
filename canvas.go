package runal

import (
	"bytes"
	"image"
	"math"
	"math/rand"
	"os"
	"strings"

	perlin "github.com/aquilax/go-perlin"
	"github.com/charmbracelet/x/ansi"
	"github.com/rahji/termenv"
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
	output  *termenv.Output
	capture *ansitoimage.Converter
	frames  []image.Image
	videoFormat
	noise  *perlin.Perlin
	random *rand.Rand

	state *state

	strokeFg, strokeBg         string
	fillFg, fillBg             string
	backgroundFg, backgroundBg string
	lastStyle                  style

	strokeText, fillText, backgroundText string

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
		output:         termenv.NewOutput(os.Stdout),
		capture:        newCapture(width, height),
		noise:          newNoise(),
		random:         newRandom(),
		scale:          1,
		strokeFg:       "#ffffff",
		strokeBg:       "#000000",
		fillFg:         "#ffffff",
		fillBg:         "#000000",
		backgroundFg:   "#ffffff",
		backgroundBg:   "#000000",
		strokeText:     defaultStrokeText,
		fillText:       defaultFillText,
		backgroundText: defaultBackgroundText,
		stroke:         true,
		autoResize:     true,
		IsLooping:      true,
	}
}

func mockCanvas(width, height int) *Canvas {
	c := newCanvas(width, height)
	c.output = termenv.NewOutput(new(bytes.Buffer))
	return c
}

func (c *Canvas) render() {
	if c.disabled {
		return
	}
	var output strings.Builder

	for y := range c.buffer {
		var line strings.Builder
		lineLen := 0
		for x := range c.buffer[y] {
			var add string
			if c.buffer[y][x].char == 0 {
				bgCell := c.backgroundCell()
				bgCellSize := ansi.StringWidth(bgCell)
				add = bgCell
				lineLen += bgCellSize
			} else {
				add = c.renderCell(c.buffer[y][x])
				lineLen += ansi.StringWidth(add)
			}
			if lineLen < c.termWidth {
				line.WriteString(add)
			}
		}
		line.WriteString(resetStyleSequence())
		forcePadding(&line, lineLen, c.termWidth, ' ')
		if y < len(c.buffer)-1 {
			line.WriteString("\r\n")
			c.lastStyle = style{}
		}
		output.WriteString(line.String())
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
	_, _ = c.output.WriteString(output.String())
}

func (c *Canvas) reset() {
	c.originX = 0
	c.originY = 0
	c.rotationAngle = 0
	c.scale = 1
	c.strokeIndex = 0
	c.backgroundIndex = 0
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
	scaledX := float64(x) * c.scale
	scaledY := float64(y) * c.scale

	rotatedX := scaledX*math.Cos(c.rotationAngle) - scaledY*math.Sin(c.rotationAngle)
	rotatedY := scaledX*math.Sin(c.rotationAngle) + scaledY*math.Cos(c.rotationAngle)

	destX := int(math.Round(rotatedX)) + c.originX
	destY := int(math.Round(rotatedY)) + c.originY

	blockSize := max(int(math.Round(c.scale)), minBlockSize)

	for dy := range blockSize {
		for dx := range blockSize {
			sx := destX + dx
			sy := destY + dy
			if c.outOfBounds(sx, sy) {
				continue
			}
			c.buffer[sy][sx] = cll

			if c.isFilling {
				// NOTE: hack to fill blank spots
				// due to rotation approx.
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

func (c *Canvas) renderCell(cll cell) string {
	currentStyle := style{foreground: cll.foreground, background: cll.background}
	var chars string

	// Text rendering when cell padding mode is enabled.
	if cll.padChar != 0 {
		chars = string([]rune{cll.char, cll.padChar})
	} else {
		switch c.cellMode {
		case cellModeCustom:
			chars = string([]rune{cll.char, c.cellModeRune})
		case cellModeDouble:
			chars = string([]rune{cll.char, cll.char})
		default:
			chars = string(cll.char)
		}
	}

	if !c.lastStyle.equals(currentStyle) {
		c.lastStyle = currentStyle
		return currentStyle.termStyle(c.output).StyledWithoutReset(chars)
	}

	// if the style hasn't changed since the previous cell, just output
	// the raw text without ANSI codes
	return chars
}

// backgroundCell tries to eliminate redundant ANSI codes in the way that renderCell does
func (c *Canvas) backgroundCell() string {
	currentStyle := style{foreground: c.backgroundFg, background: c.backgroundBg}
	var chars string
	switch c.cellMode {
	case cellModeCustom:
		chars = string([]rune{c.nextBackgroundRune(), c.cellModeRune})
	case cellModeDouble:
		next := c.nextBackgroundRune()
		chars = string([]rune{next, next})
	default:
		chars = string(c.nextBackgroundRune())
	}

	if !c.lastStyle.equals(currentStyle) {
		c.lastStyle = currentStyle
		return currentStyle.termStyle(c.output).StyledWithoutReset(chars)
	}

	return chars
}

func (c *Canvas) widthWithPadding(w int) int {
	if c.cellMode.enabled() {
		return w / 2
	}
	return w
}

func (c *Canvas) toggleFill() {
	c.isFilling = !c.isFilling
	stroke := c.strokeText
	bg := c.strokeBg
	fg := c.strokeFg
	c.strokeText = c.fillText
	c.strokeBg = c.fillBg
	c.strokeFg = c.fillFg
	c.fillText = stroke
	c.fillBg = bg
	c.fillFg = fg
}

func (c *Canvas) nextBackgroundRune() rune {
	r := []rune(c.backgroundText)[c.backgroundIndex]
	c.backgroundIndex = (c.backgroundIndex + 1) % len(c.backgroundText)
	return r
}

func (c *Canvas) nextStrokeRune() rune {
	r := []rune(c.strokeText)[c.strokeIndex]
	c.strokeIndex = (c.strokeIndex + 1) % len(c.strokeText)
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
