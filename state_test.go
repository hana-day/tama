package tama

import (
	"github.com/hyusuk/tama/types"
	"testing"
)

func TestExecuteString(t *testing.T) {
	testcases := []struct {
		stateFactory func() *State
		source       string
		resultString string
	}{
		{
			func() *State { return NewState() },
			" 1 ",
			"1",
		},
		{
			func() *State { s := NewState(); s.Global["test"] = types.Number(2); return s },
			"test",
			"2",
		},
		{
			func() *State { return NewState() },
			"(+ 1 2 3 4)",
			"10",
		},
		{
			func() *State { return NewState() },
			"(+ 1 (+ 2 3) 4)",
			"10",
		},
		{
			func() *State { return NewState() },
			"(define a 1) (+ a 1)",
			"2",
		},
		{
			func() *State { return NewState() },
			"((lambda (a) (+ 1 2 a)) 3)",
			"6",
		},
		{
			func() *State { return NewState() },
			"(begin 1 2 3 4)",
			"4",
		},
		{
			func() *State { return NewState() },
			"(((lambda (a) (lambda (b) (+ a b))) 1) 2)",
			"3",
		},
		{
			func() *State { return NewState() },
			"((((lambda (a) (lambda (b) (lambda (c) (+ a b c)))) 1) 2) 3)",
			"6",
		},
	}
	for i, tc := range testcases {
		s := tc.stateFactory()
		if err := s.ExecString(tc.source); err != nil {
			t.Fatalf("case %d: unexpected error %v", i, err)
		}
		s.CallStack.Dump()
		top := s.CallStack.Top()
		v, ok := top.(types.Object)
		if !ok {
			t.Fatalf("case %d: unsupported object %v", i, top)
		}
		if v.String() != tc.resultString {
			t.Fatalf("case %d: expected %s, but got %s", i, tc.resultString, v.String())
		}
	}
}
