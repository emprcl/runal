package main

import (
	"context"
	"os"

	"github.com/emprcl/runal"
)

func main() {
	runal.Run(context.Background(), setup, draw, onKey, nil)
}

func setup(c *runal.Canvas) {}

func draw(c *runal.Canvas) {
	c.Clear()

	c.Stroke("/", "#fffb00", "#0004ff")
	c.Fill(".", "#0004ff", "#fffb00")
	c.Rect(5, 11, 11, 5)
	c.Text("Stroke: YES", 5, 18)
	c.Text("Fill:   YES", 5, 19)

	c.Stroke("/", "#fffb00", "#0004ff")
	c.NoFill()
	c.Rect(20, 11, 11, 5)
	c.Text("Stroke: YES", 20, 18)
	c.Text("Fill:   NO", 20, 19)

	c.Fill(".", "#0004ff", "#fffb00")
	c.NoStroke()
	c.Rect(35, 11, 11, 5)
	c.Text("Stroke: NO", 35, 18)
	c.Text("Fill:   YES", 35, 19)

	c.NoStroke()
	c.NoFill()
	c.Rect(50, 11, 11, 5)
	c.Text("Stroke: NO", 50, 18)
	c.Text("Fill:   NO", 50, 19)
}

func onKey(c *runal.Canvas, e runal.KeyEvent) {
	os.Exit(0)
}
