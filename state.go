package tama

import (
	"fmt"
	"github.com/hyusuk/tama/parser"
)

var DefaultStackSize = 256 * 20

type CallInfo struct {
	Base int
	Cl   *Closure
}

type State struct {
	// call stack
	CallStack *Stack
	CallInfos *Stack
	Base      int
	Global    map[string]Value
}

func NewState() *State {
	return &State{
		CallStack: NewStack(DefaultStackSize),
		CallInfos: NewStack(DefaultStackSize),
		Global:    map[string]Value{},
	}
}

func (s *State) LoadString(source string) (*Closure, error) {
	p := &parser.Parser{}
	p.Init([]byte(source))
	f := p.ParseFile()
	return Compile(f.Exprs)
}

func (s *State) precall(clIndex int) error {
	s.Base = clIndex + 1
	cl, ok := s.CallStack.Get(clIndex).(*Closure)
	if !ok {
		return fmt.Errorf("Function is not loaded")
	}
	if cl.isGo {
		ci := &CallInfo{Base: s.Base}
		s.CallInfos.Push(ci)
		nargs := Number(s.CallStack.Sp() - clIndex - 1)
		s.CallStack.Push(nargs)

		cl.Fn(s)
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
	prevCi := s.CallInfos.Pop().(*CallInfo)
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
