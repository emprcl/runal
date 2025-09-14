package main

import (
	"context"
	"time"

	"github.com/emprcl/runal"
)

const (
	duration = 1
	scale    = 1
)

func main() {
	runal.Run(context.Background(), setup, draw, runal.WithOnKey(onKey))
}

func setup(c *runal.Canvas) {}

func draw(c *runal.Canvas) {
	c.Clear()
	c.StrokeText("THIS IS A 1D NOISE LOOP EXAMPLE")
	theta := c.LoopAngle(duration)

	for x := 0; x < c.Width; x++ {
		noise := c.NoiseLoop1D(theta, 0.1, x*scale)
		y := c.Map(noise, 0, 1, 0, float64(c.Height))
		c.Point(x, int(y))
	}
}

func onKey(c *runal.Canvas, e runal.KeyEvent) {
	if e.Key == "space" {
		c.NoiseSeed(time.Now().Unix())
	}
}
