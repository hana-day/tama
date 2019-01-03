package compiler

import (
	"github.com/hyusuk/tama"
	"github.com/hyusuk/tama/parser"
	"github.com/hyusuk/tama/scanner"
	"strconv"
)

type Compiler struct {
	Insts  []uint32
	Consts []tama.TValue
}

func (c *Compiler) add(inst uint32) {
	c.Insts = append(c.Insts, inst)
}

func (c *Compiler) addABx(op int, a int, bx int) {
	c.add(CreateABx(op, a, bx))
}

func (co *Compiler) addABC(op int, a int, b int, c int) {
	co.add(CreateABC(op, a, b, c))
}

func (c *Compiler) constIndex(v tama.TValue) int {
	for i, cs := range c.Consts {
		if cs == v {
			return i
		}
	}
	c.Consts = append(c.Consts, v)
	return len(c.Consts) - 1
}

func (c *Compiler) compilePrimitive(prim *parser.Primitive) {
	reg := 0
	if prim.Kind == scanner.INT {
		f, _ := strconv.ParseFloat(prim.Value, 64)
		v := tama.TNumber(f)
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
