package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/emprcl/runal"
)

func main() {
	runal.Run(context.Background(), setup, draw, runal.WithOnKey(onKey))
}

func setup(c *runal.Canvas) {
	c.Size(82, 41)
	c.Background("/", "237", "#000000")
	c.Stroke("0", "255", "#000000")
	c.CellModeDouble()
}

func draw(c *runal.Canvas) {
	c.Clear()
	if c.Framecount%30 < 15 {
		c.BackgroundText("\\")
	} else {
		c.BackgroundText("/")
	}
	// theta := c.LoopAngle(10)
	for i := 0; i < c.Width; i++ {
		for j := 0; j < c.Height; j++ {
			color := c.Map(
				c.Noise2D(
					float64(c.Framecount)*0.03+float64(i)*0.05,
					float64(c.Framecount)*0.03+float64(j)*0.05,
				),
				0,
				1,
				231,
				255,
			)
			colorStr := strconv.FormatFloat(color, 'f', -1, 64)
			if c.Dist(i, j, c.Width/2, c.Height/2) <= float64(c.Width)/2-4.0 {
				c.StrokeFg(colorStr)
				c.Point(i, j)
			}
		}
	}
}

func onKey(c *runal.Canvas, e runal.KeyEvent) {
	if e.Key == "c" {
		filename := fmt.Sprintf("canvas_%d.png", time.Now().Unix())
		c.SaveCanvasToPNG(filename)
	}
	if e.Key == "space" {
		c.NoiseSeed(time.Now().Unix())
		c.Redraw()
	}
}
