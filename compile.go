package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

type Compiler struct {
	Arch       Arch
	IndexReg   string
	PtrRegs    []string
	VectorRegs []string
}

// Instantiates a compiling unit for the given arch ('6'),
// with given input/output arguments and registers.
func NewCompiler(arch byte) *Compiler {
	c := new(Compiler)
	switch arch {
	case '6':
		archInfo := amd64
		c.Arch = archInfo
		c.IndexReg = archInfo.CounterReg
		c.PtrRegs = archInfo.InputRegs
		c.VectorRegs = archInfo.VectorRegs
	default:
		err := fmt.Errorf("unsupported arch %q", arch)
		panic(err)
	}
	return c
}

type Var struct {
	Name     string
	Location string
	Register bool
	Type     string
	Op       token.Token
	Left     *Var
	Right    *Var
}

type Tree map[string]*Var

func (c *Compiler) Location(v *Var) string {
	if v.Register {
		return v.Location
	}
	return fmt.Sprintf("(%s)(%s*%d)", v.Location, c.IndexReg, c.Arch.Width(v.Type))
}

func (c *Compiler) Compile(f *Function, w codeWriter) error {
	vars := make(Tree, len(f.Args))
	outv := &Var{Name: f.Args[0],
		Location: "BX",
		Register: false,
		Type:     f.ScalarType}
	vars[outv.Name] = outv
	for i, inname := range f.Args[1:] {
		inv := &Var{Name: inname,
			Location: c.PtrRegs[i],
			Register: false,
			Type:     f.ScalarType}
		vars[inv.Name] = inv
	}
	usedregs := map[string]bool{}
	for _, reg := range c.VectorRegs {
		usedregs[reg] = false
	}
	root, err := c.buildTree(f.Formula, vars, usedregs)
	if err != nil {
		return err
	}
	c.Emit(root, w)
	w.opcode("MOVUPD", c.Location(root), c.Location(outv))
	stride := c.Arch.VectorWidth / c.Arch.Width(f.ScalarType)
	w.opcode("ADDL", fmt.Sprintf("$%d", stride), c.Arch.CounterReg)
	return nil
}

// buildTree allocates registers and prepares the graph of operations.
// Register allocation is greedy.
func (c *Compiler) buildTree(expr ast.Expr, vars Tree, regs map[string]bool) (*Var, error) {
	nTemps := 0
	switch node := expr.(type) {
	case *ast.Ident:
		if v, ok := vars[node.Name]; !ok {
			return nil, fmt.Errorf("undefined variable %s", node.Name)
		} else {
			return v, nil
		}
	case *ast.BinaryExpr:
		op := node.Op
		left, err := c.buildTree(node.X, vars, regs)
		if err != nil {
			return nil, err
		}
		right, err := c.buildTree(node.Y, vars, regs)
		if err != nil {
			return nil, err
		}
		// check op
		_, ok := c.Arch.Opcode(op, left.Type)
		if !ok {
			return nil, fmt.Errorf("incompatible types %s and %s for op %s",
				left.Type, right.Type, op)
		}
		// create a temporary.
		var tmpLoc string
	regswitch:
		switch {
		case left.Register:
			// re-use register
			tmpLoc = left.Location
		case right.Register:
			// use right register for commutative operations
			switch op {
			case token.ADD, token.MUL, token.AND, token.OR:
				tmpLoc = right.Location
				left, right = right, left
				break regswitch
			}
			fallthrough
		default:
			// find another one
			for reg, used := range regs {
				if !used {
					tmpLoc = reg
					regs[reg] = true
					break
				}
			}
			if tmpLoc == "" {
				panic("out of registers")
			}
		}
		tmp := &Var{
			Name:     fmt.Sprintf("__auto_tmp_%03d", nTemps),
			Location: tmpLoc,
			Register: true,
			Type:     left.Type,
			Op:       op,
			Left:     left,
			Right:    right,
		}
		vars[tmp.Name] = tmp
		nTemps++
		return tmp, nil
	}
	return nil, fmt.Errorf("cannot handle %s", FormatNode(expr))
}

func (c *Compiler) Emit(root *Var, w codeWriter) {
	seen := make(map[*Var]bool)
	done := make(map[*Var]bool)
	var iterate func(node *Var)
	iterate = func(node *Var) {
		switch {
		case done[node]:
			return
		case seen[node]:
			err := fmt.Errorf("loop detected in optree for %s", node.Name)
			panic(err)
		}
		opcode, _ := c.Arch.Opcode(node.Op, node.Type)
		seen[node] = true
		if node.Left == nil {
			done[node] = true
			return // a leaf.
		}
		iterate(node.Left)
		if node.Right != nil {
			iterate(node.Right)
		}
		if node.Left.Register && node.Location == node.Left.Location {
			w.opcode(opcode, c.Location(node.Right), c.Location(node))
		} else {
			w.opcode("MOVUPS", c.Location(node.Left), c.Location(node))
			w.opcode(opcode, c.Location(node.Right), c.Location(node))
		}
		done[node] = true
	}
	iterate(root)
}