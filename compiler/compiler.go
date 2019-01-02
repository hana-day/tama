package compiler

import (
	"github.com/hyusuk/tama"
	"github.com/hyusuk/tama/parser"
	"github.com/hyusuk/tama/token"
	"strconv"
)

type Compiler struct {
	insts  []uint32
	consts []tama.TValue
}

func (c *Compiler) add(inst uint32) {
	c.insts = append(c.insts, inst)
}

func (c *Compiler) addABx(op int, a int, bx int) {
	c.add(createABx(op, a, bx))
}

func (co *Compiler) addABC(op int, a int, b int, c int) {
	co.add(createABC(op, a, b, c))
}

func (c *Compiler) constIndex(v tama.TValue) int {
	for i, cs := range c.consts {
		if cs == v {
			return i
		}
	}
	c.consts = append(c.consts, v)
	return len(c.consts) - 1
}

func (c *Compiler) compilePrimitive(prim *parser.Primitive) {
	reg := 0
	if prim.Kind == token.INT {
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
