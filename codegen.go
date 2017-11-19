package main

import (
	"fmt"
	"io"
	"strings"
)

type codeWriter struct {
	io.Writer
	compiler *Compiler
}

func (tr Translator) CodeGen(w io.Writer) {
	for _, f := range tr.funcs {
		f.CodeGen(codeWriter{w, NewCompiler(tr.GOARCH, "")})
	}
}

func (f *Function) CodeGen(w codeWriter) {
	fmt.Fprintf(w, "// %s\n", f.ForwardDecl())
	fmt.Fprintf(w, "TEXT ·%s(SB), 7, $0\n", f.Name)
	err := f.Compile(w)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "")
}

func frameArg(name string, offset int) string {
	return fmt.Sprintf("%s+%d(FP)", name, offset)
}

func (f *Function) Compile(w codeWriter) error {
	c := w.compiler
	ptrSize := c.Arch.PtrSize
	inputRegs := c.Arch.InputRegs
	// BX: pointer to output slice
	// CX: index counter.
	// DX: length
	// SI/DI/Rxx: pointers to inputs.
	w.comment("Load pointers.")
	outArg := f.Args[0]
	inArgs := f.Args[1:]
	if len(inArgs) > len(inputRegs) {
		panic("not enough registers")
	}
	w.opcode("MOVQ", frameArg(outArg, 0), "BX")
	for i, arg := range inArgs {
		w.opcode("MOVQ",
			frameArg(arg, ptrSize*(3*i+3)),
			inputRegs[i])
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

	stride := c.Arch.VectorWidth / c.Arch.Width(f.ScalarType)
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

	err := c.Compile(f, w)
	if err != nil {
		return err
	}

	w.opcode("ADDQ", fmt.Sprintf("$%d", stride), c.Arch.CounterReg)
	w.comment("if i >= length { break }")
	w.opcode("CMPQ", "CX", frameArg(outArg, 2*ptrSize))
	w.opcode("JGE", f.Name+"__return")
	w.opcode("JMP", f.Name+"__loop")
	w.label(f.Name, "return")
	w.opcode("RET")
	return nil
}

func (w codeWriter) comment(format string, args ...interface{}) {
	fmt.Fprintf(w, "\n\t// "+format+"\n", args...)
}

func (w codeWriter) opcode(op string, args ...string) {
	if len(args) > 0 {
		fmt.Fprintf(w, "\t%s\t%s\n", op, strings.Join(args, ", "))
	} else {
		fmt.Fprintf(w, "\t%s\n", op)
	}
}

func (w codeWriter) label(root, label string) {
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, root+"__"+label+":")
}
