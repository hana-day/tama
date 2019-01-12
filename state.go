package tama

import (
	"fmt"
	"github.com/hyusuk/tama/compiler"
	"github.com/hyusuk/tama/parser"
)

var DefaultStackSize = 256 * 20

type CallInfo struct {
	Base int
	Cl   *compiler.Closure
	Top  int
}

type State struct {
	// call stack
	CallStack *Stack
	CallInfos *Stack
	Base      int
}

func NewState() *State {
	return &State{
		CallStack: NewStack(DefaultStackSize),
		CallInfos: NewStack(DefaultStackSize),
	}
}

func (s *State) LoadString(source string) (*compiler.Closure, error) {
	p := &parser.Parser{}
	p.Init([]byte(source))
	f := p.ParseFile()
	return compiler.Compile(f.Exprs)
}

func (s *State) call(nargs int) error {
	clIndex := s.CallStack.Sp() - nargs - 1
	s.Base = clIndex + 1
	cl, ok := s.CallStack.Get(clIndex).(*compiler.Closure)
	if !ok {
		return fmt.Errorf("Function is not loaded")
	}
	ci := &CallInfo{Cl: cl, Base: s.Base, Top: cl.Proto.MaxStackSize + s.Base}
	s.CallInfos.Push(ci)
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
