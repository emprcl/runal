package main

import (
	"math"
)

var (
	t float64
)

func setup(s *state) {
	//fmt.Println(s.Width, s.Height, s.Width/10, s.Height/10)
	//os.Exit(0)
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
			k := 10.
			dx := float64(x) - float64(s.Width)/2.
			dy := float64(y) - float64(s.Height)/2.
			angle := math.Atan2(dy, dx)
			spiralPath := math.Sin(d/k + angle + t)

			df := 100.
			af := 3.
			threshold := math.Sin(d/df + af*angle)
			if spiralPath > threshold {
				s.Text("O", x, y)
			} else {
				s.Text(".", x, y)
			}
		}
	}

	t += 0.05
}
