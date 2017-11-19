package main

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/scanner"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
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
	goarch := runtime.GOARCH
	goarm := "8"
	if s := os.Getenv("GOARCH"); s != "" {
		goarch = s
	}
	if s := os.Getenv("GOARM"); s != "" {
		goarm = s
	}

	tr := NewTranslator(goarch, goarm)
	files, _ := filepath.Glob("*.vgo")
	for _, file := range files {
		err := tr.ProcessFile(file)
		if err != nil {
			scanner.PrintError(os.Stderr, err)
		}
	}
}

type Translator struct {
	fset   *token.FileSet
	goarch string
	goarm  string
	arch   Arch
	funcs  []*Function
}

func NewTranslator(goarch, goarm string) *Translator {
	return &Translator{
		fset:   token.NewFileSet(),
		goarch: goarch,
		goarm:  goarm,
		arch:   FindArch(goarch, goarm),
	}
}
