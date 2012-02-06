package main

import (
	"go/token"
)

type GoType uint

const (
	tU8 GoType = iota
	tU16
	tU32
	tU64
	tF32
	tF64
)

var types = map[string]GoType{
	"uint8":   tU8,
	"uint32":  tU32,
	"float32": tF32,
	"float64": tF64,
	"byte":    tU8,
	"uint":    tU32,
}

type Arch struct {
	CounterReg  string
	InputRegs   []string
	VectorRegs  []string
	VectorWidth int
	Types       map[GoType]Type
}

func (a *Arch) Opcode(op token.Token, typename string) (opcode string, ok bool) {
	if t, ok := a.Types[types[typename]]; ok {
		op, ok := t.Ops[op]
		return op, ok
	}
	return "", false
}

func (a *Arch) Width(typename string) int {
	return a.Types[types[typename]].Size
}

type Type struct {
	Size    int
	LogSize int // Size = 1 << LogSize
	Ops     map[token.Token]string
}

// Description of the amd64 architecture output.
var amd64 = Arch{
	CounterReg: "CX",
	InputRegs: []string{"SI", "DI",
		"R8", "R9", "R10", "R11",
		"R12", "R13", "R14", "R15"},
	VectorRegs: []string{
		"X0", "X1", "X2", "X4", "X5", "X6", "X7",
		"X8", "X9", "X10", "X11", "X12", "X13", "X14", "X15"},
	VectorWidth: 16,
	Types: map[GoType]Type{
		tU8: Type{
			Size:    1,
			LogSize: 0,
			Ops: map[token.Token]string{
				token.ADD: "PADDB",
				token.SUB: "PSUBB",
				token.AND: "PAND",
				token.OR:  "POR",
				token.XOR: "PXOR",
			},
		},
		tU16: Type{
			Size:    2,
			LogSize: 1,
			Ops: map[token.Token]string{
				token.ADD: "PADDW",
				token.MUL: "PMULLW",
				token.SUB: "PSUBW",
				token.AND: "PAND",
				token.OR:  "POR",
				token.XOR: "PXOR",
				token.SHL: "PSLLW",
				token.SHR: "PSRLW",
			},
		},
		tU32: Type{
			Size:    4,
			LogSize: 2,
			Ops: map[token.Token]string{
				token.ADD: "PADDL",
				token.SUB: "PSUBL",
				token.AND: "PAND",
				token.OR:  "POR",
				token.XOR: "PXOR",
				token.SHL: "PSLLL",
				token.SHR: "PSRLL",
			},
		},
		tU64: Type{
			Size:    8,
			LogSize: 3,
			Ops: map[token.Token]string{
				token.ADD: "PADDQ",
				token.SUB: "PSUBQ",
				token.AND: "PAND",
				token.OR:  "POR",
				token.XOR: "PXOR",
				token.SHL: "PSLLQ",
				token.SHR: "PSRLQ",
			},
		},
		tF32: Type{
			Size:    4,
			LogSize: 2,
			Ops: map[token.Token]string{
				token.ADD: "ADDPS",
				token.MUL: "MULPS",
				token.SUB: "SUBPS",
				token.QUO: "DIVPS",
			},
		},
		tF64: Type{
			Size:    8,
			LogSize: 3,
			Ops: map[token.Token]string{
				token.ADD: "ADDPD",
				token.MUL: "MULPD",
				token.SUB: "SUBPD",
				token.QUO: "DIVPD",
			},
		},
	},
}
