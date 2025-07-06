package runal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	defaultFPS = 30
)

func Run(ctx context.Context, setup, draw func(c *Canvas), onKey func(c *Canvas, key string)) {
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	Start(ctx, nil, setup, draw, onKey).Wait()
}

// func Start(ctx context.Context, done chan os.Signal, setup, draw func(c *Canvas), onKey func(c *Canvas, key string)) *sync.WaitGroup {
// 	w, h := termSize()
// 	c := newCanvas(w, h)

// 	resize := listenForResize()

// 	enterAltScreen()

// 	setup(c)
// 	render := func() {
// 		resetCursorPosition()
// 		draw(c)
// 		c.render()
// 	}
// 	render()

// 	ticker := time.NewTicker(newFramerate(defaultFPS))

// 	exit := func() {
// 		ticker.Stop()
// 		resetCursorPosition()
// 		clearScreen()
// 		showCursor()
// 	}

// 	wg := sync.WaitGroup{}
// 	wg.Add(1)
// 	go func() {
// 		defer func() {
// 			wg.Done()
// 			_ = keyboard.Close()
// 		}()
// 		keyEvent, _ := keyboard.GetKeys(1)
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				exit()
// 				return
// 			case <-resize:
// 				clearScreen()
// 				w, h := termSize()
// 				c.termWidth = w
// 				c.termHeight = h
// 				if c.autoResize {
// 					c.resize(w, h)
// 				}
// 				render()
// 			case event := <-c.bus:
// 				switch event.name {
// 				case "fps":
// 					ticker.Reset(newFramerate(event.value))
// 				case "stop":
// 					ticker.Stop()
// 				case "start":
// 					ticker.Reset(newFramerate(defaultFPS))
// 				case "render":
// 					render()
// 				}
// 			case event := <-keyEvent:
// 				// ctrl+c
// 				if event.Key == keyboard.KeyCtrlC {
// 					exit()
// 					if done != nil {
// 						done <- os.Interrupt
// 					}
// 					return
// 				}
// 				// NOTE: keyboard package has a small bug on
// 				// space key not filling the Rune attribute.
// 				if event.Key == keyboard.KeySpace {
// 					event.Rune = ' '
// 				}
// 				if onKey != nil {
// 					onKey(c, string(event.Rune))
// 				}
// 			case <-ticker.C:
// 				render()
// 			}
// 		}
// 	}()

// 	return &wg
// }

// func newFramerate(fps int) time.Duration {
// 	return time.Second / time.Duration(fps)
// }

type model struct {
	canvas      *Canvas
	setup, draw func(c *Canvas)
	onKey       func(c *Canvas, key string)
}

func newModel(setup, draw func(c *Canvas), onKey func(c *Canvas, key string)) model {
	return model{
		canvas: newCanvas(80, 40),
		setup:  setup,
		draw:   draw,
		onKey:  onKey,
	}
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(33*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	m.setup(m.canvas)
	return tea.Batch(tea.EnterAltScreen, tick())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.canvas.autoResize {
			m.canvas.resize(msg.Width, msg.Height)
		}
		return m, nil

	case tickMsg:
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	m.draw(m.canvas)
	return m.canvas.render()
}

func Start(ctx context.Context, done chan os.Signal, setup, draw func(c *Canvas), onKey func(c *Canvas, key string)) *sync.WaitGroup {
	p := tea.NewProgram(
		newModel(setup, draw, onKey),
		tea.WithContext(ctx),
		tea.WithFPS(defaultFPS),
	)

	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	}()

	return &wg
}
