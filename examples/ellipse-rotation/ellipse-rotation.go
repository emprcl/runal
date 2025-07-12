package main

import (
	"context"

	"github.com/emprcl/runal"
)

func main() {
	runal.Run(context.Background(), setup, draw, nil, nil)
}

func setup(c *runal.Canvas) {
	c.CellPadding(" ")
}

func draw(c *runal.Canvas) {
	c.Clear()
	c.Fill("ellipse", "#ffffff", "#000000")
	c.Translate(c.Width/2, c.Height/2)
	c.Rotate(float64(c.Framecount) * 0.008)
	c.Ellipse(0, 0, 15, 5)
}
