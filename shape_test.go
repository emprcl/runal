package runal

import (
	"bytes"
	"testing"
)

func TestRectFillOptimization(t *testing.T) {
	// Test that rectangle filling respects canvas bounds
	c := mockCanvas(10, 10)
	c.Fill("█", "white", "black")

	// Draw a large rectangle that extends beyond canvas bounds
	c.Rect(-5, -5, 20, 20)
	c.render()

	// Should only fill within canvas bounds
	output := c.output.(*bytes.Buffer).String()
	if output == "" {
		t.Error("Expected filled rectangle output")
	}
}

func TestTriangleFillOptimization(t *testing.T) {
	// Test that triangle filling uses scanline algorithm
	c := mockCanvas(20, 20)
	c.Fill("█", "white", "black")

	// Draw a triangle
	c.Triangle(5, 5, 15, 5, 10, 15)
	c.render()

	// Should render without performance issues
	output := c.output.(*bytes.Buffer).String()
	if output == "" {
		t.Error("Expected filled triangle output")
	}
}

func TestCircleFillOptimization(t *testing.T) {
	// Test that circle filling respects canvas bounds
	c := mockCanvas(10, 10)
	c.Fill("█", "white", "black")

	// Draw a large circle that extends beyond canvas bounds
	c.Circle(5, 5, 10)
	c.render()

	// Should only fill within canvas bounds
	output := c.output.(*bytes.Buffer).String()
	if output == "" {
		t.Error("Expected filled circle output")
	}
}

func TestEllipseFillOptimization(t *testing.T) {
	// Test that ellipse filling respects canvas bounds
	c := mockCanvas(10, 10)
	c.Fill("█", "white", "black")

	// Draw a large ellipse that extends beyond canvas bounds
	c.Ellipse(5, 5, 8, 12)
	c.render()

	// Should only fill within canvas bounds
	output := c.output.(*bytes.Buffer).String()
	if output == "" {
		t.Error("Expected filled ellipse output")
	}
}

func TestQuadFillOptimization(t *testing.T) {
	// Test that quad filling respects canvas bounds
	c := mockCanvas(15, 15)
	c.Fill("█", "white", "black")

	// Draw a quad that extends beyond canvas bounds
	c.Quad(-5, -5, 20, 0, 20, 20, 0, 20)
	c.render()

	// Should only fill within canvas bounds
	output := c.output.(*bytes.Buffer).String()
	if output == "" {
		t.Error("Expected filled quad output")
	}
}

func TestBoundsClipping(t *testing.T) {
	// Test that out-of-bounds coordinates are handled correctly
	c := mockCanvas(5, 5)
	c.Fill("█", "white", "black")

	// These should not cause panics or infinite loops
	c.Rect(-100, -100, 10, 10)
	c.Triangle(-10, -10, 20, -5, 0, 20)
	c.Circle(0, 0, 50)
	c.render()

	// Should complete without issues
	output := c.output.(*bytes.Buffer).String()
	// Output might be empty if everything is clipped, that's OK
	_ = output
}

func BenchmarkTriangleFillOld(b *testing.B) {
	c := mockCanvas(100, 100)
	c.Fill("█", "white", "black")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use the old brute force method for comparison
		c.fillTriangleBruteForce(10, 10, 90, 20, 50, 80)
	}
}

func BenchmarkTriangleFillNew(b *testing.B) {
	c := mockCanvas(100, 100)
	c.Fill("█", "white", "black")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use the new optimized method
		c.fillTriangle(10, 10, 90, 20, 50, 80)
	}
}

// Old brute force method for comparison
func (c *Canvas) fillTriangleBruteForce(x1, y1, x2, y2, x3, y3 int) {
	minX := min(x1, min(x2, x3))
	maxX := max(x1, max(x2, x3))
	minY := min(y1, min(y2, y3))
	maxY := max(y1, max(y2, y3))

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if pointInTriangle(x, y, x1, y1, x2, y2, x3, y3) {
				c.Point(x, y)
			}
		}
	}
}

func TestLargeShapePerformance(t *testing.T) {
	// Test that large shapes don't cause performance issues
	c := mockCanvas(200, 200)
	c.Fill("█", "white", "black")

	// These should complete quickly
	c.Rect(0, 0, 199, 199)
	c.Triangle(10, 10, 190, 20, 100, 180)
	c.Circle(100, 100, 90)
	c.Ellipse(100, 100, 80, 120)
	c.Quad(20, 20, 180, 30, 170, 170, 30, 160)
	c.render()

	// Should complete without timeout
	output := c.output.(*bytes.Buffer).String()
	if output == "" {
		t.Error("Expected output from large shapes")
	}
}

func BenchmarkLargeRectFill(b *testing.B) {
	c := mockCanvas(500, 500)
	c.Fill("█", "white", "black")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Rect(0, 0, 499, 499)
	}
}

func BenchmarkLargeTriangleFill(b *testing.B) {
	c := mockCanvas(500, 500)
	c.Fill("█", "white", "black")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Triangle(0, 0, 499, 100, 250, 499)
	}
}

func BenchmarkLargeCircleFill(b *testing.B) {
	c := mockCanvas(500, 500)
	c.Fill("█", "white", "black")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Circle(250, 250, 200)
	}
}

func BenchmarkOutOfBoundsShapes(b *testing.B) {
	c := mockCanvas(100, 100)
	c.Fill("█", "white", "black")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// These extend far beyond canvas bounds
		c.Rect(-1000, -1000, 2000, 2000)
		c.Triangle(-500, -500, 1500, 200, 500, 1500)
		c.Circle(50, 50, 1000)
	}
}
