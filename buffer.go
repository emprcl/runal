package runal

import "github.com/charmbracelet/lipgloss"

type Cell struct {
	Char       rune
	Foreground lipgloss.Color
	Background lipgloss.Color
}

type buffer [][]Cell

func newBuffer(width, height int) buffer {
	buff := make([][]Cell, height)
	for i := range buff {
		buff[i] = make([]Cell, width)
	}
	return buff
}
