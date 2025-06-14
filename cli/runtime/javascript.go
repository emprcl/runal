package runtime

import (
	"errors"

	"github.com/dop251/goja"
)

func parseJS(script string) (*goja.Runtime, goja.Callable, goja.Callable, goja.Callable, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	_, err := vm.RunString(script)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	setup, ok := goja.AssertFunction(vm.Get("setup"))
	if !ok {
		return nil, nil, nil, nil, errors.New("The file does not contain a setup method")
	}

	draw, ok := goja.AssertFunction(vm.Get("draw"))
	if !ok {
		return nil, nil, nil, nil, errors.New("The file does not contain a draw method")
	}

	onKey, ok := goja.AssertFunction(vm.Get("onKey"))
	if !ok {
		return vm, setup, draw, nil, nil
	}

	return vm, setup, draw, onKey, nil
}
