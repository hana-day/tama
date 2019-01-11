package parser

import (
	"github.com/hyusuk/tama/scanner"
	"testing"
)

func TestParseFile(t *testing.T) {
	p := &Parser{}
	p.Init([]byte(" 1 "))
	f := p.ParseFile()
	if len(f.Exprs) != 1 {
		t.Fatalf("expected %d, but got %d", 1, len(f.Exprs))
	}
	prim, ok := f.Exprs[0].(*Primitive)
	if !ok {
		t.Fatalf("Unexpected expression")
	}
	if prim.Kind != scanner.INT || prim.Value != "1" {
		t.Fatalf("Unexpected primitive, kind: %d, value: %s", prim.Kind, prim.Value)
	}
}

func TestParseExpr(t *testing.T) {
	p := &Parser{}

	// Parse procedure call expression
	p.Init([]byte("(+ 1 2)"))
	expr := p.parseExpr().(*CallExpr)
	name := expr.Func.(*Ident).Name
	if name != "+" {
		t.Fatalf("expected %s, but got %s", name, "+")
	}
	if len(expr.Args) != 2 {
		t.Fatalf("expected %d, but got %d", len(expr.Args), 2)
	}
	prim := expr.Args[1].(*Primitive)
	if prim.Kind != scanner.INT {
		t.Fatalf("expected %d, but got %d", prim.Kind, scanner.INT)
	}
	if prim.Value != "2" {
		t.Fatalf("expected %s, but got %s", prim.Value, "2")
	}
}
