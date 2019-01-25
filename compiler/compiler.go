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

func (fs *funcState) addASbx(op int, a int, sbx int) {
	fs.add(CreateASbx(op, a, sbx))
}

func (fs *funcState) rewriteSbx(pc int, sbx int) {
	SetArgSbx(&fs.proto.Insts[pc], sbx)
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

func (fs *funcState) nextPc() int {
	return len(fs.proto.Insts)
}

func (c *Compiler) error(format string, a ...interface{}) error {
	return fmt.Errorf("compiler: %s", fmt.Sprintf(format, a...))
}

func (c *Compiler) compileNumber(fs *funcState, num types.Number) *reg {
	r := fs.newReg()
	fs.addABx(OP_LOADK, r.n, fs.constIndex(num))
	return r
}

func (c *Compiler) compileBoolean(fs *funcState, bool types.Boolean) *reg {
	r := fs.newReg()
	fs.addABx(OP_LOADK, r.n, fs.constIndex(bool))
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
	valueR, err := c.compileObject(fs, value)
	if err != nil {
		return nil, err
	}
	fs.addABx(OP_SETGLOBAL, valueR.n, fs.constIndex(varname.Name))
	return valueR, nil
}

// Compile define syntax
//
// (define variable expression)
// (define (variable formals) body)
// (define (variable . formal) body)
func (c *Compiler) compileDefine(fs *funcState, args []types.Object) (*reg, error) {
	if len(args) < 2 {
		return nil, c.error("define: invalid syntax")
	}

	var varname *types.Symbol
	var expr types.Object

	switch first := args[0].(type) {
	case *types.Symbol: // (define variable expression)
		varname = first
		expr = args[1]
	case *types.Pair: // (define (variable formals) body)
		sym, ok := first.Car().(*types.Symbol)
		if !ok {
			return nil, c.error("define: invalid syntax")
		}
		varname = sym
		lambdaExpr := []types.Object{types.NewSymbol("lambda")}
		if first.Len() == 3 {
			second, _ := first.Second()
			sym, ok := second.(*types.Symbol)
			if !ok {
				return nil, c.error("define: invalid syntax")
			}
			// convert
			// (define (variable . formal) body)
			// =>
			// (define variable
			//   (lambda formal body))
			if sym.Name == "." {
				cddar, _ := first.Cddar()
				lambdaExpr = append(lambdaExpr, cddar)
				lambdaExpr = append(lambdaExpr, args[1:]...)
				expr = types.List(lambdaExpr...)
				break
			}
		}
		// convert
		// (define (variable formals) body)
		// =>
		// (define variable
		//   (lambda (formals) body))
		lambdaExpr = append(lambdaExpr, first.Cdr())
		lambdaExpr = append(lambdaExpr, args[1:]...)
		expr = types.List(lambdaExpr...)
	default:
		return nil, c.error("define: invalid syntax")
	}
	_, err := c.compileGlobalAssign(fs, varname, expr)
	if err != nil {
		return nil, err
	}
	r := fs.newReg()
	fs.addABC(OP_LOADUNDEF, r.n, r.n, 0)
	return r, nil
}

func (c *Compiler) lambdaForm(formals types.Object) ([]*types.Symbol, types.ArgMode, error) {
	var argSyms []*types.Symbol

	var mode types.ArgMode
	switch args := formals.(type) {
	case *types.Nil:
		mode = types.FixedArgMode
		argSyms = []*types.Symbol{}
	case *types.Symbol:
		mode = types.VArgMode
		argSyms = []*types.Symbol{args}
	case *types.Pair:
		argsArr, err := args.Slice()
		if err != nil {
			return nil, 0, err
		}
		numDots := 0
		argSyms = make([]*types.Symbol, len(argsArr))
		for i, arg := range argsArr {
			sym, ok := arg.(*types.Symbol)
			if !ok {
				return nil, 0, c.error("lambda: invalid syntax")
			}
			argSyms[i] = sym
		}
		for _, sym := range argSyms {
			if sym.Name == "." {
				numDots++
			}
		}
		switch numDots {
		case 0:
			mode = types.FixedArgMode
		case 1:
			mode = types.RestArgMode
			if len(argsArr) < 3 || argSyms[len(argSyms)-2].Name != "." {
				return nil, 0, c.error("lambda: invalid syntax")
			}
			argSyms = append(argSyms[:len(argSyms)-2], argSyms[len(argSyms)-1:]...)
		case 2:
			return nil, 0, c.error("lambda: invalid syntax")
		}
	default:
		return nil, 0, c.error("lambda: invalid syntax")
	}
	return argSyms, mode, nil
}

/// compileLambda compiles lambda syntax.
//
// (lambda (x y) ...)
// (lambda args ...)
// (lambda (x y . rest) ...)
func (c *Compiler) compileLambda(fs *funcState, lambdaArgs []types.Object) (*reg, error) {
	if len(lambdaArgs) < 2 {
		return nil, c.error("lambda: invalid syntax")
	}
	argSyms, mode, err := c.lambdaForm(lambdaArgs[0])
	if err != nil {
		return nil, err
	}

	child := newFuncState(fs)
	child.proto.Args = argSyms
	child.proto.Mode = mode

	for _, arg := range child.proto.Args {
		child.bindLocVar(arg.Name)
	}
	resultR, err := c.compileBegin(child, lambdaArgs[1:])
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

func (c *Compiler) compileBegin(fs *funcState, args []types.Object) (*reg, error) {
	if len(args) == 0 {
		return nil, c.error("begin: invalid syntax")
	}
	regs, err := c.compileObjects(fs, args)
	if err != nil {
		return nil, err
	}
	return regs[len(regs)-1], nil
}

func (c *Compiler) compileSet(fs *funcState, args []types.Object) (*reg, error) {
	if len(args) != 2 {
		return nil, c.error("set!: invalid syntax")
	}
	varname, ok := args[0].(*types.Symbol)
	if !ok {
		return nil, c.error("set!: invalid syntax")
	}
	expr := args[1]
	switch fs.getVarType(varname) {
	case varLocVar:
		index := fs.findLocVar(varname.Name)
		if index < 0 {
			index = fs.bindLocVar(varname.Name)
		}
		valueR, err := c.compileObject(fs, expr)
		if err != nil {
			return nil, err
		}
		fs.addABC(OP_MOVE, index, valueR.n, 0)
	case varGlobal:
		if _, err := c.compileGlobalAssign(fs, varname, expr); err != nil {
			return nil, err
		}
	case varUpValue:
		valueR, err := c.compileObject(fs, expr)
		if err != nil {
			return nil, err
		}
		fs.addABC(OP_SETUPVAL, valueR.n, fs.upValueIndex(varname.Name), 0)
	default:
		return nil, c.error("set!: unsupported var type")
	}
	r := fs.newReg()
	fs.addABC(OP_LOADUNDEF, r.n, r.n, 0)
	return r, nil
}

func (c *Compiler) compileQuote(fs *funcState, argsArr []types.Object) (*reg, error) {
	if len(argsArr) != 1 {
		return nil, fmt.Errorf("quote: invalid syntax")
	}
	r := fs.newReg()
	fs.addABx(OP_LOADK, r.n, fs.constIndex(argsArr[0]))
	return r, nil
}

func (c *Compiler) compileIf(fs *funcState, argsArr []types.Object) (*reg, error) {
	if len(argsArr) != 2 && len(argsArr) != 3 {
		return nil, fmt.Errorf("if: invalid syntax")
	}
	elseExists := len(argsArr) == 3

	testR, err := c.compileObject(fs, argsArr[0])
	if err != nil {
		return nil, err
	}
	resultR := fs.newReg()

	fs.addABC(OP_TEST, testR.n, 0, 1)
	thenJmpPc := fs.nextPc()
	fs.addASbx(OP_JMP, 0, 0) // jump to then expr. sbx will be set later
	elseJmpPc := fs.nextPc()
	fs.addASbx(OP_JMP, 0, 0) // jump to else expr. sbx will be set later

	thenPc := fs.nextPc()
	thenR, err := c.compileObject(fs, argsArr[1])
	fs.addABC(OP_MOVE, resultR.n, thenR.n, 0)
	lastJmpPc := fs.nextPc()
	fs.addASbx(OP_JMP, 0, 0) // jump to last expr. sbx will be set later.

	elsePc := fs.nextPc()
	if elseExists { // (if test consequent alternate)
		elseR, err := c.compileObject(fs, argsArr[2])
		if err != nil {
			return nil, err
		}
		fs.addABC(OP_MOVE, resultR.n, elseR.n, 0)
	} else { // (if test consequent)
		fs.addABC(OP_LOADUNDEF, resultR.n, resultR.n, 0)
	}
	lastPc := fs.nextPc()
	fs.rewriteSbx(thenJmpPc, thenPc-thenJmpPc-1)
	fs.rewriteSbx(elseJmpPc, elsePc-elseJmpPc-1)
	fs.rewriteSbx(lastJmpPc, lastPc-lastJmpPc-1)

	return resultR, nil
}

func (c *Compiler) compileCall(fs *funcState, proc types.Object, args []types.Object) (*reg, error) {
	procR, err := c.compileObject(fs, proc)
	if err != nil {
		return nil, err
	}

	// Arrange closure and arguments registers in the order.
	// TODO: Too verbose?
	newProcR := fs.newReg()
	fs.addABC(OP_MOVE, newProcR.n, procR.n, 0)
	regs := make([]*reg, len(args))
	for i, _ := range args {
		regs[i] = fs.newReg()
	}

	for i, arg := range args {
		r, err := c.compileObject(fs, arg)
		if err != nil {
			return nil, err
		}
		fs.addABC(OP_MOVE, regs[i].n, r.n, 0)
	}
	// Always return one value
	fs.addABC(OP_CALL, newProcR.n, 1+len(args), 2)
	return newProcR, nil
}

func (c *Compiler) compilePair(fs *funcState, pair *types.Pair) (*reg, error) {
	if pair.Len() == 0 {
		return nil, c.error("invalid syntax %s", pair.String())
	}
	cdr := pair.Cdr()
	args, ok := cdr.(types.SlicableObject)
	if !ok {
		return nil, c.error("invalid syntax %s", pair.String())
	}
	argsArr, err := args.Slice()
	if err != nil {
		return nil, err
	}
	v := pair.Car()
	switch first := v.(type) {
	case *types.Symbol:
		switch first.Name {
		case "define":
			return c.compileDefine(fs, argsArr)
		case "lambda":
			return c.compileLambda(fs, argsArr)
		case "begin":
			return c.compileBegin(fs, argsArr)
		case "set!":
			return c.compileSet(fs, argsArr)
		case "quote":
			return c.compileQuote(fs, argsArr)
		case "if":
			return c.compileIf(fs, argsArr)
		default: // (procedure-name args...)
			return c.compileCall(fs, first, argsArr)
		}
	case *types.Pair: // ((procedure-name args...) args...)
		return c.compileCall(fs, first, argsArr)
	}
	return nil, c.error("invalid procedure name %v", v)
}

func (c *Compiler) compileObject(fs *funcState, obj types.Object) (*reg, error) {
	switch o := obj.(type) {
	case types.Number:
		return c.compileNumber(fs, o), nil
	case types.Boolean:
		return c.compileBoolean(fs, o), nil
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

var DefaultSyntaxes map[string]*types.Syntax = map[string]*types.Syntax{
	// pseudo syntaxes
	"define": types.NewSyntax("define", nil),
	"lambda": types.NewSyntax("lambda", nil),
	"begin":  types.NewSyntax("begin", nil),
	"set!":   types.NewSyntax("set!", nil),
	"quote":  types.NewSyntax("quote", nil),
	"if":     types.NewSyntax("if", nil),
}
