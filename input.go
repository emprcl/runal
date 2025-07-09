package runal

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/input"
	"github.com/charmbracelet/x/term"
)

type KeyEvent struct {
	Key  string
	Code int
}

type MouseEvent struct {
	X      int
	Y      int
	Button string
}

func listenForInputEvents(ctx context.Context) chan input.Event {
	inputEvents := make(chan input.Event, 2048)
	go func() {
		defer close(inputEvents)

		state, err := term.MakeRaw(os.Stdin.Fd())
		if err != nil {
			log.Fatal(err)
		}
		defer term.Restore(os.Stdin.Fd(), state)

		reader, err := input.NewReader(os.Stdin, os.Getenv("TERM"), 0)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()

		for {
			events, err := reader.ReadEvents()
			if err != nil {
				log.Fatal(err)
				continue
			}
			for _, ev := range events {
				select {
				case inputEvents <- ev:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return inputEvents
}
