package runal

import "math"

func (c *Canvas) Distance(x1, y1, x2, y2 int) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(float64(dx*dx + dy*dy))
}

func (c *Canvas) Map(value, inputStart, inputEnd, outputStart, outputEnd float64) float64 {
	return outputStart + ((outputEnd-outputStart)/(inputEnd-inputStart))*(value-inputStart)
}
