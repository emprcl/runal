package runal

type Cell struct {
	Char       string
	Foreground string
	Background string
}

func (cll Cell) write(c *Canvas, x, y, _, _ int) {
	if c.outOfBounds(x, y) {
		return
	}
	c.write(cll.private(), x, y, 1)
}

func (cll Cell) private() cell {
	return cell{
		char:       []rune(cll.Char)[0],
		foreground: cll.Foreground,
		background: cll.Background,
	}
}

type cell struct {
	char       rune
	padChar    rune
	foreground string
	background string
}

func (c cell) public() Cell {
	return Cell{
		Char:       string(c.char),
		Foreground: c.foreground,
		Background: c.background,
	}
}
