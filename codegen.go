package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type codeWriter struct {
	w         io.Writer
	goarch    string
	gosubarch string
	arch      Arch
}

func (w codeWriter) CodeGen(f *Function) error {
	fmt.Fprintf(w.w, "// %s\n", f.ForwardDecl())
	fmt.Fprintf(w.w, "TEXT ·%s(SB), 7, $0\n", f.Name)
	err := f.Compile(w)
	if err != nil {
		return err
	}
	fmt.Fprintln(w.w, "")
	return nil
}

func frameArg(name string, offset int) string {
	return fmt.Sprintf("%s+%d(FP)", name, offset)
}

func (w codeWriter) preamble(f *Function) {
	var mov, cmp, jne, jmp, call string
	switch w.goarch {
	case "amd64":
		mov, cmp = "MOVQ", "CMPQ"
		jne, jmp, call = "JNE", "JMP", "CALL"
	case "arm":
		mov, cmp = "MOVW", "TEQ"
		jne, jmp, call = "BNE", "B", "B"
	default:
		panic("impossible")
	}

	ptrSize := w.arch.PtrSize
	w.comment("Load pointers.")
	outArg := f.Args[0]
	inArgs := f.Args[1:]
	if len(f.Args) > len(w.arch.InputRegs) {
		panic("not enough registers")
	}
	for i, arg := range f.Args {
		w.opcode(mov,
			frameArg(arg, ptrSize*(3*i)),
			w.arch.InputRegs[i])
	}

	// length and index.
	w.comment("Check lengths.")
	w.opcode(mov, frameArg(outArg, ptrSize), w.arch.LengthReg)
	for i, arg := range inArgs {
		if w.goarch == "arm" {
			const spareReg = "R0"
			w.opcode(mov, frameArg(arg, ptrSize*(3*i+3+1)), spareReg)
			w.opcode(cmp, spareReg, w.arch.LengthReg)
		} else {
			w.opcode(cmp, w.arch.LengthReg,
				frameArg(arg, ptrSize*(3*i+3+1)))
		}
		w.opcode(jne, f.Name+"__panic")
	}
	w.opcode(jmp, f.Name+"__ok")
	//
	w.label(f.Name, "panic")
	w.opcode(call, "runtime·panicindex(SB)")
	w.label(f.Name, "ok")
}

func (f *Function) Compile(w codeWriter) error {
	outArg := f.Args[0]
	ptrSize := w.arch.PtrSize

	instrs, err := Compile(f, w)
	if err != nil {
		return err
	}
	for _, ins := range instrs {
		if err := w.checkInstr(ins); err != nil {
			return err
		}
	}

	w.preamble(f)
	stride := w.arch.VectorWidth / w.arch.Width(f.ScalarType)
	// Emit code for the loop. It should look like:
	// for i := 0; ; {
	// 	if i > length-4 { i = length-4 }
	// 	process(arrays[i:i+4])
	// 	i += 4
	// 	if i >= length { break }
	// }
	regc, regl := w.arch.CounterReg, w.arch.LengthReg
	var add, sub, mov, cmp, jle, jge, jmp string
	switch w.goarch {
	case "amd64":
		add, sub = "ADDQ", "SUBQ"
		mov, cmp = "MOVQ", "CMPQ"
		jle, jge, jmp = "JLE", "JGE", "JMP"
	case "arm":
		add, sub = "ADD", "RSB"
		mov, cmp = "MOVW", "CMP" // will test LE or GE
		jle, jge, jmp = "BLE", "BGE", "B"
	default:
		panic("impossible")
	}

	w.opcode(sub, fmt.Sprintf("$%d", stride), regl) // regl = length - 4
	w.opcode(mov, "$0", regc)
	w.label(f.Name, "loop")
	w.opcode(cmp, regc, regl)
	w.comment("if i > length-%d { i = length-%d }", stride, stride)
	w.opcode(jle, f.Name+"__process")
	w.opcode(mov, regl, regc)
	w.label(f.Name, "process")

	for _, ins := range instrs {
		w.emitInstr(ins)
	}

	w.opcode(add, fmt.Sprintf("$%d", stride), regc)
	w.comment("if i >= length { break }")
	if w.goarch == "arm" {
		w.opcode(mov, frameArg(outArg, ptrSize), "R0")
		w.opcode(cmp, "R0", regc)
	} else {
		w.opcode(cmp, regc, frameArg(outArg, ptrSize))
	}
	w.opcode(jge, f.Name+"__return")
	w.opcode(jmp, f.Name+"__loop")
	w.label(f.Name, "return")
	w.opcode("RET")
	return nil
}

func (w codeWriter) checkInstr(ins Instr) error {
	if w.arch.Width(ins.Var.Type) == 0 {
		return fmt.Errorf("unsupported type %s on architecture", ins.Var.Type)
	}
	if ins.Kind == OP {
		_, ok := w.arch.Opcode(ins.Op, ins.Var.Type)
		if !ok {
			return fmt.Errorf("no instruction for %s%s%s",
				ins.Var.Type, ins.Op, ins.Var.Type)
		}
	}
	return nil
}

func (w codeWriter) emitInstr(ins Instr) {
	switch w.goarch {
	case "amd64":
		switch ins.Kind {
		case LOAD:
			width := w.arch.Width(ins.Var.Type)
			loc := fmt.Sprintf("(%s)(%s*%d)", ins.Var.AddrReg,
				w.arch.CounterReg, width)
			if w.gosubarch == "avx2" {
				w.opcode("VMOVDQU", loc, ins.Var.Location)
			} else {
				w.opcode("MOVUPS", loc, ins.Var.Location)
			}
		case STORE:
			width := w.arch.Width(ins.Var.Type)
			loc := fmt.Sprintf("(%s)(%s*%d)", ins.Var.AddrReg,
				w.arch.CounterReg, width)
			if w.gosubarch == "avx2" {
				w.opcode("VMOVDQU", ins.Var.Location, loc)
			} else {
				w.opcode("MOVUPD", ins.Var.Location, loc)
			}
		case OP:
			v := ins.Var
			w.comment("%s = %s", v.Name, v.Expr())
			opcode, ok := w.arch.Opcode(ins.Op, v.Type)
			if !ok {
				panic("unsupported operation")
			}
			if w.gosubarch == "avx2" {
				// use ternary form
				w.opcode(opcode, ins.Right.Location, ins.Left.Location, ins.Var.Location)
			} else {
				if ins.Var.Location == ins.Left.Location {
					w.opcode(opcode, ins.Right.Location, ins.Left.Location)
				} else {
					w.opcode("MOVAPS", ins.Left.Location, ins.Var.Location)
					w.opcode(opcode, ins.Right.Location, ins.Var.Location)
				}
			}
		}
	case "arm":
		const spareReg = "R0"
		switch ins.Kind {
		case LOAD:
			logwidth := w.arch.LogWidth(ins.Var.Type)
			offset := fmt.Sprintf("%s<<%d", w.arch.CounterReg, logwidth)
			w.opcode("ADD", offset, ins.Var.AddrReg, spareReg)
			if strings.HasPrefix(ins.Var.Location, "Q") {
				regd := regNr(ins.Var.Location) // encoded on bits 22, 15-12
				regs := regNr(spareReg)
				enc := assembleNEON("VLDM", false, regs, 4, 2*regd) // 4 = four words
				fmt.Fprintf(w.w, "\tWORD\t$0x%08x\t// %s %s, %s\n",
					enc, "VLDMIA", "("+spareReg+")", ins.Var.Location)
			} else {
				w.opcode("MOVD", "("+spareReg+")", ins.Var.Location)
			}
		case STORE:
			logwidth := w.arch.LogWidth(ins.Var.Type)
			offset := fmt.Sprintf("%s<<%d", w.arch.CounterReg, logwidth)
			w.opcode("ADD", offset, ins.Var.AddrReg, spareReg)
			if strings.HasPrefix(ins.Var.Location, "Q") {
				regd := regNr(ins.Var.Location) // encoded on bits 22, 15-12
				regs := regNr(spareReg)
				enc := assembleNEON("VSTM", false, regs, 4, 2*regd) // 4 = four words
				fmt.Fprintf(w.w, "\tWORD\t$0x%08x\t// %s %s, %s\n",
					enc, "VSTMIA", "("+spareReg+")", ins.Var.Location)
			} else {
				w.opcode("MOVD", ins.Var.Location, "("+spareReg+")")
			}
		case OP:
			v := ins.Var
			w.comment("%s = %s", v.Name, v.Expr())
			opcode, ok := w.arch.Opcode(ins.Op, v.Type)
			if !ok {
				panic("unsupported operation")
			}
			w.opcodeNEON(opcode, ins.Left.Location, ins.Right.Location, ins.Var.Location)
		}
	default:
		panic("not implemented")
	}
}

func (w codeWriter) comment(format string, args ...interface{}) {
	fmt.Fprintf(w.w, "\n\t// "+format+"\n", args...)
}

func (w codeWriter) opcode(op string, args ...string) {
	if len(args) > 0 {
		fmt.Fprintf(w.w, "\t%s\t%s\n", op, strings.Join(args, ", "))
	} else {
		fmt.Fprintf(w.w, "\t%s\n", op)
	}
}

func (w codeWriter) opcodeNEON(op string, reg1, reg2, regd string) {
	if regd[0] == 'Q' {
		// 128-bit neon
		encoded := assembleNEON(op, true,
			2*int(reg1[1]-'0'),
			2*int(reg2[1]-'0'),
			2*int(regd[1]-'0'))
		fmt.Fprintf(w.w, "\tWORD\t$0x%08x\t// %s %s, %s, %s\n",
			encoded, op, reg1, reg2, regd)
	} else {
		// 64-bit neon
		encoded := assembleNEON(op, false,
			int(reg1[1]-'0'),
			int(reg2[1]-'0'),
			int(regd[1]-'0'))
		fmt.Fprintf(w.w, "\tWORD\t$0x%08x\t// %s %s, %s, %s\n",
			encoded, op, reg1, reg2, regd)
	}
}

func (w codeWriter) label(root, label string) {
	fmt.Fprintln(w.w, "")
	fmt.Fprintln(w.w, root+"__"+label+":")
}

func regNr(reg string) int {
	n, err := strconv.Atoi(reg[1:])
	if err != nil {
		panic("invalid register name " + reg)
	}
	return n
}

func assembleNEON(op string, quad bool, reg1, reg2, regdest int) uint32 {
	tpl := neonTemplates[op]
	encoded := tpl |
		uint32(reg1)<<16 |
		uint32(reg2) |
		uint32(regdest)<<12
	if quad {
		encoded |= 1 << 6
	}
	return encoded
}

// ARM Architecture Reference Manual, F5.4.1
var neonTemplates = map[string]uint32{
	// Integer
	"VADD.I8":  0xf2000800,
	"VADD.I16": 0xf2100800,
	"VADD.I32": 0xf2200800,
	"VSUB.I8":  0xf3000800,
	"VSUB.I16": 0xf3100800,
	"VSUB.I32": 0xf3200800,
	"VMUL.I8":  0xf2000910,
	"VMUL.I16": 0xf2100910,
	"VMUL.I32": 0xf2200910,
	// Carryless
	"VMULCL.I8":  0xf3000910,
	"VMULCL.I16": 0xf3100910,
	"VMULCL.I32": 0xf3200910,
	// bitwise operations
	"VAND": 0xf2000110,
	"VORR": 0xf2200110,
	"VORN": 0xf2300110, // or not
	"VEOR": 0xf3000110, // xor
	// Floating point
	"VADD.F32": 0xf2000d00,
	"VSUB.F32": 0xf2200d00,
	"VMUL.F32": 0xf3000d10,
	// Multiple load/store
	"VLDM": 0xec900b00,
	"VSTM": 0xec800b00,
}
