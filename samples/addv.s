// addv(out, in1, in2 []int)
TEXT ·AddV(SB), 7, $0
    // Pointers.
    MOVQ    out+0(FP),R9
    MOVQ    in1+16(FP),R10
    MOVQ    in2+32(FP),R11
    // Length.
    MOVL    out+8(FP),R12
    MOVL    in1+24(FP),R13
    MOVL    in2+40(FP),R14

    // check lengths are equal.
    CMPL    R12, R13
    JE      ok1
    CALL    runtime·panicindex(SB)
ok1:
    CMPL    R12, R14
    JE      ok2
    CALL    runtime·panicindex(SB)

    // start the addition.
ok2:
    MOVL     $0, CX
loop:
    CMPL     CX, R12
    JGE      return
    // load in1 and in2
    MOVUPS   (R10)(CX*4), X10
    PADDL    (R11)(CX*4), X10
    MOVUPS   X10, (R9)(CX*4)

    // advance by 4 doublewords.
    ADDL     $4, CX
    JMP      loop

return:
    RET
