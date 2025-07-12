package main

import (
	"context"
	"math"

	"github.com/emprcl/runal"
)

func main() {
	runal.Run(context.Background(), setup, draw, nil, nil)
}

func setup(c *runal.Canvas) {}

func draw(c *runal.Canvas) {
	c.Clear()
	y1 := (math.Sin((float64(c.Framecount)*0.2+1000)*0.2)/2 + 0.5) * float64(c.Height) * 0.8
	x1 := (math.Cos((float64(c.Framecount)*0.2+1000)*0.2)/2 + 0.5) * float64(c.Width) * 0.8
	y2 := (math.Sin(float64(c.Framecount)*0.1)/2 + 0.5) * float64(c.Height) * 0.8
	x2 := (math.Cos(float64(c.Framecount)*0.1)/2 + 0.5) * float64(c.Width) * 0.8
	c.Line(int(x1), int(y1), int(x2), int(y2))
}
