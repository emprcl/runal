package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	w, h := termSize()
	s := NewState(w, h)
	setup(s)

	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)
	tick := time.Tick(16 * time.Millisecond)

	EnterAltScreen()

	for {
		select {
		case <-resize:
			w, h := termSize()
			s.Resize(w, h)
			ClearScreen()
		case <-tick:
			ResetCursorPosition()
			draw(s)
			s.Render()
		}
	}
}
