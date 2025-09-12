package runtime

import (
	"errors"

	"github.com/dop251/goja"
)

type callbacks struct {
	onKey          goja.Callable
	onMouseClick   goja.Callable
	onMouseRelease goja.Callable
	onMouseWheel   goja.Callable
}

func parseJS(script string) (*goja.Runtime, goja.Callable, goja.Callable, callbacks, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	_, err := vm.RunString(script)

	cb := callbacks{}

	if err != nil {
		return nil, nil, nil, cb, err
	}
	setup, ok := goja.AssertFunction(vm.Get("setup"))
	if !ok {
		return nil, nil, nil, cb, errors.New("The file does not contain a setup method")
	}

	draw, ok := goja.AssertFunction(vm.Get("draw"))
	if !ok {
		return nil, nil, nil, cb, errors.New("The file does not contain a draw method")
	}

	cb.onKey, _ = goja.AssertFunction(vm.Get("onKey"))
	cb.onMouseClick, _ = goja.AssertFunction(vm.Get("onMouseClick"))
	cb.onMouseRelease, _ = goja.AssertFunction(vm.Get("onMouseRelease"))
	cb.onMouseWheel, _ = goja.AssertFunction(vm.Get("onMouseWheel"))

	return vm, setup, draw, cb, nil
}
