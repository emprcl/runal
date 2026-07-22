package canvas

import "testing"

func TestColorHSL(t *testing.T) {
	c := mockCanvas(10, 10)

	tests := []struct {
		name    string
		h, s, l int
		want    string
	}{
		{"red", 0, 100, 50, "#ff0000"},
		{"green", 120, 100, 50, "#00ff00"},
		{"blue", 240, 100, 50, "#0000ff"},
		{"yellow", 60, 100, 50, "#ffff00"},
		{"cyan", 180, 100, 50, "#00ffff"},
		{"magenta", 300, 100, 50, "#ff00ff"},
		{"white", 0, 0, 100, "#ffffff"},
		{"black", 0, 0, 0, "#000000"},
		{"mid grey", 0, 0, 50, "#808080"},
		{"half saturation", 0, 50, 50, "#bf4040"},
		// A full turn is the same hue as no turn at all.
		{"hue wraps at 360", 360, 100, 50, "#ff0000"},
		{"hue clamped above 360", 400, 100, 50, "#ff0000"},
		{"saturation clamped", 120, 500, 50, "#00ff00"},
		{"negative lightness clamped", 120, 100, -20, "#000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.ColorHSL(tt.h, tt.s, tt.l); got != tt.want {
				t.Errorf("ColorHSL(%d, %d, %d) = %s, want %s", tt.h, tt.s, tt.l, got, tt.want)
			}
		})
	}
}

func TestColorHSV(t *testing.T) {
	c := mockCanvas(10, 10)

	tests := []struct {
		name    string
		h, s, v int
		want    string
	}{
		{"red", 0, 100, 100, "#ff0000"},
		{"green", 120, 100, 100, "#00ff00"},
		{"blue", 240, 100, 100, "#0000ff"},
		{"cyan", 180, 100, 100, "#00ffff"},
		{"magenta", 300, 100, 100, "#ff00ff"},
		{"white", 0, 0, 100, "#ffffff"},
		{"black", 0, 0, 0, "#000000"},
		{"hue wraps at 360", 360, 100, 100, "#ff0000"},
		{"value clamped", 0, 100, 500, "#ff0000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.ColorHSV(tt.h, tt.s, tt.v); got != tt.want {
				t.Errorf("ColorHSV(%d, %d, %d) = %s, want %s", tt.h, tt.s, tt.v, got, tt.want)
			}
		})
	}
}

func TestColorRGBClamps(t *testing.T) {
	c := mockCanvas(10, 10)

	tests := []struct {
		name    string
		r, g, b int
		want    string
	}{
		{"in range", 18, 52, 86, "#123456"},
		{"above range", 300, 300, 300, "#ffffff"},
		{"below range", -10, -1, -255, "#000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.ColorRGB(tt.r, tt.g, tt.b); got != tt.want {
				t.Errorf("ColorRGB(%d, %d, %d) = %s, want %s", tt.r, tt.g, tt.b, got, tt.want)
			}
		})
	}
}

func TestColorParsing(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"hex", "#ff8800", "#ff8800"},
		{"css name", "red", "#ff0000"},
		{"css name is case insensitive", "ReBeCcApUrPlE", "#663399"},
		{"ansi black", "0", "#000000"},
		{"ansi red", "1", "#800000"},
		{"ansi white", "15", "#ffffff"},
		{"ansi with whitespace", " 9 ", "#ff0000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := colorToString(color(tt.input)); got != tt.want {
				t.Errorf("color(%q) = %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}

func TestColorIndexedIsClamped(t *testing.T) {
	// An out-of-range index used to wrap around via uint8 conversion,
	// so 300 silently became colour 44.
	if got, want := color("300"), color("255"); got != want {
		t.Errorf("color(\"300\") = %v, want it clamped to %v", got, want)
	}
	if got, want := color("-5"), color("0"); got != want {
		t.Errorf("color(\"-5\") = %v, want it clamped to %v", got, want)
	}
}

func TestClamp(t *testing.T) {
	if got := clamp(5, 0, 10); got != 5 {
		t.Errorf("clamp(5, 0, 10) = %d, want 5", got)
	}
	if got := clamp(-5, 0, 10); got != 0 {
		t.Errorf("clamp(-5, 0, 10) = %d, want 0", got)
	}
	if got := clamp(50, 0, 10); got != 10 {
		t.Errorf("clamp(50, 0, 10) = %d, want 10", got)
	}
	if got := clamp(1.5, 0.0, 1.0); got != 1.0 {
		t.Errorf("clamp(1.5, 0.0, 1.0) = %v, want 1.0", got)
	}
}
