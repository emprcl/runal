package main

import (
	"context"
	"math"
	"strconv"

	"github.com/emprcl/runal"
)

var x1, y1, x2, y2 int

func main() {
	runal.Run(context.Background(), setup, draw, nil, onMouse)
}

func setup(c *runal.Canvas) {
	c.Stroke("make yourselves sheep and the wolves will eat you ", "#000000", "#000000")
	c.NoLoop()
}

func draw(c *runal.Canvas) {
	if x1 == 0 && y1 == 0 {
		return
	}
	c.Line(x1, y1, x2, y2)
}

func onMouse(c *runal.Canvas, e runal.MouseEvent) {
	// set stroke color to one of the ansi colors, but not black (1)
	c.StrokeFg(strconv.Itoa(int(math.Ceil(c.Random(1, 255)))))
	x1 = x2
	y1 = y2
	if e.Button == "left" {
		x2 = c.MouseX
		y2 = c.MouseY
	}
	c.Redraw()
}
