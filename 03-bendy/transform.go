package main

import "math"

type Vec2D struct {
	X, Y float32
}

type Transform2D struct {
	position    Vec2D
	rotation    float32
	scale       float32
	matrix      Matrix3x3
	needsUpdate bool
}

func (t *Transform2D) Pos() Vec2D {
	return t.position
}

func (t *Transform2D) Rot() float32 {
	return t.rotation
}

func (t *Transform2D) Scale() float32 {
	return t.scale
}

func (t *Transform2D) Matrix() Matrix3x3 {
	if t.needsUpdate {
		t.needsUpdate = false
		t.matrix = t.calcMatrix()
	}
	return t.matrix
}

func (t *Transform2D) SetPos(p Vec2D) {
	t.position = p
	t.needsUpdate = true
}

func (t *Transform2D) SetRot(r float32) {
	t.rotation = r
	t.needsUpdate = true
}

func (t *Transform2D) SetScale(s float32) {
	t.scale = s
	t.needsUpdate = true
}

func (t *Transform2D) calcMatrix() Matrix3x3 {
	var m Matrix3x3
	theta := float64(t.rotation)
	c := float32(math.Cos(theta))
	s := float32(math.Sin(theta))
	m[0] = c * t.scale
	m[1] = -s * t.scale
	m[2] = t.position.X
	m[3] = s * t.scale
	m[4] = c * t.scale
	m[5] = t.position.Y
	m[6] = 0
	m[7] = 0
	m[8] = 1
	return m
}
