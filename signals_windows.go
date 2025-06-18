//go:build windows
// +build windows

package runal

import "os"

func listenForResize() chan os.Signal {
	// SIGWINCH is not implemented on Windows.
	return make(chan os.Signal, 1)
}
