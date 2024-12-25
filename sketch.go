package main

import (
	"math"
)

var (
	t          float64
	background string   = "#211103"
	colors     []string = []string{
		"#f8e5ee",
		"#9f2042",
		"#7b0d1e",
		"#3d1308",
	}
)

func setup(s *state) {
	s.Background(background)
}

func draw(s *state) {
	size := 1
	cols := int(math.Round(float64(s.Width) / float64(size)))
	rows := int(math.Round(float64(s.Height) / float64(size)))
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			x := i * size
			y := j * size
			d := s.Dist(x, y, s.Width/2, s.Height/2)
			k := 0.9
			dx := float64(x) - float64(s.Width)/2.
			dy := float64(y) - float64(s.Height)/2.
			angle := math.Atan2(dy, dx)
			spiralPath := math.Sin(d/k + angle + t)
			df := 30.
			af := 3.
			threshold := math.Sin(d/df + af*angle)

			s.Foreground(colorGradient(s.Width, d))

			if spiralPath > threshold {
				s.Text("O", x, y)
			} else {
				s.Text(".", x, y)
			}
		}
	}

	t += 0.1
}

func colorGradient(width int, d float64) string {
	step := width / len(colors)
	for i := 0; i < len(colors); i++ {
		if d <= float64((i+1)*step) {
			return colors[i]
		}
	}
	return colors[len(colors)-1]
}
