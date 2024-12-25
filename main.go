package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

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
				output += " "
			} else {
				output += b[y][x]
			}
			b[y][x] = ""
		}
	}
	fmt.Print(output)
}

type state struct {
	Width, Height int
	buffer        buffer
}

func NewState(width, height int) *state {
	return &state{
		Width:  width,
		Height: height,
		buffer: NewBuffer(width, height),
	}
}

func (s *state) Resize(width, height int) {
	newBuffer := NewBuffer(width, height)

	minWidth := s.Width
	if width < s.Width {
		minWidth = width
	}

	minHeight := s.Height
	if height < s.Height {
		minHeight = height
	}

	for y := 0; y < minHeight; y++ {
		for x := 0; x < minWidth; x++ {
			newBuffer[y][x] = s.buffer[y][x]
		}
	}

	s.Width = width
	s.Height = height
	s.buffer = newBuffer
}

func (s *state) Text(str string, x, y int) {
	s.buffer[y][x] = str
}

func termSize() (int, int) {
	w, h, err := term.GetSize(0)
	if err != nil {
		log.Fatal("can't read terminal size")
	}
	return w, h
}

func main() {
	w, h := termSize()
	s := NewState(w, h)
	setup(s)

	resize := make(chan os.Signal)
	signal.Notify(resize, syscall.SIGWINCH)
	tick := time.Tick(16 * time.Millisecond)

	fmt.Print("\x1b[2J")  // clear entire screen
	fmt.Print("\x1b[25l") // hide cursor

	for {
		select {
		case <-tick:
			fmt.Print("\x1b[H") // reset cursor position
			draw(s)
			s.buffer.Render()
		case <-resize:
			w, h := termSize()
			s.Resize(w, h)
		}
	}
}

func setup(s *state) {
	//fmt.Println(s.Width, s.Height, s.Width/10, s.Height/10)
	//os.Exit(0)
}

func draw(s *state) {
	size := 10
	cols := int(math.Round(float64(s.Width) / float64(size)))
	rows := int(math.Round(float64(s.Height) / float64(size)))
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			s.Text("X", i*size, j*size)
		}
	}
}
