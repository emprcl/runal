package runal

import (
	"time"

	perlin "github.com/aquilax/go-perlin"
)

const (
	alpha = 2.
	beta  = 2.
	n     = 3
)

func newNoise() *perlin.Perlin {
	return perlin.NewPerlin(alpha, beta, n, time.Now().UnixNano())
}

func (c *Canvas) NoiseSeed(seed int64) {
	c.noise = perlin.NewPerlin(alpha, beta, n, seed)
}

func (c *Canvas) Noise1D(x float64) float64 {
	return c.noise.Noise1D(x)/2 + 0.5
}

func (c *Canvas) Noise2D(x, y float64) float64 {
	return c.noise.Noise2D(x, y)/2 + 0.5
}
