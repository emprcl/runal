package runal

import (
	"context"
	"os"
	"sync"

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

func listenForInputEvents(ctx context.Context, wg *sync.WaitGroup) chan input.Event {
	inputEvents := make(chan input.Event, 2048)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(inputEvents)

		state, err := term.MakeRaw(os.Stdin.Fd())
		if err != nil {
			log.Fatal(err)
		}
		defer term.Restore(os.Stdin.Fd(), state)

		reader, err := input.NewReader(os.Stdin, os.Getenv("TERM"), input.FlagMouseMode)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()

		readerEvents := make(chan []input.Event, 8)
		readerErrors := make(chan error, 8)
		go func() {
			for {
				events, err := reader.ReadEvents()
				if err != nil {
					readerErrors <- err
					return
				}
				readerEvents <- events
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case events := <-readerEvents:
				for _, ev := range events {
					inputEvents <- ev
				}
			case err := <-readerErrors:
				log.Fatal(err)
			}
		}
	}()
	return inputEvents
}
