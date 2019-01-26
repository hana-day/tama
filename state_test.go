package tama

import (
	"github.com/hyusuk/tama/types"
	"testing"
)

func TestTailCall(t *testing.T) {
	testcases := []struct {
		stateFactory func() *State
		source       string
		resultString string
	}{
		// if tail call doesn't work, tests below will fail due to the insufficient stack size.
		{
			func() *State { return NewState(Option{StackSize: 100}) },
			"(define (recur n) (if (= n 1) 1 (recur (- n 1)))) (recur 100)",
			"1",
		},
		{
			func() *State { return NewState(Option{StackSize: 100}) },
			"(define (recur a) (if (= a 1) 1 (begin (recur (- a 1))))) (recur 100)",
			"1",
		},
	}
	for i, tc := range testcases {
		s := tc.stateFactory()
		if err := s.ExecString(tc.source); err != nil {
			t.Fatalf("case %d: unexpected error %v ;  source: %s", i, err, tc.source)
		}
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

func TestComment(t *testing.T) {
	testcases := []struct {
		stateFactory func() *State
		source       string
		resultString string
	}{
		{
			func() *State { return NewState(Option{}) },
			"(define a 1) ; (set! a 2)\na",
			"1",
		},
		{
			func() *State { return NewState(Option{}) },
			"(define a 1) a ; (set! a 2) a",
			"1",
		},
	}

	for i, tc := range testcases {
		s := tc.stateFactory()
		if err := s.ExecString(tc.source); err != nil {
			t.Fatalf("case %d: unexpected error %v ;  source: %s", i, err, tc.source)
		}
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
