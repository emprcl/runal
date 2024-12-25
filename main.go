package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	w, h := termSize()
	s := NewState(w/2, h)
	setup(s)

	resize := make(chan os.Signal)
	signal.Notify(resize, syscall.SIGWINCH)
	tick := time.Tick(16 * time.Millisecond)

	EnterAltScreen()

	for {
		select {
		case <-tick:
			ResetCursorPosition()
			draw(s)
			s.buffer.Render()
		case <-resize:
			w, h := termSize()
			s.Resize(w/2, h)
			ClearScreen()
		}
	}
}
