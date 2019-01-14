package tama

import (
	"github.com/hyusuk/tama/types"
	"testing"
)

func TestExecuteString(t *testing.T) {
	s := NewState()
	err := s.ExecString(" 1 ")
	if err != nil {
		t.Fatal(err)
	}
	num, ok := s.CallStack.Top().(types.Number)
	if !ok {
		t.Fatalf("Invalid call stack top")
	}
	if num.String() != "1" {
		t.Fatalf("expected %s, but got %s", "1", num.String())
	}

	s = NewState()
	s.Global["test"] = types.Number(2)
	err = s.ExecString("test")
	if err != nil {
		t.Fatal(err)
	}
	num, ok = s.CallStack.Top().(types.Number)
	if !ok {
		t.Fatalf("Invalid call stack top")
	}
	if num.String() != "2" {
		t.Fatalf("expected %s, but got %s", "2", num.String())
	}

	s = NewState()
	s.OpenBase()
	err = s.ExecString("(+ 1 2 3) (+ 1 2 3 4)")
	if err != nil {
		t.Fatal(err)
	}
	num, ok = s.CallStack.Top().(types.Number)
	if !ok {
		t.Fatalf("Invalid call stack top")
	}
	if num.String() != "10" {
		t.Fatalf("expected %s, but got %s", "10", num.String())
	}
}
