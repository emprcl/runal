package canvas

import (
	"math"
	"testing"
)

func countDrawn(c *Canvas) int {
	n := 0
	for y := range c.buffer {
		for x := range c.buffer[y] {
			if c.buffer[y][x].char != 0 {
				n++
			}
		}
	}
	return n
}

func TestMat4IdentityMulVec(t *testing.T) {
	m := mat4Identity()
	x, y, z, w := m.mulVec4(1, 2, 3, 1)
	if x != 1 || y != 2 || z != 3 || w != 1 {
		t.Fatalf("identity transform changed vector: got (%v, %v, %v, %v)", x, y, z, w)
	}
}

func TestMat4TranslateMulPoint(t *testing.T) {
	p := mat4Translate(10, -5, 2).mulPoint(vec3{1, 1, 1})
	if p != (vec3{11, -4, 3}) {
		t.Fatalf("unexpected translation: got %+v", p)
	}
}

func TestProjectOriginIsCentered(t *testing.T) {
	c := mockCanvas(100, 100)
	sx, sy, _, ok := c.projectClip(vec3{0, 0, 0})
	if !ok {
		t.Fatal("origin should be in front of the default camera")
	}
	if sx != 0 || sy != 0 {
		t.Fatalf("origin should project to (0, 0), got (%d, %d)", sx, sy)
	}
}

func TestProjectBehindCameraIsCulled(t *testing.T) {
	c := mockCanvas(100, 100)
	if _, _, _, ok := c.projectClip(vec3{0, 0, 100}); ok {
		t.Fatal("point behind the camera should be culled")
	}
}

func TestProjectDepthOrdering(t *testing.T) {
	c := mockCanvas(100, 100)
	_, _, zNear, _ := c.projectClip(vec3{0, 0, 10})
	_, _, zFar, _ := c.projectClip(vec3{0, 0, -10})
	if zNear >= zFar {
		t.Fatalf("closer point should have smaller depth: near=%v far=%v", zNear, zFar)
	}
}

// TestWriteDepth directly exercises the z-buffer: nearer writes win, farther
// writes are rejected.
func TestWriteDepth(t *testing.T) {
	c := mockCanvas(10, 10)
	far := cell{char: 'F'}
	near := cell{char: 'N'}
	behind := cell{char: 'B'}

	c.writeDepth(far, 5, 5, 0.5)
	if c.buffer[5][5].char != 'F' {
		t.Fatal("first write should land")
	}
	c.writeDepth(near, 5, 5, 0.2)
	if c.buffer[5][5].char != 'N' {
		t.Fatal("nearer write should overwrite")
	}
	c.writeDepth(behind, 5, 5, 0.9)
	if c.buffer[5][5].char != 'N' {
		t.Fatal("farther write must not overwrite a nearer cell")
	}
}

func TestClearResetsDepth(t *testing.T) {
	c := mockCanvas(10, 10)
	c.writeDepth(cell{char: 'X'}, 3, 3, 0.1)
	c.Clear()
	if !math.IsInf(c.depth[3][3], 1) {
		t.Fatal("Clear should reset the depth buffer to +Inf")
	}
	// After clear, any depth should be writable again.
	c.writeDepth(cell{char: 'Y'}, 3, 3, 0.9)
	if c.buffer[3][3].char != 'Y' {
		t.Fatal("after Clear a far cell should be writable")
	}
}

func TestBoxDrawsSolid(t *testing.T) {
	c := mockCanvas(100, 100)
	c.Translate(50, 50)
	c.Box(30, 30, 30)
	// A solid box fills an area, not just an outline.
	if countDrawn(c) < 200 {
		t.Fatalf("solid box should fill many cells, got %d", countDrawn(c))
	}
}

func TestSphereDrawsSolid(t *testing.T) {
	c := mockCanvas(100, 100)
	c.Translate(50, 50)
	c.Sphere(25, 12)
	if countDrawn(c) < 200 {
		t.Fatalf("solid sphere should fill many cells, got %d", countDrawn(c))
	}
}

// TestBackfaceCulling: a triangle wound so its normal faces away from the
// camera is dropped when culling is on, and drawn when culling is off.
func TestBackfaceCulling(t *testing.T) {
	front := func(cull bool) int {
		c := mockCanvas(60, 60)
		c.Cull(cull)
		c.Translate(30, 30)
		// Facing the camera (+z normal via CCW winding seen from +z).
		c.Triangle3D(-10, -10, 0, 10, -10, 0, 0, 10, 0)
		return countDrawn(c)
	}
	if front(true) == 0 {
		t.Fatal("camera-facing triangle should render with culling on")
	}

	back := func() int {
		c := mockCanvas(60, 60)
		c.Cull(true)
		c.Translate(30, 30)
		// Reversed winding -> normal points away -> culled.
		c.Triangle3D(0, 10, 0, 10, -10, 0, -10, -10, 0)
		return countDrawn(c)
	}
	if back() != 0 {
		t.Fatal("back-facing triangle should be culled")
	}

	withoutCull := func() int {
		c := mockCanvas(60, 60)
		c.Cull(false)
		c.Translate(30, 30)
		c.Triangle3D(0, 10, 0, 10, -10, 0, -10, -10, 0)
		return countDrawn(c)
	}
	if withoutCull() == 0 {
		t.Fatal("back-facing triangle should render with culling off")
	}
}

// projectedSquareRatio returns the width/height (in cells) of a front-facing
// world square projected onto the canvas.
func projectedSquareRatio(c *Canvas) float64 {
	rx, _, _, _ := c.projectClip(vec3{10, 0, 0})
	lx, _, _, _ := c.projectClip(vec3{-10, 0, 0})
	_, ty, _, _ := c.projectClip(vec3{0, 10, 0})
	_, by, _, _ := c.projectClip(vec3{0, -10, 0})
	w := float64(rx - lx)
	h := float64(by - ty)
	return w / h
}

// TestAspectAccountsForCellShape verifies the projection compensates for the
// non-square terminal cell: a world square should span twice as many cells
// horizontally as vertically without cell mode, and equal cells with it.
func TestAspectAccountsForCellShape(t *testing.T) {
	c := mockCanvas(100, 100) // square canvas, so only cell shape matters
	if r := projectedSquareRatio(c); math.Abs(r-2.0) > 0.15 {
		t.Fatalf("without cell mode, square should span ~2x cells wide: got ratio %v", r)
	}

	c2 := mockCanvas(100, 100)
	c2.cellMode = cellModeDouble // square cells
	if r := projectedSquareRatio(c2); math.Abs(r-1.0) > 0.15 {
		t.Fatalf("with cell mode, square should span equal cells: got ratio %v", r)
	}

	c3 := mockCanvas(100, 100)
	c3.CharAspect(1.0) // pretend cells are square
	if r := projectedSquareRatio(c3); math.Abs(r-1.0) > 0.15 {
		t.Fatalf("with square char aspect, square should span equal cells: got ratio %v", r)
	}
}

func TestPushPop3DRestoresTransform(t *testing.T) {
	c := mockCanvas(100, 100)
	c.Translate3D(5, 0, 0)
	before := c.model3D
	c.Push3D()
	c.Translate3D(100, 100, 100)
	c.Pop3D()
	if c.model3D != before {
		t.Fatal("Pop3D should restore the model transform saved by Push3D")
	}
}
