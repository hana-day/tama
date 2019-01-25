package tama

import (
	"fmt"
	"github.com/hyusuk/tama/compiler"
	"github.com/hyusuk/tama/parser"
	"github.com/hyusuk/tama/types"
)

const (
	DefaultStackSize     = 256 * 20
	DefaultCallInfosSize = 256
)

type Option struct {
	StackSize     int
	CallInfosSize int
}

type State struct {
	// call stack
	CallStack *types.Stack
	CallInfos *types.Stack
	Global    map[string]types.Object
	uvhead    *types.UpValue
}

type GoFunc = func(s *State, args []types.Object) (types.Object, error)

func NewState(option Option) *State {
	if option.StackSize == 0 {
		option.StackSize = DefaultStackSize
	}
	if option.CallInfosSize == 0 {
		option.CallInfosSize = DefaultCallInfosSize
	}

	s := &State{
		CallStack: types.NewStack(option.StackSize),
		CallInfos: types.NewStack(option.CallInfosSize),
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

// popArgs pops arguements and create a slice [argument 1, ..., argument nargs].
//
//         [before]                      [after]
//
//      |            |               |            |
//      +------------+               +------------+
//      | closure    |               | closure    |
//      | argument 1 |          SP ->+------------+
//      |     ...    |               | argument 1 |
//      | argument N |               |     ...    |
// SP ->+------------+               | argument N |
//      |            |               |            |
//
func (s *State) popArgs(nargs int) []types.Object {
	topSp := s.CallStack.Sp()
	args := make([]types.Object, nargs)
	for i := 0; i < nargs; i++ {
		args[i] = s.CallStack.Get(topSp - nargs + i + 1)
	}
	s.CallStack.SetSp(topSp - nargs)
	return args
}

// precall prepares the function call.
// If the function is a scheme-function, push call information onto the stack.
// If the function is a go-function, push call information onto the stack and call it.
//
// Before precalling, the stack contents must be like below.
//
//      |            |
//      +------------+
//      | closure    |
//      | argument 1 |
//      |     ...    |
//      | argument N |
// SP ->+------------+
//      |            |
//
func (s *State) precall(clIndex int) (*types.CallInfo, error) {
	cl, ok := s.CallStack.Get(clIndex).(*types.Closure)
	if !ok {
		return nil, fmt.Errorf("function is not loaded")
	}
	nargs := s.CallStack.Sp() - clIndex
	if cl.IsGo {
		ci := &types.CallInfo{Cl: cl, Base: clIndex + 1, FuncSp: clIndex}
		s.CallInfos.Push(ci)

		fn, ok := cl.Fn.(GoFunc)
		if !ok {
			return nil, fmt.Errorf("invalid function %v", cl.Fn)
		}
		args := s.popArgs(nargs)
		retval, err := fn(s, args)
		if err != nil {
			return nil, err
		}
		s.CallStack.Push(retval)
		s.postcall(s.CallStack.Sp())
		return ci, nil
	} else {
		switch cl.Proto.Mode {
		case types.FixedArgMode:
			if nargs != len(cl.Proto.Args) {
				return nil, fmt.Errorf("invalid number of arguments")
			}
		case types.VArgMode:
			args := s.popArgs(nargs)
			s.CallStack.Push(types.List(args...))
		case types.RestArgMode:
			if nargs < len(cl.Proto.Args) {
				return nil, fmt.Errorf("insufficient number of arguments")
			}
			nrest := nargs - len(cl.Proto.Args) + 1
			rest := s.popArgs(nrest)
			s.CallStack.Push(types.List(rest...))
		}
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

func (s *State) SetGlobal(name string, obj types.Object) {
	s.Global[name] = obj
}

func (s *State) GetGlobal(name string) (obj types.Object, ok bool) {
	obj, ok = s.Global[name]
	return
}

func (s *State) RegisterFunc(name string, fn GoFunc) {
	cl := types.NewGoClosure(name, fn)
	s.SetGlobal(name, cl)
}
