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
	c.Stroke("0", "#ffffff", "#000000")
	c.Bezier(10, 10, 20, 0, 30, 20, 40, 10)
}
