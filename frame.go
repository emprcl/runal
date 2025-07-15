package runal

import "github.com/charmbracelet/lipgloss"

type Cell struct {
	Char       rune
	Foreground lipgloss.Color
	Background lipgloss.Color
}

type Frame [][]Cell

func newFrame(width, height int) Frame {
	buff := make([][]Cell, height)
	for i := range buff {
		buff[i] = make([]Cell, width)
	}
	return buff
}
