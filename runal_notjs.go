//go:build !js

package runal

import (
	"context"
	"sync"

	"github.com/emprcl/runal/internal/canvas"
)

// Run starts the runal event loop and blocks until the sketch exits.
func Run(ctx context.Context, setup, draw func(c *Canvas), opts ...CallbackOption) {
	canvas.Run(ctx, setup, draw, opts...)
}

// Start starts the runal event loop in a goroutine and returns a WaitGroup.
func Start(ctx context.Context, done chan struct{}, setup, draw func(c *Canvas), opts ...CallbackOption) *sync.WaitGroup {
	return canvas.Start(ctx, done, setup, draw, opts...)
}
