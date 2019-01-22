package tama

import (
	"fmt"
	"github.com/hyusuk/tama/compiler"
	"github.com/hyusuk/tama/parser"
	"github.com/hyusuk/tama/types"
)

var DefaultStackSize = 256 * 20

type State struct {
	// call stack
	CallStack *types.Stack
	CallInfos *types.Stack
	Global    map[string]types.Object
	uvhead    *types.UpValue
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

func (s *State) precall(clIndex int) (*types.CallInfo, error) {
	cl, ok := s.CallStack.Get(clIndex).(*types.Closure)
	if !ok {
		return nil, fmt.Errorf("function is not loaded")
	}
	if cl.IsGo {
		ci := &types.CallInfo{Cl: cl, Base: clIndex + 1, FuncSp: clIndex}
		s.CallInfos.Push(ci)
		nargs := types.Number(s.CallStack.Sp() - clIndex)
		s.CallStack.Push(nargs)

		fn, ok := cl.Fn.(func(s *State))
		if !ok {
			return nil, fmt.Errorf("invalid function %v", cl.Fn)
		}
		fn(s)
		s.postcall(s.CallStack.Sp())
		return ci, nil
	} else {
		ci := &types.CallInfo{Cl: cl, Base: clIndex + 1, FuncSp: clIndex}
		s.CallInfos.Push(ci)
		return ci, nil
	}
}

func (s *State) postcall(resultSp int) {
	curCi := s.CallInfos.Pop().(*types.CallInfo) // pop current call info
	result := s.CallStack.Get(resultSp)
	s.CallStack.Set(curCi.FuncSp, result)
	s.CallStack.SetSp(curCi.FuncSp)
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

func (s *State) findUpValue(level int) *types.UpValue {
	var prev *types.UpValue
	var next *types.UpValue
	if s.uvhead != nil {
		for uv := s.uvhead; uv != nil; uv = uv.Next {
			if uv.Index == level {
				return uv
			}
			if uv.Index > level {
				next = uv
				break
			}
			prev = uv
		}
	}
	uv := &types.UpValue{Index: level, Closed: false}
	if prev != nil {
		prev.Next = uv
	} else {
		s.uvhead = uv
	}
	if next != nil {
		uv.Next = next
	}
	return uv
}

func (s *State) closeUpValues(idx int) {
	if s.uvhead != nil {
		var prev *types.UpValue
		for uv := s.uvhead; uv != nil; uv = uv.Next {
			if uv.Index >= idx {
				if prev != nil {
					prev.Next = nil
				} else {
					s.uvhead = nil
				}
				uv.Close(s.CallStack)
			}
			prev = uv
		}
	}
}

func (s *State) RegisterFunc(name string, fn func(*State)) {
	cl := types.NewGoClosure(name, fn)
	s.Global[name] = cl
}
