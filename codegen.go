package main

import (
	"fmt"
	"io"
	"strings"
)

func (tr Translator) CodeGen(w io.Writer) {
	for _, f := range tr.funcs {
		f.CodeGen(codeWriter{w})
	}
}

func (f *Function) CodeGen(w codeWriter) {
	fmt.Fprintf(w, "// %s\n", f.ForwardDecl())
	fmt.Fprintf(w, "TEXT ·%s(SB), 7, $0\n", f.Name)
	f.WritePrologue(w)
	c := NewCompiler('6')
	err := c.Compile(f, w)
	if err != nil {
		panic(err)
	}
	f.WriteEpilogue(w)
	fmt.Fprintln(w, "")
}

func (f *Function) WritePrologue(w codeWriter) {
	inputRegs := amd64.InputRegs
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
	w.opcode("MOVQ", outArg+"+0(FP)", "BX")
	for i, arg := range inArgs {
		w.opcode("MOVQ",
			fmt.Sprintf("%s+%d(FP)", arg, 16*i+16),
			inputRegs[i])
	}

	// length and index.
	w.comment("Check lengths.")
	w.opcode("MOVL", outArg+"+8(FP)", "DX")
	for i, arg := range inArgs {
		w.opcode("CMPL", "DX",
			fmt.Sprintf("%s+%d(FP)", arg, 16*i+24))
		w.opcode("JNE", f.Name+"·panic")
	}
	w.opcode("JMP", f.Name+"·ok")
	//
	w.label(f.Name, "panic")
	w.opcode("CALL", "runtime·panicindex(SB)")
	w.label(f.Name, "ok")

	// start loop
	w.opcode("MOVL", "$0", "CX")
	w.label(f.Name, "loop")
	w.opcode("CMPL", "CX", "DX")
	w.opcode("JGE", f.Name+"·return")
}

func (f *Function) WriteEpilogue(w codeWriter) {
	w.opcode("JMP", f.Name+"·loop")
	w.label(f.Name, "return")
	w.opcode("RET")
}

type codeWriter struct{ io.Writer }

func (w codeWriter) comment(s string) {
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "\t// "+s)
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
	fmt.Fprintln(w, root+"·"+label+":")
}
