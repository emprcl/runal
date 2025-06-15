package runal

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strings"

	perlin "github.com/aquilax/go-perlin"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	ansitoimage "github.com/pavelpatrin/go-ansi-to-image"
)

const (
	defaultPaddingRune    = ' '
	defaultStrokeText     = "."
	defaultFillText       = " "
	defaultBackgroundText = " "
)

type Canvas struct {
	buffer  buffer
	output  io.Writer
	capture *ansitoimage.Converter
	noise   *perlin.Perlin
	random  *rand.Rand

	state *state

	strokeFg, strokeBg                   lipgloss.Color
	fillFg, fillBg                       lipgloss.Color
	backgroundFg, backgroundBg           lipgloss.Color
	strokeText, fillText, backgroundText string

	saveFilename string

	bus chan event

	Width, Height int
	Framecount    int

	strokeIndex, backgroundIndex int
	termWidth, termHeight        int
	originX, originY             int
	rotationAngle                float64
	scale                        float64

	cellPaddingRune rune
	cellPadding     bool
	fill            bool
	isFilling       bool
	clear           bool
	save            bool
	autoResize      bool
	disabled        bool
}

func newCanvas(width, height int) *Canvas {
	return &Canvas{
		Width:           width,
		Height:          height,
		bus:             make(chan event, 1),
		termWidth:       width,
		termHeight:      height,
		cellPaddingRune: defaultPaddingRune,
		cellPadding:     false,
		buffer:          newBuffer(width, height),
		output:          os.Stdout,
		capture:         newCapture(width, height),
		noise:           newNoise(),
		random:          newRandom(),
		scale:           1,
		strokeFg:        color("#ffffff"),
		strokeBg:        color("#000000"),
		fillFg:          color("#ffffff"),
		fillBg:          color("#000000"),
		backgroundFg:    color("#ffffff"),
		backgroundBg:    color("#000000"),
		strokeText:      defaultStrokeText,
		fillText:        defaultFillText,
		backgroundText:  defaultBackgroundText,
		autoResize:      true,
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
	bgCell := c.backgroundCell()
	bgCellSize := ansi.StringWidth(bgCell)
	for y := range c.buffer {
		var line strings.Builder
		lineLen := 0
		for x := range c.buffer[y] {
			var add string
			if c.buffer[y][x] == "" {
				add = bgCell
				lineLen += bgCellSize
			} else {
				add = c.buffer[y][x]
				lineLen += ansi.StringWidth(add)
			}
			if lineLen < c.termWidth {
				line.WriteString(add)
			}
			if c.clear {
				c.buffer[y][x] = ""
			}
		}
		forcePadding(&line, lineLen, c.termWidth, ' ')
		if y < len(c.buffer)-1 {
			line.WriteString("\r\n")
		}
		output.WriteString(line.String())
	}
	if c.save {
		c.exportCanvasToPNG(output.String())
		c.save = false
	}
	c.Framecount++
	c.reset()
	fmt.Fprint(c.output, output.String())
}

func (c *Canvas) reset() {
	c.clear = false
	c.originX = 0
	c.originY = 0
	c.rotationAngle = 0
	c.scale = 1
}

func (c *Canvas) resize(width, height int) {
	newWidth := c.widthWithPadding(width)
	newHeight := height
	newBuffer := newBuffer(newWidth, newHeight)

	minWidth := c.Width
	if newWidth < c.Width {
		minWidth = newWidth
	}

	minHeight := c.Height
	if newHeight < c.Height {
		minHeight = newHeight
	}

	for y := range minHeight {
		for x := range minWidth {
			newBuffer[y][x] = c.buffer[y][x]
		}
	}

	c.Width = newWidth
	c.Height = newHeight
	c.buffer = newBuffer
	c.capture = newCapture(c.termWidth, c.termHeight)
}

func (c *Canvas) char(char rune, x, y int) {
	formattedChar := c.formatCell(char)
	scaledX := float64(x) * c.scale
	scaledY := float64(y) * c.scale

	radians := c.rotationAngle * math.Pi / 180.0
	rotatedX := scaledX*math.Cos(radians) - scaledY*math.Sin(radians)
	rotatedY := scaledX*math.Sin(radians) + scaledY*math.Cos(radians)

	destX := int(math.Round(rotatedX)) + c.originX
	destY := int(math.Round(rotatedY)) + c.originY

	blockSize := max(int(math.Round(c.scale)), 1)

	for dy := range blockSize {
		for dx := range blockSize {
			sx := destX + dx
			sy := destY + dy
			if c.outOfBounds(sx, sy) {
				continue
			}
			c.buffer[sy][sx] = formattedChar

			if c.isFilling {
				// NOTE: hack to fill blank spots
				// due to rotation approx.
				c.forceFill(sx, sy, formattedChar)
			}
		}
	}
}

func (c *Canvas) forceFill(sx, sy int, char string) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			px := sx + dx
			py := sy + dy
			if c.outOfBounds(px, py) || (dx == 0 && dy == 0) {
				continue
			}
			if c.buffer[py][px] == "" {
				if c.inBoundsAndMatch(px+dx, py+dy, char) && c.inBoundsAndMatch(px-dx, py-dy, char) {
					c.buffer[py][px] = char
				}
			}
		}
	}
}

func (c *Canvas) inBoundsAndMatch(x, y int, char string) bool {
	return !c.outOfBounds(x, y) && c.buffer[y][x] == char
}

func (c *Canvas) style(str string) string {
	return lipgloss.NewStyle().
		Background(c.strokeBg).
		Foreground(c.strokeFg).
		Render(str)
}

func (c *Canvas) formatCell(char rune) string {
	if c.cellPadding {
		return c.style(string([]rune{char, c.cellPaddingRune}))
	}
	return c.style(string(char))
}

func (c *Canvas) backgroundCell() string {
	style := lipgloss.NewStyle().
		Background(c.backgroundBg).
		Foreground(c.backgroundFg)
	if c.cellPadding {
		return style.Render(string([]rune{c.nextBackgroundRune(), c.cellPaddingRune}))
	}
	return style.Render(string(c.cellPaddingRune))
}

func (c *Canvas) widthWithPadding(w int) int {
	if c.cellPadding {
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
	return rune(r)
}

func (c *Canvas) nextStrokeRune() rune {
	r := []rune(c.strokeText)[c.strokeIndex]
	c.strokeIndex = (c.strokeIndex + 1) % len(c.strokeText)
	return rune(r)
}

func (c *Canvas) outOfBounds(x, y int) bool {
	return x < 0 || y < 0 || x > c.Width-1 || y > c.Height-1
}
