package canvas

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"

	"golang.org/x/term"
)

func tryTermSize() (int, int, error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

func writePadding(s *strings.Builder, sLength, tLength int, padChar rune) {
	for range tLength - sLength {
		s.WriteRune(padChar)
	}
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// runeAt returns the rune at index, wrapping around the slice.
func runeAt(runes []rune, index int) rune {
	if len(runes) == 0 {
		return defaultPaddingRune
	}
	return runes[index%len(runes)]
}

func firstRune(s string, fallback rune) rune {
	for _, r := range s {
		return r
	}
	return fallback
}

func randomDir() string {
	tmp, err := os.MkdirTemp(os.TempDir(), "runal")
	if err != nil {
		log.Fatalf("error creating temporary directory: %v", err)
	}
	return tmp
}
