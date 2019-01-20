package tama

import (
	"fmt"
	"github.com/hyusuk/tama/compiler"
	"github.com/hyusuk/tama/parser"
	"github.com/hyusuk/tama/types"
)

var DefaultStackSize = 256 * 20

type CallInfo struct {
	Func int // function sp
	Base int // local sp
	Cl   *types.Closure
	Pc   int
}

type State struct {
	// call stack
	CallStack *types.Stack
	CallInfos *types.Stack
	Global    map[string]types.Object
}

func NewState() *State {
	s := &State{
		CallStack: types.NewStack(DefaultStackSize),
		CallInfos: types.NewStack(DefaultStackSize),
		Global:    map[string]types.Object{},
	}
	s.OpenBase()
	return s
}

func (s *State) LoadString(source string) (*types.Closure, error) {
	p := &parser.Parser{}
	p.Init([]byte(source))
	f, err := p.ParseFile()
	if err != nil {
		return nil, err
	}
	return compiler.Compile(f.Objs)
}

func (s *State) precall(clIndex int) (*CallInfo, error) {
	cl, ok := s.CallStack.Get(clIndex).(*types.Closure)
	if !ok {
		return nil, fmt.Errorf("function is not loaded")
	}
	if cl.IsGo {
		ci := &CallInfo{Cl: cl, Base: clIndex + 1, Func: clIndex}
		s.CallInfos.Push(ci)
		nargs := types.Number(s.CallStack.Sp() - clIndex - 1)
		s.CallStack.Push(nargs)

		fn, ok := cl.Fn.(func(s *State))
		if !ok {
			return nil, fmt.Errorf("invalid function %v", cl.Fn)
		}
		fn(s)
		s.postcall(s.CallStack.Sp())
		return ci, nil
	} else {
		ci := &CallInfo{Cl: cl, Base: clIndex + 1, Func: clIndex}
		s.CallInfos.Push(ci)
		return ci, nil
	}
}

func (s *State) postcall(resultSp int) {
	curCi := s.CallInfos.Pop().(*CallInfo) // pop current call info
	result := s.CallStack.Get(resultSp)
	s.CallStack.Set(curCi.Func, result)
	s.CallStack.SetSp(curCi.Func)
}

func (s *State) call(nargs int) error {
	clIndex := s.CallStack.Sp() - nargs
	if _, err := s.precall(clIndex); err != nil {
		return err
	}
	return runVM(s, true)
}

func (s *State) ExecString(source string) error {
	cl, err := s.LoadString(source)
	if err != nil {
		return err
	}
	s.CallStack.Push(cl)
	return s.call(0)
}
