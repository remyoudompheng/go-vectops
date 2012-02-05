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
	"float32": tF32,
	"uint":    tU32,
	"uint32":  tU32,
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
		tU32: Type{
			Size:    4,
			LogSize: 2,
			Ops: map[token.Token]string{
				token.ADD: "PADDL",
				token.MUL: "PMULLL",
				token.SUB: "PSUBL",
				token.AND: "PAND",
				token.OR:  "POR",
				token.XOR: "PXOR",
			},
		},
		tF32: Type{
			Size:    4,
			LogSize: 2,
			Ops: map[token.Token]string{
				token.ADD: "ADDPS",
				token.MUL: "MULPS",
				token.SUB: "MULPS",
				token.QUO: "MULPS",
			},
		},
	},
}
