package runal

import "fmt"

func clearScreen() {
	fmt.Print("\x1b[2J")
}

func hideCursor() {
	fmt.Print("\x1b[25l")
}

func showCursor() {
	fmt.Print("\x1b[25h")
}

func resetCursorPosition() {
	fmt.Print("\x1b[H")
}

func enterAltScreen() {
	clearScreen()
	hideCursor()
}
