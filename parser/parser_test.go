package parser

import (
	"github.com/hyusuk/tama/types"
	"testing"
)

func TestParseFile(t *testing.T) {
	p := &Parser{}
	p.Init([]byte(" 1 "))
	num, ok := p.parseObject().(types.Number)
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
	p.Init([]byte("(+ 1 2)"))
	pair, ok := p.parseObject().(*types.Pair)
	if !ok {
		t.Fatalf("expected pair")
	}
	car, _ := types.Car(pair)
	sym, ok := car.(*types.Symbol)
	if !ok {
		t.Fatalf("expected symbol")
	}
	if sym.Name != "+" {
		t.Fatalf("expected %s, but got %s", sym.Name, "+")
	}
}
