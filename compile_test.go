package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestCompile(t *testing.T) {
	file, _ := parser.ParseFile(token.NewFileSet(), "", "package p; "+testDecl, 0)
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
	for _, ins := range instrs {
		t.Log(ins)
	}
}
