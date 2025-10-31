function setup(c) {
	c.noLoop()
}

function draw(c) {
	for (i = 0; i < 256; i++) {
		c.push()
		c.stroke(" ", "0", i)
		c.translate((i % 16) * 10, Math.floor(i / 16) * 3)
		c.line(2, 1, 6, 1)
		c.line(2, 2, 6, 2)
		c.stroke(" ", "15", "0")
		c.text(i, 8, 1)
		c.pop()
	}
}
