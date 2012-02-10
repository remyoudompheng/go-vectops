package main

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/scanner"
	"go/token"
	"os"
	"path/filepath"
)

var (
	fset        = token.NewFileSet()
	printconfig = printer.Config{Mode: printer.SourcePos}
)

func FormatNode(node ast.Node) string {
	buf := new(bytes.Buffer)
	printer.Fprint(buf, fset, node)
	return string(buf.Bytes())
}

func main() {
	files, _ := filepath.Glob("*.vgo")
	for _, file := range files {
		err := ProcessFile(fset, file)
		if err != nil {
			scanner.PrintError(os.Stderr, err)
		}
	}
}
