package tama

import (
	"testing"
)

func TestExecuteString(t *testing.T) {
	s := NewState()
	err := s.ExecString(" 1 ")
	if err != nil {
		t.Fatal(err)
	}
	num, ok := s.CallStack.Top().(Number)
	if !ok {
		t.Fatalf("Invalid call stack top")
	}
	if num.String() != "1" {
		t.Fatalf("expected %s, but got %s", "1", num.String())
	}

	s = NewState()
	s.Global["test"] = Number(2)
	err = s.ExecString("test")
	if err != nil {
		t.Fatal(err)
	}
	num, ok = s.CallStack.Top().(Number)
	if !ok {
		t.Fatalf("Invalid call stack top")
	}
	if num.String() != "2" {
		t.Fatalf("expected %s, but got %s", "2", num.String())
	}
}
