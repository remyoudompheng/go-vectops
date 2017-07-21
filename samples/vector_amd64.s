// func NormFloat32s(out, x, y []float32)
TEXT ·NormFloat32s(SB), 7, $0

	// Load pointers.
	MOVQ	out+0(FP), BX
	MOVQ	x+24(FP), SI
	MOVQ	y+48(FP), DI

	// Check lengths.
	MOVL	out+8(FP), DX
	CMPL	DX, x+32(FP)
	JNE	NormFloat32s__panic
	CMPL	DX, y+56(FP)
	JNE	NormFloat32s__panic
	JMP	NormFloat32s__ok

NormFloat32s__panic:
	CALL	runtime·panicindex(SB)

NormFloat32s__ok:
	SUBL	$4, DX
	XORL	CX, CX

NormFloat32s__loop:
	CMPL	CX, DX

	// if i > length-4 { i = length-4 }
	JLE	NormFloat32s__process
	MOVL	DX, CX

NormFloat32s__process:
	MOVUPS	(SI)(CX*4), X0

	// __auto_tmp_000 = x * x
	MOVAPS	X0, X2
	MULPS	X0, X2
	MOVUPS	(DI)(CX*4), X1

	// __auto_tmp_001 = y * y
	MOVAPS	X1, X4
	MULPS	X1, X4

	// __auto_tmp_002 = __auto_tmp_000 + __auto_tmp_001
	ADDPS	X4, X2
	MOVUPD	X2, (BX)(CX*4)
	ADDL	$4, CX

	// if i >= length { break }
	CMPL	CX, out+16(FP)
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
	MOVL	out+8(FP), DX
	CMPL	DX, in1+32(FP)
	JNE	AddUints__panic
	CMPL	DX, in2+56(FP)
	JNE	AddUints__panic
	JMP	AddUints__ok

AddUints__panic:
	CALL	runtime·panicindex(SB)

AddUints__ok:
	SUBL	$2, DX
	XORL	CX, CX

AddUints__loop:
	CMPL	CX, DX

	// if i > length-2 { i = length-2 }
	JLE	AddUints__process
	MOVL	DX, CX

AddUints__process:
	MOVUPS	(SI)(CX*8), X0
	MOVUPS	(DI)(CX*8), X1

	// __auto_tmp_000 = in1 + in2
	PADDQ	X1, X0
	MOVUPD	X0, (BX)(CX*8)
	ADDL	$2, CX

	// if i >= length { break }
	CMPL	CX, out+16(FP)
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
	MOVL	out+8(FP), DX
	CMPL	DX, x+32(FP)
	JNE	SomeFormula__panic
	CMPL	DX, y+56(FP)
	JNE	SomeFormula__panic
	CMPL	DX, z+80(FP)
	JNE	SomeFormula__panic
	JMP	SomeFormula__ok

SomeFormula__panic:
	CALL	runtime·panicindex(SB)

SomeFormula__ok:
	SUBL	$4, DX
	XORL	CX, CX

SomeFormula__loop:
	CMPL	CX, DX

	// if i > length-4 { i = length-4 }
	JLE	SomeFormula__process
	MOVL	DX, CX

SomeFormula__process:
	MOVUPS	(SI)(CX*4), X0
	MOVUPS	(DI)(CX*4), X1

	// __auto_tmp_000 = x * y
	MOVAPS	X0, X4
	MULPS	X1, X4
	MOVUPS	(R8)(CX*4), X2

	// __auto_tmp_001 = y * z
	MOVAPS	X1, X5
	MULPS	X2, X5

	// __auto_tmp_002 = __auto_tmp_000 + __auto_tmp_001
	ADDPS	X5, X4

	// __auto_tmp_003 = z * x
	MOVAPS	X2, X6
	MULPS	X0, X6

	// __auto_tmp_004 = __auto_tmp_002 + __auto_tmp_003
	ADDPS	X6, X4

	// __auto_tmp_005 = x * x
	MOVAPS	X0, X7
	MULPS	X0, X7

	// __auto_tmp_006 = y * y
	MOVAPS	X1, X8
	MULPS	X1, X8

	// __auto_tmp_007 = __auto_tmp_005 + __auto_tmp_006
	ADDPS	X8, X7

	// __auto_tmp_008 = z * z
	MOVAPS	X2, X9
	MULPS	X2, X9

	// __auto_tmp_009 = __auto_tmp_007 + __auto_tmp_008
	ADDPS	X9, X7

	// __auto_tmp_010 = __auto_tmp_004 / __auto_tmp_009
	DIVPS	X7, X4
	MOVUPD	X4, (BX)(CX*4)
	ADDL	$4, CX

	// if i >= length { break }
	CMPL	CX, out+16(FP)
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
	MOVL	out+8(FP), DX
	CMPL	DX, a+32(FP)
	JNE	subByte__panic
	CMPL	DX, b+56(FP)
	JNE	subByte__panic
	JMP	subByte__ok

subByte__panic:
	CALL	runtime·panicindex(SB)

subByte__ok:
	SUBL	$16, DX
	XORL	CX, CX

subByte__loop:
	CMPL	CX, DX

	// if i > length-16 { i = length-16 }
	JLE	subByte__process
	MOVL	DX, CX

subByte__process:
	MOVUPS	(SI)(CX*1), X0
	MOVUPS	(DI)(CX*1), X1

	// __auto_tmp_000 = a - b
	PSUBB	X1, X0
	MOVUPD	X0, (BX)(CX*1)
	ADDL	$16, CX

	// if i >= length { break }
	CMPL	CX, out+16(FP)
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
	MOVL	out+8(FP), DX
	CMPL	DX, a+32(FP)
	JNE	subuint__panic
	CMPL	DX, b+56(FP)
	JNE	subuint__panic
	JMP	subuint__ok

subuint__panic:
	CALL	runtime·panicindex(SB)

subuint__ok:
	SUBL	$2, DX
	XORL	CX, CX

subuint__loop:
	CMPL	CX, DX

	// if i > length-2 { i = length-2 }
	JLE	subuint__process
	MOVL	DX, CX

subuint__process:
	MOVUPS	(SI)(CX*8), X0
	MOVUPS	(DI)(CX*8), X1

	// __auto_tmp_000 = a - b
	PSUBQ	X1, X0
	MOVUPD	X0, (BX)(CX*8)
	ADDL	$2, CX

	// if i >= length { break }
	CMPL	CX, out+16(FP)
	JGE	subuint__return
	JMP	subuint__loop

subuint__return:
	RET

// func DetF64(det, a, b, c, d []float64)
TEXT ·DetF64(SB), 7, $0

	// Load pointers.
	MOVQ	det+0(FP), BX
	MOVQ	a+24(FP), SI
	MOVQ	b+48(FP), DI
	MOVQ	c+72(FP), R8
	MOVQ	d+96(FP), R9

	// Check lengths.
	MOVL	det+8(FP), DX
	CMPL	DX, a+32(FP)
	JNE	DetF64__panic
	CMPL	DX, b+56(FP)
	JNE	DetF64__panic
	CMPL	DX, c+80(FP)
	JNE	DetF64__panic
	CMPL	DX, d+104(FP)
	JNE	DetF64__panic
	JMP	DetF64__ok

DetF64__panic:
	CALL	runtime·panicindex(SB)

DetF64__ok:
	SUBL	$2, DX
	XORL	CX, CX

DetF64__loop:
	CMPL	CX, DX

	// if i > length-2 { i = length-2 }
	JLE	DetF64__process
	MOVL	DX, CX

DetF64__process:
	MOVUPS	(SI)(CX*8), X0
	MOVUPS	(R9)(CX*8), X1

	// __auto_tmp_000 = a * d
	MULPD	X1, X0
	MOVUPS	(DI)(CX*8), X2
	MOVUPS	(R8)(CX*8), X4

	// __auto_tmp_001 = b * c
	MULPD	X4, X2

	// __auto_tmp_002 = __auto_tmp_000 - __auto_tmp_001
	SUBPD	X2, X0
	MOVUPD	X0, (BX)(CX*8)
	ADDL	$2, CX

	// if i >= length { break }
	CMPL	CX, det+16(FP)
	JGE	DetF64__return
	JMP	DetF64__loop

DetF64__return:
	RET

