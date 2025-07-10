package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	"github.com/fsnotify/fsnotify"

	"github.com/emprcl/runal/cli/runtime"
)

//go:embed VERSION
var AppVersion string

//go:embed demo.js
var Demo string

const (
	logo = `
      ___           ___           ___           ___           ___
     /\  \         /\__\         /\__\         /\  \         /\__\
    /::\  \       /:/  /        /::|  |       /::\  \       /:/  /
   /:/\:\  \     /:/  /        /:|:|  |      /:/\:\  \     /:/  /
  /::\~\:\  \   /:/  /  ___   /:/|:|  |__   /::\~\:\  \   /:/  /
 /:/\:\ \:\__\ /:/__/  /\__\ /:/ |:| /\__\ /:/\:\ \:\__\ /:/__/
 \/_|::\/:/  / \:\  \ /:/  / \/__|:|/:/  / \/__\:\/:/  / \:\  \
    |:|::/  /   \:\  /:/  /      |:/:/  /       \::/  /   \:\  \
    |:|\/__/     \:\/:/  /       |::/  /        /:/  /     \:\  \
    |:|  |        \::/  /        /:/  /        /:/  /       \:\__\
     \|__|         \/__/         \/__/         \/__/         \/__/`

	defaultLogFile = "console.log"
)

func main() {
	file := flag.String("f", "", "sketch file (.js)")
	demo := flag.Bool("demo", false, "demo mode")
	flag.Parse()

	if *demo {
		r := runtime.New("", nil, nil)
		r.RunDemo(Demo)
		return
	}

	if *file == "" {
		displayHelp()
		return
	}

	if _, err := os.Stat(*file); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("sketch file %s does not exist", *file)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	consoleLogFile, err := os.Create(defaultLogFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		consoleLogFile.Close()
		err := os.Remove(defaultLogFile)
		if err != nil {
			log.Fatal(err)
		}
	}()

	r := runtime.New(*file, watcher, consoleLogFile)
	r.Run()
}

func displayHelp() {
	fmt.Println(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("201")).Render(logo),
			lipgloss.NewStyle().MarginLeft(5).Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.NewStyle().Foreground(lipgloss.Color("79")).Render(fmt.Sprintf("%s - https://empr.cl/runal/", strings.TrimSuffix(AppVersion, "\n"))),
					lipgloss.NewStyle().MarginTop(2).Bold(true).Foreground(lipgloss.Color("81")).Render("USAGE"),
					lipgloss.NewStyle().MarginLeft(2).Render(
						lipgloss.JoinVertical(
							lipgloss.Left,
							lipgloss.JoinHorizontal(
								lipgloss.Left,
								lipgloss.NewStyle().Width(15).Render("-f [FILE]"),
								lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render("the javascript sketch file (.js) to watch"),
							),
							lipgloss.JoinHorizontal(
								lipgloss.Left,
								lipgloss.NewStyle().Width(15).Render("-demo"),
								lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render("demo sketch (press space to reseed, c to capture png)"),
							),
						),
					),
					lipgloss.NewStyle().MarginTop(2).Bold(true).Foreground(lipgloss.Color("81")).Render("EXAMPLE"),
					lipgloss.NewStyle().MarginLeft(2).Render(
						lipgloss.JoinVertical(
							lipgloss.Left,
							"runal -f my_sketch.js",
							"runal -demo",
						),
					),
					lipgloss.NewStyle().MarginTop(2).Bold(true).Foreground(lipgloss.Color("81")).Render("LOGS"),
					lipgloss.NewStyle().MarginLeft(2).MarginBottom(2).Render(
						lipgloss.JoinVertical(
							lipgloss.Left,
							lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render(
								lipgloss.JoinVertical(lipgloss.Left,
									"console.log() messages are written to a console.log file",
									"that is deleted upon exit. You can watch these logs with",
									"tail is another terminal window/pane/tab.",
								),
							),
							"tail -f console.log",
						),
					),
				),
			),
		),
	)
}
