package parser

import (
	"github.com/hyusuk/tama/types"
	"testing"
)

func TestParseFile(t *testing.T) {
	p := &Parser{}
	if err := p.Init([]byte(" 1 ")); err != nil {
		t.Fatal(err)
	}
	obj, err := p.parseObject()
	if err != nil {
		t.Fatal(err)
	}
	num, ok := obj.(types.Number)
	if !ok {
		t.Fatalf("exected number")
	}
	if num.String() != "1" {
		t.Fatalf("expected %s, but got %s", "1", num.String())
	}
}

func TestParsePair(t *testing.T) {
	p := &Parser{}

	// Parse procedure call expression
	if err := p.Init([]byte("(+ 1 2)")); err != nil {
		t.Fatal(err)
	}
	obj, err := p.parseObject()
	if err != nil {
		t.Fatal(err)
	}
	pair, ok := obj.(*types.Pair)
	if !ok {
		t.Fatalf("expected pair")
	}
	sym, ok := pair.Car().(*types.Symbol)
	if !ok {
		t.Fatalf("expected symbol")
	}
	if sym.Name != "+" {
		t.Fatalf("expected %s, but got %s", sym.Name, "+")
	}
}
