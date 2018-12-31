package tama

import (
	"testing"
)

func TestExecString(t *testing.T) {
	testcases := []struct {
		source   string
		expected string
	}{
		{"1", "1"},
	}

	state := NewState()
	for i, tc := range testcases {
		result, _ := state.ExecString(tc.source)
		if result != tc.expected {
			t.Fatalf("case: %d, expected: %s, but got %s", i, tc.expected, result)
		}
	}
}
