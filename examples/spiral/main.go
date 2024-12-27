package main

import (
	"math"

	"github.com/emprcl/runal"
)

var (
	t          float64
	background string   = "#000000"
	colors     []string = []string{
		"#fcf6bd",
		"#d0f4de",
		"#a9def9",
		"#e4c1f9",
		"#ff99c8",
	}
)

func main() {
	runal.Run(setup, draw, runal.WithFPS(60))
}

func setup(c *runal.Canvas) {
	c.Background(background)
}

func draw(c *runal.Canvas) {
	size := 1
	cols := int(math.Round(float64(c.Width()) / float64(size)))
	rows := int(math.Round(float64(c.Height()) / float64(size)))
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			x := i * size
			y := j * size
			d := c.Dist(x, y, c.Width()/2, c.Height()/2)
			k := .6
			dx := float64(x) - float64(c.Width())/2.
			dy := float64(y) - float64(c.Height())/2.
			angle := math.Atan2(dy, dx)
			spiralPath := math.Sin(d/k + angle + t)
			df := 2.
			af := 2.
			threshold := math.Sin(d/df + af*angle)

			c.Foreground(colorGradient(c.Width(), d))

			if spiralPath > threshold {
				c.Text("⬤", x, y)
			}
		}
	}

	t += 0.5
}

func colorGradient(width int, d float64) string {
	step := width / len(colors)
	for i := 0; i < len(colors); i++ {
		if d <= float64((i+1)*step) {
			return colors[i]
		}
	}
	return colors[len(colors)-1]
}
