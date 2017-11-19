package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
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

func forwardDecl(decl *ast.FuncDecl) string {
	fwd := *decl
	fwd.Body = nil
	return FormatNode(&fwd)
}

func (f Function) String() string {
	buf := bytes.NewBuffer(nil)
	printer.Fprint(buf, token.NewFileSet(), f.Formula)
	return fmt.Sprintf(
		"%s (%v :: [%s]) -> %s",
		f.Name, f.Args, f.ScalarType, buf.Bytes())
}

// Visit implements the ast.Visitor interface.
func (t *Translator) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if op, err := IsVectorOp(n); err == nil {
			t.funcs = append(t.funcs, op)
		} else {
			fmt.Fprintln(os.Stderr, forwardDecl(n))
			fmt.Fprintf(os.Stderr, "\tskipping: %s\n", err)
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
func IsVectorOp(decl *ast.FuncDecl) (f *Function, err error) {
	switch {
	case decl.Recv != nil:
		return nil, fmt.Errorf("is method")
	case decl.Type.Results != nil:
		return nil, fmt.Errorf("has return values")
	case len(decl.Type.Params.List) != 1:
		return nil, fmt.Errorf("many parameter lists")
	case len(decl.Body.List) != 1:
		return nil, fmt.Errorf("more than 1 statement")
	}
	// Now the function declaration has the form:
	//     func F(a1, a2, a3 T) { stmt; }
	paramType := decl.Type.Params.List[0].Type
	var scalarType string
	if t, ok := paramType.(*ast.ArrayType); !ok || t.Len != nil {
		// only process slice types.
		return nil, fmt.Errorf("non-slice type %s", FormatNode(paramType))
	} else {
		// the slice element type must be simple.
		elemType := t.Elt
		if ident, ok := elemType.(*ast.Ident); !ok {
			return nil, fmt.Errorf("unsupported type %s", FormatNode(paramType))
		} else {
			// an identifier
			if _, ok := types[ident.Name]; ok || ident.Name == "uint" {
				scalarType = ident.Name
			} else {
				return nil, fmt.Errorf("unsupported type %s", FormatNode(paramType))
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
		return nil, fmt.Errorf("statement %s is not an assignment", FormatNode(body))
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
		Formula:    expr}, nil
}

// ProcessFile processes an input file and write a go and an assembly
// source file.
func (t *Translator) ProcessFile(filename string) (err error) {
	baseName := filename[:len(filename)-len(".vgo")]
	goFile := baseName + "_" + t.goarch + ".go"
	asmFile := baseName + "_" + t.goarch + ".s"

	// Parse and preprocess.
	goInput, err := parser.ParseFile(t.fset, filename, nil, parser.ParseComments)
	if err != nil {
		return
	}
	ast.Walk(t, goInput)

	// Write modified Go file
	goF, err := os.Create(goFile)
	if err != nil {
		return fmt.Errorf("error creating %s: %s", goFile, err)
	}
	printconfig.Fprint(goF, t.fset, goInput)
	err = goF.Close()
	if err != nil {
		return fmt.Errorf("error creating %s: %s", goFile, err)
	}

	// Write assembly.
	asmBuf := new(bytes.Buffer)
	w := codeWriter{
		w:      asmBuf,
		goarch: t.goarch,
		arch:   t.arch,
	}
	for _, f := range t.funcs {
		err = w.CodeGen(f)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(asmFile, asmBuf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error creating %s: %s", asmFile, err)
	}
	return nil
}
