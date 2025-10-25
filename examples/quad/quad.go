package main

import (
	"context"
	"math"

	"github.com/emprcl/runal"
)

func main() {
	runal.Run(context.Background(), setup, draw)
}

func setup(c *runal.Canvas) {
	c.CellModeCustom(" ")
}

func draw(c *runal.Canvas) {
	c.Clear()
	c.Stroke(".", "#ffffff", "#555555")
	c.Fill("quad", "#ffffff", "#000000")
	c.Quad(
		int(c.Map(math.Sin(float64(c.Framecount)*0.1), -1, 1, 1, 35)),
		1,
		int(c.Map(math.Cos(float64(c.Framecount)*0.1), -1, 1, 1, 35)),
		3,
		16,
		12,
		2,
		18,
	)
}
