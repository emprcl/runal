package main

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
)

const (
	padChar = " "
)

type state struct {
	Width, Height          int
	termWidth, termHeight  int
	buffer                 buffer
	strokeColor, fillColor lipgloss.Color
}

func NewState(width, height int) *state {
	return &state{
		Width:       width / 2,
		Height:      height,
		termWidth:   width,
		termHeight:  height,
		buffer:      NewBuffer(width, height),
		strokeColor: lipgloss.Color("#ffffff"),
		fillColor:   lipgloss.Color("#000000"),
	}
}

func (s *state) Render() {
	output := ""
	for y := range s.buffer {
		line := ""
		for x := range s.buffer[y] {
			add := ""
			if s.buffer[y][x] == "" {
				add = s.Style("  ")
			} else {
				add = s.buffer[y][x]
			}
			if lipgloss.Width(line+add) <= s.termWidth {
				line += add
			}
			s.buffer[y][x] = ""
		}
		output += forcePadding(line, s.termWidth, ' ')
	}
	fmt.Print(output)
}

func (s *state) Resize(width, height int) {
	newWidth := width / 2
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
	s.buffer[y][x] = s.Style(str + padChar)
}

func (s *state) Style(str string) string {
	return lipgloss.NewStyle().
		Background(s.fillColor).
		Foreground(s.strokeColor).
		Render(str)
}

func (s *state) Dist(x1, y1, x2, y2 int) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(float64(dx*dx + dy*dy))
}

func (s *state) Background(color string) {
	s.fillColor = lipgloss.Color(color)
}

func (s *state) Foreground(color string) {
	s.strokeColor = lipgloss.Color(color)
}
