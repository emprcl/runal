//go:build !js

package js

import "github.com/charmbracelet/log"

func logError(a any)                      { log.Error(a) }
func logErrorf(format string, a ...any)   { log.Errorf(format, a...) }
func logFatal(a any)                      { log.Fatal(a) }
