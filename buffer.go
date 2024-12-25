package main

import "fmt"

type buffer [][]string

func NewBuffer(width, height int) buffer {
	buff := make([][]string, height)
	for i := range buff {
		buff[i] = make([]string, width)
	}
	return buff
}

func (b buffer) Render() {
	output := ""
	for y := range b {
		for x := range b[y] {
			if b[y][x] == "" {
				output += "  "
			} else {
				output += b[y][x]
			}
			b[y][x] = ""
		}
	}
	fmt.Print(output)
}
