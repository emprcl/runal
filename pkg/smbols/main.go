package smbols

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"smbols/internal/ansi"
	"smbols/internal/util"
)

func Run(setup, draw func(c *Canvas)) {
	w, h := util.TermSize()
	c := newCanvas(w, h)
	setup(c)

	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)
	tick := time.Tick(16 * time.Millisecond)

	ansi.EnterAltScreen()

	for {
		select {
		case <-resize:
			w, h := util.TermSize()
			c.resize(w, h)
			ansi.ClearScreen()
		case <-tick:
			ansi.ResetCursorPosition()
			draw(c)
			c.render()
		}
	}
}
