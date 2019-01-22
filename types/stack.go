package types

import "fmt"

type Stack struct {
	arr []Object
	sp  int
	len int
}

func NewStack(defaultLen int) *Stack {
	return new(Stack).Init(defaultLen)
}

func (s *Stack) Init(defaultLen int) *Stack {
	s.arr = make([]Object, defaultLen)
	s.sp = -1
	s.len = defaultLen
	return s
}

func (s *Stack) Top() Object {
	return s.Get(s.sp)
}

func (s *Stack) Sp() int {
	return s.sp
}

func (s *Stack) SetSp(i int) {
	s.sp = i
}

func (s *Stack) Len() int {
	return s.len
}

func (s *Stack) Push(obj Object) {
	s.sp++
	s.arr[s.sp] = obj
}

func (s *Stack) Pop() Object {
	if s.sp < 0 {
		return nil
	}
	obj := s.Get(s.sp)
	s.arr[s.sp] = nil
	s.sp--
	return obj
}

func (s *Stack) Get(i int) Object {
	if i >= 0 && i < s.len {
		obj := s.arr[i]
		if obj == nil {
			return nil
		}
		return obj
	}
	return nil
}

func (s *Stack) Set(i int, obj Object) {
	s.arr[i] = obj
}

func (s *Stack) Dump() {
	fmt.Printf("SP = %d, LEN = %d\n", s.sp, s.len)
	for i := 0; i < 20; i++ {
		fmt.Printf("%d => %v\n", i, s.arr[i])
	}
}
