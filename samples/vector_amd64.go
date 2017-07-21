//line vector.vgo:1
package samples

//line vector.vgo:4

//line vector.vgo:3
func NormFloat32s(out, x, y []float32)

//line vector.vgo:8

//line vector.vgo:7
func AddUints(out, in1, in2 []uint)

//line vector.vgo:12

//line vector.vgo:11
func SomeFormula(out, x, y, z []float32)

//line vector.vgo:16

//line vector.vgo:15
func subByte(out, a, b []byte)

//line vector.vgo:20

//line vector.vgo:19
func Diff(a []byte) []byte {
	result := make([]byte, len(a)-1)
	subByte(result, a[1:], a[:len(a)-1])
	return result
}

//line vector.vgo:26

//line vector.vgo:25
func subuint(out, a, b []uint)

//line vector.vgo:30

//line vector.vgo:29
func DiffInt(a []uint) []uint {
	result := make([]uint, len(a)-1)
	subuint(result, a[:len(a)-1], a[1:])
	return result
}

//line vector.vgo:36

//line vector.vgo:35
func DetF64(det, a, b, c, d []float64)

//line vector.vgo:39

// vim: set ft=go:
