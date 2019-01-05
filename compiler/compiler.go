package compiler

import (
	"github.com/hyusuk/tama/parser"
	"github.com/hyusuk/tama/scanner"
	"github.com/hyusuk/tama/types"
	"strconv"
)

type Closure struct {
	Insts        []uint32
	Consts       []types.Value
	MaxStackSize int
}

func newClosure() *Closure {
	return &Closure{
		Insts:        []uint32{},
		Consts:       []types.Value{},
		MaxStackSize: 256,
	}
}

type Compiler struct {
	Cl *Closure
}

func (c *Compiler) add(inst uint32) {
	c.Cl.Insts = append(c.Cl.Insts, inst)
}

func (c *Compiler) addABx(op int, a int, bx int) {
	c.add(CreateABx(op, a, bx))
}

func (co *Compiler) addABC(op int, a int, b int, c int) {
	co.add(CreateABC(op, a, b, c))
}

func (c *Compiler) constIndex(v types.Value) int {
	for i, cs := range c.Cl.Consts {
		if cs == v {
			return i
		}
	}
	c.Cl.Consts = append(c.Cl.Consts, v)
	return len(c.Cl.Consts) - 1
}

func (c *Compiler) compilePrimitive(prim *parser.Primitive) {
	reg := 0
	if prim.Kind == scanner.INT {
		f, _ := strconv.ParseFloat(prim.Value, 64)
		v := types.Number(f)
		c.addABx(LOADK, reg, c.constIndex(v))
	}
}

func (c *Compiler) compileExpr(expr parser.Expr) {
	if prim, ok := expr.(*parser.Primitive); ok {
		c.compilePrimitive(prim)
	}
}

func (c *Compiler) compileExprs(exprs []parser.Expr) {
	for _, expr := range exprs {
		c.compileExpr(expr)
		c.addABC(RETURN, 0, 2, 0)
	}
}

func Compile(exprs []parser.Expr) (*Closure, error) {
	c := Compiler{
		Cl: newClosure(),
	}
	c.compileExprs(exprs)
	return c.Cl, nil
}
