package tama

import (
	"github.com/hyusuk/tama/types"
)

func (s *State) OpenBase() *State {
	cl := types.NewGoClosure()
	cl.Fn = fnAdd
	s.Global["+"] = cl
	return s
}

func fnAdd(s *State) {
	nargs := s.CallStack.Pop().(types.Number)
	var result types.Number = 0
	for i := 0; i < int(nargs); i++ {
		result += s.CallStack.Pop().(types.Number)
	}
	s.CallStack.Push(result)
}
