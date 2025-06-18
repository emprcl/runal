package runal

import "math"

// Dist returns the Euclidean distance between two points.
func (c *Canvas) Dist(x1, y1, x2, y2 int) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(float64(dx*dx + dy*dy))
}

// Map linearly maps a value from one range to another.
func (c *Canvas) Map(value, inputStart, inputEnd, outputStart, outputEnd float64) float64 {
	return outputStart + ((outputEnd-outputStart)/(inputEnd-inputStart))*(value-inputStart)
}
