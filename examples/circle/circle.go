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
	c.CellModeCustom(".")
}

func draw(c *runal.Canvas) {
	c.Clear()

	radius1 := ((math.Sin(float64(c.Framecount)*0.1)*0.5 + 0.5) * float64(c.Width)) / 2
	radius2 := ((math.Sin(float64(c.Framecount)*0.2)*0.5 + 0.5) * float64(c.Width)) / 3

	c.Stroke("COUCOU", "#ffffff", "#000000")
	c.Fill("i", "#ffffff", "#000000")
	c.Circle(c.Width/2, c.Height/2, int(radius1))

	c.Stroke("C", "#ffffff", "#000000")
	c.Fill("vvvv", "#ffffff", "#000000")
	c.Circle(c.Width/2, c.Height/2, int(radius2))
}
