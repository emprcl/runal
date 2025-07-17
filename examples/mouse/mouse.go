package main

import (
	"context"
	"strconv"

	"github.com/emprcl/runal"
)

func main() {
	runal.Run(context.Background(), setup, draw, onKey, nil)
}

func setup(c *runal.Canvas) {}

func draw(c *runal.Canvas) {
	c.Clear()
	c.Circle(c.MouseX, c.MouseY, 5)
}

func onKey(c *runal.Canvas, e runal.KeyEvent) {
	if e.Key == "space" {
		color := c.Random(0, 255)
		colorStr := strconv.FormatFloat(color, 'f', -1, 64)
		c.BackgroundBg(colorStr)
	}
}

func onMouse(c *runal.Canvas, e runal.MouseEvent) {
	if e.Button == "left" {
		color := c.Random(0, 255)
		colorStr := strconv.FormatFloat(color, 'f', -1, 64)
		c.BackgroundBg(colorStr)
	}
}
