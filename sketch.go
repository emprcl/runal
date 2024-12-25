package main

import "math"

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
			s.Text("X", i*size, j*size)
		}
	}
}
