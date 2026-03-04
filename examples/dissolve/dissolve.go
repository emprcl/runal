// dissolve — noise-driven glyph landscape
//
// Maps 2D noise values to character pools of increasing visual density.
// The glyph choice itself is the visual — no shapes, no lines, just
// character weight shifting across the screen as the noise field
// drifts through time.
package main

import (
	"context"
	"math"
	"strconv"

	"github.com/emprcl/runal"
)

// Character pools ordered by visual density.
// Curating these is where the aesthetic lives —
// different pools make completely different landscapes
// from the same noise field.
var pools = [][]rune{
	[]rune(" "),
	[]rune("·.:,"),
	[]rune("░╌┆╍"),
	[]rune("▒▌▐╎"),
	[]rune("▓█▀▄"),
}

func main() {
	runal.Run(context.Background(), setup, draw)
}

func setup(c *runal.Canvas) {
	c.Fps(10)
}

func draw(c *runal.Canvas) {
	t := float64(c.Framecount) * 0.008

	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			n := c.Noise2D(float64(x)*0.04+t, float64(y)*0.06+t*0.3)

			// Select pool by noise band
			band := int(math.Floor(c.Map(n, 0, 1, 0, float64(len(pools)))))
			if band >= len(pools) {
				band = len(pools) - 1
			}
			pool := pools[band]

			// Random character from the pool
			ch := string(pool[int(c.Random(0, len(pool)))])

			// Brightness follows density
			fg := strconv.Itoa(int(math.Floor(c.Map(n, 0, 1, 236, 255))))

			c.Stroke(ch, fg, "0")
			c.Point(x, y)
		}
	}
}
