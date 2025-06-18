package runal

import "math"

// Bezier draws a Bézier curve using four control points.
func (c *Canvas) Bezier(x1, y1, x2, y2, x3, y3, x4, y4 int) {
	steps := 50

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		u := 1 - t

		x := u*u*u*float64(x1) +
			3*u*u*t*float64(x2) +
			3*u*t*t*float64(x3) +
			t*t*t*float64(x4)

		y := u*u*u*float64(y1) +
			3*u*u*t*float64(y2) +
			3*u*t*t*float64(y3) +
			t*t*t*float64(y4)

		c.Point(int(math.Round(x)), int(math.Round(y)))
	}
}
