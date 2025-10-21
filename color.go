package runal

import (
	"fmt"
	col "image/color"
)

// func color(color string) lipgloss.Color {
// 	if strings.HasPrefix(color, "#") {
// 		return lipgloss.Color(color)
// 	}
// 	c, err := strconv.ParseFloat(strings.TrimSpace(color), 64)
// 	if err != nil {
// 		return lipgloss.Color(color)
// 	}
// 	return lipgloss.Color(strconv.Itoa(int(math.Round(c))))
// }

// func colorFromImage(c col.Color) lipgloss.Color {
// 	rgba := col.RGBAModel.Convert(c).(col.RGBA)
// 	return lipgloss.Color(fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B))
// }

func colorStringFromImage(c col.Color) string {
	rgba := col.RGBAModel.Convert(c).(col.RGBA)
	return fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}
