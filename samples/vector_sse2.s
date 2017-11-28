// +build sse2

// func NormFloat32s(out, x, y []float32)
TEXT ·NormFloat32s(SB), 7, $0

	// Load pointers.
	MOVQ	out+0(FP), BX
	MOVQ	x+24(FP), SI
	MOVQ	y+48(FP), DI

	// Check lengths.
	MOVQ	out+8(FP), DX
	CMPQ	DX, x+32(FP)
	JNE	NormFloat32s__panic
	CMPQ	DX, y+56(FP)
	JNE	NormFloat32s__panic
	JMP	NormFloat32s__ok

NormFloat32s__panic:
	CALL	runtime·panicindex(SB)

NormFloat32s__ok:
	SUBQ	$4, DX
	MOVQ	$0, CX

NormFloat32s__loop:
	CMPQ	CX, DX

	// if i > length-4 { i = length-4 }
	JLE	NormFloat32s__process
	MOVQ	DX, CX

NormFloat32s__process:
	MOVUPS	(SI)(CX*4), X0

	// __auto_tmp_000 = x * x
	MOVAPS	X0, X2
	MULPS	X0, X2
	MOVUPS	(DI)(CX*4), X1

	// __auto_tmp_001 = y * y
	MOVAPS	X1, X4
	MULPS	X1, X4

	// out = __auto_tmp_000 + __auto_tmp_001
	ADDPS	X4, X2
	MOVUPD	X2, (BX)(CX*4)
	ADDQ	$4, CX

	// if i >= length { break }
	CMPQ	CX, out+8(FP)
	JGE	NormFloat32s__return
	JMP	NormFloat32s__loop

NormFloat32s__return:
	RET

// func AddUints(out, in1, in2 []uint)
TEXT ·AddUints(SB), 7, $0

	// Load pointers.
	MOVQ	out+0(FP), BX
	MOVQ	in1+24(FP), SI
	MOVQ	in2+48(FP), DI

	// Check lengths.
	MOVQ	out+8(FP), DX
	CMPQ	DX, in1+32(FP)
	JNE	AddUints__panic
	CMPQ	DX, in2+56(FP)
	JNE	AddUints__panic
	JMP	AddUints__ok

AddUints__panic:
	CALL	runtime·panicindex(SB)

AddUints__ok:
	SUBQ	$2, DX
	MOVQ	$0, CX

AddUints__loop:
	CMPQ	CX, DX

	// if i > length-2 { i = length-2 }
	JLE	AddUints__process
	MOVQ	DX, CX

AddUints__process:
	MOVUPS	(SI)(CX*8), X0
	MOVUPS	(DI)(CX*8), X1

	// out = in1 + in2
	PADDQ	X1, X0
	MOVUPD	X0, (BX)(CX*8)
	ADDQ	$2, CX

	// if i >= length { break }
	CMPQ	CX, out+8(FP)
	JGE	AddUints__return
	JMP	AddUints__loop

AddUints__return:
	RET

// func SomeFormula(out, x, y, z []float32)
TEXT ·SomeFormula(SB), 7, $0

	// Load pointers.
	MOVQ	out+0(FP), BX
	MOVQ	x+24(FP), SI
	MOVQ	y+48(FP), DI
	MOVQ	z+72(FP), R8

	// Check lengths.
	MOVQ	out+8(FP), DX
	CMPQ	DX, x+32(FP)
	JNE	SomeFormula__panic
	CMPQ	DX, y+56(FP)
	JNE	SomeFormula__panic
	CMPQ	DX, z+80(FP)
	JNE	SomeFormula__panic
	JMP	SomeFormula__ok

SomeFormula__panic:
	CALL	runtime·panicindex(SB)

SomeFormula__ok:
	SUBQ	$4, DX
	MOVQ	$0, CX

SomeFormula__loop:
	CMPQ	CX, DX

	// if i > length-4 { i = length-4 }
	JLE	SomeFormula__process
	MOVQ	DX, CX

SomeFormula__process:
	MOVUPS	(SI)(CX*4), X0

	// __auto_tmp_000 = x * x
	MOVAPS	X0, X4
	MULPS	X0, X4
	MOVUPS	(DI)(CX*4), X1

	// __auto_tmp_001 = y * y
	MOVAPS	X1, X5
	MULPS	X1, X5

	// __auto_tmp_002 = __auto_tmp_000 + __auto_tmp_001
	ADDPS	X5, X4
	MOVUPS	(R8)(CX*4), X2

	// __auto_tmp_003 = z * z
	MOVAPS	X2, X6
	MULPS	X2, X6

	// __auto_tmp_004 = __auto_tmp_002 + __auto_tmp_003
	ADDPS	X6, X4

	// __auto_tmp_005 = x * y
	MOVAPS	X0, X7
	MULPS	X1, X7

	// __auto_tmp_006 = y * z
	MOVAPS	X1, X8
	MULPS	X2, X8

	// __auto_tmp_007 = __auto_tmp_005 + __auto_tmp_006
	ADDPS	X8, X7

	// __auto_tmp_008 = z * x
	MOVAPS	X2, X9
	MULPS	X0, X9

	// __auto_tmp_009 = __auto_tmp_007 + __auto_tmp_008
	ADDPS	X9, X7

	// out = __auto_tmp_004 - __auto_tmp_009
	SUBPS	X7, X4
	MOVUPD	X4, (BX)(CX*4)
	ADDQ	$4, CX

	// if i >= length { break }
	CMPQ	CX, out+8(FP)
	JGE	SomeFormula__return
	JMP	SomeFormula__loop

SomeFormula__return:
	RET

// func subByte(out, a, b []byte)
TEXT ·subByte(SB), 7, $0

	// Load pointers.
	MOVQ	out+0(FP), BX
	MOVQ	a+24(FP), SI
	MOVQ	b+48(FP), DI

	// Check lengths.
	MOVQ	out+8(FP), DX
	CMPQ	DX, a+32(FP)
	JNE	subByte__panic
	CMPQ	DX, b+56(FP)
	JNE	subByte__panic
	JMP	subByte__ok

subByte__panic:
	CALL	runtime·panicindex(SB)

subByte__ok:
	SUBQ	$16, DX
	MOVQ	$0, CX

subByte__loop:
	CMPQ	CX, DX

	// if i > length-16 { i = length-16 }
	JLE	subByte__process
	MOVQ	DX, CX

subByte__process:
	MOVUPS	(SI)(CX*1), X0
	MOVUPS	(DI)(CX*1), X1

	// out = a - b
	PSUBB	X1, X0
	MOVUPD	X0, (BX)(CX*1)
	ADDQ	$16, CX

	// if i >= length { break }
	CMPQ	CX, out+8(FP)
	JGE	subByte__return
	JMP	subByte__loop

subByte__return:
	RET

// func subuint(out, a, b []uint)
TEXT ·subuint(SB), 7, $0

	// Load pointers.
	MOVQ	out+0(FP), BX
	MOVQ	a+24(FP), SI
	MOVQ	b+48(FP), DI

	// Check lengths.
	MOVQ	out+8(FP), DX
	CMPQ	DX, a+32(FP)
	JNE	subuint__panic
	CMPQ	DX, b+56(FP)
	JNE	subuint__panic
	JMP	subuint__ok

subuint__panic:
	CALL	runtime·panicindex(SB)

subuint__ok:
	SUBQ	$2, DX
	MOVQ	$0, CX

subuint__loop:
	CMPQ	CX, DX

	// if i > length-2 { i = length-2 }
	JLE	subuint__process
	MOVQ	DX, CX

subuint__process:
	MOVUPS	(SI)(CX*8), X0
	MOVUPS	(DI)(CX*8), X1

	// out = a - b
	PSUBQ	X1, X0
	MOVUPD	X0, (BX)(CX*8)
	ADDQ	$2, CX

	// if i >= length { break }
	CMPQ	CX, out+8(FP)
	JGE	subuint__return
	JMP	subuint__loop

subuint__return:
	RET

// func DetF32(det, a, b, c, d []float32)
TEXT ·DetF32(SB), 7, $0

	// Load pointers.
	MOVQ	det+0(FP), BX
	MOVQ	a+24(FP), SI
	MOVQ	b+48(FP), DI
	MOVQ	c+72(FP), R8
	MOVQ	d+96(FP), R9

	// Check lengths.
	MOVQ	det+8(FP), DX
	CMPQ	DX, a+32(FP)
	JNE	DetF32__panic
	CMPQ	DX, b+56(FP)
	JNE	DetF32__panic
	CMPQ	DX, c+80(FP)
	JNE	DetF32__panic
	CMPQ	DX, d+104(FP)
	JNE	DetF32__panic
	JMP	DetF32__ok

DetF32__panic:
	CALL	runtime·panicindex(SB)

DetF32__ok:
	SUBQ	$4, DX
	MOVQ	$0, CX

DetF32__loop:
	CMPQ	CX, DX

	// if i > length-4 { i = length-4 }
	JLE	DetF32__process
	MOVQ	DX, CX

DetF32__process:
	MOVUPS	(SI)(CX*4), X0
	MOVUPS	(R9)(CX*4), X1

	// __auto_tmp_000 = a * d
	MULPS	X1, X0
	MOVUPS	(DI)(CX*4), X2
	MOVUPS	(R8)(CX*4), X4

	// __auto_tmp_001 = b * c
	MULPS	X4, X2

	// det = __auto_tmp_000 - __auto_tmp_001
	SUBPS	X2, X0
	MOVUPD	X0, (BX)(CX*4)
	ADDQ	$4, CX

	// if i >= length { break }
	CMPQ	CX, det+8(FP)
	JGE	DetF32__return
	JMP	DetF32__loop

DetF32__return:
	RET

