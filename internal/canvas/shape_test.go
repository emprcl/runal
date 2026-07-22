package canvas

import (
	"strings"
	"testing"
)

// grid renders the canvas buffer as one string per row, using a dot for
// empty cells, so expectations can be written as ascii art.
func grid(c *Canvas) []string {
	rows := make([]string, len(c.buffer))
	for y, row := range c.buffer {
		var b strings.Builder
		for _, cll := range row {
			if cll.char == 0 {
				b.WriteRune('.')
				continue
			}
			b.WriteRune(cll.char)
		}
		rows[y] = b.String()
	}
	return rows
}

func assertGrid(t *testing.T, c *Canvas, want []string) {
	t.Helper()
	got := grid(c)
	if len(got) != len(want) {
		t.Fatalf("got %d rows, want %d", len(got), len(want))
	}
	for y := range want {
		if got[y] != want[y] {
			t.Errorf("row %d:\n got %q\nwant %q", y, got[y], want[y])
		}
	}
}

func TestLine(t *testing.T) {
	tests := []struct {
		name           string
		x1, y1, x2, y2 int
		want           []string
	}{
		{
			name: "horizontal",
			x1:   1, y1: 1, x2: 4, y2: 1,
			want: []string{
				".....",
				".####",
				".....",
				".....",
				".....",
			},
		},
		{
			name: "vertical",
			x1:   2, y1: 0, x2: 2, y2: 3,
			want: []string{
				"..#..",
				"..#..",
				"..#..",
				"..#..",
				".....",
			},
		},
		{
			name: "diagonal",
			x1:   0, y1: 0, x2: 4, y2: 4,
			want: []string{
				"#....",
				".#...",
				"..#..",
				"...#.",
				"....#",
			},
		},
		{
			name: "reversed endpoints draw the same line",
			x1:   4, y1: 4, x2: 0, y2: 0,
			want: []string{
				"#....",
				".#...",
				"..#..",
				"...#.",
				"....#",
			},
		},
		{
			name: "single point",
			x1:   2, y1: 2, x2: 2, y2: 2,
			want: []string{
				".....",
				".....",
				"..#..",
				".....",
				".....",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := mockCanvas(5, 5)
			c.StrokeText("#")
			c.Line(tt.x1, tt.y1, tt.x2, tt.y2)
			assertGrid(t, c, tt.want)
		})
	}
}

func TestLineClipsOutOfBounds(t *testing.T) {
	c := mockCanvas(5, 5)
	c.StrokeText("#")
	// Must not panic when the line runs off the canvas.
	c.Line(-10, 2, 10, 2)
	assertGrid(t, c, []string{
		".....",
		".....",
		"#####",
		".....",
		".....",
	})
}

func TestRectOutline(t *testing.T) {
	c := mockCanvas(6, 6)
	c.StrokeText("#")
	c.NoFill()
	c.Rect(1, 1, 3, 3)
	assertGrid(t, c, []string{
		"......",
		".####.",
		".#..#.",
		".#..#.",
		".####.",
		"......",
	})
}

func TestRectFilled(t *testing.T) {
	c := mockCanvas(6, 6)
	c.StrokeText("#")
	c.FillText("*")
	c.Rect(1, 1, 3, 3)
	assertGrid(t, c, []string{
		"......",
		".####.",
		".#**#.",
		".#**#.",
		".####.",
		"......",
	})
}

func TestTriangleFilled(t *testing.T) {
	c := mockCanvas(7, 5)
	c.StrokeText("#")
	c.FillText("*")
	c.Triangle(3, 0, 0, 4, 6, 4)
	got := grid(c)

	// The interior scanline must be filled, not left blank.
	if !strings.Contains(got[3], "*") {
		t.Errorf("expected the triangle interior to be filled, got %q", got[3])
	}
	// And the outline must still be drawn on top.
	if !strings.Contains(got[0], "#") {
		t.Errorf("expected the apex to be stroked, got %q", got[0])
	}
}

// Multi-byte stroke characters used to be indexed by byte, which sliced
// them into garbage. Circle is the shape that walks the stroke text.
func TestCircleWithMultiByteStroke(t *testing.T) {
	c := mockCanvas(11, 11)
	c.StrokeText("█▓▒")
	c.NoFill()
	c.Circle(5, 5, 4)

	allowed := map[rune]bool{'█': true, '▓': true, '▒': true}
	drawn := 0
	for _, row := range c.buffer {
		for _, cll := range row {
			if cll.char == 0 {
				continue
			}
			drawn++
			if !allowed[cll.char] {
				t.Fatalf("circle drew %q, which is not one of the stroke runes", cll.char)
			}
		}
	}
	if drawn == 0 {
		t.Fatal("circle drew nothing")
	}
}

func TestTextWritesRunesNotBytes(t *testing.T) {
	c := mockCanvas(6, 1)
	c.StrokeText("#")
	c.Text("héllo", 0, 0)
	assertGrid(t, c, []string{"héllo."})
}

func TestRuneAt(t *testing.T) {
	runes := []rune("█▓▒")

	tests := []struct {
		index int
		want  rune
	}{
		{0, '█'},
		{1, '▓'},
		{2, '▒'},
		{3, '█'},
		{7, '▓'},
	}

	for _, tt := range tests {
		if got := runeAt(runes, tt.index); got != tt.want {
			t.Errorf("runeAt(%d) = %q, want %q", tt.index, got, tt.want)
		}
	}

	if got := runeAt(nil, 0); got != defaultPaddingRune {
		t.Errorf("runeAt(nil, 0) = %q, want the padding rune", got)
	}
}

func TestFirstRune(t *testing.T) {
	if got := firstRune("abc", 'x'); got != 'a' {
		t.Errorf("firstRune(\"abc\") = %q, want 'a'", got)
	}
	if got := firstRune("█▓", 'x'); got != '█' {
		t.Errorf("firstRune(\"█▓\") = %q, want '█'", got)
	}
	if got := firstRune("", 'x'); got != 'x' {
		t.Errorf("firstRune(\"\") = %q, want the fallback", got)
	}
}

// An empty string from a sketch used to panic on []rune(s)[0].
func TestEmptyStringsDoNotPanic(t *testing.T) {
	c := mockCanvas(5, 5)
	c.CellModeCustom("")
	c.Set(0, 0, Cell{Char: "", Foreground: "red", Background: "black"})
}
