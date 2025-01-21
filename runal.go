package runal

func (c *Canvas) Size(w, h int) {
	c.autoResize = false
	if c.widthPadding {
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

func (c *Canvas) DisableRendering() {
	c.disabled = true
	c.NoLoop()
}

func (c *Canvas) WidthPadding(char string) {
	previousValue := c.widthPadding
	c.widthPadding = true
	c.widthPaddingChar = rune(char[0])

	if c.autoResize && !previousValue {
		c.resize(c.Width, c.Height)
	} else if !previousValue {
		c.resize(c.Width*2, c.Height)
	}
}

func (c *Canvas) Fps(fps int) {
	c.bus <- newFPSEvent(fps)
}
