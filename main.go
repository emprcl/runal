package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

type state struct {
	Width, Height int
}

func size() (int, int) {
	w, h, err := term.GetSize(0)
	if err != nil {
		log.Fatal("can't read terminal size")
	}
	return w, h
}

func main() {
	w, h := size()
	s := &state{
		Width:  w,
		Height: h,
	}
	setup(s)

	resize := make(chan os.Signal)
	signal.Notify(resize, syscall.SIGWINCH)
	tick := time.Tick(16 * time.Millisecond)

	for {
		select {
		case <-tick:
			draw(s)
		case <-resize:
			w, h := size()
			s.Width = w
			s.Height = h
		}
	}
}

func setup(s *state) {
	fmt.Println(s.Width, s.Height)
}

func draw(s *state) {
	fmt.Println(s.Width, s.Height)
}
