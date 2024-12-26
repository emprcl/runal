package util

import (
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func TermSize() (int, int) {
	w, h, err := term.GetSize(0)
	if err != nil {
		log.Fatal("can't read terminal size")
	}
	return w, h
}

func ForcePadding(s string, length int, padChar rune) string {
	if lipgloss.Width(s) < length {
		padding := strings.Repeat(string(padChar), length-lipgloss.Width(s))
		return s + padding
	}
	return s
}
