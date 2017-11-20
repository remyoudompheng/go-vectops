package main

import (
	"fmt"
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

type Op int

const (
	ADD = iota
	SUB
	MUL
	DIV
	AND
	XOR
	OR
	SHL
	SHR
)

var opstring = [...]string{
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
	AND: "&",
	OR:  "|",
	XOR: "^",
	SHL: "<<",
	SHR: ">>",
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

func (a *Arch) Opcode(op Op, typename string) (opcode string, ok bool) {
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
	Ops     map[Op]string
}

func FindArch(goarch, gosubarch string) Arch {
	switch goarch {
	case "amd64":
		if gosubarch == "avx2" {
			return avx2
		}
		return amd64
	case "arm":
		return armv7
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
			Ops: map[Op]string{
				ADD: "PADDB",
				SUB: "PSUBB",
				AND: "PAND",
				OR:  "POR",
				XOR: "PXOR",
			},
		},
		tU16: Type{
			Size:    2,
			LogSize: 1,
			Ops: map[Op]string{
				ADD: "PADDW",
				MUL: "PMULLW",
				SUB: "PSUBW",
				AND: "PAND",
				OR:  "POR",
				XOR: "PXOR",
				SHL: "PSLLW",
				SHR: "PSRLW",
			},
		},
		tU32: Type{
			Size:    4,
			LogSize: 2,
			Ops: map[Op]string{
				ADD: "PADDL",
				SUB: "PSUBL",
				AND: "PAND",
				OR:  "POR",
				XOR: "PXOR",
				SHL: "PSLLL",
				SHR: "PSRLL",
			},
		},
		tU64: Type{
			Size:    8,
			LogSize: 3,
			Ops: map[Op]string{
				ADD: "PADDQ",
				SUB: "PSUBQ",
				AND: "PAND",
				OR:  "POR",
				XOR: "PXOR",
				SHL: "PSLLQ",
				SHR: "PSRLQ",
			},
		},
		tF32: Type{
			Size:    4,
			LogSize: 2,
			Ops: map[Op]string{
				ADD: "ADDPS",
				MUL: "MULPS",
				SUB: "SUBPS",
				DIV: "DIVPS",
			},
		},
		tF64: Type{
			Size:    8,
			LogSize: 3,
			Ops: map[Op]string{
				ADD: "ADDPD",
				MUL: "MULPD",
				SUB: "SUBPD",
				DIV: "DIVPD",
			},
		},
	},
}

var avx2 = Arch{
	PtrSize:    8,
	UintType:   tU64,
	CounterReg: "CX",
	LengthReg:  "DX",
	InputRegs: []string{"BX", "SI", "DI",
		"R8", "R9", "R10", "R11",
		"R12", "R13", "R14", "R15"},
	VectorRegs: []string{
		"Y0", "Y1", "Y2", "Y4", "Y5", "Y6", "Y7",
		"Y8", "Y9", "Y10", "Y11", "Y12", "Y13", "Y14", "Y15"},
	VectorWidth: 32, // 256-bit AVX registers
	Types: map[GoType]Type{
		tU8: Type{
			Size:    1,
			LogSize: 0,
			Ops: map[Op]string{
				ADD: "VPADDB",
				SUB: "VPSUBB",
				AND: "VPAND",
				OR:  "VPOR",
				XOR: "VPXOR",
			},
		},
		tU16: Type{
			Size:    2,
			LogSize: 1,
			Ops: map[Op]string{
				ADD: "VPADDW",
				MUL: "VPMULLW",
				SUB: "VPSUBW",
				AND: "VPAND",
				OR:  "VPOR",
				XOR: "VPXOR",
				SHL: "VPSLLW",
				SHR: "VPSRLW",
			},
		},
		tU32: Type{
			Size:    4,
			LogSize: 2,
			Ops: map[Op]string{
				ADD: "VPADDL",
				SUB: "VPSUBL",
				AND: "VPAND",
				OR:  "VPOR",
				XOR: "VPXOR",
				SHL: "VPSLLL",
				SHR: "VPSRLL",
			},
		},
		tU64: Type{
			Size:    8,
			LogSize: 3,
			Ops: map[Op]string{
				ADD: "VPADDQ",
				SUB: "VPSUBQ",
				AND: "VPAND",
				OR:  "VPOR",
				XOR: "VPXOR",
				SHL: "VPSLLQ",
				SHR: "VPSRLQ",
			},
		},
		tF32: Type{
			Size:    4,
			LogSize: 2,
			Ops: map[Op]string{
				ADD: "VADDPS",
				MUL: "VMULPS",
				SUB: "VSUBPS",
				DIV: "VDIVPS",
			},
		},
		tF64: Type{
			Size:    8,
			LogSize: 3,
			Ops: map[Op]string{
				ADD: "VADDPD",
				MUL: "VMULPD",
				SUB: "VSUBPD",
				DIV: "VDIVPD",
			},
		},
	},
}

// Description of the ARM NEON SIMD instructions
// Looks good on Cortex-A9 and Cortex-A53.
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
		"Q0", "Q1", "Q2", "Q3", "Q4", "Q5", "Q6", "Q7"},
	VectorWidth: 16, // 128-bit registers
	Types:       neonTypes,
}

var neonTypes = map[GoType]Type{
	tU8: Type{
		Size:    1,
		LogSize: 0,
		Ops: map[Op]string{
			ADD: "VADD.I8",
			SUB: "VSUB.I8",
			MUL: "VMUL.I8",
			AND: "VAND",
			OR:  "VORR",
			XOR: "VEOR",
		},
	},
	tU16: Type{
		Size:    2,
		LogSize: 1,
		Ops: map[Op]string{
			ADD: "VADD.I16",
			SUB: "VSUB.I16",
			MUL: "VMUL.I16",
			AND: "VAND",
			OR:  "VORR",
			XOR: "VEOR",
		},
	},
	tU32: Type{
		Size:    4,
		LogSize: 2,
		Ops: map[Op]string{
			ADD: "VADD.I32",
			SUB: "VSUB.I32",
			MUL: "VMUL.I32",
			AND: "VAND",
			OR:  "VORR",
			XOR: "VEOR",
		},
	},
	tF32: Type{
		Size:    4,
		LogSize: 2,
		Ops: map[Op]string{
			ADD: "VADD.F32",
			MUL: "VMUL.F32",
			SUB: "VSUB.F32",
		},
	},
}
