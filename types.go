package main

import (
	"fmt"
	"go/token"
)

type GoType uint

const (
	tUndefined GoType = iota
	tU8
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
}

type Arch struct {
	PtrSize     int
	UintType    GoType // the type of unsized uint
	CounterReg  string
	InputRegs   []string
	VectorRegs  []string
	VectorWidth int
	Types       map[GoType]Type
}

func (a *Arch) Opcode(op token.Token, typename string) (opcode string, ok bool) {
	typ := types[typename]
	if typename == "uint" {
		typ = a.UintType
	}
	if t, ok := a.Types[typ]; ok {
		op, ok := t.Ops[op]
		return op, ok
	}
	return "", false
}

func IsCommutative(op token.Token) bool {
	switch op {
	case token.ADD, token.AND, token.OR, token.XOR, token.MUL:
		return true
	}
	return false
}

func (a *Arch) Width(typename string) int {
	typ := types[typename]
	if typename == "uint" {
		typ = a.UintType
	}
	return a.Types[typ].Size
}

type Type struct {
	Size    int
	LogSize int // Size = 1 << LogSize
	Ops     map[token.Token]string
}

func FindArch(goarch, goarm string) Arch {
	switch goarch {
	case "amd64":
		return amd64
	default:
		err := fmt.Errorf("unsupported GOARCH=%q", goarch)
		panic(err)
	}
}

// Description of the amd64 architecture output.
var amd64 = Arch{
	PtrSize:    8,
	UintType:   tU64,
	CounterReg: "CX",
	InputRegs: []string{"BX", "SI", "DI",
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
