package main

import (
	"context"

	"github.com/emprcl/runal"
)

func main() {
	runal.Run(context.Background(), setup, draw)
}

var img runal.Image

const (
	glitchPoints = 80
	glitchSize   = 2
	glitchLines  = 5
)

func setup(c *runal.Canvas) {
	img = c.LoadImage("mona-lisa.jpg")
	c.Fps(5)
}

func draw(c *runal.Canvas) {
	c.Image(img, int(c.Random(-2, 2)), int(c.Random(-2, 2)), c.Width, c.Height)

	for range glitchPoints {
		part := c.Get(int(c.Random(0, c.Width)), int(c.Random(0, c.Height)), glitchSize, int(c.Random(1, glitchSize+1)))
		c.Set(int(c.Random(0, c.Width)), int(c.Random(0, c.Height)), part)
	}

	for range glitchLines {
		part := c.Get(0, int(c.Random(0, c.Height)), c.Width, 1)
		c.Set(0, int(c.Random(0, c.Height)), part)
	}
}
