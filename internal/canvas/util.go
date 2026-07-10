package canvas

import (
	"strings"
)

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

func strIndex(str string, index int) rune {
	return rune(str[index%len(str)])
}
