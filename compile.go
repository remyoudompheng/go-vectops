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

type Var struct {
	Name     string
	Location string
	AddrReg  string
	Type     string
	ReadOnly bool
	Op       token.Token
	Left     *Var
	Right    *Var
}

type Tree map[string]*Var

func (c *Compiler) MemLocation(v *Var) string {
	return fmt.Sprintf("(%s)(%s*%d)", v.AddrReg, c.IndexReg, c.Arch.Width(v.Type))
}

func Compile(f *Function, w codeWriter) error {
	c := &Compiler{
		Arch:       w.arch,
		IndexReg:   w.arch.CounterReg,
		PtrRegs:    w.arch.InputRegs,
		VectorRegs: w.arch.VectorRegs,
	}
	vars := make(Tree, len(f.Args))
	outv := &Var{Name: f.Args[0],
		AddrReg: "BX",
		Type:    f.ScalarType}
	vars[outv.Name] = outv
	for i, inname := range f.Args[1:] {
		inv := &Var{Name: inname,
			AddrReg: c.PtrRegs[i],
			Type:    f.ScalarType}
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
	w.opcode("MOVUPD", root.Location, c.MemLocation(outv))
	return nil
}

func (c *Compiler) getFreeReg(regs map[string]bool) string {
	for _, reg := range c.VectorRegs {
		if !regs[reg] {
			regs[reg] = true
			// fmt.Printf("using register %s\n", reg)
			return reg
		}
	}
	panic("out of registers")
}

// buildTree allocates registers and prepares the graph of operations.
// Register allocation is greedy.
func (c *Compiler) buildTree(expr ast.Expr, vars Tree, regs map[string]bool) (*Var, error) {
	var pass1 func(ast.Expr) error
	var pass2 func(ast.Expr) (*Var, error)
	// Pass 1 allocates registers for input array elements.
	pass1 = func(e ast.Expr) error {
		switch node := e.(type) {
		case *ast.Ident:
			if v, ok := vars[node.Name]; !ok {
				return fmt.Errorf("undefined variable %s", node.Name)
			} else {
				if v.Location != "" {
					// already chosen a register: it's shared, mark it read only
					v.ReadOnly = true
				} else {
					v.Location = c.getFreeReg(regs)
				}
				return nil
			}
		case *ast.ParenExpr:
			return pass1(node.X)
		case *ast.BinaryExpr:
			if err := pass1(node.X); err != nil {
				return err
			}
			if err := pass1(node.Y); err != nil {
				return err
			}
		}
		return nil
	}
	// Pass 2 build temporary variables for intermediate nodes.
	nTemps := 0
	pass2 = func(e ast.Expr) (*Var, error) {
		switch node := e.(type) {
		case *ast.Ident:
			return vars[node.Name], nil
		case *ast.ParenExpr:
			return pass2(node.X)
		case *ast.BinaryExpr:
			op := node.Op
			left, err := pass2(node.X)
			if err != nil {
				return nil, err
			}
			right, err := pass2(node.Y)
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
			if left.ReadOnly && !right.ReadOnly && IsCommutative(op) {
				left, right = right, left
			}
			tmp := &Var{
				Name:     fmt.Sprintf("__auto_tmp_%03d", nTemps),
				Location: left.Location,
				Type:     left.Type,
				Op:       op,
				Left:     left,
				Right:    right,
			}
			if left.ReadOnly {
				// must allocate a register
				tmp.Location = c.getFreeReg(regs)
			}
			vars[tmp.Name] = tmp
			nTemps++
			return tmp, nil
		}
		return nil, fmt.Errorf("cannot handle %s", FormatNode(expr))
	}

	// run passes
	err := pass1(expr)
	if err != nil {
		return nil, err
	}
	return pass2(expr)
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
			if node.AddrReg != "" {
				w.opcode("MOVUPS", c.MemLocation(node), node.Location)
				node.AddrReg = ""
			}
			done[node] = true
			return // a leaf.
		}
		iterate(node.Left)
		if node.Right != nil {
			iterate(node.Right)
		}
		w.comment("%s = %s %s %s", node.Name, node.Left.Name, node.Op, node.Right.Name)
		if node.Location != node.Left.Location {
			w.opcode("MOVAPS", node.Left.Location, node.Location)
		}
		w.opcode(opcode, node.Right.Location, node.Location)
		done[node] = true
	}
	iterate(root)
}
