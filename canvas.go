package runal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

const (
	defaultPaddingRune    = ' '
	defaultStrokeText     = "."
	defaultFillText       = " "
	defaultBackgroundText = " "
)

type Canvas struct {
	buffer buffer
	output io.Writer

	strokeFg, strokeBg                   lipgloss.Color
	fillFg, fillBg                       lipgloss.Color
	backgroundFg, backgroundBg           lipgloss.Color
	strokeText, fillText, backgroundText string

	bus chan event

	Width, Height int
	Framecount    int

	strokeIndex, backgroundIndex int
	termWidth, termHeight        int
	cellPaddingRune              rune
	cellPadding                  bool
	fill                         bool
	clear                        bool
	autoResize                   bool
	disabled                     bool
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
		strokeFg:        lipgloss.Color("#ffffff"),
		strokeBg:        lipgloss.Color("#000000"),
		fillFg:          lipgloss.Color("#ffffff"),
		fillBg:          lipgloss.Color("#000000"),
		backgroundFg:    lipgloss.Color("#ffffff"),
		backgroundBg:    lipgloss.Color("#000000"),
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
			line.WriteString("\n")
		}
		output.WriteString(line.String())
	}
	c.Framecount++
	c.clear = false
	fmt.Fprint(c.output, output.String())
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

	for y := 0; y < minHeight; y++ {
		for x := 0; x < minWidth; x++ {
			newBuffer[y][x] = c.buffer[y][x]
		}
	}

	c.Width = newWidth
	c.Height = newHeight
	c.buffer = newBuffer
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

func (c *Canvas) OutOfBounds(x, y int) bool {
	return x < 0 || y < 0 || x > c.Width-1 || y > c.Height-1
}
