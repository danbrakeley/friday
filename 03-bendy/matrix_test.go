package main

import "testing"

func assertMatrix(t *testing.T, m, n Matrix3x3) {
	t.Helper()
	for i := 0; i < 9; i++ {
		if m[i] != n[i] {
			t.Errorf("matrix mismatch: %v != %v", m, n)
			return
		}
	}
}

func TestMatrix3x3_Multiply(t *testing.T) {
	iden := Identity3x3()
	a := Matrix3x3{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	}
	assertMatrix(t, iden.Multiply(a), a)
	assertMatrix(t, a.Multiply(iden), a)

	b := Matrix3x3{
		9, 8, 7,
		6, 5, 4,
		3, 2, 1,
	}
	assertMatrix(t, a.Multiply(b), Matrix3x3{
		30, 24, 18,
		84, 69, 54,
		138, 114, 90,
	})
	assertMatrix(t, b.Multiply(a), Matrix3x3{
		90, 114, 138,
		54, 69, 84,
		18, 24, 30,
	})
}
