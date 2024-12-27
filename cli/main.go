package main

import (
	"context"
	"errors"
	"flag"
	"os"

	"github.com/charmbracelet/log"

	"github.com/dop251/goja"
	"github.com/emprcl/runal"
	"github.com/fsnotify/fsnotify"
)

func main() {
	file := flag.String("f", "", "sketch file (js)")
	flag.Parse()

	if *file == "" {
		log.Fatal("sketch file (-f) argument is mandatory")
	}

	if _, err := os.Stat(*file); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("sketch file %s does not exit", *file)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	ctx, cancel := context.WithCancel(context.Background())
	content, err := os.ReadFile(*file)
	vm, setup, draw, err := parseJS(string(content))
	if err != nil {
		log.Error("error:", err)
	} else {
		runSketch(ctx, vm, setup, draw)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Infof("%s updated. Reloading...\n", event.Name)
					content, err := os.ReadFile(event.Name)
					if err != nil {
						log.Error("error:", err)
						return
					}
					vm, setup, draw, err := parseJS(string(content))
					if err != nil {
						log.Error("error:", err)
						return
					}
					cancel()
					ctx, cancel = context.WithCancel(context.Background())
					runSketch(ctx, vm, setup, draw)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(*file)
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}

func parseJS(script string) (*goja.Runtime, goja.Callable, goja.Callable, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	_, err := vm.RunString(script)
	if err != nil {
		return nil, nil, nil, err
	}
	setup, ok := goja.AssertFunction(vm.Get("setup"))
	if !ok {
		return nil, nil, nil, errors.New("The file does not contain a setup method.")
	}

	draw, ok := goja.AssertFunction(vm.Get("draw"))
	if !ok {
		return nil, nil, nil, errors.New("The file does not contain a draw method.")
	}

	return vm, setup, draw, nil
}

type console struct{}

func (c console) Log(msg string) {
	log.Info(msg)
}

func runSketch(ctx context.Context, vm *goja.Runtime, setup, draw goja.Callable) {
	runal.Run(
		ctx,
		func(c *runal.Canvas) {
			vm.Set("console", console{})
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
