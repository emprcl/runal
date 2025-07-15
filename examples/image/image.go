package main

import (
	"context"

	"github.com/emprcl/runal"
)

func main() {
	runal.Run(context.Background(), setup, draw, nil, nil)
}

var img runal.Image

func setup(c *runal.Canvas) {
	img = c.LoadImage("wish.png")
}

func draw(c *runal.Canvas) {
	c.Clear()
	c.Translate(c.Width/2, c.Height/2)
	c.Rotate(float64(c.Framecount) * 0.08)
	c.Image(img, 0, 0, 40, 40)
}
