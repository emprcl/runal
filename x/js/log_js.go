//go:build js

package js

import (
	"fmt"
	"os"
)

// On the web build, logging goes to stderr (surfaced in the browser console)
// rather than through charmbracelet/log, which would pull lipgloss/termenv and
// their transitive dependencies into the wasm binary. logFatal does not exit.

func logError(a any)                    { fmt.Fprintln(os.Stderr, a) }
func logErrorf(format string, a ...any) { fmt.Fprintf(os.Stderr, format+"\n", a...) }
func logFatal(a any)                    { fmt.Fprintln(os.Stderr, a) }
