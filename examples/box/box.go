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
	c.Light(-0.5, 0.8, 1)
}

func draw(c *runal.Canvas) {
	c.Clear()
	c.Fill("", "#00ffcc", "#000000")
	c.Translate(c.Width/2, c.Height/2)
	c.RotateX(float64(c.Framecount) * 0.02)
	c.RotateY(float64(c.Framecount) * 0.03)
	c.RotateZ(float64(c.Framecount) * 0.01)
	c.Box(20, 20, 20)
}
