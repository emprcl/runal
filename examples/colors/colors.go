package main

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/emprcl/runal"
)

func main() {
	runal.Run(context.Background(), setup, draw)
}

func setup(c *runal.Canvas) {
	c.NoLoop()
}

func draw(c *runal.Canvas) {
	for i := range 256 {
		c.Push()
		c.Stroke(" ", "0", strconv.Itoa(i))
		c.Translate((i%16)*10, int(math.Floor(float64(i/16)))*3)
		c.Line(2, 1, 6, 1)
		c.Line(2, 2, 6, 2)
		c.Stroke(" ", "15", "0")
		text := fmt.Sprintf("%d", i)
		c.Text(text, 8, 1)
		c.Pop()
	}
}
