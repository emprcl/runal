package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	"github.com/fsnotify/fsnotify"

	"runal-cli/runtime"
)

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
	infile := flag.String("f", "", "sketch file (.js)")
	outfile := flag.String("o", "", "output executable file")
	demo := flag.Bool("demo", false, "demo mode")
	flag.Parse()

	embedded, err := readEmbeddedSketch()
	if err != nil {
		log.Fatalf("%v reading embedded sketch", err)
		return
	}
	if embedded != "" {
		r := runtime.New("", nil)
		r.RunInternal(embedded)
		return
	}

	if *demo {
		r := runtime.New("", nil)
		r.RunInternal(Demo)
		return
	}

	if *infile == "" {
		displayHelp()
		return
	}

	if *outfile != "" {
		err := createEmbeddedExecutable(*outfile, *infile)
		if err != nil {
			log.Fatalf("%v creating the embedded executable", err)
		}
		return
	}

	log.SetOutput(os.Stdout)

	if _, err := os.Stat(*infile); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("sketch file %s does not exist", *infile)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	r := runtime.New(*infile, watcher)
	r.Run()
}

func getVersion() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "dirty"
	}

	if buildInfo.Main.Version == "" {
		return "dirty"
	}

	return buildInfo.Main.Version
}

func displayHelp() {
	fmt.Println(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("201")).Render(logo),
			lipgloss.NewStyle().MarginLeft(5).Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.NewStyle().Foreground(lipgloss.Color("79")).Render(fmt.Sprintf("%s - https://empr.cl/runal/", getVersion())),
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
							lipgloss.JoinHorizontal(
								lipgloss.Left,
								lipgloss.NewStyle().Width(15).Render("-o [FILE]"),
								lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render("creates a standalone executable from the -f [FILE]"),
							),
						),
					),
					lipgloss.NewStyle().MarginTop(2).Bold(true).Foreground(lipgloss.Color("81")).Render("EXAMPLE"),
					lipgloss.NewStyle().MarginLeft(2).Render(
						lipgloss.JoinVertical(
							lipgloss.Left,
							"runal -f my_sketch.js",
							"runal -demo",
							"runal -f my_sketch.js -o my_executable",
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
