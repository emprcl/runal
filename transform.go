package runal

// Translate offsets the drawing context by (x, y).
func (c *Canvas) Translate(x, y int) {
	c.originX = x
	c.originY = y
}

// Rotate rotates the drawing context by the given angle in radians.
func (c *Canvas) Rotate(angle float64) {
	c.rotationAngle = angle
}

// Scale scales the drawing context by the given factor.
func (c *Canvas) Scale(scale float64) {
	c.scale = scale
}
