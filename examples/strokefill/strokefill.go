package main

import (
	"context"

	"github.com/emprcl/runal"
)

var xList = []int{1, 21, 31, 41}

func main() {

	runal.Run(context.Background(), setup, draw)
}

func setup(c *runal.Canvas) {}

func draw(c *runal.Canvas) {
	c.Clear()

	// rects

	c.Stroke("1234567890", "#fffb00", "#0004ff")
	c.Fill("1234567890", "#0004ff", "#fffb00")
	c.Rect(1, 1, 11, 5)
	c.Stroke("1234567890", "255", "#ee0000")
	c.Text("RECT", 5, 2)
	c.Text("Stroke: YES", 5, 3)
	c.Text("Fill:   YES", 5, 4)

	c.Stroke("1234567890", "#fffb00", "#0004ff")
	c.NoFill()
	c.Rect(21, 1, 11, 5)
	c.Stroke("1234567890", "255", "#ee0000")
	c.Text("RECT", 25, 2)
	c.Text("Stroke: YES", 25, 3)
	c.Text("Fill:   NO ", 25, 4)

	c.Fill("1234567890", "#0004ff", "#fffb00")
	c.NoStroke()
	c.Rect(41, 1, 11, 5)
	c.Stroke("1234567890", "255", "#ee0000")
	c.Text("RECT", 45, 2)
	c.Text("Stroke: NO ", 45, 3)
	c.Text("Fill:   YES", 45, 4)

	// circles

	c.Stroke("1234567890", "#fffb00", "#0004ff")
	c.Fill("1234567890", "#0004ff", "#fffb00")
	c.Circle(7, 13, 5)
	c.Stroke("1234567890", "255", "#ee0000")
	c.Text("CIRCLE", 5, 12)
	c.Text("Stroke: YES", 5, 13)
	c.Text("Fill:   YES", 5, 14)

	c.Stroke("1234567890", "#fffb00", "#0004ff")
	c.NoFill()
	c.Circle(27, 13, 5)
	c.Stroke("1234567890", "255", "#ee0000")
	c.Text("CIRCLE", 25, 12)
	c.Text("Stroke: YES", 25, 13)
	c.Text("Fill:   NO ", 25, 14)

	c.Fill("1234567890", "#0004ff", "#fffb00")
	c.NoStroke()
	c.Circle(47, 13, 5)
	c.Stroke("1234567890", "255", "#ee0000")
	c.Text("CIRCLE", 45, 12)
	c.Text("Stroke: NO ", 45, 13)
	c.Text("Fill:   YES", 45, 14)

	// ellipses

	c.Stroke("1234567890", "#fffb00", "#0004ff")
	c.Fill("1234567890", "#0004ff", "#fffb00")
	c.Ellipse(7, 25, 5, 5)
	c.Stroke("1234567890", "255", "#ee0000")
	c.Text("ELLIPSE", 5, 24)
	c.Text("Stroke: YES", 5, 25)
	c.Text("Fill:   YES", 5, 26)

	c.Stroke("1234567890", "#fffb00", "#0004ff")
	c.NoFill()
	c.Ellipse(27, 25, 5, 5)
	c.Stroke("1234567890", "255", "#ee0000")
	c.Text("ELLIPSE", 25, 24)
	c.Text("Stroke: YES", 25, 25)
	c.Text("Fill:   NO ", 25, 26)

	c.Fill("1234567890", "#0004ff", "#fffb00")
	c.NoStroke()
	c.Ellipse(47, 25, 5, 5)
	c.Stroke("1234567890", "255", "#ee0000")
	c.Text("ELLIPSE", 45, 24)
	c.Text("Stroke: NO ", 45, 25)
	c.Text("Fill:   YES", 45, 26)

	// quads

	c.Stroke("1234567890", "#fffb00", "#0004ff")
	c.Fill("1234567890", "#0004ff", "#fffb00")
	c.Quad(64, 1, 74, 1, 74, 6, 64, 6)
	c.Stroke("1234567890", "255", "#ee0000")
	c.Text("QUAD", 70, 2)
	c.Text("Stroke: YES", 70, 3)
	c.Text("Fill:   YES", 70, 4)

	c.Stroke("1234567890", "#fffb00", "#0004ff")
	c.NoFill()
	c.Quad(84, 1, 94, 1, 94, 6, 84, 6)
	c.Stroke("1234567890", "255", "#ee0000")
	c.Text("QUAD", 88, 2)
	c.Text("Stroke: YES", 88, 3)
	c.Text("Fill:   NO ", 88, 4)

	c.Fill("1234567890", "#0004ff", "#fffb00")
	c.NoStroke()
	c.Quad(104, 1, 114, 1, 114, 6, 104, 6)
	c.Stroke("1234567890", "255", "#ee0000")
	c.Text("QUAD", 108, 2)
	c.Text("Stroke: NO ", 108, 3)
	c.Text("Fill:   YES", 108, 4)

	// triangles

	c.Stroke("1234567890", "#fffb00", "#0004ff")
	c.Fill("1234567890", "#0004ff", "#fffb00")
	c.Triangle(64, 9, 74, 9, 74, 18)

	c.Stroke("1234567890", "#fffb00", "#0004ff")
	c.NoFill()
	c.Triangle(84, 9, 94, 9, 94, 18)

	c.Fill("1234567890", "#0004ff", "#fffb00")
	c.NoStroke()
	c.Triangle(104, 9, 114, 9, 114, 18)

	// line

	c.Stroke("1234567890", "#fffb00", "#0004ff")
	c.Fill("1234567890", "#0004ff", "#fffb00")
	c.Line(64, 25, 74, 25)

	c.Stroke("1234567890", "#fffb00", "#0004ff")
	c.NoFill()
	c.Line(84, 25, 94, 25)

	c.Fill("1234567890", "#0004ff", "#fffb00")
	c.NoStroke()
	c.Line(104, 25, 114, 25)

}
