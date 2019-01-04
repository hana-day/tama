package tama

import (
	"testing"
)

func TestNewStack(t *testing.T) {
	s := NewStack()
	if s.Top() != nil {
		t.Fatalf("expected nil")
	}
	if s.Len() != 0 {
		t.Fatalf("expected %d, but got %d", 0, s.Len())
	}
}

func TestStackPushAndPop(t *testing.T) {
	s := NewStack()
	s.Push(1)
	if i := s.Pop().(int); i != 1 {
		t.Fatalf("expected %d, but got %d", 1, i)
	}
	if s.Pop() != nil {
		t.Fatalf("expected nil")
	}
}
