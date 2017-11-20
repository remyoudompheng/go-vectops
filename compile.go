package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

type Compiler struct {
	PtrRegs    []string
	VectorRegs []string
	Vars       map[string]*Var
}

type Var struct {
	Name     string
	Location string
	AddrReg  string
	Type     string
	ReadOnly bool
	Op       Op
	Left     *Var
	Right    *Var
}

func (v *Var) Expr() string {
	return v.Left.Name + " " + opstring[v.Op] + " " + v.Right.Name
}

type Instr struct {
	Kind    int
	RegDest string
	// For LOAD, STORE
	Var *Var
	// For OP
	Op       Op
	RegLeft  string
	RegRight string
}

const (
	LOAD = iota
	STORE
	OP
)

func (ins Instr) String() string {
	switch ins.Kind {
	case LOAD:
		return fmt.Sprintf("LOAD %s+%s(*), %s",
			ins.Var.Name, ins.Var.AddrReg, ins.RegDest)
	case STORE:
		return fmt.Sprintf("STORE %s, %s+%s(*)",
			ins.RegDest, ins.Var.Name, ins.Var.AddrReg)
	case OP:
		return fmt.Sprintf("OP%s %s, %s, %s",
			ins.Op, ins.RegLeft, ins.RegRight, ins.RegDest)
	default:
		panic("impossible")
	}
}

func Compile(f *Function, w codeWriter) ([]Instr, error) {
	c := &Compiler{
		PtrRegs:    w.arch.InputRegs,
		VectorRegs: w.arch.VectorRegs,
		Vars:       make(map[string]*Var),
	}
	for i, name := range f.Args {
		v := &Var{Name: name,
			AddrReg: c.PtrRegs[i],
			Type:    f.ScalarType}
		c.Vars[v.Name] = v
	}
	usedregs := map[string]bool{}
	for _, reg := range c.VectorRegs {
		usedregs[reg] = false
	}
	root, err := c.buildTree(f.Formula, usedregs)
	if err != nil {
		return nil, err
	}
	// root is actually f.Args[0]
	outvar := *root
	delete(c.Vars, root.Name)
	outvar.Name = f.Args[0]
	outvar.AddrReg = w.arch.InputRegs[0]
	root = &outvar
	c.Vars[outvar.Name] = root
	instrs := c.Emit(root)
	return instrs, nil
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
func (c *Compiler) buildTree(expr ast.Expr, regs map[string]bool) (*Var, error) {
	var pass1 func(ast.Expr) error
	var pass2 func(ast.Expr) (*Var, error)
	// Pass 1 allocates registers for input array elements.
	pass1 = func(e ast.Expr) error {
		switch node := e.(type) {
		case *ast.Ident:
			if v, ok := c.Vars[node.Name]; !ok {
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
			return c.Vars[node.Name], nil
		case *ast.ParenExpr:
			return pass2(node.X)
		case *ast.BinaryExpr:
			op := TokenOp(node.Op)
			left, err := pass2(node.X)
			if err != nil {
				return nil, err
			}
			right, err := pass2(node.Y)
			if err != nil {
				return nil, err
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
			c.Vars[tmp.Name] = tmp
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

func (c *Compiler) Emit(root *Var) []Instr {
	var instrs []Instr
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
		seen[node] = true
		if node.Left == nil {
			if node.AddrReg == "" {
				// A leaf necessarily comes from memory.
				panic("impossible")
			}
			instr := Instr{Kind: LOAD, Var: node, RegDest: node.Location}
			instrs = append(instrs, instr)
			done[node] = true
			return // a leaf.
		}
		iterate(node.Left)
		if node.Right != nil {
			iterate(node.Right)
		}
		instr := Instr{Kind: OP, Op: node.Op, Var: node,
			RegLeft:  node.Left.Location,
			RegRight: node.Right.Location,
			RegDest:  node.Location,
		}
		instrs = append(instrs, instr)
		done[node] = true
	}
	iterate(root)
	instrs = append(instrs,
		Instr{Kind: STORE, Var: root, RegDest: root.Location})
	return instrs
}

func TokenOp(tok token.Token) Op {
	switch tok {
	case token.ADD:
		return ADD
	case token.SUB:
		return SUB
	case token.MUL:
		return MUL
	case token.QUO:
		return DIV
	case token.AND:
		return AND
	case token.OR:
		return OR
	case token.XOR:
		return XOR
	case token.SHL:
		return SHL
	case token.SHR:
		return SHR
	default:
		panic("unsupported operator " + tok.String())
	}
}

func IsCommutative(op Op) bool {
	switch op {
	case ADD, AND, OR, XOR, MUL:
		return true
	}
	return false
}
