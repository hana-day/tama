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
		{
			func() *State { return NewState() },
			"((lambda (a) (+ a 1) (+ a 2)) 1)",
			"3",
		},
		{
			func() *State { return NewState() },
			"((lambda (a) (+ (car a) (cdr a))) (cons 1 2))",
			"3",
		},
		{
			func() *State { return NewState() },
			"(set! a 1) a",
			"1",
		},
	}
	for i, tc := range testcases {
		s := tc.stateFactory()
		if err := s.ExecString(tc.source); err != nil {
			t.Fatalf("case %d: unexpected error %v ;  source: %s", i, err, tc.source)
		}
		s.CallStack.Dump()
		top := s.CallStack.Top()
		v, ok := top.(types.Object)
		if !ok {
			t.Fatalf("case %d: unsupported object %v ; source: %s", i, top, tc.source)
		}
		if v.String() != tc.resultString {
			t.Fatalf("case %d: expected %s, but got %s ; source %s", i, tc.resultString, v.String(), tc.source)
		}
	}
}
