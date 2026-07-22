package canvas

import (
	"fmt"
	"io"

	"github.com/charmbracelet/x/ansi"
)

func exitAltScreen(w io.Writer) {
	fmt.Fprint(w, ansi.ResetModeAltScreen)
}

func eraseScreen(w io.Writer) {
	fmt.Fprint(w, ansi.EraseEntireScreen)
}

func hideCursor(w io.Writer) {
	fmt.Fprint(w, ansi.HideCursor)
}

func showCursor(w io.Writer) {
	fmt.Fprint(w, ansi.ShowCursor)
}

func resetCursorPosition(w io.Writer) {
	fmt.Fprint(w, ansi.CursorHomePosition)
}

func enterAltScreen(w io.Writer) {
	fmt.Fprint(w, ansi.SetModeAltScreen)
	hideCursor(w)
}

func enableMouse(w io.Writer) {
	fmt.Fprint(w, ansi.SetModeMouseAnyEvent)
	fmt.Fprint(w, ansi.SetModeMouseExtSgr)
}

func disableMouse(w io.Writer) {
	fmt.Fprint(w, ansi.ResetModeMouseAnyEvent)
	fmt.Fprint(w, ansi.ResetModeMouseExtSgr)
}

const clearLineStr = "\r\n\x1b[2K"
