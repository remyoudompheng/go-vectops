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
	WORD	$0xec900b04	// VLDMIA (R0), Q0

	// __auto_tmp_000 = x * x
	WORD	$0xf3000d50	// VMUL.F32 Q0, Q0, Q0
	ADD	R11<<2, R3, R0
	WORD	$0xec902b04	// VLDMIA (R0), Q1

	// __auto_tmp_001 = y * y
	WORD	$0xf3022d52	// VMUL.F32 Q1, Q1, Q1

	// out = __auto_tmp_000 + __auto_tmp_001
	WORD	$0xf2022d40	// VADD.F32 Q1, Q0, Q1
	ADD	R11<<2, R1, R0
	WORD	$0xec802b04	// VSTMIA (R0), Q1
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
	WORD	$0xec900b04	// VLDMIA (R0), Q0
	ADD	R11<<2, R3, R0
	WORD	$0xec902b04	// VLDMIA (R0), Q1

	// out = in1 + in2
	WORD	$0xf2222840	// VADD.I32 Q1, Q0, Q1
	ADD	R11<<2, R1, R0
	WORD	$0xec802b04	// VSTMIA (R0), Q1
	ADD	$4, R11

	// if i >= length { break }
	MOVW	out+4(FP), R0
	CMP	R0, R11
	BGE	AddUints__return
	B	AddUints__loop

AddUints__return:
	RET

// func SomeFormula(out, x, y, z []float32)
TEXT ·SomeFormula(SB), 7, $0

	// Load pointers.
	MOVW	out+0(FP), R1
	MOVW	x+12(FP), R2
	MOVW	y+24(FP), R3
	MOVW	z+36(FP), R4

	// Check lengths.
	MOVW	out+4(FP), R12
	MOVW	x+16(FP), R0
	TEQ	R0, R12
	BNE	SomeFormula__panic
	MOVW	y+28(FP), R0
	TEQ	R0, R12
	BNE	SomeFormula__panic
	MOVW	z+40(FP), R0
	TEQ	R0, R12
	BNE	SomeFormula__panic
	B	SomeFormula__ok

SomeFormula__panic:
	B	runtime·panicindex(SB)

SomeFormula__ok:
	RSB	$4, R12
	MOVW	$0, R11

SomeFormula__loop:
	CMP	R11, R12

	// if i > length-4 { i = length-4 }
	BLE	SomeFormula__process
	MOVW	R12, R11

SomeFormula__process:
	ADD	R11<<2, R2, R0
	WORD	$0xec900b04	// VLDMIA (R0), Q0

	// __auto_tmp_000 = x * x
	WORD	$0xf3002d50	// VMUL.F32 Q0, Q0, Q1
	ADD	R11<<2, R3, R0
	WORD	$0xec904b04	// VLDMIA (R0), Q2

	// __auto_tmp_001 = y * y
	WORD	$0xf3046d54	// VMUL.F32 Q2, Q2, Q3

	// __auto_tmp_002 = __auto_tmp_000 + __auto_tmp_001
	WORD	$0xf2066d42	// VADD.F32 Q3, Q1, Q3
	ADD	R11<<2, R4, R0
	WORD	$0xec902b04	// VLDMIA (R0), Q1

	// __auto_tmp_003 = z * z
	WORD	$0xf3028d52	// VMUL.F32 Q1, Q1, Q4

	// __auto_tmp_004 = __auto_tmp_002 + __auto_tmp_003
	WORD	$0xf2088d46	// VADD.F32 Q4, Q3, Q4

	// __auto_tmp_005 = x * y
	WORD	$0xf3006d54	// VMUL.F32 Q0, Q2, Q3

	// __auto_tmp_006 = y * z
	WORD	$0xf3044d52	// VMUL.F32 Q2, Q1, Q2

	// __auto_tmp_007 = __auto_tmp_005 + __auto_tmp_006
	WORD	$0xf2044d46	// VADD.F32 Q2, Q3, Q2

	// __auto_tmp_008 = z * x
	WORD	$0xf3000d52	// VMUL.F32 Q0, Q1, Q0

	// __auto_tmp_009 = __auto_tmp_007 + __auto_tmp_008
	WORD	$0xf2000d44	// VADD.F32 Q0, Q2, Q0

	// out = __auto_tmp_004 - __auto_tmp_009
	WORD	$0xf2282d40	// VSUB.F32 Q4, Q0, Q1
	ADD	R11<<2, R1, R0
	WORD	$0xec802b04	// VSTMIA (R0), Q1
	ADD	$4, R11

	// if i >= length { break }
	MOVW	out+4(FP), R0
	CMP	R0, R11
	BGE	SomeFormula__return
	B	SomeFormula__loop

SomeFormula__return:
	RET

// func subByte(out, a, b []byte)
TEXT ·subByte(SB), 7, $0

	// Load pointers.
	MOVW	out+0(FP), R1
	MOVW	a+12(FP), R2
	MOVW	b+24(FP), R3

	// Check lengths.
	MOVW	out+4(FP), R12
	MOVW	a+16(FP), R0
	TEQ	R0, R12
	BNE	subByte__panic
	MOVW	b+28(FP), R0
	TEQ	R0, R12
	BNE	subByte__panic
	B	subByte__ok

subByte__panic:
	B	runtime·panicindex(SB)

subByte__ok:
	RSB	$16, R12
	MOVW	$0, R11

subByte__loop:
	CMP	R11, R12

	// if i > length-16 { i = length-16 }
	BLE	subByte__process
	MOVW	R12, R11

subByte__process:
	ADD	R11<<0, R2, R0
	WORD	$0xec900b04	// VLDMIA (R0), Q0
	ADD	R11<<0, R3, R0
	WORD	$0xec902b04	// VLDMIA (R0), Q1

	// out = a - b
	WORD	$0xf3000842	// VSUB.I8 Q0, Q1, Q0
	ADD	R11<<0, R1, R0
	WORD	$0xec800b04	// VSTMIA (R0), Q0
	ADD	$16, R11

	// if i >= length { break }
	MOVW	out+4(FP), R0
	CMP	R0, R11
	BGE	subByte__return
	B	subByte__loop

subByte__return:
	RET

// func subuint(out, a, b []uint)
TEXT ·subuint(SB), 7, $0

	// Load pointers.
	MOVW	out+0(FP), R1
	MOVW	a+12(FP), R2
	MOVW	b+24(FP), R3

	// Check lengths.
	MOVW	out+4(FP), R12
	MOVW	a+16(FP), R0
	TEQ	R0, R12
	BNE	subuint__panic
	MOVW	b+28(FP), R0
	TEQ	R0, R12
	BNE	subuint__panic
	B	subuint__ok

subuint__panic:
	B	runtime·panicindex(SB)

subuint__ok:
	RSB	$4, R12
	MOVW	$0, R11

subuint__loop:
	CMP	R11, R12

	// if i > length-4 { i = length-4 }
	BLE	subuint__process
	MOVW	R12, R11

subuint__process:
	ADD	R11<<2, R2, R0
	WORD	$0xec900b04	// VLDMIA (R0), Q0
	ADD	R11<<2, R3, R0
	WORD	$0xec902b04	// VLDMIA (R0), Q1

	// out = a - b
	WORD	$0xf3200842	// VSUB.I32 Q0, Q1, Q0
	ADD	R11<<2, R1, R0
	WORD	$0xec800b04	// VSTMIA (R0), Q0
	ADD	$4, R11

	// if i >= length { break }
	MOVW	out+4(FP), R0
	CMP	R0, R11
	BGE	subuint__return
	B	subuint__loop

subuint__return:
	RET

// func DetF32(det, a, b, c, d []float32)
TEXT ·DetF32(SB), 7, $0

	// Load pointers.
	MOVW	det+0(FP), R1
	MOVW	a+12(FP), R2
	MOVW	b+24(FP), R3
	MOVW	c+36(FP), R4
	MOVW	d+48(FP), R5

	// Check lengths.
	MOVW	det+4(FP), R12
	MOVW	a+16(FP), R0
	TEQ	R0, R12
	BNE	DetF32__panic
	MOVW	b+28(FP), R0
	TEQ	R0, R12
	BNE	DetF32__panic
	MOVW	c+40(FP), R0
	TEQ	R0, R12
	BNE	DetF32__panic
	MOVW	d+52(FP), R0
	TEQ	R0, R12
	BNE	DetF32__panic
	B	DetF32__ok

DetF32__panic:
	B	runtime·panicindex(SB)

DetF32__ok:
	RSB	$4, R12
	MOVW	$0, R11

DetF32__loop:
	CMP	R11, R12

	// if i > length-4 { i = length-4 }
	BLE	DetF32__process
	MOVW	R12, R11

DetF32__process:
	ADD	R11<<2, R2, R0
	WORD	$0xec900b04	// VLDMIA (R0), Q0
	ADD	R11<<2, R5, R0
	WORD	$0xec902b04	// VLDMIA (R0), Q1

	// __auto_tmp_000 = a * d
	WORD	$0xf3022d50	// VMUL.F32 Q1, Q0, Q1
	ADD	R11<<2, R3, R0
	WORD	$0xec900b04	// VLDMIA (R0), Q0
	ADD	R11<<2, R4, R0
	WORD	$0xec904b04	// VLDMIA (R0), Q2

	// __auto_tmp_001 = b * c
	WORD	$0xf3044d50	// VMUL.F32 Q2, Q0, Q2

	// det = __auto_tmp_000 - __auto_tmp_001
	WORD	$0xf2220d44	// VSUB.F32 Q1, Q2, Q0
	ADD	R11<<2, R1, R0
	WORD	$0xec800b04	// VSTMIA (R0), Q0
	ADD	$4, R11

	// if i >= length { break }
	MOVW	det+4(FP), R0
	CMP	R0, R11
	BGE	DetF32__return
	B	DetF32__loop

DetF32__return:
	RET

