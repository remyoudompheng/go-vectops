package main

import (
	"fmt"
	"io"
	"strings"
)

type codeWriter struct {
	w    io.Writer
	arch Arch
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

func (f *Function) Compile(w codeWriter) error {
	if w.arch.Width(f.ScalarType) == 0 {
		return fmt.Errorf("unsupported data type %s", f.ScalarType)
	}
	instrs, err := Compile(f, w)
	if err != nil {
		return err
	}
	for _, ins := range instrs {
		if err := w.checkInstr(ins); err != nil {
			return err
		}
	}

	ptrSize := w.arch.PtrSize
	// BX: pointer to output slice
	// CX: index counter.
	// DX: length
	// SI/DI/Rxx: pointers to inputs.
	w.comment("Load pointers.")
	outArg := f.Args[0]
	inArgs := f.Args[1:]
	if len(inArgs) > len(w.arch.InputRegs) {
		panic("not enough registers")
	}
	for i, arg := range f.Args {
		w.opcode("MOVQ",
			frameArg(arg, ptrSize*(3*i)),
			w.arch.InputRegs[i])
	}

	// length and index.
	w.comment("Check lengths.")
	w.opcode("MOVQ", outArg+"+8(FP)", "DX")
	for i, arg := range inArgs {
		w.opcode("CMPQ", "DX",
			frameArg(arg, ptrSize*(3*i+3+1)))
		w.opcode("JNE", f.Name+"__panic")
	}
	w.opcode("JMP", f.Name+"__ok")
	//
	w.label(f.Name, "panic")
	w.opcode("CALL", "runtime·panicindex(SB)")
	w.label(f.Name, "ok")

	stride := w.arch.VectorWidth / w.arch.Width(f.ScalarType)
	// Emit code for the loop. It should look like:
	// for i := 0; ; {
	// 	if i > length-4 { i = length-4 }
	// 	process(arrays[i:i+4])
	// 	i += 4
	// 	if i >= length { break }
	// }
	w.opcode("SUBQ", fmt.Sprintf("$%d", stride), "DX")
	w.opcode("XORQ", "CX", "CX")
	w.label(f.Name, "loop")
	w.opcode("CMPQ", "CX", "DX")
	w.comment("if i > length-%d { i = length-%d }", stride, stride)
	w.opcode("JLE", f.Name+"__process")
	w.opcode("MOVQ", "DX", "CX")
	w.label(f.Name, "process")

	for _, ins := range instrs {
		w.emitInstr(ins)
	}

	w.opcode("ADDQ", fmt.Sprintf("$%d", stride), w.arch.CounterReg)
	w.comment("if i >= length { break }")
	w.opcode("CMPQ", "CX", frameArg(outArg, 2*ptrSize))
	w.opcode("JGE", f.Name+"__return")
	w.opcode("JMP", f.Name+"__loop")
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
	if true {
		// amd64
		switch ins.Kind {
		case LOAD:
			width := w.arch.Width(ins.Var.Type)
			loc := fmt.Sprintf("(%s)(%s*%d)", ins.Var.AddrReg,
				w.arch.CounterReg, width)
			w.opcode("MOVUPS", loc, ins.RegDest)
		case STORE:
			width := w.arch.Width(ins.Var.Type)
			loc := fmt.Sprintf("(%s)(%s*%d)", ins.Var.AddrReg,
				w.arch.CounterReg, width)
			w.opcode("MOVUPD", ins.RegDest, loc)
		case OP:
			v := ins.Var
			w.comment("%s = %s %s %s", v.Name, v.Left.Name, v.Op, v.Right.Name)
			opcode, ok := w.arch.Opcode(ins.Op, v.Type)
			if !ok {
				panic("unsupported operation")
			}
			if ins.RegDest == ins.RegLeft {
				w.opcode(opcode, ins.RegRight, ins.RegLeft)
			} else {
				w.opcode("MOVAPS", ins.RegLeft, ins.RegDest)
				w.opcode(opcode, ins.RegRight, ins.RegDest)
			}
		}
	} else {
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

func (w codeWriter) label(root, label string) {
	fmt.Fprintln(w.w, "")
	fmt.Fprintln(w.w, root+"__"+label+":")
}
