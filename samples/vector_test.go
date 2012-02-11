package vectops

import (
	"testing"
)

func TestUints(t *testing.T) {
	var a, b, c [512]uint
	for i := range a {
		a[i] = uint(i)
		b[i] = 2 * uint(i)
	}
	AddUints(c[:], a[:], b[:])
	for i, x := range c {
		if x != 3*uint(i) {
			t.Errorf("c[%d] = %d, expected %d", i, x, 3*i)
		}
	}
}

// Code that should panic.
func TestPanic(t *testing.T) {
	defer func() { recover() }()
	var a, b [512]uint
	var c [510]uint
	AddUints(c[:], a[:], b[:])
	t.Errorf("should have panicked")
}

func TestNormF32(t *testing.T) {
	var a, b, c [512]float32
	for i := range c {
		a[i] = float32(i) / 12
		b[i] = float32((i + 3) / 4)
	}
	NormFloat32s(c[:], a[:], b[:])
	for i, n := range c {
		if n != a[i]*a[i]+b[i]*b[i] {
			t.Errorf("c[%d] = %-1g, expected %-1g", i, n, a[i]*a[i]+b[i]*b[i])
		}
	}
}

func TestFormula(t *testing.T) {
	var a, b, c, out [512]float32
	for i := range out {
		a[i] = 1 + float32(i)/12
		b[i] = 1 + float32(i)/4
		c[i] = 1 + float32(i)/2
	}
	SomeFormula(out[:], a[:], b[:], c[:])
	for i, xout := range out {
		x, y, z := a[i], b[i], c[i]
		expect := (x*y + y*z + z*x) / (x*x + y*y + z*z)
		if xout != expect {
			t.Errorf("c[%d] = %-1g, expected %-1g", i, xout, expect)
		}
	}
}

func TestDiff(t *testing.T) {
	var a [512]byte
	for i := range a {
		a[i] = 'a' + byte(i/10)
	}
	// input of substract won't be multiple of 128-bit.
	c := Diff(a[:])
	for i, x := range c {
		if x != a[i+1]-a[i] {
			t.Errorf("got c[%d] = %d, expected %d", i, x, a[i+1]-a[i])
		}
	}
}

func TestDiffInt(t *testing.T) {
	var a [257]uint
	for i := range a {
		a[i] = uint(i * i)
	}
	c := DiffInt(a[:])
	for i, x := range c {
		if x != a[i]-a[i+1] {
			t.Errorf("got c[%d] = %d, expected %d", i, x, a[i+1]-a[i])
		}
	}
}

func BenchmarkDiff(b *testing.B) {
	const length = 1 << 16
	// a[1:] must have length multiple of 16.
	var a [length]byte
	for i := 0; i < b.N; i++ {
		Diff(a[:])
	}
	b.SetBytes(length)
}

func BenchmarkDiffNoAlloc(b *testing.B) {
	const length = 1 << 16
	var a [length + 1]byte
	var out [length]byte
	for i := 0; i < b.N; i++ {
		subByte(out[:], a[1:], a[:length])
	}
	b.SetBytes(length)
}

func BenchmarkDiffNoSIMD(b *testing.B) {
	const length = 1 << 16
	var a [length + 1]byte
	var out [length]byte
	for i := 0; i < b.N; i++ {
		for i, x := range a[1:] {
			out[i] = x - a[i]
		}
	}
	b.SetBytes(length)
}
