package canvas

import "math"

// depthBuffer holds a per-cell depth value used for hidden-surface removal
// when drawing solid 3D shapes. Values are NDC depth; smaller is closer.
type depthBuffer [][]float64

func newDepthBuffer(width, height int) depthBuffer {
	values := make([]float64, width*height)
	for i := range values {
		values[i] = math.Inf(1)
	}
	buff := make([][]float64, height)
	for i := range buff {
		buff[i] = values[i*width : (i+1)*width]
	}
	return buff
}

// clearDepth resets every cell to the far value. Called whenever the color
// buffer is cleared so the two stay in sync.
func (c *Canvas) clearDepth() {
	for y := range c.depth {
		for x := range c.depth[y] {
			c.depth[y][x] = math.Inf(1)
		}
	}
}

// writeDepth writes a cell only if z is closer than what is already there,
// updating the depth buffer. This is the funnel used by every solid 3D shape.
func (c *Canvas) writeDepth(cll cell, x, y int, z float64) {
	if c.outOfBounds(x, y) {
		return
	}
	if z <= c.depth[y][x] {
		c.depth[y][x] = z
		c.buffer[y][x] = cll
	}
}
