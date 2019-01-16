package compiler

import (
	"fmt"
	"github.com/hyusuk/tama/types"
)

type Compiler struct {
	Proto *types.ClosureProto // current function header
	nreg  int                 // number of registers
}

type Reg struct {
	N int // register number
}

func (c *Compiler) error(msg string) error {
	return fmt.Errorf("compiler: %s", msg)
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
func (c *Compiler) compileDefine(pair *types.Pair) (*Reg, error) {
	errobj := c.error(fmt.Sprintf("invalid define"))
	cdar, err := types.Cdar(pair)
	if err != nil {
		return nil, errobj
	}
	cddr, err := types.Cddr(pair)
	if err != nil {
		return nil, errobj
	}
	rest, err := types.Car(cddr)
	if err != nil {
		return nil, errobj
	}
	sym, ok := cdar.(*types.Symbol)
	if !ok {
		return nil, c.error("The first argument of define must be a symbol")
	}
	nameReg := c.newReg()
	c.addABx(LOADK, nameReg.N, c.constIndex(sym.Name))
	valueReg, err := c.compileObject(rest)
	if err != nil {
		return nil, err
	}
	c.addABx(SETGLOBAL, valueReg.N, nameReg.N)
	return valueReg, nil
}

func (c *Compiler) compilePair(pair *types.Pair) (*Reg, error) {
	v, err := types.Car(pair)
	if err != nil {
		return nil, err
	}
	first, ok := v.(*types.Symbol)
	if !ok {
		return nil, c.error("invalid procedure/define name")
	}
	switch first.Name {
	case "define":
		return c.compileDefine(pair)
	default:
		r1 := c.compileSymbol(first)
		cdr, err := types.Cdr(pair)
		if err != nil {
			return nil, err
		}
		args := cdr.(*types.Pair)
		argsArr := args.ListToArray()
		argRegs := make([]*Reg, len(argsArr))
		for i := 0; i < len(argsArr); i++ {
			argRegs[i] = c.newReg()
		}
		for i, arg := range argsArr {
			r, err := c.compileObject(arg)
			if err != nil {
				return nil, err
			}
			c.addABC(MOVE, argRegs[i].N, r.N, 0)
		}
		// Always return one value
		c.addABC(CALL, r1.N, 1+len(argsArr), 2)
		return r1, nil
	}
}

func (c *Compiler) compileObject(obj types.Object) (*Reg, error) {
	switch o := obj.(type) {
	case types.Number:
		return c.compileNumber(o), nil
	case *types.Symbol:
		return c.compileSymbol(o), nil
	case *types.Pair:
		return c.compilePair(o)
	default:
		return nil, c.error(fmt.Sprintf("Unknown type of object %v", o))
	}
}

func (c *Compiler) compileObjects(objs []types.Object) ([]*Reg, error) {
	regs := make([]*Reg, len(objs))
	for i, obj := range objs {
		reg, err := c.compileObject(obj)
		if err != nil {
			return regs, err
		}
		regs[i] = reg
	}
	return regs, nil
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
	regs, err := c.compileObjects(objs)
	if err != nil {
		return nil, err
	}
	lastReg := regs[len(regs)-1]
	c.addABC(RETURN, lastReg.N, 2, 0)

	cl := types.NewScmClosure()
	cl.Proto = c.Proto
	return cl, nil
}
