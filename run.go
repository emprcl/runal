package runal

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

const (
	defaultFPS = 30
)

func Run(ctx context.Context, setup, draw func(c *Canvas), onKey func(c *Canvas, key string)) {
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	Start(ctx, nil, setup, draw, onKey).Wait()
}

type model struct {
	canvas      *Canvas
	done        chan struct{}
	setup, draw func(c *Canvas)
	onKey       func(c *Canvas, key string)
	framerate   time.Duration
}

func newModel(done chan struct{}, setup, draw func(c *Canvas), onKey func(c *Canvas, key string)) model {
	w, h := termSize()
	return model{
		canvas:    newCanvas(w, h),
		done:      done,
		setup:     setup,
		draw:      draw,
		onKey:     onKey,
		framerate: newFramerate(defaultFPS),
	}
}

type tickMsg time.Time

func tick(m model) tea.Cmd {
	return tea.Tick(m.framerate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func newFramerate(fps int) time.Duration {
	return time.Second / time.Duration(fps)
}

type renderMsg struct{}

func render() tea.Cmd {
	return func() tea.Msg {
		return struct{}{}
	}
}

func (m model) Init() tea.Cmd {
	m.setup(m.canvas)
	return tea.Batch(tea.EnterAltScreen, tick(m))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// handle events coming from the canvas
	select {
	case event := <-m.canvas.bus:
		switch event.name {
		case "fps":
			m.framerate = newFramerate(event.value)
		}
	default:
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.canvas.termWidth = msg.Width
		m.canvas.termHeight = msg.Height
		if m.canvas.autoResize {
			m.canvas.resize(msg.Width, msg.Height)
		}
		return m, nil

	case tickMsg:
		if !m.canvas.IsLooping {
			return m, nil
		}
		return m, tick(m)

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			if m.done != nil {
				m.done <- struct{}{}
			}
			return m, tea.Quit
		default:
			if m.onKey != nil {
				m.onKey(m.canvas, msg.String())
			}
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	m.draw(m.canvas)
	return m.canvas.render()
}

func Start(ctx context.Context, done chan struct{}, setup, draw func(c *Canvas), onKey func(c *Canvas, key string)) *sync.WaitGroup {
	p := tea.NewProgram(
		newModel(done, setup, draw, onKey),
		tea.WithContext(ctx),
		tea.WithFPS(defaultFPS),
	)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := p.Run(); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Errorf("Error: %v", err)
		}
	}()

	return &wg
}
