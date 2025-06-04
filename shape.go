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

func (c *Canvas) char(char rune, x, y int) {
	destX := c.originX + x
	destY := c.originY + y

	if c.OutOfBounds(destX, destY) {
		return
	}
	c.buffer[destY][destX] = c.formatCell(char)
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

func (c *Canvas) Circle(xCenter, yCenter, r int) {
	x := 0
	y := r
	d := 1 - r
	char := 0

	for x < y {
		char = char + 8
		c.plotCircle(c.strokeText, char, xCenter, yCenter, x, y)

		x++
		if d < 0 {
			d += 2*x + 1
		} else {
			y--
			d += 2*(x-y) + 1
		}
	}
}

func (c *Canvas) plotCircle(borderText string, char, xCenter, yCenter, x, y int) {
	if c.fill {
		c.toggleFill()
		c.Line(xCenter-x, yCenter+y, xCenter+x, yCenter+y)
		c.Line(xCenter-x, yCenter-y, xCenter+x, yCenter-y)
		c.Line(xCenter-y, yCenter+x, xCenter+y, yCenter+x)
		c.Line(xCenter-y, yCenter-x, xCenter+y, yCenter-x)
		c.toggleFill()
	}

	c.char(strIndex(borderText, char), xCenter+x, yCenter+y)
	c.char(strIndex(borderText, char+1), xCenter-x, yCenter+y)
	c.char(strIndex(borderText, char+2), xCenter+x, yCenter-y)
	c.char(strIndex(borderText, char+3), xCenter-x, yCenter-y)
	c.char(strIndex(borderText, char+4), xCenter+y, yCenter+x)
	c.char(strIndex(borderText, char+5), xCenter-y, yCenter+x)
	c.char(strIndex(borderText, char+6), xCenter+y, yCenter-x)
	c.char(strIndex(borderText, char+7), xCenter-y, yCenter-x)
}
