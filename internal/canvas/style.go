package canvas

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

type style struct {
	foreground ansi.Color
	background ansi.Color
}

func (s style) equals(s2 style) bool {
	return s.background == s2.background && s.foreground == s2.foreground
}

func (s style) render(str string) string {
	fr, fg, fb, _ := s.foreground.RGBA()
	br, bg, bb, _ := s.background.RGBA()
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%s",
		fr>>8, fg>>8, fb>>8,
		br>>8, bg>>8, bb>>8,
		str)
}

func (s style) writeTo(b *strings.Builder) {
	fr, fg, fb, _ := s.foreground.RGBA()
	br, bg, bb, _ := s.background.RGBA()
	fmt.Fprintf(b, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm",
		fr>>8, fg>>8, fb>>8,
		br>>8, bg>>8, bb>>8)
}

const resetStyleStr = "\x1b[0m"
