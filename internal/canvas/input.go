package canvas

import (
	"context"
	"os"
	"sync"

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

// inputReader owns the raw terminal state and the stdin event reader.
// It must be closed to restore the terminal and unblock the pump goroutine.
type inputReader struct {
	reader    *input.Reader
	termState *term.State
	events    chan input.Event
	done      chan struct{}
	listening bool
	closeOnce sync.Once
}

func newInputReader() (*inputReader, error) {
	termState, err := term.MakeRaw(os.Stdin.Fd())
	if err != nil {
		return nil, err
	}

	reader, err := input.NewReader(os.Stdin, os.Getenv("TERM"), input.FlagMouseMode)
	if err != nil {
		term.Restore(os.Stdin.Fd(), termState) // nolint: errcheck
		return nil, err
	}

	return &inputReader{
		reader:    reader,
		termState: termState,
		events:    make(chan input.Event, 2048),
		done:      make(chan struct{}),
	}, nil
}

func (i *inputReader) listen(ctx context.Context, wg *sync.WaitGroup) {
	i.listening = true
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(i.events)
		defer close(i.done)
		for {
			events, err := i.reader.ReadEvents()
			if err != nil {
				return
			}
			for _, ev := range events {
				select {
				case i.events <- ev:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
}

// close cancels any pending read and restores the terminal mode.
// It is safe to call more than once.
func (i *inputReader) close() {
	i.closeOnce.Do(func() {
		i.reader.Cancel()
		if i.listening {
			<-i.done
		}
		i.reader.Close()                         // nolint: errcheck
		term.Restore(os.Stdin.Fd(), i.termState) // nolint: errcheck
	})
}
