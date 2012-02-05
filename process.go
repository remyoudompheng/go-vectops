package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

type Function struct {
	Name       string
	Args       []string // argument names.
	ScalarType string   // the scalar type ("int", "uint", "float64"...)

	// AST information.
	Decl    *ast.FuncDecl
	Body    *ast.BlockStmt
	Formula ast.Expr // main formula.
}

func (f Function) ForwardDecl() string {
	return FormatNode(f.Decl)
}

func (f Function) String() string {
	buf := bytes.NewBuffer(nil)
	printer.Fprint(buf, token.NewFileSet(), f.Formula)
	return fmt.Sprintf(
		"%s (%v :: [%s]) -> %s",
		f.Name, f.Args, f.ScalarType, buf.Bytes())
}

type Translator struct {
	funcs []*Function
}

// Visit implements the ast.Visitor interface.
func (t *Translator) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if op, ok := IsVectorOp(n); ok {
			t.funcs = append(t.funcs, op)
		}
	}
	return t
}

// IsVectorOp returns true if the given node is a vectorizable function
// declaration. A vectorizable function must have the form:
//
// 	func Op(out, in1, in2, ..., in_k []type) {
// 		out = arithExpr(in1, in2, ..., in_k) 
// 	}
//
// where arithExpr is a simple arithmetic expression.
func IsVectorOp(decl *ast.FuncDecl) (f *Function, vectOk bool) {
	switch {
	case
		// Don't process methods.
		decl.Recv != nil,
		// Don't process functions with results.
		decl.Type.Results != nil,
		// Only one parameter list.
		len(decl.Type.Params.List) != 1,
		// Only one statement.
		len(decl.Body.List) != 1:
		return
	}
	// Now the function declaration has the form:
	//     func F(a1, a2, a3 T) { stmt; }
	paramType := decl.Type.Params.List[0].Type
	var scalarType string
	if t, ok := paramType.(*ast.ArrayType); !ok || t.Len != nil {
		// only process slice types.
		return
	} else {
		// the slice element type must be simple.
		elemType := t.Elt
		if ident, ok := elemType.(*ast.Ident); !ok {
			return
		} else {
			// an identifier
			switch ident.Name {
			case "float32", "float64", "uint", "uint32":
				// Ok.
				scalarType = ident.Name
			default:
				return
			}
		}
	}
	// Now check the body is a single assignment.
	paramNames := make([]string, len(decl.Type.Params.List[0].Names))
	for i, paramNameNode := range decl.Type.Params.List[0].Names {
		paramNames[i] = paramNameNode.Name
	}
	body, isAssign := decl.Body.List[0].(*ast.AssignStmt)
	switch {
	case !isAssign, len(body.Lhs) != 1, len(body.Rhs) != 1:
		return
	}
	expr := body.Rhs[0]
	// Save function body.
	savebody := decl.Body
	decl.Body = nil
	return &Function{
		Name:       decl.Name.Name,
		Decl:       decl,
		Body:       savebody,
		Args:       paramNames,
		ScalarType: scalarType,
		Formula:    expr}, true
}

// ProcessFile processes an input file and write a go and an assembly
// source file.
func ProcessFile(fset *token.FileSet, filename string) (err error) {
	baseName := filename[:len(filename)-len(".vgo")]
	goFile := baseName + "_amd64.go"
	asmFile := baseName + "_amd64.s"

	// Parse and preprocess.
	goInput, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return
	}
	tr := new(Translator)
	ast.Walk(tr, goInput)

	// Write modified Go file
	goF, err := os.Create(goFile)
	if err != nil {
		return fmt.Errorf("error creating %s: %s", goFile, err)
	}
	defer goF.Close()
	printer.Fprint(goF, fset, goInput)

	// Write assembly.
	asmF, err := os.Create(asmFile)
	if err != nil {
		return fmt.Errorf("error creating %s: %s", goFile, err)
	}
	defer asmF.Close()
	tr.CodeGen(asmF)
	return
}
