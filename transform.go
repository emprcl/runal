package runal

import (
	"math"
)

func (c *Canvas) Translate(x, y int) {
	c.originX = x
	c.originY = y
}

// TODO: fix weird empty characters after rotation (ex: scale.js)
func (c *Canvas) Rotate(angle float64) {
	rotated := newBuffer(c.Width, c.Height)

	radians := angle * math.Pi / 180.0

	// TODO: change origin of rotation
	centerX := float64(c.Width-1) / 2.0
	centerY := float64(c.Height-1) / 2.0

	for y := range c.Height {
		for x := range c.Width {
			xPrime := float64(x) - centerX
			yPrime := float64(y) - centerY
			xRot := xPrime*math.Cos(radians) - yPrime*math.Sin(radians)
			yRot := xPrime*math.Sin(radians) + yPrime*math.Cos(radians)
			xFinal := int(math.Round(xRot + centerX))
			yFinal := int(math.Round(yRot + centerY))
			if c.OutOfBounds(xFinal, yFinal) {
				continue
			}
			rotated[yFinal][xFinal] = c.buffer[y][x]
		}
	}

	c.buffer = rotated
}

func (c *Canvas) Scale(scale float64) {
	if scale < 0 {
		return
	}

	scaled := newBuffer(c.Width, c.Height)

	for y := range c.Height {
		for x := range c.Width {
			if c.buffer[y][x] == "" {
				continue
			}
			destX := int(float64(x) * scale)
			destY := int(float64(y) * scale)
			if c.OutOfBounds(destX, destY) {
				continue
			}

			for fillY := destY; fillY < destY+int(math.Round(scale)); fillY++ {
				for fillX := destX; fillX < destX+int(math.Round(scale)); fillX++ {
					if c.OutOfBounds(fillX, fillY) {
						continue
					}
					scaled[fillY][fillX] = c.buffer[y][x]
				}
			}
		}
	}

	c.buffer = scaled
}
