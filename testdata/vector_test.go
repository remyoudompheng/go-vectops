package vectops

import (
	"testing"
)

func TestUints(t *testing.T) {
	a := make([]uint, 512)
	b := make([]uint, 512)
	c := make([]uint, 512)
	for i := range a {
		a[i] = uint(i)
		b[i] = 2 * uint(i)
	}
	AddUints(c, a, b)
	for i, x := range c {
		if x != 3*uint(i) {
			t.Errorf("c[%d] = %d, expected %d", i, x, 3*i)
		}
	}
}

func TestNormF32(t *testing.T) {
	a := make([]float32, 512)
	b := make([]float32, 512)
	c := make([]float32, 512)
	for i := range c {
		a[i] = float32(i) / 12
		b[i] = float32((i + 3) / 4)
	}
	NormFloat32s(c, a, b)
	for i, n := range c {
		if n != a[i]*a[i]+b[i]*b[i] {
			t.Errorf("c[%d] = %-1g, expected %-1g", i, n, a[i]*a[i]+b[i]*b[i])
		}
	}
}
