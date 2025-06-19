package runal

import (
	"math"
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

// NoiseSeed sets the random seed for noise generation.
func (c *Canvas) NoiseSeed(seed int64) {
	c.noise = perlin.NewPerlin(alpha, beta, n, seed)
}

// Noise1D generates 1D Perlin noise for a given input.
func (c *Canvas) Noise1D(x float64) float64 {
	return c.noise.Noise1D(x)/2 + 0.5
}

// Noise2D generates 2D Perlin noise for a given (x, y) coordinate.
func (c *Canvas) Noise2D(x, y float64) float64 {
	return c.noise.Noise2D(x, y)/2 + 0.5
}

// LoopAngle returns the angular progress (in radians) through a looping cycle
// of given duration in seconds.
func (c *Canvas) LoopAngle(duration int) float64 {
	totalFrames := c.fps * duration
	frame := c.Framecount % totalFrames
	return c.Map(float64(frame), 0, float64(totalFrames), 0, 2*math.Pi)
}

// NoiseLoop returns a noise value by sampling the noise on a circular
// path of the given radius.
// This is useful for creating cyclic animations or evolving patterns
// that repeat perfectly after one full loop.
//
// Parameters:
//   - angle: the loop angle in radians, from 0 to 2π.
//   - radius: the radius of the circular path in noise space.
//     the higher the radius, the more instable it gets.
func (c *Canvas) NoiseLoop(angle, radius float64) float64 {
	x := c.Map(math.Cos(angle), -1, 1, 0, float64(radius))
	y := c.Map(math.Sin(angle), -1, 1, 0, float64(radius))
	return c.Noise2D(x, y)
}

func (c *Canvas) NoiseLoop1D(angle, radius float64, x int) float64 {
	nx := c.Map(math.Cos(angle), -1, 1, 0, float64(radius))
	ny := c.Map(math.Sin(angle), -1, 1, 0, float64(radius))
	return c.Noise2D(nx+float64(x), ny)
}

func (c *Canvas) NoiseLoop2D(angle, radius float64, x, y int) float64 {
	nx := c.Map(math.Cos(angle), -1, 1, 0, float64(radius))
	ny := c.Map(math.Sin(angle), -1, 1, 0, float64(radius))
	return c.Noise2D(nx+float64(x), ny+float64(y))
}
