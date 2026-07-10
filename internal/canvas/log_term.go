//go:build !js

package canvas

import "github.com/charmbracelet/log"

func logErrorf(format string, a ...any) {
	log.Errorf(format, a...)
}
