function setup(c) { }

function draw(c) {
  c.clear();

	c.stroke("/", "#fffb00", "#0004ff")
	c.fill(".", "#0004ff", "#fffb00")
	c.rect(5, 11, 11, 5)
	c.text("Stroke: YES", 5, 18)
	c.text("Fill:   YES", 5, 19)

	c.stroke("/", "#fffb00", "#0004ff")
	c.noFill()
	c.rect(20, 11, 11, 5)
	c.text("Stroke: YES", 20, 18)
	c.text("Fill:   NO", 20, 19)

	c.fill(".", "#0004ff", "#fffb00")
	c.noStroke()
	c.rect(35, 11, 11, 5)
	c.text("Stroke: NO", 35, 18)
	c.text("Fill:   YES", 35, 19)

	c.noStroke()
	c.noFill()
	c.rect(50, 11, 11, 5)
	c.text("Stroke: NO", 50, 18)
	c.text("Fill:   NO", 50, 19)
}
