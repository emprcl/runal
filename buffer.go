package runal

import "github.com/charmbracelet/lipgloss"

type cell struct {
	char       rune
	foreground lipgloss.Color
	background lipgloss.Color
}

type buffer [][]cell

func newBuffer(width, height int) buffer {
	buff := make([][]cell, height)
	for i := range buff {
		buff[i] = make([]cell, width)
	}
	return buff
}
