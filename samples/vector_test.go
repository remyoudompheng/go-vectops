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

func TestDiff(t *testing.T) {
	// 257 is for a[1:] to have length multiple of 16.
	a := make([]byte, 257)
	for i := range a {
		a[i] = 'a' + byte(i/10)
	}
	c := Diff(a)
	for i, x := range c {
		if x != a[i+1]-a[i] {
			t.Errorf("got c[%d] = %d, expected %d", i, x, a[i+1]-a[i])
		}
	}
}

func BenchmarkDiff(b *testing.B) {
	const length = 1 << 16
	// a[1:] must have length multiple of 16.
	a := make([]byte, length+1)
	for i := 0; i < b.N; i++ {
		Diff(a)
	}
	b.SetBytes(length)
}

func BenchmarkDiffNoAlloc(b *testing.B) {
	const length = 1 << 16
	// a[1:] must have length multiple of 16.
	a := make([]byte, length+1)
	out := make([]byte, length)
	for i := 0; i < b.N; i++ {
		subByte(out, a[1:], a[:length])
	}
	b.SetBytes(length)
}

func BenchmarkDiffNoSIMD(b *testing.B) {
	const length = 1 << 16
	a := make([]byte, length+1)
	out := make([]byte, length)
	for i := 0; i < b.N; i++ {
		for i, x := range a[1:] {
			out[i] = x - a[i]
		}
	}
	b.SetBytes(length)
}
