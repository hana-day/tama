package types

import "fmt"

type valueType interface{}

type Element struct {
	Value valueType

	// The stack to which this element belongs.
	stack *Stack
}

type Stack struct {
	arr []*Element
	sp  int
	len int
}

func NewStack(defaultLen int) *Stack {
	return new(Stack).Init(defaultLen)
}

func (s *Stack) Init(defaultLen int) *Stack {
	s.arr = make([]*Element, defaultLen)
	s.sp = -1
	s.len = defaultLen
	return s
}

func (s *Stack) Top() valueType {
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

func (s *Stack) Push(value valueType) {
	s.sp++
	s.arr[s.sp] = &Element{
		Value: value,
		stack: s,
	}
}

func (s *Stack) Pop() valueType {
	if s.sp < 0 {
		return nil
	}
	v := s.Get(s.sp)
	s.arr[s.sp] = nil
	s.sp--
	return v
}

func (s *Stack) Get(i int) valueType {
	if i >= 0 && i < s.len {
		v := s.arr[i]
		if v == nil {
			return nil
		}
		return v.Value
	}
	return nil
}

func (s *Stack) Set(i int, v valueType) {
	s.arr[i] = &Element{
		Value: v,
		stack: s,
	}
}

func (s *Stack) Dump() {
	fmt.Printf("SP = %d, LEN = %d\n", s.sp, s.len)
	for i := 0; i < 20; i++ {
		fmt.Printf("%d => %v\n", i, s.arr[i])
	}
}
