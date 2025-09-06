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
	pointsNb = 26
	a        = 5
	b        = 2
)

type point struct {
	x, y     float64
	char     string
	duration int
}

func (p *point) update(c *runal.Canvas) {
	c.StrokeText(p.char)
	c.Point(int(p.x), int(p.y))
	p.duration--
}

var points = []point{}

func main() {
	runal.Run(context.Background(), setup, draw, runal.WithOnKey(onKey))
}

func setup(c *runal.Canvas) {
	c.Size(40, 20)
	c.BackgroundBg("197")
	c.Stroke(".", "255", "197")
}

func draw(c *runal.Canvas) {
	c.Clear()
	theta := c.LoopAngle(duration)
	x := c.Map(math.Sin(a*theta), -1, 1, margin, float64(c.Width-margin))
	y := c.Map(math.Sin(b*theta), -1, 1, margin, float64(c.Height-margin))
	x2 := c.Map(math.Sin(b*theta), -1, 1, margin, float64(c.Width-margin))
	y2 := c.Map(math.Sin(a*theta), -1, 1, margin, float64(c.Height-margin))
	if len(points) <= pointsNb*2 {
		points = append(points, point{x, y, "0", pointsNb})
		points = append(points, point{x2, y2, "#", pointsNb})
	}

	for i := range points {
		points[i].update(c)
	}

	// remove points with a duration of zero
	filtered := points[:0]
	for _, p := range points {
		if p.duration > 0 {
			filtered = append(filtered, p)
		}
	}
	points = filtered
}

func onKey(c *runal.Canvas, e runal.KeyEvent) {
	if e.Key == "space" {
		c.NoiseSeed(time.Now().Unix())
		c.Redraw()
	}
}
