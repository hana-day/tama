package tama

import (
	"fmt"
	"github.com/hyusuk/tama/compiler"
	"github.com/hyusuk/tama/parser"
	"github.com/hyusuk/tama/types"
)

var DefaultStackSize = 256 * 20

type CallInfo struct {
	Base int
	Cl   *types.Closure
}

type State struct {
	// call stack
	CallStack *types.Stack
	CallInfos *types.Stack
	Base      int
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

func (s *State) precall(clIndex int) error {
	s.Base = clIndex + 1
	cl, ok := s.CallStack.Get(clIndex).(*types.Closure)
	if !ok {
		return fmt.Errorf("Function is not loaded")
	}
	if cl.IsGo {
		ci := &CallInfo{Base: s.Base}
		s.CallInfos.Push(ci)
		nargs := types.Number(s.CallStack.Sp() - clIndex - 1)
		s.CallStack.Push(nargs)

		fn, ok := cl.Fn.(func(s *State))
		if !ok {
			return fmt.Errorf("invalid function %v", cl.Fn)
		}
		fn(s)
		s.postcall(clIndex)
		return nil
	} else {
		ci := &CallInfo{Cl: cl, Base: s.Base}
		s.CallInfos.Push(ci)
		return nil
	}
}

func (s *State) postcall(resultSp int) {
	_ = s.CallInfos.Pop() // pop current call info
	prevCi := s.CallInfos.Top().(*CallInfo)
	s.Base = prevCi.Base
	result := s.CallStack.Pop()
	s.CallStack.Set(resultSp, result)
}

func (s *State) call(nargs int) error {
	clIndex := s.CallStack.Sp() - nargs - 1
	if err := s.precall(clIndex); err != nil {
		return err
	}
	runVM(s)
	return nil
}

func (s *State) ExecString(source string) error {
	cl, err := s.LoadString(source)
	if err != nil {
		return err
	}
	s.CallStack.Push(cl)
	return s.call(0)
}
