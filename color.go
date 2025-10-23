package runal

import (
	"fmt"
	col "image/color"
	"math"
	"strconv"
	"strings"

	"github.com/rahji/termenv"
)

type style struct {
	foreground string
	background string
}

func (s style) equals(s2 style) bool {
	return s.background == s2.background && s.foreground == s2.foreground
}

func (s style) termStyle(t *termenv.Output) termenv.Style {
	return termenv.Style{}.
		Foreground(t.Color(color(s.foreground))).
		Background(t.Color(color(s.background)))
}

func color(color string) string {
	if strings.HasPrefix(color, "#") {
		return color
	}
	c, err := strconv.ParseFloat(strings.TrimSpace(color), 64)
	if err != nil {
		return color
	}
	return strconv.Itoa(int(math.Round(c)))
}

func colorFromImage(c col.Color) string {
	rgba := col.RGBAModel.Convert(c).(col.RGBA)
	return fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}
