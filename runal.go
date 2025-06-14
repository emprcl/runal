package runal

func (c *Canvas) Size(w, h int) {
	c.autoResize = false
	if c.cellPadding {
		c.resize(w*2, h)
	} else {
		c.resize(w, h)
	}
}

func (c *Canvas) Clear() {
	c.clear = true
}

func (c *Canvas) NoLoop() {
	c.bus <- newStopEvent()
}

func (c *Canvas) Redraw() {
	c.bus <- newRenderEvent()
}

func (c *Canvas) DisableRendering() {
	c.disabled = true
	c.NoLoop()
}

func (c *Canvas) CellPadding(char string) {
	previousValue := c.cellPadding
	c.cellPadding = true
	c.cellPaddingRune = rune(char[0])

	if c.autoResize && !previousValue {
		c.resize(c.Width, c.Height)
	} else if !previousValue {
		c.resize(c.Width*2, c.Height)
	}
}

func (c *Canvas) FPS(fps int) {
	c.bus <- newFPSEvent(fps)
}
