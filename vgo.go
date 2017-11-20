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
	gosubarch := ""
	if s := os.Getenv("GOARCH"); s != "" {
		goarch = s
	}
	if s := os.Getenv("GOSUBARCH"); s != "" {
		gosubarch = s
	}

	tr := NewTranslator(goarch, gosubarch)
	files, _ := filepath.Glob("*.vgo")
	for _, file := range files {
		err := tr.ProcessFile(file)
		if err != nil {
			scanner.PrintError(os.Stderr, err)
		}
	}
}

type Translator struct {
	fset      *token.FileSet
	goarch    string
	gosubarch string
	arch      Arch
	funcs     []*Function
}

func NewTranslator(goarch, gosubarch string) *Translator {
	return &Translator{
		fset:      token.NewFileSet(),
		goarch:    goarch,
		gosubarch: gosubarch,
		arch:      FindArch(goarch, gosubarch),
	}
}
