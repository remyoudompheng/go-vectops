package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"
)

func TestCompile(t *testing.T) {
	file, _ := parser.ParseFile(token.NewFileSet(), "", "package p; "+testDecl, 0)
	decl := file.Decls[0].(*ast.FuncDecl)
	f, ok := IsVectorOp(decl)
	if !ok {
		t.Fatalf("could not check %s", FormatNode(decl))
	}
	c := NewCompiler('6')
	err := c.Compile(f, codeWriter{os.Stderr})
	if err != nil {
		t.Error(err)
	}
}
