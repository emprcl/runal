package main

import (
	"errors"
	"flag"
	"os"

	"github.com/charmbracelet/log"

	"github.com/fsnotify/fsnotify"

	"runal-cli/runtime"
)

func main() {
	// p := perlin.NewPerlin(2., 2., 3, time.Now().UnixNano())
	// for x := 0.; x < 3; x++ {
	// 	for y := 0.; y < 3; y++ {
	// 		fmt.Printf("%0.0f;%0.0f;%0.4f\n", x, y, p.Noise2D(x*0.05, y*0.05))
	// 	}

	// }
	// os.Exit(0)
	file := flag.String("f", "", "sketch file (.js)")
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

	consoleLogFile, err := os.Create("console.log")
	if err != nil {
		log.Fatal(err)
	}
	defer consoleLogFile.Close()

	r := runtime.New(*file, watcher, consoleLogFile)
	r.Run()
}
