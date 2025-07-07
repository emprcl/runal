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
}

// NoLoop disables automatic canvas redrawing.
func (c *Canvas) NoLoop() {
	c.IsLooping = false
}

// Redraw triggers a manual rendering pass when
// canvas is not redrawing automatically.
func (c *Canvas) Redraw() {
	c.lastFrame = ""
}

// DisableRendering disables all rendering updates.
// Used when an error is rendered.
func (c *Canvas) DisableRendering() {
	c.disabled = true
}

// CellPadding sets a character used for cell spacing between elements.
func (c *Canvas) CellPadding(char string) {
	prev := c.cellPadding
	c.cellPadding = cellPaddingCustom
	c.cellPaddingRune = rune(char[0])

	if prev.enabled() {
		c.resize(c.Width*2, c.Height)
	} else {
		c.resize(c.Width, c.Height)
	}
}

// CellPaddingDouble makes every cell duplicated.
func (c *Canvas) CellPaddingDouble() {
	prev := c.cellPadding
	c.cellPadding = cellPaddingDouble
	if prev.enabled() {
		c.resize(c.Width*2, c.Height)
	} else {
		c.resize(c.Width, c.Height)
	}
}

// Fps sets the rendering framerate in frames per second.
func (c *Canvas) Fps(fps int) {
	c.fps = fps
	c.bus <- newFPSEvent(fps)
}
