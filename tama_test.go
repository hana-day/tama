package main

import (
	"testing"
)

func TestTama(t *testing.T) {
	var testcases = []struct {
		code     string
		expected string
	}{
		{"1", "1"},
	}
	for i, tc := range testcases {
		result := eval(tc.code)
		if result != tc.expected {
			t.Fatalf("case: %d, expected %s, but got %s", i, tc.expected, result)
		}
	}
}
