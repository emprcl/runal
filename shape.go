package runal

import (
	"math"
	"sort"
)

// Text renders a string at the given canvas coordinates.
func (c *Canvas) Text(text string, x, y int) {
	for i, r := range text {
		c.char(r, x+i, y)
	}
}

// Point draws a single point at the given position.
func (c *Canvas) Point(x, y int) {
	c.char(c.nextStrokeRune(), x, y)
}

// Line draws a straight line between two points.
func (c *Canvas) Line(x1, y1, x2, y2 int) {
	// Bresenham algorithm
	// https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
	dx := absInt(x2 - x1)
	dy := absInt(y2 - y1)
	sx := 1
	sy := 1

	if x1 > x2 {
		sx = -1
	}
	if y1 > y2 {
		sy = -1
	}

	d := dx - dy

	char := 0
	for {
		c.char(rune(c.strokeText[char]), x1, y1)
		if x1 == x2 && y1 == y2 {
			break
		}
		char = (char + 1) % len(c.strokeText)
		e2 := 2 * d
		if e2 > -dy {
			d -= dy
			x1 += sx
		}
		if e2 < dx {
			d += dx
			y1 += sy
		}
	}
}

// Square draws a square with the given top-left corner and side length.
func (c *Canvas) Square(x, y, size int) {
	c.Rect(x, y, size, size)
}

// Rect draws a rectangle starting at (x, y) with width w and height h.
func (c *Canvas) Rect(x, y, w, h int) {
	if c.fill {
		c.toggleFill()
		for dy := range h {
			c.Line(x, y+dy, x+w, y+dy)
		}
		c.toggleFill()
	}
	if c.stroke {
		c.Line(x, y, x+w, y)
		c.Line(x+w, y, x+w, y+h)
		c.Line(x+w, y+h, x, y+h)
		c.Line(x, y+h, x, y)
	}
}

// Quad draws a quadrilateral defined by four points.
func (c *Canvas) Quad(x1, y1, x2, y2, x3, y3, x4, y4 int) {
	if c.fill {
		vertices := [][2]int{{x1, y1}, {x2, y2}, {x3, y3}, {x4, y4}}
		scanlineIntersections := map[int][]int{}
		for i := 0; i < 4; i++ {
			xStart, yStart := vertices[i][0], vertices[i][1]
			xEnd, yEnd := vertices[(i+1)%4][0], vertices[(i+1)%4][1]

			if yStart == yEnd {
				y := yStart
				xmin := xStart
				xmax := xEnd
				if xmin > xmax {
					xmin, xmax = xmax, xmin
				}
				scanlineIntersections[y] = append(scanlineIntersections[y], xmin, xmax)
			} else {
				if yStart > yEnd {
					yStart, yEnd = yEnd, yStart
					xStart, xEnd = xEnd, xStart
				}
				for y := yStart; y <= yEnd; y++ {
					t := float64(y-yStart) / float64(yEnd-yStart)
					x := int(math.Round(float64(xStart) + t*float64(xEnd-xStart)))
					scanlineIntersections[y] = append(scanlineIntersections[y], x)
				}
			}
		}

		c.toggleFill()
		for y, xs := range scanlineIntersections {
			if len(xs) < 2 {
				continue
			}
			sort.Ints(xs)
			for i := 0; i < len(xs); i += 2 {
				if i+1 < len(xs) {
					c.Line(xs[i], y, xs[i+1], y)
				}
			}
		}
		c.toggleFill()
	}

	// Draw outline lines
	if c.stroke {
		c.Line(x1, y1, x2, y2)
		c.Line(x2, y2, x3, y3)
		c.Line(x3, y3, x4, y4)
		c.Line(x4, y4, x1, y1)
	}
}

// Ellipse draws an ellipse centered at (x, y) with radiuses rx and ry.
func (c *Canvas) Ellipse(xCenter, yCenter, rx, ry int) {
	steps := 360
	points := make([][2]int, 0, steps)

	for i := range steps {
		theta := 2 * math.Pi * float64(i) / float64(steps)
		x := int(math.Round(float64(rx) * math.Cos(theta)))
		y := int(math.Round(float64(ry) * math.Sin(theta)))
		points = append(points, [2]int{xCenter + x, yCenter + y})
	}

	if c.fill {
		c.toggleFill()

		rows := map[int][]int{}
		for _, p := range points {
			y := p[1]
			x := p[0]
			rows[y] = append(rows[y], x)
		}

		for y, xs := range rows {
			if len(xs) < 2 {
				continue
			}
			sort.Ints(xs)
			c.Line(xs[0], y, xs[len(xs)-1], y)
		}

		c.toggleFill()
	}

	if c.stroke {
		for _, p := range points {
			c.Point(p[0], p[1])
		}
	}
}

// Circle draws a circle centered at (x, y) with the given radius.
func (c *Canvas) Circle(xCenter, yCenter, r int) {
	x := 0
	y := r
	d := 1 - r
	char := 0

	for x <= y {
		char = char + 8
		if c.fill {
			c.toggleFill()
			if y <= r-1 {
				c.Line(xCenter-x, yCenter+y, xCenter+x, yCenter+y)
				c.Line(xCenter-x, yCenter-y, xCenter+x, yCenter-y)
			}
			c.Line(xCenter-y, yCenter+x, xCenter+y, yCenter+x)
			c.Line(xCenter-y, yCenter-x, xCenter+y, yCenter-x)
			c.toggleFill()
		}

		if c.stroke {
			c.char(strIndex(c.strokeText, char), xCenter+x, yCenter+y)
			c.char(strIndex(c.strokeText, char+1), xCenter-x, yCenter+y)
			c.char(strIndex(c.strokeText, char+2), xCenter+x, yCenter-y)
			c.char(strIndex(c.strokeText, char+3), xCenter-x, yCenter-y)
			c.char(strIndex(c.strokeText, char+4), xCenter+y, yCenter+x)
			c.char(strIndex(c.strokeText, char+5), xCenter-y, yCenter+x)
			c.char(strIndex(c.strokeText, char+6), xCenter+y, yCenter-x)
			c.char(strIndex(c.strokeText, char+7), xCenter-y, yCenter-x)
		}

		x++
		if d < 0 {
			d += 2*x + 1
		} else {
			y--
			d += 2*(x-y) + 1
		}
	}
}

// Triangle draws a triangle using three vertex points.
func (c *Canvas) Triangle(x1, y1, x2, y2, x3, y3 int) {
	if c.fill {
		c.toggleFill()
		c.fillTriangle(x1, y1, x2, y2, x3, y3)
		c.toggleFill()
	}

	if c.stroke {
		c.Line(x1, y1, x2, y2)
		c.Line(x2, y2, x3, y3)
		c.Line(x3, y3, x1, y1)
	}
}

func (c *Canvas) fillTriangle(x1, y1, x2, y2, x3, y3 int) {
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

func pointInTriangle(px, py, x1, y1, x2, y2, x3, y3 int) bool {
	ax, ay := float64(x1), float64(y1)
	bx, by := float64(x2), float64(y2)
	cx, cy := float64(x3), float64(y3)
	pxf, pyf := float64(px), float64(py)

	v0x, v0y := cx-ax, cy-ay
	v1x, v1y := bx-ax, by-ay
	v2x, v2y := pxf-ax, pyf-ay

	d00 := v0x*v0x + v0y*v0y
	d01 := v0x*v1x + v0y*v1y
	d02 := v0x*v2x + v0y*v2y
	d11 := v1x*v1x + v1y*v1y
	d12 := v1x*v2x + v1y*v2y

	denom := d00*d11 - d01*d01
	if denom == 0 {
		return false
	}
	invDenom := 1 / denom
	u := (d11*d02 - d01*d12) * invDenom
	v := (d00*d12 - d01*d02) * invDenom

	return u >= 0 && v >= 0 && u+v <= 1
}
