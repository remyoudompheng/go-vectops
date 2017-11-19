package main

import (
	"testing"
)

func TestAsmNEON(t *testing.T) {
	for _, x := range testsNEON {
		op := assembleNEON(x.Opcode, x.Quad, x.Reg1, x.Reg2, x.RegDest)
		if op != x.Encoding {
			t.Errorf("encoding %s D%d, D%d, D%d = 0x%08x expected 0x%08x",
				x.Opcode, x.Reg1, x.Reg2, x.RegDest, op, x.Encoding)
		}
	}
}

type testArmOp struct {
	Opcode   string
	Quad     bool
	Reg1     int
	Reg2     int
	RegDest  int
	Encoding uint32
}

var testsNEON = []testArmOp{
	{"VADD.F32", false, 1, 2, 3, 0xf2013d02},
	{"VADD.F32", true, 0, 2, 4, 0xf2004d42},
	{"VADD.I32", false, 1, 2, 3, 0xf2213802},
	{"VADD.I16", false, 1, 2, 3, 0xf2113802},
}
