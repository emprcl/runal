//go:build js

// Command runal-wasm is the WebAssembly entry point for runal. It exposes two
// globals to JavaScript:
//
//	runalStart(source, canvasEl[, fontSize]) -> ""|errorString
//	runalStop()
//
// runalStart mounts the given <canvas> element as the render target and runs
// the provided JavaScript sketch source. Starting a new sketch stops the
// previous one.
package main

import (
	"syscall/js"

	"github.com/emprcl/runal"
	rjs "github.com/emprcl/runal/x/js"
)

var stopCurrent func()

func start(_ js.Value, args []js.Value) any {
	if len(args) < 2 {
		return "runalStart(source, canvasEl[, fontSize]) requires 2 arguments"
	}
	source := args[0].String()
	el := args[1]
	fontSize := 16
	if len(args) > 2 && args[2].Type() == js.TypeNumber {
		fontSize = args[2].Int()
	}

	if stopCurrent != nil {
		stopCurrent()
		stopCurrent = nil
	}

	runal.SetDisplayCanvas(el, fontSize)

	stop, err := rjs.New("", nil).Start(source)
	if err != nil {
		return err.Error()
	}
	stopCurrent = stop
	return ""
}

func stop(_ js.Value, _ []js.Value) any {
	if stopCurrent != nil {
		stopCurrent()
		stopCurrent = nil
	}
	return nil
}

func main() {
	js.Global().Set("runalStart", js.FuncOf(start))
	js.Global().Set("runalStop", js.FuncOf(stop))
	select {} // keep the Go runtime alive for the registered callbacks
}
