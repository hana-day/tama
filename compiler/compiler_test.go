package compiler

import (
	"github.com/hyusuk/tama/parser"
	"github.com/hyusuk/tama/token"
	"testing"
)

func TestCompileExprs(t *testing.T) {
	exprs := []parser.Expr{
		&parser.Primitive{
			Kind:  token.INT,
			Value: "1",
		},
	}
	c := &Compiler{}
	c.compileExprs(exprs)
	if len(c.insts) != 2 {
		t.Fatalf("expected %d, but got %d", 2, len(c.insts))
	}
	if opcode := getOpCode(c.insts[0]); opcode != LOADK {
		t.Fatalf("expected %d, but got %d", LOADK, opcode)
	}
	if opcode := getOpCode(c.insts[1]); opcode != RETURN {
		t.Fatalf("expected %d, but got %d", RETURN, opcode)
	}
}
