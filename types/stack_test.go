package types

import (
	"testing"
)

func TestNewStack(t *testing.T) {
	s := NewStack(100)
	if s.Len() != 100 {
		t.Fatalf("expected %d, but got %d", 100, s.Len())
	}
	if s.Top() != nil {
		t.Fatalf("expected nil")
	}
}

func TestStackPushAndPop(t *testing.T) {
	s := NewStack(100)
	s.Push(1)
	i, ok := s.Pop().(int)
	if !ok {
		t.Fatalf("unexpected value")
	}
	if i != 1 {
		t.Fatalf("expected %d, but got %d", 1, i)
	}
	if s.Pop() != nil {
		t.Fatalf("expected nil")
	}
	if s.Top() != nil {
		t.Fatalf("expected nil")
	}
}
