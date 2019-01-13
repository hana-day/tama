package tama

func (s *State) OpenBase() {
	cl := NewGoClosure()
	cl.Fn = fnAdd
	s.Global["+"] = cl
}

func fnAdd(s *State) {
	nargs := s.CallStack.Pop().(Number)
	var result Number = 0
	for i := 0; i < int(nargs); i++ {
		result += s.CallStack.Pop().(Number)
	}
	s.CallStack.Push(result)
}
