package samples

import "testing"

func AddV(out, in1, in2 []int)

func SlowAddV(out, in1, in2 []int) {
	for i := range out {
		out[i] = in1[i] + in2[i]
	}
}

func TestAddV(t *testing.T) {
	out := make([]int, 256)
	in1 := make([]int, 256)
	in2 := make([]int, 256)
	for i := range out {
		in1[i] = i * i
		in2[i] = i + 1
	}
	AddV(out, in1, in2)
	for i, v := range out {
		if v != i*i+i+1 {
			t.Errorf("wrong output out[%d] != in1[%d] + in2[%d] (%d != %d + %d)",
				i, i, i, out[i], in1[i], in2[i])
		}
	}
}

func TestAddV253(t *testing.T) {
	out := make([]int, 253)
	in1 := make([]int, 253)
	in2 := make([]int, 253)
	for i := range out {
		in1[i] = i * i
		in2[i] = i + 1
	}
	AddV(out, in1, in2)
	for i, v := range out {
		if v != i*i+i+1 {
			t.Errorf("wrong output out[%d] != in1[%d] + in2[%d] (%d != %d + %d)",
				i, i, i, out[i], in1[i], in2[i])
		}
	}
}

func BenchmarkSlowAddV(b *testing.B) {
	out := make([]int, 256)
	in1 := make([]int, 256)
	in2 := make([]int, 256)
	b.SetBytes(256)
	for i := range out {
		in1[i] = i * i
		in2[i] = i + 1
	}
	for n := 0; n < b.N; n++ {
		SlowAddV(out, in1, in2)
	}
}

func BenchmarkAddV(b *testing.B) {
	out := make([]int, 256)
	in1 := make([]int, 256)
	in2 := make([]int, 256)
	b.SetBytes(256)
	for i := range out {
		in1[i] = i * i
		in2[i] = i + 1
	}
	for n := 0; n < b.N; n++ {
		AddV(out, in1, in2)
	}
}
