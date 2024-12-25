package main

import "fmt"

func ClearScreen() {
	fmt.Print("\x1b[2J")
}

func HideCursor() {
	fmt.Print("\x1b[25l")
}

func ResetCursorPosition() {
	fmt.Print("\x1b[H")
}

func EnterAltScreen() {
	ClearScreen()
	HideCursor()
}
