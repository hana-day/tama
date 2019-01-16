package compiler

import (
	"github.com/hyusuk/tama/types"
	"log"
)

type Compiler struct {
	Proto *types.ClosureProto // current function header
	nreg  int                 // number of registers
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

func (c *Compiler) constIndex(v types.Object) int {
	for i, cs := range c.Proto.Consts {
		if cs == v {
			return i
		}
	}
	c.Proto.Consts = append(c.Proto.Consts, v)
	return len(c.Proto.Consts) - 1
}

func (c *Compiler) compileNumber(num types.Number) *Reg {
	reg := c.newReg()
	c.addABx(LOADK, reg.N, c.constIndex(num))
	return reg
}

func (c *Compiler) compileSymbol(sym *types.Symbol) *Reg {
	r1 := c.newReg()
	c.addABx(GETGLOBAL, r1.N, c.constIndex(types.String(sym.Name)))
	return r1
}

// Compile define syntax
//
// (define a 1)
func (c *Compiler) compileDefine(pair *types.Pair) *Reg {
	cdr, _ := types.Cdr(pair)
	cdar, _ := types.Car(cdr)
	cddr, _ := types.Cdr(cdr)
	rest, _ := types.Car(cddr)
	sym, ok := cdar.(*types.Symbol)
	if !ok {
		log.Fatalf("The first argument of define must be a symbol")
	}
	nameReg := c.newReg()
	c.addABx(LOADK, nameReg.N, c.constIndex(sym.Name))
	valueReg := c.compileObject(rest)
	c.addABx(SETGLOBAL, valueReg.N, nameReg.N)
	return valueReg
}

func (c *Compiler) compilePair(pair *types.Pair) *Reg {
	v, _ := types.Car(pair)
	first, ok := v.(*types.Symbol)
	if !ok {
		log.Fatalf("Invalid function name")
	}
	switch first.Name {
	case "define":
		return c.compileDefine(pair)
	default:
		r1 := c.compileSymbol(first)
		cdr, _ := types.Cdr(pair)
		args := cdr.(*types.Pair)
		argsArr := args.ListToArray()
		argRegs := make([]*Reg, len(argsArr))
		for i := 0; i < len(argsArr); i++ {
			argRegs[i] = c.newReg()
		}
		var r *Reg
		for i, arg := range argsArr {
			r = c.compileObject(arg)
			c.addABC(MOVE, argRegs[i].N, r.N, 0)
		}
		// Always return one value
		c.addABC(CALL, r1.N, 1+len(argsArr), 2)
		return r1
	}
}

func (c *Compiler) compileObject(obj types.Object) *Reg {
	switch o := obj.(type) {
	case types.Number:
		return c.compileNumber(o)
	case *types.Symbol:
		return c.compileSymbol(o)
	case *types.Pair:
		return c.compilePair(o)
	default:
		log.Fatalf("Unknown type of object %v", o)
	}
	return nil
}

func (c *Compiler) compileObjects(objs []types.Object) []*Reg {
	regs := make([]*Reg, len(objs))
	for i, obj := range objs {
		regs[i] = c.compileObject(obj)
	}
	return regs
}

func newClosureProto() *types.ClosureProto {
	return &types.ClosureProto{
		Insts:        []uint32{},
		Consts:       []types.Object{},
		MaxStackSize: 256,
	}
}

func Compile(objs []types.Object) (*types.Closure, error) {
	c := Compiler{
		Proto: newClosureProto(),
		nreg:  0,
	}
	regs := c.compileObjects(objs)
	lastReg := regs[len(regs)-1]
	c.addABC(RETURN, lastReg.N, 2, 0)

	cl := types.NewScmClosure()
	cl.Proto = c.Proto
	return cl, nil
}
