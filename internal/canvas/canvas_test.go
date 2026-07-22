package canvas

import (
	"fmt"
	"math"
	"testing"
	"time"
)

var benchmarks = []struct {
	size int
}{
	{size: 20},
	{size: 100},
	{size: 200},
	{size: 300},
	{size: 400},
	{size: 500},
}

func BenchmarkCanvasRender(b *testing.B) {
	for _, v := range benchmarks {
		b.Run(fmt.Sprintf("canvas_size_%dx%d", v.size, v.size), func(b *testing.B) {
			canvas := mockCanvas(v.size, v.size)
			canvas.Rect(5, 5, 5, 5)
			canvas.Clear()
			for i := 0; i < b.N; i++ {
				canvas.render()
			}
		})
	}
}

// Push/Pop is a stack: a nested pair must restore the outer state, not
// collapse it into the innermost one.
func TestPushPopNests(t *testing.T) {
	c := mockCanvas(10, 10)

	c.Translate(1, 1)
	c.Push()

	c.Translate(10, 10)
	c.Push()

	c.Translate(100, 100)
	if c.originX != 111 {
		t.Fatalf("originX = %d, want 111", c.originX)
	}

	c.Pop()
	if c.originX != 11 || c.originY != 11 {
		t.Errorf("after first Pop: origin = (%d, %d), want (11, 11)", c.originX, c.originY)
	}

	c.Pop()
	if c.originX != 1 || c.originY != 1 {
		t.Errorf("after second Pop: origin = (%d, %d), want (1, 1)", c.originX, c.originY)
	}
}

func TestPopWithoutPushIsNoop(t *testing.T) {
	c := mockCanvas(10, 10)
	c.Translate(3, 4)
	c.Pop()
	if c.originX != 3 || c.originY != 4 {
		t.Errorf("origin = (%d, %d), want it untouched at (3, 4)", c.originX, c.originY)
	}
}

func TestPushRestoresStyle(t *testing.T) {
	c := mockCanvas(10, 10)
	c.StrokeText("#")
	c.Push()
	c.StrokeText("@")
	c.Rotate(math.Pi)
	c.Scale(3)
	c.Pop()

	if c.strokeText != "#" {
		t.Errorf("strokeText = %q, want %q", c.strokeText, "#")
	}
	if c.rotationAngle != 0 {
		t.Errorf("rotationAngle = %v, want 0", c.rotationAngle)
	}
	if c.scale != 1 {
		t.Errorf("scale = %v, want 1", c.scale)
	}
}

// An unbalanced Push in draw() must not grow the stack every frame.
func TestStateStackResetsEachFrame(t *testing.T) {
	c := mockCanvas(5, 5)
	for range 10 {
		c.Push()
		c.render()
	}
	if len(c.stateStack) != 0 {
		t.Errorf("stateStack has %d entries after rendering, want 0", len(c.stateStack))
	}
}

// A zero or negative framerate used to divide by zero, and ticker.Reset
// panics on a non-positive duration.
func TestNewFramerate(t *testing.T) {
	tests := []struct {
		fps  int
		want time.Duration
	}{
		{30, time.Second / 30},
		{60, time.Second / 60},
		{0, time.Second},
		{-10, time.Second},
		{100000, time.Second / maxFPS},
	}

	for _, tt := range tests {
		got := newFramerate(tt.fps)
		if got != tt.want {
			t.Errorf("newFramerate(%d) = %v, want %v", tt.fps, got, tt.want)
		}
		if got <= 0 {
			t.Errorf("newFramerate(%d) = %v, which would panic ticker.Reset", tt.fps, got)
		}
	}
}

func TestFpsIsClamped(t *testing.T) {
	c := mockCanvas(5, 5)

	c.Fps(0)
	if c.fps < minFPS {
		t.Errorf("Fps(0) stored %d, want at least %d", c.fps, minFPS)
	}
	c.Fps(100000)
	if c.fps > maxFPS {
		t.Errorf("Fps(100000) stored %d, want at most %d", c.fps, maxFPS)
	}
}

// LoopAngle(0) used to divide by zero via Framecount % totalFrames.
func TestLoopAngleZeroDuration(t *testing.T) {
	c := mockCanvas(5, 5)
	c.Framecount = 7
	got := c.LoopAngle(0)
	if math.IsNaN(got) {
		t.Errorf("LoopAngle(0) = NaN")
	}
}

func TestGifDelay(t *testing.T) {
	tests := []struct {
		fps  int
		want int
	}{
		{30, 3},
		{25, 4},
		{50, 2},
		// Above 100fps the delay rounds to zero, which makes viewers
		// fall back to their own default speed.
		{200, 1},
		{0, 100},
	}

	for _, tt := range tests {
		if got := gifDelay(tt.fps); got != tt.want {
			t.Errorf("gifDelay(%d) = %d, want %d", tt.fps, got, tt.want)
		}
		if gifDelay(tt.fps) < 1 {
			t.Errorf("gifDelay(%d) must never be zero", tt.fps)
		}
	}
}

func TestMap(t *testing.T) {
	c := mockCanvas(5, 5)
	tests := []struct {
		name  string
		value float64
		want  float64
	}{
		{"start", 0, 0},
		{"middle", 5, 50},
		{"end", 10, 100},
		{"beyond range extrapolates", 20, 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Map(tt.value, 0, 10, 0, 100); got != tt.want {
				t.Errorf("Map(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestDist(t *testing.T) {
	c := mockCanvas(5, 5)
	if got := c.Dist(0, 0, 3, 4); got != 5 {
		t.Errorf("Dist(0,0,3,4) = %v, want 5", got)
	}
	if got := c.Dist(2, 2, 2, 2); got != 0 {
		t.Errorf("Dist of a point to itself = %v, want 0", got)
	}
}

func TestFrameSizeOnEmptyFrame(t *testing.T) {
	// A zero-sized frame used to panic on f[0].
	w, h := newFrame(0, 0).size()
	if w != 0 || h != 0 {
		t.Errorf("empty frame size = (%d, %d), want (0, 0)", w, h)
	}
}

func TestGetOutOfBoundsRegion(t *testing.T) {
	c := mockCanvas(5, 5)
	c.StrokeText("#")
	c.Point(0, 0)

	// Must clip rather than panic.
	img := c.Get(3, 3, 10, 10)
	if img == nil {
		t.Fatal("Get returned nil")
	}
	// Reading a cell that was never drawn used to nil-deref, because a
	// zero cell carries no colors.
	if got := img.Cell(0, 0); got.Char != "" {
		t.Errorf("undrawn cell = %q, want an empty Cell", got.Char)
	}
	if got := img.Cell(100, 100); got.Char != "" {
		t.Errorf("out-of-range cell = %q, want an empty Cell", got.Char)
	}
}

// A non-positive duration used to make a slice with a negative capacity,
// which panics inside the sketch rather than reporting anything useful.
func TestCaptureFrames(t *testing.T) {
	tests := []struct {
		duration, fps int
		want          int
	}{
		{2, 30, 60},
		{1, 30, 30},
		{0, 30, 30},
		{-5, 30, 30},
		{1, 0, 1},
		{1, 10000, maxFPS},
	}
	for _, tt := range tests {
		if got := captureFrames(tt.duration, tt.fps); got != tt.want {
			t.Errorf("captureFrames(%d, %d) = %d, want %d",
				tt.duration, tt.fps, got, tt.want)
		}
	}
}

// A zero or negative duration must not panic when it reaches make().
func TestSaveCanvasToGIFNonPositiveDuration(t *testing.T) {
	for _, duration := range []int{0, -1} {
		c := mockCanvas(5, 5)
		c.SaveCanvasToGIF("out.gif", duration)
		if cap(c.frames) < 1 {
			t.Errorf("duration %d gave a %d-frame buffer, want at least 1",
				duration, cap(c.frames))
		}
	}
}

// A terminal narrower than the capture padding must still build a
// converter: a nil one nil-panics on the next export.
func TestCaptureConfigStaysPositive(t *testing.T) {
	for _, size := range []int{0, 1, 2} {
		config := newCaptureConfig(size, size)
		if config.PageCols < 1 || config.PageRows < 1 {
			t.Errorf("size %d gave PageCols=%d PageRows=%d, want both >= 1",
				size, config.PageCols, config.PageRows)
		}
		if _, err := newCapture(size, size); err != nil {
			t.Errorf("newCapture(%d, %d) failed: %v", size, size, err)
		}
	}
}

// Dropping an unbalanced Push silently hides the missing Pop.
func TestUnbalancedPushWarnsOnce(t *testing.T) {
	c := mockCanvas(5, 5)
	for range 5 {
		c.Push()
		c.render()
	}
	if len(c.debugBuffer) != 1 {
		t.Errorf("debugBuffer has %d entries, want exactly 1", len(c.debugBuffer))
	}
}

// A balanced sketch must stay quiet.
func TestBalancedPushDoesNotWarn(t *testing.T) {
	c := mockCanvas(5, 5)
	for range 5 {
		c.Push()
		c.Pop()
		c.render()
	}
	if len(c.debugBuffer) != 0 {
		t.Errorf("debugBuffer = %v, want empty", c.debugBuffer)
	}
}
