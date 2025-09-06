package main

import (
	"context"
	"strconv"
	"time"

	"github.com/emprcl/runal"
)

const (
	duration = 3
	scale    = 0.3
)

func main() {
	runal.Run(context.Background(), setup, draw, runal.WithOnKey(onKey))
}

func setup(c *runal.Canvas) {}

func draw(c *runal.Canvas) {
	c.Clear()
	theta := c.LoopAngle(duration)

	for x := 0; x < c.Width; x++ {
		for y := 0; y < c.Height; y++ {
			noise := c.NoiseLoop2D(
				theta,
				1,
				int(float64(x)*scale),
				int(float64(y)*scale),
			)
			color := c.Map(noise, 0, 1, 232, 255)
			colorStr := strconv.FormatFloat(color, 'f', -1, 64)
			c.Stroke("§", colorStr, "0")
			c.Point(x, y)
		}
	}
}

func onKey(c *runal.Canvas, e runal.KeyEvent) {
	if e.Key == "space" {
		c.NoiseSeed(time.Now().Unix())
	}
}
