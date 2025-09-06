package main

import (
	"context"
	"math"
	"time"

	"github.com/emprcl/runal"
)

const (
	duration = 5
	margin   = 3
)

func main() {
	runal.Run(context.Background(), setup, draw, runal.WithOnKey(onKey))
}

func setup(c *runal.Canvas) {
	c.Size(40, 21)
	c.BackgroundBg("197")
	c.SaveCanvasToGIF("canvas.gif", duration)
}

func draw(c *runal.Canvas) {
	c.Clear()
	c.Stroke("RUNAL", "255", "197")
	theta := c.LoopAngle(duration)

	for y := margin; y < c.Height-margin; y++ {
		x := c.Map(math.Sin(theta+float64(y)), -1, 1, margin, float64(c.Width-margin))
		c.Point(int(x), y)
	}
}

func onKey(c *runal.Canvas, e runal.KeyEvent) {
	if e.Key == "space" {
		c.NoiseSeed(time.Now().Unix())
	}
}
