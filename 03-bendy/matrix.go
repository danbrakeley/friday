package main

// Matrix3x3 is a 3x3 matrix stored in row-major order.
// [0 1 2]
// [3 4 5]
// [6 7 8]
type Matrix3x3 [9]float32

func Identity3x3() Matrix3x3 {
	return Matrix3x3{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
}

// Multiply returns M∙N. In practical terms, this means that the
// matrix on the right (N) is applied first.
func (m Matrix3x3) Multiply(n Matrix3x3) Matrix3x3 {
	var r Matrix3x3
	r[0] = m[0]*n[0] + m[1]*n[3] + m[2]*n[6]
	r[1] = m[0]*n[1] + m[1]*n[4] + m[2]*n[7]
	r[2] = m[0]*n[2] + m[1]*n[5] + m[2]*n[8]
	r[3] = m[3]*n[0] + m[4]*n[3] + m[5]*n[6]
	r[4] = m[3]*n[1] + m[4]*n[4] + m[5]*n[7]
	r[5] = m[3]*n[2] + m[4]*n[5] + m[5]*n[8]
	r[6] = m[6]*n[0] + m[7]*n[3] + m[8]*n[6]
	r[7] = m[6]*n[1] + m[7]*n[4] + m[8]*n[7]
	r[8] = m[6]*n[2] + m[7]*n[5] + m[8]*n[8]
	return r
}

func (m Matrix3x3) MultiplyVec2D(v Vec2D) Vec2D {
	return Vec2D{
		m[0]*v.X + m[1]*v.Y + m[2],
		m[3]*v.X + m[4]*v.Y + m[5],
	}
}
