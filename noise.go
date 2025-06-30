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

// Noise1D generates 1D Perlin noise (range [0, 1]) for a given input.
func (c *Canvas) Noise1D(x float64) float64 {
	return c.noise.Noise1D(x)/2 + 0.5
}

// Noise2D generates 2D Perlin noise (range [0, 1]) for a given (x, y) coordinate.
func (c *Canvas) Noise2D(x, y float64) float64 {
	return c.noise.Noise2D(x, y)/2 + 0.5
}

// LoopAngle returns the angular progress (in radians, range [0, 2π]) through a looping cycle
// of given duration in seconds, with a given frame offset.
func (c *Canvas) LoopAngle(duration, offset int) float64 {
	totalFrames := c.fps * duration
	frame := (c.Framecount + offset) % totalFrames
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
//
// Returns:
//   - A float64 noise value in the range [0, 1].
func (c *Canvas) NoiseLoop(angle, radius float64) float64 {
	x := c.Map(math.Cos(angle), -1, 1, 0, radius)
	y := c.Map(math.Sin(angle), -1, 1, 0, radius)
	return c.Noise2D(x, y)
}

// NoiseLoop1D returns a 1D noise value that loops as the angle progresses.
// It samples a 2D noise space using the given radius and combines it with a horizontal offset.
//
// Parameters:
//   - angle: the loop angle in radians, typically from 0 to 2π.
//   - radius: the radius of the circular path in noise space.
//   - x: horizontal position to offset the noise sampling.
//
// Returns:
//   - A float64 noise value in the range [0, 1].
func (c *Canvas) NoiseLoop1D(angle, radius float64, x int) float64 {
	nx := c.Map(math.Cos(angle), -1, 1, 0, radius)
	ny := c.Map(math.Sin(angle), -1, 1, 0, radius)
	return c.Noise2D(nx+float64(x), ny)
}

// NoiseLoop2D returns a 2D noise value that loops as the angle progresses.
// It samples a circular path in 2D noise space, offset by the (x, y) coordinates.
//
// Parameters:
//   - angle: the loop angle in radians, typically from 0 to 2π.
//   - radius: the radius of the circular path in noise space.
//   - x, y: coordinates used to offset the sampled position in the noise field.
//
// Returns:
//   - A float64 noise value in the range [0, 1].
func (c *Canvas) NoiseLoop2D(angle, radius float64, x, y int) float64 {
	nx := c.Map(math.Cos(angle), -1, 1, 0, radius)
	ny := c.Map(math.Sin(angle), -1, 1, 0, radius)
	return c.Noise2D(nx+float64(x), ny+float64(y))
}
