package runal

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func clearScreen() {
	fmt.Print(ansi.ResetModeAltScreen)
}

func hideCursor() {
	fmt.Print(ansi.HideCursor)
}

func showCursor() {
	fmt.Print(ansi.ShowCursor)
}

func resetCursorPosition() {
	fmt.Print(ansi.CursorHomePosition)
}

func enterAltScreen() {
	fmt.Print(ansi.SetModeAltScreen)
	hideCursor()
}

func enableMouse() {
	fmt.Print(ansi.SetModeMouseAnyEvent)
	fmt.Print(ansi.SetModeMouseExtSgr)
}

func disableMouse() {
	fmt.Print(ansi.ResetModeMouseAnyEvent)
	fmt.Print(ansi.ResetModeMouseExtSgr)
}

func clearLineSequence() string {
	return "\r\n\x1b[2K"
}
