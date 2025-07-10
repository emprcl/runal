package runal

import (
	"math"
	"sort"
)

// Text renders a string at the given canvas coordinates.
func (c *Canvas) Text(text string, x, y int) {
	destX := c.originX + x
	destY := c.originY + y

	if c.outOfBounds(destX, destY) {
		return
	}
	for i, r := range text {
		if x+i < len(c.buffer[destY])-1 {
			c.buffer[destY][destX+i] = c.formatCell(r)
		} else if x+i == len(c.buffer[destY])-1 {
			c.buffer[destY][destX+i] = c.style(string(r))
		}
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
		// Clip to canvas bounds for filling
		fillStartY := max(0, y-c.originY)
		fillEndY := min(c.Height-1, y+h-c.originY)
		fillStartX := max(0, x-c.originX)
		fillEndX := min(c.Width-1, x+w-c.originX)

		for dy := fillStartY; dy <= fillEndY; dy++ {
			actualY := dy + c.originY
			if actualY >= y && actualY <= y+h {
				c.Line(max(fillStartX+c.originX, x), actualY, min(fillEndX+c.originX, x+w), actualY)
			}
		}
		c.toggleFill()
	}
	c.Line(x, y, x+w, y)
	c.Line(x+w, y, x+w, y+h)
	c.Line(x+w, y+h, x, y+h)
	c.Line(x, y+h, x, y)
}

// Quad draws a quadrilateral defined by four points.
func (c *Canvas) Quad(x1, y1, x2, y2, x3, y3, x4, y4 int) {
	if c.fill {
		vertices := [][2]int{{x1, y1}, {x2, y2}, {x3, y3}, {x4, y4}}
		scanlineIntersections := map[int][]int{}

		// Find bounding box and clip to canvas
		minY := min(y1, min(y2, min(y3, y4)))
		maxY := max(y1, max(y2, max(y3, y4)))
		minY = max(minY, -c.originY)
		maxY = min(maxY, c.Height-1-c.originY)

		for i := 0; i < 4; i++ {
			xStart, yStart := vertices[i][0], vertices[i][1]
			xEnd, yEnd := vertices[(i+1)%4][0], vertices[(i+1)%4][1]

			if yStart == yEnd {
				y := yStart
				if y < minY || y > maxY {
					continue
				}
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
				for y := max(yStart, minY); y <= min(yEnd, maxY); y++ {
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
					// Clip horizontal lines to canvas bounds
					x1Clipped := max(xs[i], -c.originX)
					x2Clipped := min(xs[i+1], c.Width-1-c.originX)
					if x1Clipped <= x2Clipped {
						c.Line(x1Clipped, y, x2Clipped, y)
					}
				}
			}
		}
		c.toggleFill()
	}

	// Draw outline lines
	c.Line(x1, y1, x2, y2)
	c.Line(x2, y2, x3, y3)
	c.Line(x3, y3, x4, y4)
	c.Line(x4, y4, x1, y1)
}

// Ellipse draws an ellipse centered at (x, y) with radiuses rx and ry.
func (c *Canvas) Ellipse(xCenter, yCenter, rx, ry int) {
	steps := 360
	points := make([][2]int, 0, steps)

	for i := 0; i < steps; i++ {
		theta := 2 * math.Pi * float64(i) / float64(steps)
		x := int(math.Round(float64(rx) * math.Cos(theta)))
		y := int(math.Round(float64(ry) * math.Sin(theta)))
		points = append(points, [2]int{xCenter + x, yCenter + y})
	}

	if c.fill {
		c.toggleFill()

		// Find bounding box and clip to canvas
		minY := max(yCenter-ry, -c.originY)
		maxY := min(yCenter+ry, c.Height-1-c.originY)

		rows := map[int][]int{}
		for _, p := range points {
			y := p[1]
			x := p[0]
			// Only add points within Y bounds
			if y >= minY && y <= maxY {
				rows[y] = append(rows[y], x)
			}
		}

		for y, xs := range rows {
			if len(xs) < 2 {
				continue
			}
			sort.Ints(xs)
			// Clip horizontal lines to canvas bounds
			x1Clipped := max(xs[0], -c.originX)
			x2Clipped := min(xs[len(xs)-1], c.Width-1-c.originX)
			if x1Clipped <= x2Clipped {
				c.Line(x1Clipped, y, x2Clipped, y)
			}
		}

		c.toggleFill()
	}

	for _, p := range points {
		c.Point(p[0], p[1])
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
			// Clip horizontal lines to canvas bounds
			if y <= r-1 {
				// Top and bottom horizontal lines
				if yCenter+y >= -c.originY && yCenter+y <= c.Height-1-c.originY {
					x1 := max(xCenter-x, -c.originX)
					x2 := min(xCenter+x, c.Width-1-c.originX)
					if x1 <= x2 {
						c.Line(x1, yCenter+y, x2, yCenter+y)
					}
				}
				if yCenter-y >= -c.originY && yCenter-y <= c.Height-1-c.originY {
					x1 := max(xCenter-x, -c.originX)
					x2 := min(xCenter+x, c.Width-1-c.originX)
					if x1 <= x2 {
						c.Line(x1, yCenter-y, x2, yCenter-y)
					}
				}
			}
			// Left and right horizontal lines
			if yCenter+x >= -c.originY && yCenter+x <= c.Height-1-c.originY {
				x1 := max(xCenter-y, -c.originX)
				x2 := min(xCenter+y, c.Width-1-c.originX)
				if x1 <= x2 {
					c.Line(x1, yCenter+x, x2, yCenter+x)
				}
			}
			if yCenter-x >= -c.originY && yCenter-x <= c.Height-1-c.originY {
				x1 := max(xCenter-y, -c.originX)
				x2 := min(xCenter+y, c.Width-1-c.originX)
				if x1 <= x2 {
					c.Line(x1, yCenter-x, x2, yCenter-x)
				}
			}
			c.toggleFill()
		}

		c.char(strIndex(c.strokeText, char), xCenter+x, yCenter+y)
		c.char(strIndex(c.strokeText, char+1), xCenter-x, yCenter+y)
		c.char(strIndex(c.strokeText, char+2), xCenter+x, yCenter-y)
		c.char(strIndex(c.strokeText, char+3), xCenter-x, yCenter-y)
		c.char(strIndex(c.strokeText, char+4), xCenter+y, yCenter+x)
		c.char(strIndex(c.strokeText, char+5), xCenter-y, yCenter+x)
		c.char(strIndex(c.strokeText, char+6), xCenter+y, yCenter-x)
		c.char(strIndex(c.strokeText, char+7), xCenter-y, yCenter-x)

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
	c.Line(x1, y1, x2, y2)
	c.Line(x2, y2, x3, y3)
	c.Line(x3, y3, x1, y1)
}

func (c *Canvas) fillTriangle(x1, y1, x2, y2, x3, y3 int) {
	// Use scanline algorithm instead of brute force
	minY := min(y1, min(y2, y3))
	maxY := max(y1, max(y2, y3))

	// Clip to canvas bounds
	minY = max(minY, -c.originY)
	maxY = min(maxY, c.Height-1-c.originY)

	for y := minY; y <= maxY; y++ {
		intersections := []int{}

		// Check intersections with each edge
		edges := [][4]int{{x1, y1, x2, y2}, {x2, y2, x3, y3}, {x3, y3, x1, y1}}

		for _, edge := range edges {
			x1e, y1e, x2e, y2e := edge[0], edge[1], edge[2], edge[3]

			if y1e == y2e {
				// Horizontal edge
				if y1e == y && x1e != x2e {
					intersections = append(intersections, min(x1e, x2e), max(x1e, x2e))
				}
			} else {
				// Non-horizontal edge
				if (y1e <= y && y < y2e) || (y2e <= y && y < y1e) {
					// Calculate intersection
					t := float64(y-y1e) / float64(y2e-y1e)
					x := int(math.Round(float64(x1e) + t*float64(x2e-x1e)))
					intersections = append(intersections, x)
				}
			}
		}

		if len(intersections) >= 2 {
			sort.Ints(intersections)
			// Fill between pairs of intersections
			for i := 0; i < len(intersections)-1; i += 2 {
				if i+1 < len(intersections) {
					x1Fill := max(intersections[i], -c.originX)
					x2Fill := min(intersections[i+1], c.Width-1-c.originX)
					if x1Fill <= x2Fill {
						c.Line(x1Fill, y, x2Fill, y)
					}
				}
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
