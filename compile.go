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
	Kind int
	Var  *Var
	// For OP
	Op    Op
	Left  *Var
	Right *Var
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
			ins.Var.Name, ins.Var.AddrReg, ins.Var.Location)
	case STORE:
		return fmt.Sprintf("STORE %s, %s+%s(*)",
			ins.Var.Location, ins.Var.Name, ins.Var.AddrReg)
	case OP:
		return fmt.Sprintf("OP%s %s, %s, %s",
			ins.Op, ins.Left.Location, ins.Right.Location, ins.Var.Location)
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
	c.AllocRegs(instrs)
	return instrs, nil
}

// buildTree transforms the original expression in a SSA-like
// set of variables.
func (c *Compiler) buildTree(expr ast.Expr, regs map[string]bool) (*Var, error) {
	var do func(ast.Expr) (*Var, error)
	// build temporary variables for intermediate nodes.
	nTemps := 0
	do = func(e ast.Expr) (*Var, error) {
		switch node := e.(type) {
		case *ast.Ident:
			if v, ok := c.Vars[node.Name]; ok {
				return v, nil
			} else {
				return nil, fmt.Errorf("undefined variable %s", node.Name)
			}
		case *ast.ParenExpr:
			return do(node.X)
		case *ast.BinaryExpr:
			op := TokenOp(node.Op)
			left, err := do(node.X)
			if err != nil {
				return nil, err
			}
			right, err := do(node.Y)
			if err != nil {
				return nil, err
			}
			tmp := &Var{
				Name:     fmt.Sprintf("__auto_tmp_%03d", nTemps),
				Location: left.Location,
				Type:     left.Type,
				Op:       op,
				Left:     left,
				Right:    right,
			}
			c.Vars[tmp.Name] = tmp
			nTemps++
			return tmp, nil
		}
		return nil, fmt.Errorf("cannot handle %s", FormatNode(expr))
	}
	return do(expr)
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
			instr := Instr{Kind: LOAD, Var: node}
			instrs = append(instrs, instr)
			done[node] = true
			return // a leaf.
		}
		iterate(node.Left)
		if node.Right != nil {
			iterate(node.Right)
		}
		instr := Instr{Kind: OP, Op: node.Op, Var: node,
			Left:  node.Left,
			Right: node.Right,
		}
		instrs = append(instrs, instr)
		done[node] = true
	}
	iterate(root)
	instrs = append(instrs,
		Instr{Kind: STORE, Var: root})
	return instrs
}

func (c *Compiler) AllocRegs(prog []Instr) {
	lastRef := make(map[*Var]int)
	// lastRef[v] == idx if prog[idx] is the last reference to v
	for idx, ins := range prog {
		if ins.Var != nil {
			lastRef[ins.Var] = idx
		}
		if ins.Left != nil {
			lastRef[ins.Left] = idx
		}
		if ins.Right != nil {
			lastRef[ins.Right] = idx
		}
	}

	regs := make(map[string]bool)
	getReg := func() string {
		for _, reg := range c.VectorRegs {
			if !regs[reg] {
				regs[reg] = true
				return reg
			}
		}
		panic("out of registers")
	}

	for idx, ins := range prog {
		switch ins.Kind {
		case LOAD:
			if ins.Var.Location == "" {
				ins.Var.Location = getReg()
			} else {
				panic("cannot LOAD " + ins.Var.Name + " twice")
			}
		case STORE:
			if ins.Var.Location == "" {
				panic("STORE before assignment of " + ins.Var.Name)
			}
		case OP:
			if ins.Left.Location == "" {
				panic("use before assignment of " + ins.Left.Name)
			}
			if ins.Right.Location == "" {
				panic("use before assignment of " + ins.Right.Name)
			}
			if ins.Var.Location != "" {
				panic("assigned " + ins.Var.Name + " twice")
			}
			// It's forbidden to reuse Right for the result
			// of the instruction (on amd64, we will do
			// MOV Left, Var; OP Right, Var)
			if IsCommutative(ins.Op) && lastRef[ins.Right] == idx {
				ins.Left, ins.Right = ins.Right, ins.Left
				prog[idx] = ins
			}
			if lastRef[ins.Left] == idx {
				regs[ins.Left.Location] = false // free register
			}
			ins.Var.Location = getReg()
			if lastRef[ins.Right] == idx && ins.Right.Location != ins.Var.Location {
				regs[ins.Right.Location] = false // free register
			}
		}
	}
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
