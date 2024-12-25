package main

import (
	"log"

	"golang.org/x/term"
)

func termSize() (int, int) {
	w, h, err := term.GetSize(0)
	if err != nil {
		log.Fatal("can't read terminal size")
	}
	return w, h
}
