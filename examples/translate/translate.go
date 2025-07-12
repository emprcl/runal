package main

import (
	"context"

	"github.com/emprcl/runal"
)

func main() {
	runal.Run(context.Background(), setup, draw, nil, nil)
}

func setup(c *runal.Canvas) {
	c.NoLoop()
	c.CellPadding(" ")
}

func draw(c *runal.Canvas) {
	c.Clear()
	c.Circle(0, 0, 5)
	c.Translate(c.Width/2, c.Height/2)
	c.Circle(0, 0, 5)
}
