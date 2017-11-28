// +build avx2

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
	SUBQ	$8, DX
	MOVQ	$0, CX

NormFloat32s__loop:
	CMPQ	CX, DX

	// if i > length-8 { i = length-8 }
	JLE	NormFloat32s__process
	MOVQ	DX, CX

NormFloat32s__process:
	VMOVDQU	(SI)(CX*4), Y0

	// __auto_tmp_000 = x * x
	VMULPS	Y0, Y0, Y0
	VMOVDQU	(DI)(CX*4), Y1

	// __auto_tmp_001 = y * y
	VMULPS	Y1, Y1, Y1

	// out = __auto_tmp_000 + __auto_tmp_001
	VADDPS	Y0, Y1, Y1
	VMOVDQU	Y1, (BX)(CX*4)
	ADDQ	$8, CX

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
	SUBQ	$4, DX
	MOVQ	$0, CX

AddUints__loop:
	CMPQ	CX, DX

	// if i > length-4 { i = length-4 }
	JLE	AddUints__process
	MOVQ	DX, CX

AddUints__process:
	VMOVDQU	(SI)(CX*8), Y0
	VMOVDQU	(DI)(CX*8), Y1

	// out = in1 + in2
	VPADDQ	Y0, Y1, Y1
	VMOVDQU	Y1, (BX)(CX*8)
	ADDQ	$4, CX

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
	SUBQ	$8, DX
	MOVQ	$0, CX

SomeFormula__loop:
	CMPQ	CX, DX

	// if i > length-8 { i = length-8 }
	JLE	SomeFormula__process
	MOVQ	DX, CX

SomeFormula__process:
	VMOVDQU	(SI)(CX*4), Y0

	// __auto_tmp_000 = x * x
	VMULPS	Y0, Y0, Y1
	VMOVDQU	(DI)(CX*4), Y2

	// __auto_tmp_001 = y * y
	VMULPS	Y2, Y2, Y4

	// __auto_tmp_002 = __auto_tmp_000 + __auto_tmp_001
	VADDPS	Y1, Y4, Y4
	VMOVDQU	(R8)(CX*4), Y1

	// __auto_tmp_003 = z * z
	VMULPS	Y1, Y1, Y5

	// __auto_tmp_004 = __auto_tmp_002 + __auto_tmp_003
	VADDPS	Y4, Y5, Y5

	// __auto_tmp_005 = x * y
	VMULPS	Y2, Y0, Y4

	// __auto_tmp_006 = y * z
	VMULPS	Y1, Y2, Y2

	// __auto_tmp_007 = __auto_tmp_005 + __auto_tmp_006
	VADDPS	Y4, Y2, Y2

	// __auto_tmp_008 = z * x
	VMULPS	Y1, Y0, Y0

	// __auto_tmp_009 = __auto_tmp_007 + __auto_tmp_008
	VADDPS	Y2, Y0, Y0

	// out = __auto_tmp_004 - __auto_tmp_009
	VSUBPS	Y0, Y5, Y1
	VMOVDQU	Y1, (BX)(CX*4)
	ADDQ	$8, CX

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
	SUBQ	$32, DX
	MOVQ	$0, CX

subByte__loop:
	CMPQ	CX, DX

	// if i > length-32 { i = length-32 }
	JLE	subByte__process
	MOVQ	DX, CX

subByte__process:
	VMOVDQU	(SI)(CX*1), Y0
	VMOVDQU	(DI)(CX*1), Y1

	// out = a - b
	VPSUBB	Y1, Y0, Y0
	VMOVDQU	Y0, (BX)(CX*1)
	ADDQ	$32, CX

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
	SUBQ	$4, DX
	MOVQ	$0, CX

subuint__loop:
	CMPQ	CX, DX

	// if i > length-4 { i = length-4 }
	JLE	subuint__process
	MOVQ	DX, CX

subuint__process:
	VMOVDQU	(SI)(CX*8), Y0
	VMOVDQU	(DI)(CX*8), Y1

	// out = a - b
	VPSUBQ	Y1, Y0, Y0
	VMOVDQU	Y0, (BX)(CX*8)
	ADDQ	$4, CX

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
	SUBQ	$8, DX
	MOVQ	$0, CX

DetF32__loop:
	CMPQ	CX, DX

	// if i > length-8 { i = length-8 }
	JLE	DetF32__process
	MOVQ	DX, CX

DetF32__process:
	VMOVDQU	(SI)(CX*4), Y0
	VMOVDQU	(R9)(CX*4), Y1

	// __auto_tmp_000 = a * d
	VMULPS	Y0, Y1, Y1
	VMOVDQU	(DI)(CX*4), Y0
	VMOVDQU	(R8)(CX*4), Y2

	// __auto_tmp_001 = b * c
	VMULPS	Y0, Y2, Y2

	// det = __auto_tmp_000 - __auto_tmp_001
	VSUBPS	Y2, Y1, Y0
	VMOVDQU	Y0, (BX)(CX*4)
	ADDQ	$8, CX

	// if i >= length { break }
	CMPQ	CX, det+8(FP)
	JGE	DetF32__return
	JMP	DetF32__loop

DetF32__return:
	RET

