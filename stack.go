package tama

type Element struct {
	Value interface{}
	next  *Element

	// The stack to which this element belongs.
	stack *Stack
}

type Stack struct {
	top *Element
	len int
}

func NewStack() *Stack {
	return new(Stack).Init()
}

func (s *Stack) Init() *Stack {
	s.top = nil
	s.len = 0
	return s
}

func (s *Stack) Top() interface{} {
	if s.Len() > 0 {
		return s.top.Value
	}
	return nil
}

func (s *Stack) Len() int {
	return s.len
}

func (s *Stack) Push(value interface{}) {
	s.top = &Element{
		Value: value,
		next:  s.top,
	}
	s.len++
}

func (s *Stack) Pop() (value interface{}) {
	if s.Len() > 0 {
		value = s.top.Value
		s.top = s.top.next
		s.len--
		return
	}
	return nil
}
