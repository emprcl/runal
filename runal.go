package runal

// Size sets the dimensions of the canvas.
func (c *Canvas) Size(w, h int) {
	c.autoResize = false
	if c.cellPadding.enabled() {
		c.resize(w*2, h)
	} else {
		c.resize(w, h)
	}
}

// Clear clears the canvas contents.
func (c *Canvas) Clear() {
	c.clear = true
}

// Loop enables continuous redrawing of the canvas.
func (c *Canvas) Loop() {
	c.IsLooping = true
	c.bus <- newStartEvent()
}

// NoLoop disables automatic canvas redrawing.
func (c *Canvas) NoLoop() {
	c.IsLooping = false
	c.bus <- newStopEvent()
}

// Redraw triggers a manual rendering pass.
func (c *Canvas) Redraw() {
	c.bus <- newRenderEvent()
}

// DisableRendering disables all rendering updates.
// Used when an error is rendered.
func (c *Canvas) DisableRendering() {
	c.disabled = true
	c.NoLoop()
}

// CellPadding sets a character used for cell spacing between elements.
func (c *Canvas) CellPadding(char string) {
	previousValue := c.cellPadding.enabled()
	c.cellPadding = cellPaddingCustom
	c.cellPaddingRune = rune(char[0])

	if c.autoResize && !previousValue {
		c.resize(c.Width, c.Height)
	} else if !previousValue {
		c.resize(c.Width*2, c.Height)
	}
}

// CellPaddingDouble makes every cell duplicated.
func (c *Canvas) CellPaddingDouble() {
	previousValue := c.cellPadding.enabled()
	c.cellPadding = cellPaddingDouble

	if c.autoResize && !previousValue {
		c.resize(c.Width, c.Height)
	} else if !previousValue {
		c.resize(c.Width*2, c.Height)
	}
}

// Fps sets the rendering framerate in frames per second.
func (c *Canvas) Fps(fps int) {
	c.fps = fps
	c.bus <- newFPSEvent(fps)
}
