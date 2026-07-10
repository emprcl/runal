//go:build js

package canvas

import (
	"fmt"
	"syscall/js"
)

func logErrorf(format string, a ...any) {
	js.Global().Get("console").Call("error", fmt.Sprintf(format, a...))
}
