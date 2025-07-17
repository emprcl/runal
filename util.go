package runal

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"

	"golang.org/x/term"
)

func termSize() (int, int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal("can't read terminal size")
	}
	return w, h
}

func forcePadding(s *strings.Builder, sLength, tLength int, padChar rune) {
	if sLength >= tLength {
		return
	}
	padding := strings.Repeat(string(padChar), tLength-sLength)
	s.WriteString(padding)
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func strIndex(str string, index int) rune {
	return rune(str[index%len(str)])
}

func randomDir() string {
	tmp, err := os.MkdirTemp(os.TempDir(), "runal")
	if err != nil {
		log.Fatalf("error creating temporary directory: %v", err)
	}
	return tmp
}
