package main

import (
	"github.com/dop251/goja"
	"github.com/emprcl/runal"
)

func main() {
	const SCRIPT = `
	function setup() {

	}

	function draw(runal) {
		runal.Char('G', 10, 10);

		//return runal
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
			_, err := setup(goja.Undefined())
			if err != nil {
				panic(err)
			}
		},
		func(c *runal.Canvas) {
			_, err := draw(goja.Undefined(), goja.New().ToValue(c))
			if err != nil {
				panic(err)
			}
			//can := result.Export().(*runal.Canvas)
			//fmt.Println(&c, &can)
		},
		runal.WithFPS(60),
	)
}

func setup(c *runal.Canvas) {}

func draw(c *runal.Canvas) {}
