package canvas

import (
	"fmt"
)

const maxDebugBufferSize = 100

func (c *Canvas) Debug(messages ...any) {
	msg := fmt.Sprint(messages...)
	c.debugBuffer = append(c.debugBuffer, msg)

	if len(c.debugBuffer) > min(maxDebugBufferSize, c.termHeight) {
		c.debugBuffer = c.debugBuffer[len(c.debugBuffer)-min(maxDebugBufferSize, c.termHeight):]
	}
}
