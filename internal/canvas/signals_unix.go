//go:build darwin || linux

package canvas

import (
	"os"
	"os/signal"
	"syscall"
)

func listenForResize() chan os.Signal {
	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)
	return resize
}
