package main

import (
	"github.com/dop251/goja"
	"github.com/emprcl/runal"
)

func main() {
	const SCRIPT = `
	let text = "";
	function setup() {
		text = "coucou";
	}

	function draw() {
		runal.text(text, 10, 10);
	}
	`

	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
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
			vm.Set("runal", c)
			_, err := setup(goja.Undefined())
			if err != nil {
				panic(err)
			}
		},
		func(c *runal.Canvas) {
			vm.Set("runal", c)
			_, err := draw(goja.Undefined())
			if err != nil {
				panic(err)
			}
		},
		runal.WithFPS(60),
	)
}
