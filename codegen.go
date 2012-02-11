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
	err := f.Compile(w)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "")
}

func (f *Function) Compile(w codeWriter) error {
	c := NewCompiler('6')
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

	// Emit code for the loop. It should look like:
	// for i := 0; ; {
	// 	if i > length-4 { i = length-4 }
	// 	process(arrays[i:i+4])
	// 	i += 4
	// 	if i >= length { break }
	// }
	w.opcode("SUBL", "$4", "DX")
	w.opcode("XORL", "CX", "CX")
	w.label(f.Name, "loop")
	w.opcode("CMPL", "CX", "DX")
	w.comment("if i > length-4 { i = length-4 }")
	w.opcode("JLE", f.Name+"·process")
	w.opcode("MOVL", "DX", "CX")
	w.label(f.Name, "process")

	err := c.Compile(f, w)
	if err != nil {
		return err
	}

	w.comment("if i >= length { break }")
	w.opcode("CMPL", "CX", outArg+"+8(FP)")
	w.opcode("JGE", f.Name+"·return")
	w.opcode("JMP", f.Name+"·loop")
	w.label(f.Name, "return")
	w.opcode("RET")
	return nil
}

type codeWriter struct{ io.Writer }

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
	fmt.Fprintln(w, root+"·"+label+":")
}
