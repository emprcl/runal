package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
