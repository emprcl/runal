package runal

func (c *Canvas) Background(text, fg, bg string) {
	c.BackgroundText(text)
	c.BackgroundBg(bg)
	c.BackgroundFg(fg)
}

func (c *Canvas) BackgroundText(text string) {
	c.backgroundIndex = 0
	if len(text) == 0 {
		c.backgroundText = defaultBackgroundText
		return
	}
	c.backgroundText = text
}

func (c *Canvas) BackgroundFg(fg string) {
	c.backgroundFg = color(fg)
}

func (c *Canvas) BackgroundBg(bg string) {
	c.backgroundBg = color(bg)
}

func (c *Canvas) Fill(text, fg, bg string) {
	c.FillText(text)
	c.FillBg(bg)
	c.FillFg(fg)
}

func (c *Canvas) FillText(text string) {
	c.fill = true
	if len(text) == 0 {
		c.fillText = defaultFillText
		return
	}
	c.fillText = text
}

func (c *Canvas) FillFg(fg string) {
	c.fill = true
	c.fillFg = color(fg)
}

func (c *Canvas) FillBg(bg string) {
	c.fill = true
	c.fillBg = color(bg)
}

func (c *Canvas) Stroke(text, fg, bg string) {
	c.StrokeText(text)
	c.StrokeBg(bg)
	c.StrokeFg(fg)
}

func (c *Canvas) StrokeText(text string) {
	c.strokeIndex = 0
	if len(text) == 0 {
		c.strokeText = defaultStrokeText
	}
	c.strokeText = text
}

func (c *Canvas) StrokeFg(fg string) {
	c.strokeFg = color(fg)
}

func (c *Canvas) StrokeBg(bg string) {
	c.strokeBg = color(bg)
}

func (c *Canvas) NoFill() {
	c.fill = false
}
