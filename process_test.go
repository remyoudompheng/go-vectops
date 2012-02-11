package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"reflect"
	"testing"
)

const testDecl = `
func F(out, in1, in2, in3 []float32) {
	out = in1 + in2 * in3
}
`

func TestVectorFunc(t *testing.T) {
	file, _ := parser.ParseFile(token.NewFileSet(), "", "package p; "+testDecl, 0)
	decl := file.Decls[0].(*ast.FuncDecl)
	f, err := IsVectorOp(decl)
	switch {
	case err != nil:
		t.Errorf("F should be vectorizable but %s", err)
	case !reflect.DeepEqual(f.Args, []string{"out", "in1", "in2", "in3"}):
		t.Errorf("wrong args for F: got %v", f.Args)
	case f.ScalarType != "float32":
		t.Errorf("wrong scalar type %s, expected float32", f.ScalarType)
	}
	t.Logf("function info: %+v", f)
	buf := new(bytes.Buffer)
	printer.Fprint(buf, token.NewFileSet(), decl)
	t.Logf("processed declaration: %s", buf.Bytes())
}
