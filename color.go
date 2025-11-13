package runal

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"golang.org/x/exp/constraints"
)

func (c *Canvas) ColorRGB(r, g, b int) string {
	r = clamp(r, 0, 255)
	g = clamp(g, 0, 255)
	b = clamp(b, 0, 255)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func (c *Canvas) ColorHSL(h, s, l float64) string {
	h = clamp[float64](h, 0, 360)
	s = clamp[float64](s, 0., 1.)
	l = clamp[float64](l, 0., 1.)

	C := (1 - math.Abs((2*l)-1)) * s
	X := C * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - (C / 2)
	var Rnot, Gnot, Bnot float64

	switch {
	case 0 <= h && h < 60:
		Rnot, Gnot, Bnot = C, X, 0
	case 60 <= h && h < 120:
		Rnot, Gnot, Bnot = X, C, 0
	case 120 <= h && h < 180:
		Rnot, Gnot, Bnot = 0, C, X
	case 180 <= h && h < 240:
		Rnot, Gnot, Bnot = 0, X, C
	case 240 <= h && h < 300:
		Rnot, Gnot, Bnot = X, 0, C
	case 300 <= h && h < 360:
		Rnot, Gnot, Bnot = C, 0, X
	}
	r := int(math.Round((Rnot + m) * 255))
	g := int(math.Round((Gnot + m) * 255))
	b := int(math.Round((Bnot + m) * 255))

	return c.ColorRGB(r, g, b)
}

func (c *Canvas) ColorHSV(h, s, v float64) string {
	h = clamp[float64](h, 0, 360)
	s = clamp[float64](s, 0., 1.)
	v = clamp[float64](v, 0., 1.)

	C := v * s
	X := C * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - C
	var Rnot, Gnot, Bnot float64
	switch {
	case 0 <= h && h < 60:
		Rnot, Gnot, Bnot = C, X, 0
	case 60 <= h && h < 120:
		Rnot, Gnot, Bnot = X, C, 0
	case 120 <= h && h < 180:
		Rnot, Gnot, Bnot = 0, C, X
	case 180 <= h && h < 240:
		Rnot, Gnot, Bnot = 0, X, C
	case 240 <= h && h < 300:
		Rnot, Gnot, Bnot = X, 0, C
	case 300 <= h && h < 360:
		Rnot, Gnot, Bnot = C, 0, X
	}
	r := int(math.Round((Rnot + m) * 255))
	g := int(math.Round((Gnot + m) * 255))
	b := int(math.Round((Bnot + m) * 255))

	return c.ColorRGB(r, g, b)
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

type Number interface {
	constraints.Integer | constraints.Float
}

func clamp[T Number](value, minimum, maximum T) T {
	return max(min(value, maximum), minimum)
}
