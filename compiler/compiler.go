package compiler

import (
	"fmt"
	"github.com/hyusuk/tama/types"
)

type Compiler struct {
}

type varType int

const (
	varGlobal varType = iota
	varUpValue
	varLocVar
)

type reg struct {
	n int // register number
}

type nameStorage struct {
	names    []types.String
	len      int
	capacity int
}

func newNameStorage(cap int) *nameStorage {
	return &nameStorage{
		names:    make([]types.String, cap),
		len:      0,
		capacity: cap,
	}
}

func (ns *nameStorage) Len() int {
	return ns.len
}

func (ns *nameStorage) Capacity() int {
	return ns.capacity
}

func (ns *nameStorage) Name(index int) types.String {
	return ns.names[index]
}

func (ns *nameStorage) Find(name types.String) int {
	for i, nm := range ns.names {
		if nm == name {
			return i
		}
	}
	return -1
}

func (ns *nameStorage) grow() {
	if ns.len >= ns.capacity {
		ns.capacity = (ns.capacity + 1) * 2
		newOne := make([]types.String, ns.capacity)
		copy(newOne, ns.names)
		ns.names = newOne
	}
}

func (ns *nameStorage) Register(name types.String) int {
	i := ns.Find(name)
	if i >= 0 {
		return i
	}
	ns.grow()
	l := ns.len
	ns.names[l] = name
	ns.len++
	return l
}

type funcState struct {
	proto         *types.ClosureProto // current function header
	nreg          int                 // number of registers
	prev          *funcState          // enclosing function
	locVars       *nameStorage
	upVals        *nameStorage
	closeRequired bool // whether upvalues inside the function need to be closed or
}

func newFuncState(prev *funcState) *funcState {
	return &funcState{
		proto:         types.NewClosureProto(),
		nreg:          0,
		prev:          prev,
		locVars:       newNameStorage(16),
		upVals:        newNameStorage(16),
		closeRequired: false,
	}
}

func (fs *funcState) newReg() *reg {
	r := &reg{n: fs.nreg}
	fs.nreg++
	return r
}

func (fs *funcState) add(inst uint32) {
	fs.proto.Insts = append(fs.proto.Insts, inst)
}

func (fs *funcState) addABx(op int, a int, bx int) {
	fs.add(CreateABx(op, a, bx))
}

func (fs *funcState) addABC(op int, a int, b int, c int) {
	fs.add(CreateABC(op, a, b, c))
}

func (fs *funcState) constIndex(v types.Object) int {
	for i, cs := range fs.proto.Consts {
		if cs == v {
			return i
		}
	}
	fs.proto.Consts = append(fs.proto.Consts, v)
	return len(fs.proto.Consts) - 1
}

func (fs *funcState) bindLocVar(name types.String) int {
	fs.locVars.Register(name)
	fs.nreg++
	return fs.nreg - 1
}

func (fs *funcState) findLocVar(name types.String) int {
	return fs.locVars.Find(name)
}

func (fs *funcState) upValueIndex(name types.String) int {
	i := fs.upVals.Find(name)
	if i < 0 {
		return fs.upVals.Register(name)
	}
	return i
}

func (fs *funcState) getVarType(sym *types.Symbol) varType {
	for cur := fs; cur != nil; cur = cur.prev {
		if index := cur.findLocVar(sym.Name); index > -1 {
			if cur == fs {
				return varLocVar
			}
			return varUpValue
		}
	}
	return varGlobal
}

func (c *Compiler) error(format string, a ...interface{}) error {
	return fmt.Errorf("compiler: %s", fmt.Sprintf(format, a...))
}

func (c *Compiler) compileNumber(fs *funcState, num types.Number) *reg {
	r := fs.newReg()
	fs.addABx(OP_LOADK, r.n, fs.constIndex(num))
	return r
}

func (c *Compiler) compileSymbol(fs *funcState, sym *types.Symbol) *reg {
	switch fs.getVarType(sym) {
	case varLocVar:
		index := fs.findLocVar(sym.Name)
		return &reg{n: index}
	case varGlobal:
		r := fs.newReg()
		fs.addABx(OP_GETGLOBAL, r.n, fs.constIndex(sym.Name))
		return r
	case varUpValue:
		r := fs.newReg()
		fs.addABC(OP_GETUPVAL, r.n, fs.upValueIndex(sym.Name), 0)
		return r
	default:
		return nil
	}
}

func (c *Compiler) compileGlobalAssign(fs *funcState, varname *types.Symbol, value types.Object) (*reg, error) {
	nameR := fs.newReg()
	fs.addABx(OP_LOADK, nameR.n, fs.constIndex(varname.Name))
	valueR, err := c.compileObject(fs, value)
	if err != nil {
		return nil, err
	}
	fs.addABx(OP_SETGLOBAL, valueR.n, nameR.n)
	return valueR, nil
}

// Compile define syntax
//
// (define a 1)
func (c *Compiler) compileDefine(fs *funcState, pair *types.Pair) (*reg, error) {
	if pair.Len() != 3 {
		return nil, c.error("define: invalid syntax")
	}
	cdar, _ := types.Cdar(pair)
	varname, ok := cdar.(*types.Symbol)
	if !ok {
		return nil, c.error("define: invalid syntax")
	}
	cddr, _ := types.Cddr(pair)
	expr, _ := types.Car(cddr)
	r, err := c.compileGlobalAssign(fs, varname, expr)
	if err != nil {
		return nil, err
	}
	return r, nil
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
	child.proto.Args = argSyms
	for _, arg := range child.proto.Args {
		child.bindLocVar(arg.Name)
	}
	cddr, _ := types.Cddr(pair)
	body := types.Cons(&types.Symbol{Name: "begin"}, cddr)
	resultR, err := c.compileObject(child, body)
	if err != nil {
		return nil, err
	}
	if child.closeRequired {
		child.addABC(OP_CLOSE, child.locVars.Len()-1, 0, 0)
	}
	child.addABC(OP_RETURN, resultR.n, 2, 0)

	child.proto.NUpVals = child.upVals.Len()
	protoIndex := len(fs.proto.Protos)
	fs.proto.Protos = append(fs.proto.Protos, child.proto)
	r := fs.newReg()
	fs.addABx(OP_CLOSURE, r.n, protoIndex)

	for i := 0; i < child.upVals.Len(); i++ {
		uvName := child.upVals.Name(i)
		locIndex := fs.findLocVar(uvName)
		if locIndex > -1 {
			fs.addABC(OP_MOVE, 0, locIndex, 0)
			fs.closeRequired = true
			continue
		}
		uvIndex := fs.upValueIndex(uvName)
		if uvIndex < 0 {
			uvIndex = fs.upVals.Register(uvName)
		}
		fs.addABC(OP_GETUPVAL, 0, uvIndex, 0)
	}
	return r, nil
}

func (c *Compiler) compileBegin(fs *funcState, pair *types.Pair) (*reg, error) {
	if pair.Len() < 2 {
		return nil, c.error("invalid begin %s", pair.String())
	}
	cdr, _ := types.Cdr(pair)
	exprs, err := cdr.(*types.Pair).Slice()
	if err != nil {
		return nil, err
	}
	regs, err := c.compileObjects(fs, exprs)
	if err != nil {
		return nil, err
	}
	return regs[len(regs)-1], nil
}

func (c *Compiler) compileSet(fs *funcState, pair *types.Pair) (*reg, error) {
	if pair.Len() != 3 {
		return nil, c.error("set!: invalid syntax")
	}
	cdar, _ := types.Cdar(pair)
	varname, ok := cdar.(*types.Symbol)
	cddr, _ := types.Cddr(pair)
	expr, _ := types.Car(cddr)
	if !ok {
		return nil, c.error("set!: invalid syntax")
	}
	switch fs.getVarType(varname) {
	case varLocVar:
		return nil, nil
	case varGlobal:
		return c.compileGlobalAssign(fs, varname, expr)
	case varUpValue:
		return nil, nil
	}
	return nil, c.error("set!: unsupported var type")
}

func (c *Compiler) compileCall(fs *funcState, proc types.Object, args types.SlicableObject) (*reg, error) {
	procR, err := c.compileObject(fs, proc)
	if err != nil {
		return nil, err
	}
	argsArr, err := args.Slice()
	if err != nil {
		return nil, err
	}

	// Arrange closure and arguments registers in the order.
	// TODO: Too verbose?
	newProcR := fs.newReg()
	fs.addABC(OP_MOVE, newProcR.n, procR.n, 0)
	regs := make([]*reg, len(argsArr))
	for i, _ := range argsArr {
		regs[i] = fs.newReg()
	}

	for i, arg := range argsArr {
		r, err := c.compileObject(fs, arg)
		if err != nil {
			return nil, err
		}
		fs.addABC(OP_MOVE, regs[i].n, r.n, 0)
	}
	// Always return one value
	fs.addABC(OP_CALL, newProcR.n, 1+len(argsArr), 2)
	return newProcR, nil
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
		case "begin":
			return c.compileBegin(fs, pair)
		case "set!":
			return c.compileSet(fs, pair)
		default: // (procedure-name args...)
			return c.compileCall(fs, first, args)
		}
	case *types.Pair: // ((procedure-name args...) args...)
		return c.compileCall(fs, first, args)
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
	fs.addABC(OP_RETURN, lastR.n, 2, 0)

	cl := types.NewScmClosure(fs.proto, 0)
	return cl, nil
}
