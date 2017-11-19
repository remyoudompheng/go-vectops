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
	LengthReg   string
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

func (a *Arch) LogWidth(typename string) int {
	typ := types[typename]
	if typename == "uint" {
		typ = a.UintType
	}
	return a.Types[typ].LogSize
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
	case "arm":
		switch goarm {
		case "7":
			return armv7
		case "8":
			return armv8
		default:
			panic("unsupported goarm=" + goarm)
		}
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
	LengthReg:  "DX",
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

// Description of the ARM NEON SIMD instructions
var armv7 = Arch{
	PtrSize:  4,
	UintType: tU32,
	// R10 is g, R13 is sp, R14 is lr, R15 is pc.
	CounterReg: "R11",
	LengthReg:  "R12",
	InputRegs: []string{
		"R1", "R2", "R3", "R4",
		"R5", "R6", "R7", "R8"},
	VectorRegs: []string{
		"F0", "F1", "F2", "F3", "F4", "F5", "F6", "F7",
		"F8", "F9", "F10", "F11", "F12", "F13", "F14", "F15"},
	VectorWidth: 8,
	Types:       neonTypes,
}

// ARMv8 in 32-bit mode (Aarch32)
var armv8 = Arch{
	PtrSize:    4,
	UintType:   tU32,
	CounterReg: "R11",
	LengthReg:  "R12",
	InputRegs: []string{
		"R1", "R2", "R3", "R4",
		"R5", "R6", "R7", "R8"},
	VectorRegs: []string{
		"Q0", "Q1", "Q2", "Q3", "Q4", "Q5", "Q6", "Q7"},
	VectorWidth: 16, // 128-bit registers on Aarch32
	Types:       neonTypes,
}

var neonTypes = map[GoType]Type{
	tU8: Type{
		Size:    1,
		LogSize: 0,
		Ops: map[token.Token]string{
			token.ADD: "VADD.I8",
			token.SUB: "VSUB.I8",
			token.MUL: "VMUL.I8",
			token.AND: "VAND",
			token.OR:  "VORR",
			token.XOR: "VEOR",
		},
	},
	tU16: Type{
		Size:    2,
		LogSize: 1,
		Ops: map[token.Token]string{
			token.ADD: "VADD.I16",
			token.SUB: "VSUB.I16",
			token.MUL: "VMUL.I16",
			token.AND: "VAND",
			token.OR:  "VORR",
			token.XOR: "VEOR",
		},
	},
	tU32: Type{
		Size:    4,
		LogSize: 2,
		Ops: map[token.Token]string{
			token.ADD: "VADD.I32",
			token.SUB: "VSUB.I32",
			token.MUL: "VMUL.I32",
			token.AND: "VAND",
			token.OR:  "VORR",
			token.XOR: "VEOR",
		},
	},
	tF32: Type{
		Size:    4,
		LogSize: 2,
		Ops: map[token.Token]string{
			token.ADD: "VADD.F32",
			token.MUL: "VMUL.F32",
			token.SUB: "VSUB.F32",
		},
	},
}
