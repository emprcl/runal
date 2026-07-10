//go:build js

package runal

import (
	"syscall/js"

	"github.com/emprcl/runal/internal/canvas"
)

// SetDisplayCanvas registers the DOM <canvas> element the sketch renders into,
// along with the font size (in CSS pixels) used for cell metrics. It must be
// called before Run/Start. Only available on the js/wasm build.
func SetDisplayCanvas(el js.Value, fontSize int) {
	canvas.SetDisplayCanvas(el, fontSize)
}
