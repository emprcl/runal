package main

import (
	"context"
	"strconv"
	"time"

	"github.com/emprcl/runal"
)

const (
	duration = 5
	radius   = 4
)

var seed1, seed2 int64

func main() {
	seed1 = time.Now().Unix()
	seed2 = seed1 + 1000
	runal.Run(context.Background(), setup, draw, runal.WithOnKey(onKey))
}

func setup(c *runal.Canvas) {
	c.CellModeDouble()
}

func draw(c *runal.Canvas) {
	c.Clear()

	theta := c.LoopAngle(duration)
	c.NoiseSeed(seed1)
	noise := c.NoiseLoop(theta, 1)
	c.NoiseSeed(seed2)
	noise2 := c.NoiseLoop(theta, 1)
	x := c.Map(noise, 0, 1, radius, float64(c.Width-radius))
	y := c.Map(noise2, 0, 1, radius, float64(c.Height-radius))
	color := c.Map(noise, 0, 1, 0, 10)
	colorStr := strconv.FormatFloat(color, 'f', -1, 64)
	colorBg := c.Map(noise2, 0, 1, 200, 210)
	colorBgStr := strconv.FormatFloat(colorBg, 'f', -1, 64)

	c.Background("#", colorBgStr, colorBgStr)
	c.Stroke(" ", colorStr, colorStr)
	c.Fill(" ", colorStr, colorStr)
	c.Circle(int(x), int(y), radius)
}

func onKey(c *runal.Canvas, e runal.KeyEvent) {
	if e.Key == "space" {
		seed1 = time.Now().Unix()
		seed2 = seed1 + 1000
	}
	if e.Key == "c" {
		c.SaveCanvasToMP4("flash.mp4", duration)
	}
}
