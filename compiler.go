package tama

import (
	"github.com/hyusuk/tama/parser"
	"github.com/hyusuk/tama/scanner"
	"log"
	"strconv"
)

type Compiler struct {
	Proto *ClosureProto // current function header
	nreg  int           // number of registers
}

type Reg struct {
	N int // register number
}

func (c *Compiler) newReg() *Reg {
	reg := &Reg{
		N: c.nreg,
	}
	c.nreg++
	return reg
}

func (c *Compiler) add(inst uint32) {
	c.Proto.Insts = append(c.Proto.Insts, inst)
}

func (c *Compiler) addABx(op int, a int, bx int) {
	c.add(CreateABx(op, a, bx))
}

func (co *Compiler) addABC(op int, a int, b int, c int) {
	co.add(CreateABC(op, a, b, c))
}

func (c *Compiler) constIndex(v Value) int {
	for i, cs := range c.Proto.Consts {
		if cs == v {
			return i
		}
	}
	c.Proto.Consts = append(c.Proto.Consts, v)
	return len(c.Proto.Consts) - 1
}

func (c *Compiler) compilePrimitive(prim *parser.Primitive) *Reg {
	reg := c.newReg()
	if prim.Kind == scanner.INT {
		f, _ := strconv.ParseFloat(prim.Value, 64)
		v := Number(f)
		c.addABx(LOADK, reg.N, c.constIndex(v))
	}
	return reg
}

func (c *Compiler) compileIdent(ident *parser.Ident) *Reg {
	r1 := c.newReg()
	c.addABx(LOADK, r1.N, c.constIndex(String(ident.Name)))
	r2 := c.newReg()
	c.addABx(GETGLOBAL, r2.N, r1.N)
	return r2
}

func (c *Compiler) compileExpr(expr parser.Expr) *Reg {
	switch ex := expr.(type) {
	case *parser.Primitive:
		return c.compilePrimitive(ex)
	case *parser.Ident:
		return c.compileIdent(ex)
	default:
		log.Fatalf("Unknown expression %v", ex)
	}
	return nil
}

func (c *Compiler) compileExprs(exprs []parser.Expr) {
	for _, expr := range exprs {
		c.compileExpr(expr)
	}
}

func newClosureProto() *ClosureProto {
	return &ClosureProto{
		Insts:        []uint32{},
		Consts:       []Value{},
		MaxStackSize: 256,
	}
}

func Compile(exprs []parser.Expr) (*Closure, error) {
	c := Compiler{
		Proto: newClosureProto(),
		nreg:  0,
	}
	c.compileExprs(exprs)
	rn := c.nreg - 1
	c.addABC(RETURN, rn, 2, 0)

	cl := NewScmClosure()
	cl.Proto = c.Proto
	return cl, nil
}
