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
	c.SavedCanvasFontSize(24)
	c.CellPaddingDouble()
}

func draw(c *runal.Canvas) {
	for i := 0; i < c.Width; i++ {
		for j := 0; j < c.Height; j++ {
			color := c.Map(
				c.Noise2D(
					float64(i)*0.009+float64(c.Framecount)/1000,
					float64(j)*0.009+float64(c.Framecount)/1000,
				),
				0,
				1,
				150,
				231,
			)
			colorStr := strconv.FormatFloat(color, 'f', -1, 64)
			c.Stroke("§", colorStr, "#000000")
			c.Point(i, j)
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

func onMouseClick(c *runal.Canvas, e runal.MouseEvent) {
	c.NoiseSeed(time.Now().Unix())
	c.Redraw()
}
