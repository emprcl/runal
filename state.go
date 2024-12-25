package main

import (
	"fmt"
	"math"
)

type state struct {
	Width, Height         int
	termWidth, termHeight int
	buffer                buffer
}

func NewState(width, height int) *state {
	return &state{
		Width:      width/2 + 1,
		Height:     height,
		termWidth:  width,
		termHeight: height,
		buffer:     NewBuffer(width, height),
	}
}

func (s *state) Render() {
	output := ""
	for y := range s.buffer {
		line := ""
		for x := range s.buffer[y] {
			if s.buffer[y][x] == "" {
				line += " "
			} else {
				line += s.buffer[y][x]
			}
			line += " "
			s.buffer[y][x] = ""
		}
		output += forceLength(line, s.termWidth, ' ')
	}
	fmt.Print(output)
}

func (s *state) Resize(width, height int) {
	newWidth := width/2 + 1
	newHeight := height
	newBuffer := NewBuffer(newWidth, newHeight)

	minWidth := s.Width
	if newWidth < s.Width {
		minWidth = newWidth
	}

	minHeight := s.Height
	if newHeight < s.Height {
		minHeight = newHeight
	}

	for y := 0; y < minHeight; y++ {
		for x := 0; x < minWidth; x++ {
			newBuffer[y][x] = s.buffer[y][x]
		}
	}

	s.Width = newWidth
	s.Height = newHeight
	s.termWidth = width
	s.termHeight = height
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
