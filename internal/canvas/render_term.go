//go:build !js

package canvas

import (
	"fmt"
	"strings"
)

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

var debugStyle = style{
	background: color("9"),
	foreground: color("15"),
}

func (c *Canvas) renderDebug() {
	if len(c.debugBuffer) == 0 {
		return
	}
	resetCursorPosition()
	for y, msg := range c.debugBuffer {
		if y >= c.termHeight-1 {
			continue
		}
		fmt.Fprint(c.output, debugStyle.render(fmt.Sprintf("%s\r\n", msg)))
	}
}
