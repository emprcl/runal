//go:build js

// Command runal-wasm is the WebAssembly entry point for runal.
//
// The sketch itself runs in the browser's native JS engine; this module only
// hosts the Go canvas engine, exposed to JS through the //go:wasmexport
// functions in internal/canvas (see web_js.go). The blank import pulls in the
// engine and its exported functions, and main blocks so the Go runtime stays
// resident and those exports remain callable.
package main

import (
	_ "github.com/emprcl/runal/internal/canvas"
)

func main() { select {} }
