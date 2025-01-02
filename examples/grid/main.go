package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/emprcl/runal"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	runal.Run(ctx, setup, draw).Wait()
}

func setup(c *runal.Canvas) {}

func draw(c *runal.Canvas) {
	c.Flush()
	for i := 0; i < c.Width; i++ {
		for j := 0; j < c.Height; j++ {
			c.Text(".", i, j)
		}
	}
}
