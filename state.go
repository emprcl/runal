package main

import (
	"math"
)

type state struct {
	Width, Height int
	buffer        buffer
}

func NewState(width, height int) *state {
	return &state{
		Width:  width,
		Height: height,
		buffer: NewBuffer(width, height),
	}
}

func (s *state) Resize(width, height int) {
	newBuffer := NewBuffer(width, height)

	minWidth := s.Width
	if width < s.Width {
		minWidth = width
	}

	minHeight := s.Height
	if height < s.Height {
		minHeight = height
	}

	for y := 0; y < minHeight; y++ {
		for x := 0; x < minWidth; x++ {
			newBuffer[y][x] = s.buffer[y][x]
		}
	}

	s.Width = width
	s.Height = height
	s.buffer = newBuffer
}

func (s *state) Text(str string, x, y int) {
	s.buffer[y][x] = str
}

func (s *state) Dist(x1, y1, x2, y2 int) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(float64(dx*dx + dy*dy))
}
