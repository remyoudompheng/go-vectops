vgo is a compiler for array arithmetic expressions.

It is meant for use by Go programs, and currently targets SSE2 on amd64.
It compiles functions written in Go syntax: see samples/vector.vgo
for examples.

Usage
-----

vgo's language is Go-like (hence gofmt-friendly) and very limited
and simple. Functions must take arguments of the same type,
the first type being the output slice.

All arguments must have the same length and have size multiple of 128
bits, otherwise an index runtime-panic will be triggered.

    func SomeFormula(out, x, y, z []float32) {
    	out = (x*y + y*z + z*x) / (x*x + y*y + z*z)
    }

Functions may only be made of a single assignment statement
with an arithmetic expression. Not all operators are supported
by all types.

A vgo source file may mix normal Go functions and vgo functions.
It is named with the `.vgo` extension. The vgo invocation

    $ vgo

will produce from a vgo source file `source.vgo` a Go source file
`source_amd64.go` and a Plan 9-like assembly file `source_amd64.s` suitable for
use by 6a.

TODOs
-----

* support multiple statements.
* support integer quotients (via magic constants?).
* support scalar/vector operations and constants.

