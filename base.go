package tama

func (s *State) OpenBase() {
	cl := NewGoClosure()
	cl.Fn = fnAdd
	s.Global["+"] = cl
}

func fnAdd(s *State) int {
	// number of arguments
	_ = s.CallStack.Pop().(Number)
	a1 := s.CallStack.Pop().(Number)
	a2 := s.CallStack.Pop().(Number)
	var result Number = a1 + a2
	s.CallStack.Push(result)
	return 0
}
