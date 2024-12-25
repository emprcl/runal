package main

import (
	"log"
	"strings"

	"golang.org/x/term"
)

func termSize() (int, int) {
	w, h, err := term.GetSize(0)
	if err != nil {
		log.Fatal("can't read terminal size")
	}
	return w, h
}

func forceLength(s string, length int, padChar rune) string {
	if len(s) > length {
		return s[:length]
	} else if len(s) < length {
		padding := strings.Repeat(string(padChar), length-len(s))
		return s + padding
	}
	return s
}
