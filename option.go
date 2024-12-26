package runal

import (
	"log"
	"time"
)

const (
	defaultFrameDuration = 16 * time.Millisecond
)

type option func(*options)

type options struct {
	frameDuration time.Duration
}

func newOptions() *options {
	return &options{
		frameDuration: defaultFrameDuration,
	}
}

func WithFPS(fps int) option {
	if fps <= 0 {
		log.Fatal("FPS must be greater than 0")
	}
	return func(opts *options) {
		opts.frameDuration = time.Second / time.Duration(fps)
	}
}
