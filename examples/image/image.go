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
	c.Image(img, 0, 0, c.Width, c.Height)

	fullCanvas := c.Get(0, 0, c.Width, c.Height)

	c.Image(fullCanvas, c.Width/2, 0, c.Width/2, c.Height)
}
