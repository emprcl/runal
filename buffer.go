package runal

type buffer [][]string

func newBuffer(width, height int) buffer {
	buff := make([][]string, height)
	for i := range buff {
		buff[i] = make([]string, width)
	}
	return buff
}
