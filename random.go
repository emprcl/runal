package runal

import (
	"math/rand"
	"time"
)

func newRandom() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Random returns a random float between minimum and maximum.
func (c *Canvas) Random(minimum, maximum int) float64 {
	return c.Map(c.random.Float64(), 0, 1, float64(minimum), float64(maximum))
}

// RandomSeed sets the random number generator seed.
func (c *Canvas) RandomSeed(seed int64) {
	c.random = rand.New(rand.NewSource(seed))
}
