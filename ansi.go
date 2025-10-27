package runal

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func clearScreen() {
	fmt.Print(ansi.ResetAltScreenMode)
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
	fmt.Print(ansi.SetAltScreenMode)
	hideCursor()
}

func enableMouse() {
	fmt.Print(ansi.SetAnyEventMouseMode)
	fmt.Print(ansi.SetSgrExtMouseMode)
}

func disableMouse() {
	fmt.Print(ansi.ResetAnyEventMouseMode)
	fmt.Print(ansi.ResetSgrExtMouseMode)
}

func clearLineSequence() string {
	return "\r\n\x1b[2K"
}
