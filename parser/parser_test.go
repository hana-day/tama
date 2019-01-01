package parser

import (
	"github.com/hyusuk/tama/token"
	"testing"
)

func TestParseFile(t *testing.T) {
	p := &Parser{}
	p.init([]byte(" 1 "))
	f := p.ParseFile()
	if len(f.Exprs) != 1 {
		t.Fatalf("expected %d, but got %d", 1, len(f.Exprs))
	}
	prim, ok := f.Exprs[0].(*Primitive)
	if !ok {
		t.Fatalf("Unexpected expression")
	}
	if prim.Kind != token.INT || prim.Value != "1" {
		t.Fatalf("Unexpected primitive, kind: %d, value: %s", prim.Kind, prim.Value)
	}
}
