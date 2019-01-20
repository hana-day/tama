package compiler

import (
	"fmt"
	"github.com/hyusuk/tama/types"
)

type Compiler struct {
}

type reg struct {
	N int // register number
}

type funcState struct {
	Proto *types.ClosureProto // current function header
	nreg  int                 // number of registers
	prev  *funcState          // enclosing function
}

func newFuncState(prev *funcState) *funcState {
	return &funcState{
		Proto: types.NewClosureProto(),
		nreg:  0,
		prev:  prev,
	}
}

func (fs *funcState) newReg() *reg {
	r := &reg{N: fs.nreg}
	fs.nreg++
	return r
}

func (fs *funcState) add(inst uint32) {
	fs.Proto.Insts = append(fs.Proto.Insts, inst)
}

func (fs *funcState) addABx(op int, a int, bx int) {
	fs.add(CreateABx(op, a, bx))
}

func (fs *funcState) addABC(op int, a int, b int, c int) {
	fs.add(CreateABC(op, a, b, c))
}

func (fs *funcState) constIndex(v types.Object) int {
	for i, cs := range fs.Proto.Consts {
		if cs == v {
			return i
		}
	}
	fs.Proto.Consts = append(fs.Proto.Consts, v)
	return len(fs.Proto.Consts) - 1
}

func (fs *funcState) bindLocVar(sym *types.Symbol) {
	index := fs.nreg
	v := &types.LocVar{Name: sym.Name, Index: index}
	fs.Proto.LocVars[sym.Name] = v
	fs.nreg++
}

func (fs *funcState) findLocVar(sym *types.Symbol) *types.LocVar {
	loc, ok := fs.Proto.LocVars[sym.Name]
	if !ok {
		return nil
	}
	return loc
}

func (c *Compiler) error(format string, a ...interface{}) error {
	return fmt.Errorf("compiler: %s", fmt.Sprintf(format, a...))
}

func (c *Compiler) compileNumber(fs *funcState, num types.Number) *reg {
	r := fs.newReg()
	fs.addABx(LOADK, r.N, fs.constIndex(num))
	return r
}

func (c *Compiler) compileSymbol(fs *funcState, sym *types.Symbol) *reg {
	if loc := fs.findLocVar(sym); loc != nil {
		return &reg{N: loc.Index}
	}
	r1 := fs.newReg()
	fs.addABx(GETGLOBAL, r1.N, fs.constIndex(types.String(sym.Name)))
	return r1
}

// Compile define syntax
//
// (define a 1)
func (c *Compiler) compileDefine(fs *funcState, pair *types.Pair) (*reg, error) {
	errobj := c.error("invalid define")
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
	nameR := fs.newReg()
	fs.addABx(LOADK, nameR.N, fs.constIndex(sym.Name))
	valueR, err := c.compileObject(fs, rest)
	if err != nil {
		return nil, err
	}
	fs.addABx(SETGLOBAL, valueR.N, nameR.N)
	return valueR, nil
}

// (lambda (x y) ...)
func (c *Compiler) compileLambda(fs *funcState, pair *types.Pair) (*reg, error) {
	if pair.Len() < 3 {
		return nil, c.error("invalid lambda %s", pair.String())
	}
	cdar, _ := types.Cdar(pair)
	args, ok := cdar.(types.SlicableObject)
	if !ok {
		return nil, c.error("invalid lambda %s", pair.String())
	}
	child := newFuncState(fs)
	argsArr, err := args.Slice()
	if err != nil {
		return nil, err
	}
	argSyms := make([]*types.Symbol, len(argsArr))
	for i, arg := range argsArr {
		sym, ok := arg.(*types.Symbol)
		if !ok {
			return nil, c.error("invalid lambda %s", pair.String())
		}
		argSyms[i] = sym
	}
	child.Proto.Args = argSyms
	for _, arg := range child.Proto.Args {
		child.bindLocVar(arg)
	}
	cddr, _ := types.Cddr(pair)
	body, _ := types.Car(cddr)
	resultR, err := c.compileObject(child, body)
	if err != nil {
		return nil, err
	}
	child.addABC(RETURN, resultR.N, 2, 0)

	protoIndex := len(fs.Proto.Protos)
	fs.Proto.Protos = append(fs.Proto.Protos, child.Proto)
	r := fs.newReg()
	fs.addABx(CLOSURE, r.N, protoIndex)
	return r, nil
}

func (c *Compiler) compileCall(fs *funcState, proc *reg, args types.SlicableObject) (*reg, error) {
	argsArr, err := args.Slice()
	if err != nil {
		return nil, err
	}
	for i, arg := range argsArr {
		r, err := c.compileObject(fs, arg)
		if proc.N+i+1 != r.N {
			fs.addABC(MOVE, proc.N+i+1, r.N, 0)
		}
		if err != nil {
			return nil, err
		}
	}
	// Always return one value
	fs.addABC(CALL, proc.N, 1+len(argsArr), 2)
	return proc, nil
}

func (c *Compiler) compilePair(fs *funcState, pair *types.Pair) (*reg, error) {
	if pair.Len() == 0 {
		return nil, c.error("invalid syntax %s", pair.String())
	}
	cdr, _ := types.Cdr(pair)
	args, ok := cdr.(types.SlicableObject)
	if !ok {
		return nil, c.error("invalid syntax %s", pair.String())
	}
	v, _ := types.Car(pair)
	switch first := v.(type) {
	case *types.Symbol:
		switch first.Name {
		case "define":
			return c.compileDefine(fs, pair)
		case "lambda":
			return c.compileLambda(fs, pair)
		default: // (procedure-name args...)
			proc := c.compileSymbol(fs, first)
			return c.compileCall(fs, proc, args)
		}
	case *types.Pair: // ((procedure-name args...) args...)
		proc, err := c.compilePair(fs, first)
		if err != nil {
			return nil, err
		}
		return c.compileCall(fs, proc, args)
	}
	return nil, c.error("invalid procedure name %v", v)
}

func (c *Compiler) compileObject(fs *funcState, obj types.Object) (*reg, error) {
	switch o := obj.(type) {
	case types.Number:
		return c.compileNumber(fs, o), nil
	case *types.Symbol:
		return c.compileSymbol(fs, o), nil
	case *types.Pair:
		return c.compilePair(fs, o)
	default:
		return nil, c.error("Unknown type of object %v", o)
	}
}

func (c *Compiler) compileObjects(fs *funcState, objs []types.Object) ([]*reg, error) {
	regs := make([]*reg, len(objs))
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
	lastR := regs[len(regs)-1]
	fs.addABC(RETURN, lastR.N, 2, 0)

	cl := types.NewScmClosure()
	cl.Proto = fs.Proto
	return cl, nil
}
