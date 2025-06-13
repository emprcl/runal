package runal

import (
	"math"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func color(color string) lipgloss.Color {
	if strings.HasPrefix(color, "#") {
		return lipgloss.Color(color)
	}
	c, err := strconv.ParseFloat(strings.TrimSpace(color), 64)
	if err != nil {
		return lipgloss.Color(color)
	}
	return lipgloss.Color(strconv.Itoa(int(math.Round(c))))
}
