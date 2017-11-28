package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	gotypes "go/types"
	"testing"
)

func TestCompile(t *testing.T) {
	for _, expr := range testExprs {
		pkg := "package p; func F(out, in1, in2, in3 []float32) {" + expr + "}"
		file, _ := parser.ParseFile(token.NewFileSet(), "", pkg, 0)
		decl := file.Decls[0].(*ast.FuncDecl)
		f, err := IsVectorOp(decl)
		if err != nil {
			t.Fatalf("could not check %s", FormatNode(decl))
		}
		w := codeWriter{arch: amd64}
		instrs, err := Compile(f, w)
		if err != nil {
			t.Error(err)
		}
		t.Logf("compiled %s", gotypes.ExprString(f.Formula))
		for _, ins := range instrs {
			t.Log(ins)
		}
	}
}

var testExprs = []string{
	"out = in1 + in2 * in3",
	"out = in1 + in2 - in3",
	"out = in1 - in2 - in3",
	"out = in1 / in2 - in3",
	"out = in1 * in2 - in1",
	"out = in1 * in1 + in2 * in2",
}
