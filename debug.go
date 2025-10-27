package runal

import (
	"fmt"
)

const maxDebugBufferSize = 100

var debugStyle = style{
	background: color("9"),
	foreground: color("15"),
}

func (c *Canvas) Debug(messages ...any) {
	msg := fmt.Sprint(messages...)
	c.debugBuffer = append(c.debugBuffer, msg)

	if len(c.debugBuffer) > min(maxDebugBufferSize, c.termHeight) {
		c.debugBuffer = c.debugBuffer[len(c.debugBuffer)-min(maxDebugBufferSize, c.termHeight):]
	}
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
