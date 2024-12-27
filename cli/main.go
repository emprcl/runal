package main

import (
	"github.com/dop251/goja"
	"github.com/emprcl/runal"
)

func main() {
	const SCRIPT = `
	let text = "";
	function setup(runal) {
		text = "coucou";
	}

	function draw(runal) {
		runal.Text(text, 10, 10);
	}
	`

	vm := goja.New()
	_, err := vm.RunString(SCRIPT)
	if err != nil {
		panic(err)
	}
	setup, ok := goja.AssertFunction(vm.Get("setup"))
	if !ok {
		panic("Not a function")
	}

	draw, ok := goja.AssertFunction(vm.Get("draw"))
	if !ok {
		panic("Not a function")
	}

	runal.Run(
		func(c *runal.Canvas) {
			_, err := setup(goja.Undefined(), goja.New().ToValue(c))
			if err != nil {
				panic(err)
			}
		},
		func(c *runal.Canvas) {
			_, err := draw(goja.Undefined(), goja.New().ToValue(c))
			if err != nil {
				panic(err)
			}
		},
		runal.WithFPS(60),
	)
}
