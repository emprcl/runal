package runal

func (c *Canvas) Text(text string, x, y int) {
	destX := c.originX + x
	destY := c.originY + y

	if c.OutOfBounds(destX, destY) {
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

func (c *Canvas) Point(x, y int) {
	c.char(c.nextStrokeRune(), x, y)
}

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

func (c *Canvas) Square(x, y, size int) {
	c.Rect(x, y, size, size)
}

func (c *Canvas) Rect(x, y, w, h int) {
	if c.fill {
		c.toggleFill()
		for dy := range h {
			c.Line(x, y+dy, x+w, y+dy)
		}
		c.toggleFill()
	}
	c.Line(x, y, x+w, y)
	c.Line(x+w, y, x+w, y+h)
	c.Line(x+w, y+h, x, y+h)
	c.Line(x, y+h, x, y)
}

func (c *Canvas) Ellipse(xCenter, yCenter, rx, ry int) {
	x := 0
	y := ry

	rx2 := rx * rx
	ry2 := ry * ry
	twory2 := 2 * ry2
	tworx2 := 2 * rx2

	px := 0
	py := tworx2 * y

	char := 0

	// Region 1
	d := ry2 - (rx2 * ry) + (rx2 / 4)
	for px < py {
		char += 8
		c.plotEllipsePoints(xCenter, yCenter, x, y, char)

		if c.fill && y < ry {
			c.toggleFill()
			c.Line(xCenter-x, yCenter+y, xCenter+x, yCenter+y)
			c.Line(xCenter-x, yCenter-y, xCenter+x, yCenter-y)
			c.toggleFill()
		}

		x++
		px += twory2
		if d < 0 {
			d += ry2 + px
		} else {
			y--
			py -= tworx2
			d += ry2 + px - py
		}
	}

	// Region 2
	d = 4*ry2*(x+1)*(x+1) + rx2*(2*y-1)*(2*y-1) - 4*rx2*ry2
	for y >= 0 {
		char += 8
		c.plotEllipsePoints(xCenter, yCenter, x, y, char)

		if c.fill && y < ry {
			c.toggleFill()
			c.Line(xCenter-x, yCenter+y, xCenter+x, yCenter+y)
			c.Line(xCenter-x, yCenter-y, xCenter+x, yCenter-y)
			c.toggleFill()
		}

		y--
		py -= tworx2
		if d > 0 {
			d += rx2 - py
		} else {
			x++
			px += twory2
			d += rx2 - py + px
		}
	}
}

func (c *Canvas) plotEllipsePoints(cx, cy, x, y, char int) {
	c.char(strIndex(c.strokeText, char+0), cx+x, cy+y)
	c.char(strIndex(c.strokeText, char+1), cx-x, cy+y)
	c.char(strIndex(c.strokeText, char+2), cx+x, cy-y)
	c.char(strIndex(c.strokeText, char+3), cx-x, cy-y)
	c.char(strIndex(c.strokeText, char+4), cx+y, cy+x)
	c.char(strIndex(c.strokeText, char+5), cx-y, cy+x)
	c.char(strIndex(c.strokeText, char+6), cx+y, cy-x)
	c.char(strIndex(c.strokeText, char+7), cx-y, cy-x)
}

func (c *Canvas) Circle(xCenter, yCenter, r int) {
	c.Ellipse(xCenter, yCenter, r, r)
}

// func (c *Canvas) Circle(xCenter, yCenter, r int) {
// 	x := 0
// 	y := r
// 	d := 1 - r
// 	char := 0

// 	for x <= y {
// 		char = char + 8
// 		if c.fill {
// 			c.toggleFill()
// 			if y <= r-1 {
// 				c.Line(xCenter-x, yCenter+y, xCenter+x, yCenter+y)
// 				c.Line(xCenter-x, yCenter-y, xCenter+x, yCenter-y)
// 			}
// 			c.Line(xCenter-y, yCenter+x, xCenter+y, yCenter+x)
// 			c.Line(xCenter-y, yCenter-x, xCenter+y, yCenter-x)
// 			c.toggleFill()
// 		}

// 		c.char(strIndex(c.strokeText, char), xCenter+x, yCenter+y)
// 		c.char(strIndex(c.strokeText, char+1), xCenter-x, yCenter+y)
// 		c.char(strIndex(c.strokeText, char+2), xCenter+x, yCenter-y)
// 		c.char(strIndex(c.strokeText, char+3), xCenter-x, yCenter-y)
// 		c.char(strIndex(c.strokeText, char+4), xCenter+y, yCenter+x)
// 		c.char(strIndex(c.strokeText, char+5), xCenter-y, yCenter+x)
// 		c.char(strIndex(c.strokeText, char+6), xCenter+y, yCenter-x)
// 		c.char(strIndex(c.strokeText, char+7), xCenter-y, yCenter-x)

// 		x++
// 		if d < 0 {
// 			d += 2*x + 1
// 		} else {
// 			y--
// 			d += 2*(x-y) + 1
// 		}
// 	}
// }

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
