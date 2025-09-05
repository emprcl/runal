package main

import (
	"context"
	"math/rand"

	"github.com/emprcl/runal"
)

func main() {
	runal.Run(context.Background(), setup, draw, onKey, onMouse)
}

func setup(c *runal.Canvas) {
	c.NoLoop()
}

func draw(c *runal.Canvas) {
	c.Clear()
	for i := 0; i < c.Width; i++ {
		for j := 0; j < c.Height; j++ {
			if rand.Intn(100) < 80 {
				c.Text(".", i, j)
			}
		}
	}
}

func onKey(c *runal.Canvas, e runal.KeyEvent) {
	if e.Key == "space" {
		c.Redraw()
	}
}

func onMouse(c *runal.Canvas, e runal.MouseEvent) {
	if e.Type == "click" && e.Button == "left" {
		c.Redraw()
	}
}
