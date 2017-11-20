//line vector.vgo:1
package samples

func NormFloat32s(out, x, y []float32)

//line vector.vgo:7
func AddUints(out, in1, in2 []uint)

//line vector.vgo:11
func SomeFormula(out, x, y, z []float32)

//line vector.vgo:15
func subByte(out, a, b []byte)

//line vector.vgo:19
func Diff(a []byte) []byte {
	result := make([]byte, len(a)-1)
	subByte(result, a[1:], a[:len(a)-1])
	return result
}

func subuint(out, a, b []uint)

//line vector.vgo:29
func DiffInt(a []uint) []uint {
	result := make([]uint, len(a)-1)
	subuint(result, a[:len(a)-1], a[1:])
	return result
}

func DetF32(det, a, b, c, d []float32)

//line vector.vgo:39

// vim: set ft=go:
