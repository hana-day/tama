package scanner

import (
	"testing"
)

func TestScan(t *testing.T) {
	var s Scanner
	type expect struct {
		tok Token
		lit string
	}
	testcases := []struct {
		src     []byte
		expects []expect
	}{
		{
			src: []byte(" 123 "),
			expects: []expect{
				{tok: INT, lit: "123"},
				{tok: EOF, lit: ""},
			},
		},
	}
	for _, tc := range testcases {
		s.Init(tc.src)
		for _, expect := range tc.expects {
			tok, lit := s.Scan()
			if tok != expect.tok {
				t.Fatalf("expected %d, but got %d", expect.tok, tok)
			}
			if lit != expect.lit {
				t.Fatalf("expected %s, but got %s", expect.lit, lit)
			}
		}
	}
}
