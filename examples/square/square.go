package main

import (
	"context"

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
	c.Stroke("BORDER", "#ffffff", "#555555")
	c.Fill("square", "#ffffff", "#000000")
	c.Translate(c.Width/2, c.Height/2)
	c.Rotate(float64(c.Framecount) * 0.08)
	c.Scale(1)
	c.Square(0, 0, 10)
}
