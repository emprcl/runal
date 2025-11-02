package runal

// Size sets the dimensions of the canvas.
func (c *Canvas) Size(w, h int) {
	c.autoResize = false
	if c.cellMode.enabled() {
		c.resize(w*2, h)
	} else {
		c.resize(w, h)
	}
}

// Clear clears the canvas contents.
func (c *Canvas) Clear() {
	c.buffer = newFrame(c.Width, c.Height)
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

// CellModeCustom sets a specific character used for cell spacing between elements.
func (c *Canvas) CellModeCustom(char string) {
	prev := c.cellMode
	c.cellMode = cellModeCustom
	c.cellModeRune = []rune(char)[0]

	if prev.enabled() {
		c.resize(c.Width*2, c.Height)
	} else {
		c.resize(c.Width, c.Height)
	}
}

// CellModeDouble makes every cell duplicated.
func (c *Canvas) CellModeDouble() {
	prev := c.cellMode
	c.cellMode = cellModeDouble
	if prev.enabled() {
		c.resize(c.Width*2, c.Height)
	} else {
		c.resize(c.Width, c.Height)
	}
}

// CellModeDefault disables cell mode.
func (c *Canvas) CellModeDefault() {
	prev := c.cellMode
	c.cellMode = cellModeDisabled
	c.cellModeRune = 0
	if prev.enabled() {
		c.resize(c.Width*2, c.Height)
	} else {
		c.resize(c.Width, c.Height)
	}
}

// DEPRECATED: Use CellModeCustom() instead.
func (c *Canvas) CellPadding(char string) {
	c.CellModeCustom(char)
}

// DEPRECATED: Use CellModeDouble() instead.
func (c *Canvas) CellPaddingDouble(char string) {
	c.CellModeDouble()
}

// DEPRECATED: Use CellModeDefault() instead.
func (c *Canvas) NoCellPadding() {
	c.CellModeDefault()
}

// Fps sets the rendering framerate in frames per second.
func (c *Canvas) Fps(fps int) {
	c.fps = fps
	c.bus <- newFPSEvent(fps)
}

// Exit ends the program execution.
func (c *Canvas) Exit() {
	c.bus <- newExitEvent()
}
