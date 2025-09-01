package runal

import (
	"fmt"
	"testing"
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
