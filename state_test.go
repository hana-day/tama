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
			"(+ 1 2 3) (+ 1 2 3 4)",
			"10",
		},
	}
	for i, tc := range testcases {
		s := tc.stateFactory()
		if err := s.ExecString(tc.source); err != nil {
			t.Fatalf("case %d: unexpected error %v", i, err)
		}
		v, ok := s.CallStack.Top().(types.Object)
		if !ok {
			t.Fatalf("case %d: unsupported object was found", i)
		}
		if v.String() != tc.resultString {
			t.Fatalf("case %d: expected %s, but got %s", i, v.String(), tc.resultString)
		}
	}
}
