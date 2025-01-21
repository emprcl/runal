package runal

import "math"


func (c *Canvas) Rotate(angle float64) {
	rows := len(c.buffer)
	cols := len(c.buffer[0])

	rotated := newBuffer(cols, rows)

	radians := angle * math.Pi / 180.0

	centerX := float64(cols-1) / 2.0
	centerY := float64(rows-1) / 2.0

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			xPrime := float64(x) - centerX
			yPrime := float64(y) - centerY

			xRot := xPrime*math.Cos(radians) - yPrime*math.Sin(radians)
			yRot := xPrime*math.Sin(radians) + yPrime*math.Cos(radians)

			xFinal := int(math.Round(xRot + centerX))
			yFinal := int(math.Round(yRot + centerY))

			if xFinal >= 0 && xFinal < cols && yFinal >= 0 && yFinal < rows {
				rotated[yFinal][xFinal] = c.buffer[y][x]
			}
		}
	}

	c.buffer = rotated
}
