package runal

import tea "github.com/charmbracelet/bubbletea"

type event struct {
	name  string
	value int
}

func newFPSEvent(fps int) event {
	return event{
		name:  "fps",
		value: fps,
	}
}

func newRedrawEvent() event {
	return event{
		name: "redraw",
	}
}

type MouseEvent tea.MouseEvent

type KeyEvent struct {
	Key     string
	KeyCode int
}
