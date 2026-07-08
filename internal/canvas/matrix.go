package canvas

import "math"

// vec3 is a 3D vector used by the 3D drawing pipeline.
type vec3 struct {
	x, y, z float64
}

func (v vec3) add(o vec3) vec3 { return vec3{v.x + o.x, v.y + o.y, v.z + o.z} }
func (v vec3) sub(o vec3) vec3 { return vec3{v.x - o.x, v.y - o.y, v.z - o.z} }

func (v vec3) cross(o vec3) vec3 {
	return vec3{
		v.y*o.z - v.z*o.y,
		v.z*o.x - v.x*o.z,
		v.x*o.y - v.y*o.x,
	}
}

func (v vec3) dot(o vec3) float64 {
	return v.x*o.x + v.y*o.y + v.z*o.z
}

func (v vec3) normalize() vec3 {
	l := math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
	if l == 0 {
		return v
	}
	return vec3{v.x / l, v.y / l, v.z / l}
}

// mat4 is a 4x4 matrix stored row-major. Points are treated as column
// vectors and transformed as m * v.
type mat4 [4][4]float64

func mat4Identity() mat4 {
	return mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// mul returns m * o.
func (m mat4) mul(o mat4) mat4 {
	var r mat4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			r[i][j] = m[i][0]*o[0][j] +
				m[i][1]*o[1][j] +
				m[i][2]*o[2][j] +
				m[i][3]*o[3][j]
		}
	}
	return r
}

// mulVec4 transforms the homogeneous point (x, y, z, w) by m.
func (m mat4) mulVec4(x, y, z, w float64) (float64, float64, float64, float64) {
	return m[0][0]*x + m[0][1]*y + m[0][2]*z + m[0][3]*w,
		m[1][0]*x + m[1][1]*y + m[1][2]*z + m[1][3]*w,
		m[2][0]*x + m[2][1]*y + m[2][2]*z + m[2][3]*w,
		m[3][0]*x + m[3][1]*y + m[3][2]*z + m[3][3]*w
}

// mulPoint transforms an affine point (w = 1) and returns the xyz result,
// discarding w. Used for model-space transforms that never project.
func (m mat4) mulPoint(p vec3) vec3 {
	x, y, z, _ := m.mulVec4(p.x, p.y, p.z, 1)
	return vec3{x, y, z}
}

func mat4Translate(x, y, z float64) mat4 {
	return mat4{
		{1, 0, 0, x},
		{0, 1, 0, y},
		{0, 0, 1, z},
		{0, 0, 0, 1},
	}
}

func mat4Scale(x, y, z float64) mat4 {
	return mat4{
		{x, 0, 0, 0},
		{0, y, 0, 0},
		{0, 0, z, 0},
		{0, 0, 0, 1},
	}
}

func mat4RotX(a float64) mat4 {
	s, c := math.Sincos(a)
	return mat4{
		{1, 0, 0, 0},
		{0, c, -s, 0},
		{0, s, c, 0},
		{0, 0, 0, 1},
	}
}

func mat4RotY(a float64) mat4 {
	s, c := math.Sincos(a)
	return mat4{
		{c, 0, s, 0},
		{0, 1, 0, 0},
		{-s, 0, c, 0},
		{0, 0, 0, 1},
	}
}

func mat4RotZ(a float64) mat4 {
	s, c := math.Sincos(a)
	return mat4{
		{c, -s, 0, 0},
		{s, c, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// mat4Perspective builds a right-handed perspective projection matrix.
// fovY is the vertical field of view in radians.
func mat4Perspective(fovY, aspect, near, far float64) mat4 {
	f := 1 / math.Tan(fovY/2)
	nf := 1 / (near - far)
	return mat4{
		{f / aspect, 0, 0, 0},
		{0, f, 0, 0},
		{0, 0, (far + near) * nf, 2 * far * near * nf},
		{0, 0, -1, 0},
	}
}

// mat4LookAt builds a right-handed view matrix looking from eye to center.
func mat4LookAt(eye, center, up vec3) mat4 {
	f := center.sub(eye).normalize()
	s := f.cross(up).normalize()
	u := s.cross(f)
	return mat4{
		{s.x, s.y, s.z, -s.dot(eye)},
		{u.x, u.y, u.z, -u.dot(eye)},
		{-f.x, -f.y, -f.z, f.dot(eye)},
		{0, 0, 0, 1},
	}
}
