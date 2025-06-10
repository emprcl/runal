package runal

func (c *Canvas) Translate(x, y int) {
	c.originX = x
	c.originY = y
}

// TODO: check behavior with scale.js
func (c *Canvas) Rotate(angle float64) {
	c.rotationAngle = angle
}

func (c *Canvas) Scale(scale float64) {
	c.scale = scale
}
