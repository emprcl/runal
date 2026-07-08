package canvas

import "math"

// defaultCharAspect is the assumed height-to-width ratio of a terminal
// character cell. Most terminals render cells about twice as tall as wide.
const defaultCharAspect = 2.0

// camera holds the view and projection parameters used to project 3D
// points onto the 2D canvas.
type camera struct {
	eye, center, up vec3
	fov, near, far  float64
}

func defaultCamera() camera {
	return camera{
		eye:    vec3{0, 0, 30},
		center: vec3{0, 0, 0},
		up:     vec3{0, 1, 0},
		fov:    math.Pi / 3,
		near:   0.1,
		far:    1000,
	}
}

// Camera positions the eye, sets the point it looks at and the up vector.
func (c *Canvas) Camera(eyeX, eyeY, eyeZ, centerX, centerY, centerZ, upX, upY, upZ float64) {
	c.cam.eye = vec3{eyeX, eyeY, eyeZ}
	c.cam.center = vec3{centerX, centerY, centerZ}
	c.cam.up = vec3{upX, upY, upZ}
}

// Perspective sets the vertical field of view (in radians), and the near and
// far clipping planes. The aspect ratio is derived from the canvas size.
func (c *Canvas) Perspective(fov, near, far float64) {
	c.cam.fov = fov
	c.cam.near = near
	c.cam.far = far
}

// Translate3D offsets the 3D model transform by (x, y, z).
func (c *Canvas) Translate3D(x, y, z float64) {
	c.model3D = c.model3D.mul(mat4Translate(x, y, z))
}

// RotateX rotates the 3D model transform around the x axis (radians).
func (c *Canvas) RotateX(angle float64) {
	c.model3D = c.model3D.mul(mat4RotX(angle))
}

// RotateY rotates the 3D model transform around the y axis (radians).
func (c *Canvas) RotateY(angle float64) {
	c.model3D = c.model3D.mul(mat4RotY(angle))
}

// RotateZ rotates the 3D model transform around the z axis (radians).
func (c *Canvas) RotateZ(angle float64) {
	c.model3D = c.model3D.mul(mat4RotZ(angle))
}

// Scale3D scales the 3D model transform uniformly by the given factor.
func (c *Canvas) Scale3D(scale float64) {
	c.model3D = c.model3D.mul(mat4Scale(scale, scale, scale))
}

// Push3D saves the current 3D model transform onto the stack.
func (c *Canvas) Push3D() {
	c.model3DStack = append(c.model3DStack, c.model3D)
}

// Pop3D restores the last saved 3D model transform.
func (c *Canvas) Pop3D() {
	n := len(c.model3DStack)
	if n == 0 {
		return
	}
	c.model3D = c.model3DStack[n-1]
	c.model3DStack = c.model3DStack[:n-1]
}

// Light sets the direction the light comes from. The vector points from the
// scene toward the light source and is normalized internally.
func (c *Canvas) Light(x, y, z float64) {
	d := vec3{x, y, z}.normalize()
	if d != (vec3{}) {
		c.lightDir = d
	}
}

// Ambient sets the ambient light level (0..1) applied to every face
// regardless of orientation.
func (c *Canvas) Ambient(level float64) {
	c.ambient = clamp(level, 0, 1)
}

// ShadeChars sets the character ramp used to shade solid 3D faces, ordered
// from darkest to brightest.
func (c *Canvas) ShadeChars(chars string) {
	if len(chars) == 0 {
		return
	}
	c.shadeRunes = []rune(chars)
}

// Cull enables or disables back-face culling for solid 3D shapes.
func (c *Canvas) Cull(enabled bool) {
	c.cull = enabled
}

// CharAspect sets the assumed height-to-width ratio of a terminal character
// cell, used to keep 3D geometry from looking squished. The default of 2.0
// suits most terminals; lower it if shapes look stretched, raise it if they
// look squashed.
func (c *Canvas) CharAspect(ratio float64) {
	if ratio > 0 {
		c.charAspect = ratio
	}
}

// reset3D restores the per-frame 3D transform state. Lighting, culling and
// the shade ramp are settings and persist across frames.
func (c *Canvas) reset3D() {
	c.cam = defaultCamera()
	c.model3D = mat4Identity()
	c.model3DStack = c.model3DStack[:0]
}

// worldPoint applies the current 3D model transform to a point.
func (c *Canvas) worldPoint(p vec3) vec3 {
	return c.model3D.mulPoint(p)
}

// projectClip transforms a world-space point through the view and projection
// matrices and maps it to canvas coordinates. It returns the depth (NDC z,
// smaller is closer) and ok is false when the point is behind the camera.
//
// The returned coordinates are centered on the origin; callers apply the 2D
// transform afterwards so Translate can move the result to the canvas center.
func (c *Canvas) projectClip(w vec3) (int, int, float64, bool) {
	view := mat4LookAt(c.cam.eye, c.cam.center, c.cam.up)

	aspect := 1.0
	if c.Height != 0 {
		aspect = c.cellPhysAspect() * float64(c.Width) / float64(c.Height)
	}
	proj := mat4Perspective(c.cam.fov, aspect, c.cam.near, c.cam.far)

	vp := proj.mul(view)
	x, y, z, wc := vp.mulVec4(w.x, w.y, w.z, 1)
	if wc <= 1e-6 {
		return 0, 0, 0, false
	}

	ndcX := x / wc
	ndcY := y / wc
	ndcZ := z / wc

	// Per-axis viewport: NDC [-1,1] maps across the full canvas. The aspect
	// above (which folds in the character-cell shape) keeps geometry
	// undistorted, so a world cube renders as a cube on screen.
	sx := int(math.Round(ndcX * float64(c.Width) / 2))
	sy := int(math.Round(-ndcY * float64(c.Height) / 2))
	return sx, sy, ndcZ, true
}

// cellPhysAspect returns the physical width-to-height ratio of a single
// logical canvas cell. A bare terminal cell is about twice as tall as it is
// wide; the cell modes render each cell two characters wide, making it square.
func (c *Canvas) cellPhysAspect() float64 {
	charsWide := 1.0
	if c.cellMode.enabled() {
		charsWide = 2.0
	}
	ca := c.charAspect
	if ca == 0 {
		ca = defaultCharAspect
	}
	return charsWide / ca
}

// projectVertex transforms a model-space point all the way to final canvas
// coordinates (view, projection and the 2D transform) plus its depth.
func (c *Canvas) projectVertex(p vec3) (int, int, float64, bool) {
	sx, sy, z, ok := c.projectClip(c.worldPoint(p))
	if !ok {
		return 0, 0, 0, false
	}
	tx, ty := c.applyTransform2D(sx, sy)
	return tx, ty, z, true
}
