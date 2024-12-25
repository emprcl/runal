package main

type buffer [][]string

func NewBuffer(width, height int) buffer {
	buff := make([][]string, height)
	for i := range buff {
		buff[i] = make([]string, width)
	}
	return buff
}
