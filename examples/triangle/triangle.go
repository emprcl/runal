package main

import (
	"context"

	"github.com/emprcl/runal"
)

func main() {
	runal.Run(context.Background(), setup, draw)
}

func setup(c *runal.Canvas) {
	c.CellPadding(" ")
}

func draw(c *runal.Canvas) {
	c.Clear()
	c.Stroke(".", "#ffffff", "#000000")
	c.Fill("triangle", "#ffffff", "#000000")
	c.Translate(c.Width/2, c.Height/2)
	c.Rotate(float64(c.Framecount) * 0.008)
	c.Scale(1)
	c.Triangle(5, 5, 15, 15, 2, 15)
}
