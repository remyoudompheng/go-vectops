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

ok2:
    // last index is len(out)-4
    SUBL     $4, R12
    // start the addition.
    MOVL     $0, CX    // i = 0
loop:
    CMPL     CX, R12
    JLE      process   // if i > len(out)-4 { i = len(out)-4 }
    MOVL     R12, CX
process:
    MOVUPS   (R10)(CX*4), X10 // x = in1[i:i+4]
    MOVUPS   (R11)(CX*4), X11
    PADDL    X11, X10         // x += in2[i:i+4]
    MOVUPS   X10, (R9)(CX*4)  // out[i:i+4] = x

    ADDL     $4, CX   // i += 4
    CMPL     CX, out+8(FP)
    JGE      return
    JMP      loop

return:
    RET
