package token

import (
	"testing"
)

func TestScan(t *testing.T) {
	var s Scanner
	s.Init([]byte("  123 "))
	token, lit := s.Scan()
	if token != INT {
		t.Fatalf("expected %d, but got %d", INT, token)
	}
	if lit != "123" {
		t.Fatalf("expected %s, but got %s", "123", lit)
	}
	token, _ = s.Scan()
	if token != EOF {
		t.Fatalf("expected %d, but got %d", EOF, token)
	}
}
