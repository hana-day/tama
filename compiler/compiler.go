package compiler

import (
	"github.com/hyusuk/tama"
	"github.com/hyusuk/tama/parser"
	"github.com/hyusuk/tama/scanner"
	"strconv"
)

type FuncProto struct {
	Insts  []uint32
	Consts []tama.TValue
}

func newFuncProto() *FuncProto {
	return &FuncProto{
		Insts:  []uint32{},
		Consts: []tama.TValue{},
	}
}

type Compiler struct {
	Proto *FuncProto
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

func (c *Compiler) constIndex(v tama.TValue) int {
	for i, cs := range c.Proto.Consts {
		if cs == v {
			return i
		}
	}
	c.Proto.Consts = append(c.Proto.Consts, v)
	return len(c.Proto.Consts) - 1
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

func Compile(exprs []parser.Expr) (*FuncProto, error) {
	c := Compiler{
		Proto: newFuncProto(),
	}
	c.compileExprs(exprs)
	return c.Proto, nil
}
