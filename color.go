package runal

import (
	"fmt"
	col "image/color"
	"math"
	"strconv"
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
	return ansi.NewStyle().
		BackgroundColor(s.background).
		ForegroundColor(s.foreground).
		String() + str
}

func resetStyle() string {
	return "\x1b[0m"
}

func color(color string) ansi.Color {
	if strings.HasPrefix(color, "#") {
		return ansi.HexColor(color)
	}
	c, err := strconv.ParseFloat(strings.TrimSpace(color), 64)
	if err != nil {
		return ansi.HexColor(color)
	}
	return ansi.IndexedColor(uint8(math.Round(c)))
}

func colorFromImage(c col.Color) ansi.Color {
	rgba := col.RGBAModel.Convert(c).(col.RGBA)
	return ansi.RGBColor{
		R: rgba.R,
		G: rgba.G,
		B: rgba.B,
	}
}

func colorToString(c ansi.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%d%d%d", r, g, b)
}
