package tama

import (
	"github.com/hyusuk/tama/types"
)

func (s *State) OpenBase() *State {
	s.RegisterFunc("+", fnAdd)
	s.RegisterFunc("cons", fnCons)
	s.RegisterFunc("car", fnCar)
	s.RegisterFunc("cdr", fnCdr)
	return s
}

func fnAdd(s *State) {
	// TODO: handle error cases
	nargs := s.CallStack.Pop().(types.Number)
	var result types.Number = 0
	for i := 0; i < int(nargs); i++ {
		result += s.CallStack.Pop().(types.Number)
	}
	s.CallStack.Push(result)
}

func fnCons(s *State) {
	// TODO: handle error cases
	_ = s.CallStack.Pop().(types.Number)
	cdr := s.CallStack.Pop().(types.Object)
	car := s.CallStack.Pop().(types.Object)
	s.CallStack.Push(types.Cons(car, cdr))
}

func fnCar(s *State) {
	// TODO: handle error cases
	_ = s.CallStack.Pop().(types.Number)
	pair := s.CallStack.Pop().(*types.Pair)
	s.CallStack.Push(pair.Car)
}

func fnCdr(s *State) {
	// TODO: handle error cases
	_ = s.CallStack.Pop().(types.Number)
	pair := s.CallStack.Pop().(*types.Pair)
	s.CallStack.Push(pair.Cdr)
}
