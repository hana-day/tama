package compiler

import (
	"fmt"
	"github.com/hyusuk/tama/types"
)

type Compiler struct {
}

type Reg struct {
	N int // register number
}

func newReg(n int) *Reg {
	return &Reg{N: n}
}

type FuncState struct {
	Proto *types.ClosureProto // current function header
	nreg  int                 // number of registers
	prev  *FuncState          // enclosing function
}

func newFuncState(prev *FuncState) *FuncState {
	return &FuncState{
		Proto: types.NewClosureProto(),
		nreg:  0,
		prev:  prev,
	}
}

func (c *Compiler) error(msg string) error {
	return fmt.Errorf("compiler: %s", msg)
}

func (fs *FuncState) newReg() *Reg {
	reg := newReg(fs.nreg)
	fs.nreg++
	return reg
}

func (fs *FuncState) add(inst uint32) {
	fs.Proto.Insts = append(fs.Proto.Insts, inst)
}

func (fs *FuncState) addABx(op int, a int, bx int) {
	fs.add(CreateABx(op, a, bx))
}

func (fs *FuncState) addABC(op int, a int, b int, c int) {
	fs.add(CreateABC(op, a, b, c))
}

func (fs *FuncState) constIndex(v types.Object) int {
	for i, cs := range fs.Proto.Consts {
		if cs == v {
			return i
		}
	}
	fs.Proto.Consts = append(fs.Proto.Consts, v)
	return len(fs.Proto.Consts) - 1
}

func (fs *FuncState) bindLocVar(sym *types.Symbol) {
	index := fs.nreg
	v := &types.LocVar{Name: sym.Name, Index: index}
	fs.Proto.LocVars[sym.Name] = v
	fs.nreg++
}

func (fs *FuncState) findLocVar(sym *types.Symbol) *types.LocVar {
	loc, ok := fs.Proto.LocVars[sym.Name]
	if !ok {
		return nil
	}
	return loc
}

func (c *Compiler) compileNumber(fs *FuncState, num types.Number) *Reg {
	reg := fs.newReg()
	fs.addABx(LOADK, reg.N, fs.constIndex(num))
	return reg
}

func (c *Compiler) compileSymbol(fs *FuncState, sym *types.Symbol) *Reg {
	if loc := fs.findLocVar(sym); loc != nil {
		return newReg(loc.Index)
	}
	r1 := fs.newReg()
	fs.addABx(GETGLOBAL, r1.N, fs.constIndex(types.String(sym.Name)))
	return r1
}

// Compile define syntax
//
// (define a 1)
func (c *Compiler) compileDefine(fs *FuncState, pair *types.Pair) (*Reg, error) {
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
	nameReg := fs.newReg()
	fs.addABx(LOADK, nameReg.N, fs.constIndex(sym.Name))
	valueReg, err := c.compileObject(fs, rest)
	if err != nil {
		return nil, err
	}
	fs.addABx(SETGLOBAL, valueReg.N, nameReg.N)
	return valueReg, nil
}

// (lambda (x y) ...)
func (c *Compiler) compileLambda(fs *FuncState, pair *types.Pair) (*Reg, error) {
	if pair.Len() < 3 {
		return nil, c.error(fmt.Sprintf("invalid lambda %s", pair.String()))
	}
	cdar, _ := types.Cdar(pair)
	args, ok := cdar.(types.ArrayableObject)
	if !ok {
		return nil, c.error(fmt.Sprintf("invalid lambda %s", pair.String()))
	}
	child := newFuncState(fs)
	argsArr := args.Array()
	argSyms := make([]*types.Symbol, len(argsArr))
	for i, arg := range argsArr {
		sym, ok := arg.(*types.Symbol)
		if !ok {
			return nil, c.error(fmt.Sprintf("invalid lambda %s", pair.String()))
		}
		argSyms[i] = sym
	}
	child.Proto.Args = argSyms
	for _, arg := range child.Proto.Args {
		child.bindLocVar(arg)
	}
	cddr, _ := types.Cddr(pair)
	body, _ := types.Car(cddr)
	resultReg, err := c.compileObject(child, body)
	if err != nil {
		return nil, err
	}
	child.addABC(RETURN, resultReg.N, 2, 0)

	protoIndex := len(fs.Proto.Protos)
	fs.Proto.Protos = append(fs.Proto.Protos, child.Proto)
	r := fs.newReg()
	fs.addABx(CLOSURE, r.N, protoIndex)
	return r, nil
}

func (c *Compiler) compileCall(fs *FuncState, proc *Reg, args types.ArrayableObject) (*Reg, error) {
	argsArr := args.Array()
	argRegs := make([]*Reg, len(argsArr))
	for i := 0; i < len(argsArr); i++ {
		argRegs[i] = fs.newReg()
	}
	for i, arg := range argsArr {
		r, err := c.compileObject(fs, arg)
		if err != nil {
			return nil, err
		}
		fs.addABC(MOVE, argRegs[i].N, r.N, 0)
	}
	// Always return one value
	fs.addABC(CALL, proc.N, 1+len(argsArr), 2)
	return proc, nil
}

func (c *Compiler) compilePair(fs *FuncState, pair *types.Pair) (*Reg, error) {
	if pair.Len() == 0 {
		return nil, c.error(fmt.Sprintf("invalid syntax %s", pair.String()))
	}
	cdr, _ := types.Cdr(pair)
	args, ok := cdr.(types.ArrayableObject)
	if !ok {
		return nil, c.error(fmt.Sprintf("invalid syntax %s", pair.String()))
	}
	v, _ := types.Car(pair)
	switch first := v.(type) {
	case *types.Symbol:
		switch first.Name {
		case "define":
			return c.compileDefine(fs, pair)
		case "lambda":
			return c.compileLambda(fs, pair)
		default:
			proc := c.compileSymbol(fs, first)
			return c.compileCall(fs, proc, args)
		}
	case *types.Pair:
		proc, err := c.compilePair(fs, first)
		if err != nil {
			return nil, err
		}
		return c.compileCall(fs, proc, args)
	}
	return nil, c.error(fmt.Sprintf("invalid procedure name %v", v))
}

func (c *Compiler) compileObject(fs *FuncState, obj types.Object) (*Reg, error) {
	switch o := obj.(type) {
	case types.Number:
		return c.compileNumber(fs, o), nil
	case *types.Symbol:
		return c.compileSymbol(fs, o), nil
	case *types.Pair:
		return c.compilePair(fs, o)
	default:
		return nil, c.error(fmt.Sprintf("Unknown type of object %v", o))
	}
}

func (c *Compiler) compileObjects(fs *FuncState, objs []types.Object) ([]*Reg, error) {
	regs := make([]*Reg, len(objs))
	for i, obj := range objs {
		reg, err := c.compileObject(fs, obj)
		if err != nil {
			return regs, err
		}
		regs[i] = reg
	}
	return regs, nil
}

func Compile(objs []types.Object) (*types.Closure, error) {
	c := Compiler{}
	fs := newFuncState(nil)
	regs, err := c.compileObjects(fs, objs)
	if err != nil {
		return nil, err
	}
	lastReg := regs[len(regs)-1]
	fs.addABC(RETURN, lastReg.N, 2, 0)

	cl := types.NewScmClosure()
	cl.Proto = fs.Proto
	return cl, nil
}
