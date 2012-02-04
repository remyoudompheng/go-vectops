package main

import (
	"go/scanner"
	"go/token"
	"os"
	"path/filepath"
)

var (
	fset = token.NewFileSet()
)

func main() {
	files, _ := filepath.Glob("*.vgo")
	for _, file := range files {
		err := ProcessFile(fset, file)
		if err != nil {
			scanner.PrintError(os.Stderr, err)
		}
	}
}
