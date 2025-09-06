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
	//c.NoLoop();
	c.CellPadding(" ")
}

func draw(c *runal.Canvas) {
	c.Clear()
	c.Stroke(".", "#ffffff", "#ffffff")
	c.Fill(".", "#ffffff", "#000000")
	c.Translate(c.Width/2, c.Height/2)
	c.Scale(c.Map(math.Sin(float64(c.Framecount)*0.1), -1, 1, 1, 4))
	c.Rotate(float64(c.Framecount) * 0.008)
	c.Circle(0, 0, 5)
	c.Circle(10, 10, 5)
}
