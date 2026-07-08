package canvas

import (
	"fmt"
	"math"

	"github.com/charmbracelet/x/ansi"
)

// vertex2D is a rasterizer vertex: integer screen position plus interpolated
// depth.
type vertex2D struct {
	x, y int
	z    float64
}

// Point3D draws a single depth-tested point at the given 3D position using the
// stroke color.
func (c *Canvas) Point3D(x, y, z float64) {
	sx, sy, depth, ok := c.projectVertex(vec3{x, y, z})
	if !ok {
		return
	}
	c.writeDepth(cell{char: c.nextStrokeRune(), foreground: c.strokeFg, background: c.strokeBg}, sx, sy, depth)
}

// Line3D draws a depth-tested straight line between two 3D points using the
// stroke color.
func (c *Canvas) Line3D(x1, y1, z1, x2, y2, z2 float64) {
	ax, ay, az, ok1 := c.projectVertex(vec3{x1, y1, z1})
	bx, by, bz, ok2 := c.projectVertex(vec3{x2, y2, z2})
	if !ok1 || !ok2 {
		return
	}
	c.rasterLine(vertex2D{ax, ay, az}, vertex2D{bx, by, bz})
}

// Triangle3D draws a solid, shaded, depth-tested triangle from three
// model-space vertices.
func (c *Canvas) Triangle3D(x1, y1, z1, x2, y2, z2, x3, y3, z3 float64) {
	w0 := c.worldPoint(vec3{x1, y1, z1})
	w1 := c.worldPoint(vec3{x2, y2, z2})
	w2 := c.worldPoint(vec3{x3, y3, z3})

	normal := w1.sub(w0).cross(w2.sub(w0)).normalize()

	// Back-face culling: skip faces whose normal points away from the eye.
	if c.cull && normal.dot(w0.sub(c.cam.eye)) > 0 {
		return
	}

	intensity := c.faceIntensity(normal)
	shaded := c.shadedCell(intensity)

	s0x, s0y, z0, ok0 := c.projectClipAnd2D(w0)
	s1x, s1y, z1v, ok1 := c.projectClipAnd2D(w1)
	s2x, s2y, z2v, ok2 := c.projectClipAnd2D(w2)
	if !ok0 || !ok1 || !ok2 {
		return
	}

	c.rasterTriangle(
		vertex2D{s0x, s0y, z0},
		vertex2D{s1x, s1y, z1v},
		vertex2D{s2x, s2y, z2v},
		shaded,
	)
}

// Quad3D draws a solid quadrilateral as two triangles.
func (c *Canvas) Quad3D(x1, y1, z1, x2, y2, z2, x3, y3, z3, x4, y4, z4 float64) {
	c.Triangle3D(x1, y1, z1, x2, y2, z2, x3, y3, z3)
	c.Triangle3D(x1, y1, z1, x3, y3, z3, x4, y4, z4)
}

// Box draws a solid box centered on the model origin with the given width,
// height and depth. Faces are wound counter-clockwise when seen from outside
// so back-face culling works.
func (c *Canvas) Box(w, h, d float64) {
	hw, hh, hd := w/2, h/2, d/2
	// 8 corners
	p := [8]vec3{
		{-hw, -hh, -hd}, // 0
		{hw, -hh, -hd},  // 1
		{hw, hh, -hd},   // 2
		{-hw, hh, -hd},  // 3
		{-hw, -hh, hd},  // 4
		{hw, -hh, hd},   // 5
		{hw, hh, hd},    // 6
		{-hw, hh, hd},   // 7
	}
	quad := func(a, b, cc, d int) {
		c.Quad3D(
			p[a].x, p[a].y, p[a].z,
			p[b].x, p[b].y, p[b].z,
			p[cc].x, p[cc].y, p[cc].z,
			p[d].x, p[d].y, p[d].z,
		)
	}
	quad(4, 5, 6, 7) // front  (+z)
	quad(1, 0, 3, 2) // back   (-z)
	quad(0, 4, 7, 3) // left   (-x)
	quad(5, 1, 2, 6) // right  (+x)
	quad(3, 7, 6, 2) // top    (+y)
	quad(0, 1, 5, 4) // bottom (-y)
}

// Sphere draws a solid sphere of the given radius. detail controls the number
// of latitude and longitude segments (minimum 3).
func (c *Canvas) Sphere(r float64, detail int) {
	if detail < 3 {
		detail = 3
	}
	lat := detail
	lon := detail

	at := func(i, j int) vec3 {
		theta := math.Pi*float64(i)/float64(lat) - math.Pi/2
		phi := 2 * math.Pi * float64(j) / float64(lon)
		cosTheta := math.Cos(theta)
		return vec3{
			r * cosTheta * math.Cos(phi),
			r * math.Sin(theta),
			r * cosTheta * math.Sin(phi),
		}
	}

	for i := 0; i < lat; i++ {
		for j := 0; j < lon; j++ {
			a := at(i, j)
			b := at(i+1, j)
			cc := at(i+1, j+1)
			d := at(i, j+1)
			c.Quad3D(
				a.x, a.y, a.z,
				b.x, b.y, b.z,
				cc.x, cc.y, cc.z,
				d.x, d.y, d.z,
			)
		}
	}
}

// projectClipAnd2D projects a world point and applies the 2D transform.
func (c *Canvas) projectClipAnd2D(w vec3) (int, int, float64, bool) {
	sx, sy, z, ok := c.projectClip(w)
	if !ok {
		return 0, 0, 0, false
	}
	tx, ty := c.applyTransform2D(sx, sy)
	return tx, ty, z, true
}

// faceIntensity returns the light intensity for a face with the given normal.
func (c *Canvas) faceIntensity(normal vec3) float64 {
	diffuse := normal.dot(c.lightDir)
	if !c.cull {
		// With culling off, both sides are visible; light either face.
		diffuse = math.Abs(diffuse)
	} else if diffuse < 0 {
		diffuse = 0
	}
	return clamp(c.ambient+(1-c.ambient)*diffuse, 0, 1)
}

// shadedCell builds a fill cell for the given light intensity, picking a glyph
// from the shade ramp and scaling the fill color.
func (c *Canvas) shadedCell(intensity float64) cell {
	runes := c.shadeRunes
	idx := int(intensity*float64(len(runes)-1) + 0.5)
	idx = clamp(idx, 0, len(runes)-1)
	return cell{
		char:       runes[idx],
		foreground: shadeColor(c.fillFg, intensity),
		background: c.fillBg,
	}
}

func shadeColor(base ansi.Color, intensity float64) ansi.Color {
	r, g, b, _ := base.RGBA()
	s := clamp(intensity, 0, 1)
	rr := int(float64(r>>8) * s)
	gg := int(float64(g>>8) * s)
	bb := int(float64(b>>8) * s)
	return color(fmt.Sprintf("#%02x%02x%02x", clamp(rr, 0, 255), clamp(gg, 0, 255), clamp(bb, 0, 255)))
}

func edgeFn(a, b vertex2D, x, y int) int {
	return (b.x-a.x)*(y-a.y) - (b.y-a.y)*(x-a.x)
}

// rasterTriangle fills a triangle in screen space, interpolating depth per
// cell and writing through the depth test.
func (c *Canvas) rasterTriangle(a, b, cc vertex2D, cll cell) {
	minX := max(min(a.x, min(b.x, cc.x)), 0)
	maxX := min(max(a.x, max(b.x, cc.x)), c.Width-1)
	minY := max(min(a.y, min(b.y, cc.y)), 0)
	maxY := min(max(a.y, max(b.y, cc.y)), c.Height-1)

	area := edgeFn(a, b, cc.x, cc.y)
	if area == 0 {
		return
	}
	invArea := 1 / float64(area)

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			w0 := edgeFn(b, cc, x, y)
			w1 := edgeFn(cc, a, x, y)
			w2 := edgeFn(a, b, x, y)
			inside := (w0 >= 0 && w1 >= 0 && w2 >= 0) || (w0 <= 0 && w1 <= 0 && w2 <= 0)
			if !inside {
				continue
			}
			l0 := float64(w0) * invArea
			l1 := float64(w1) * invArea
			l2 := float64(w2) * invArea
			z := l0*a.z + l1*b.z + l2*cc.z
			c.writeDepth(cll, x, y, z)
		}
	}
}

// rasterLine draws a depth-tested line using Bresenham, interpolating depth
// along the way, with the stroke color.
func (c *Canvas) rasterLine(a, b vertex2D) {
	dx := absInt(b.x - a.x)
	dy := absInt(b.y - a.y)
	steps := max(dx, dy)

	sx := 1
	sy := 1
	if a.x > b.x {
		sx = -1
	}
	if a.y > b.y {
		sy = -1
	}
	d := dx - dy

	x, y := a.x, a.y
	runes := c.strokeRunes
	ri := 0
	for i := 0; ; i++ {
		t := 0.0
		if steps > 0 {
			t = float64(i) / float64(steps)
		}
		z := a.z + (b.z-a.z)*t
		c.writeDepth(cell{char: runes[ri], foreground: c.strokeFg, background: c.strokeBg}, x, y, z)
		if x == b.x && y == b.y {
			break
		}
		ri = (ri + 1) % len(runes)
		e2 := 2 * d
		if e2 > -dy {
			d -= dy
			x += sx
		}
		if e2 < dx {
			d += dx
			y += sy
		}
	}
}
