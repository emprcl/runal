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
	file := flag.String("f", "", "sketch file (.js)")
	fps := flag.Int("fps", 60, "frame per seconds")
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

	r := runtime.New(*file, watcher, consoleLogFile, *fps)
	r.Run()
}
