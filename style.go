package runal

import "github.com/charmbracelet/x/ansi"

type style struct {
	foreground ansi.Color
	background ansi.Color
}

func (s style) equals(s2 style) bool {
	return s.background == s2.background && s.foreground == s2.foreground
}

func (s style) render(str string) string {
	return ansi.NewStyle().
		BackgroundColor(s.background).
		ForegroundColor(s.foreground).
		String() + str
}

func resetStyle() string {
	return "\x1b[0m"
}
