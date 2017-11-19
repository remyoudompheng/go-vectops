// func NormFloat32s(out, x, y []float32)
TEXT ·NormFloat32s(SB), 7, $0

	// Load pointers.
	MOVW	out+0(FP), R1
	MOVW	x+12(FP), R2
	MOVW	y+24(FP), R3

	// Check lengths.
	MOVW	out+4(FP), R12
	MOVW	x+16(FP), R0
	TEQ	R0, R12
	BNE	NormFloat32s__panic
	MOVW	y+28(FP), R0
	TEQ	R0, R12
	BNE	NormFloat32s__panic
	B	NormFloat32s__ok

NormFloat32s__panic:
	B	runtime·panicindex(SB)

NormFloat32s__ok:
	RSB	$4, R12
	MOVW	$0, R11

NormFloat32s__loop:
	CMP	R11, R12

	// if i > length-4 { i = length-4 }
	BLE	NormFloat32s__process
	MOVW	R12, R11

NormFloat32s__process:
	ADD	R11<<2, R2, R0
	WORD	$0xec900b04	// VLDM (R0), Q0

	// __auto_tmp_000 = x * x
	WORD	$0xf3004d50	// VMUL.F32 Q0, Q0, Q2
	ADD	R11<<2, R3, R0
	WORD	$0xec902b04	// VLDM (R0), Q1

	// __auto_tmp_001 = y * y
	WORD	$0xf3026d52	// VMUL.F32 Q1, Q1, Q3

	// out = __auto_tmp_000 + __auto_tmp_001
	WORD	$0xf2044d46	// VADD.F32 Q2, Q3, Q2
	ADD	R11<<2, R1, R0
	WORD	$0xec804b04	// VSTM (R0), Q2
	ADD	$4, R11

	// if i >= length { break }
	MOVW	out+4(FP), R0
	CMP	R0, R11
	BGE	NormFloat32s__return
	B	NormFloat32s__loop

NormFloat32s__return:
	RET

// func AddUints(out, in1, in2 []uint)
TEXT ·AddUints(SB), 7, $0

	// Load pointers.
	MOVW	out+0(FP), R1
	MOVW	in1+12(FP), R2
	MOVW	in2+24(FP), R3

	// Check lengths.
	MOVW	out+4(FP), R12
	MOVW	in1+16(FP), R0
	TEQ	R0, R12
	BNE	AddUints__panic
	MOVW	in2+28(FP), R0
	TEQ	R0, R12
	BNE	AddUints__panic
	B	AddUints__ok

AddUints__panic:
	B	runtime·panicindex(SB)

AddUints__ok:
	RSB	$4, R12
	MOVW	$0, R11

AddUints__loop:
	CMP	R11, R12

	// if i > length-4 { i = length-4 }
	BLE	AddUints__process
	MOVW	R12, R11

AddUints__process:
	ADD	R11<<2, R2, R0
	WORD	$0xec900b04	// VLDM (R0), Q0
	ADD	R11<<2, R3, R0
	WORD	$0xec902b04	// VLDM (R0), Q1

	// out = in1 + in2
	WORD	$0xf2200842	// VADD.I32 Q0, Q1, Q0
	ADD	R11<<2, R1, R0
	WORD	$0xec800b04	// VSTM (R0), Q0
	ADD	$4, R11

	// if i >= length { break }
	MOVW	out+4(FP), R0
	CMP	R0, R11
	BGE	AddUints__return
	B	AddUints__loop

AddUints__return:
	RET

