package compiler

import (
	"github.com/hyusuk/tama/parser"
	"github.com/hyusuk/tama/scanner"
	"testing"
)

func TestCompileExprs(t *testing.T) {
	exprs := []parser.Expr{
		&parser.Primitive{
			Kind:  scanner.INT,
			Value: "1",
		},
	}
	c := &Compiler{}
	c.compileExprs(exprs)
	if len(c.Insts) != 2 {
		t.Fatalf("expected %d, but got %d", 2, len(c.Insts))
	}
	if opcode := GetOpCode(c.Insts[0]); opcode != LOADK {
		t.Fatalf("expected %d, but got %d", LOADK, opcode)
	}
	if opcode := GetOpCode(c.Insts[1]); opcode != RETURN {
		t.Fatalf("expected %d, but got %d", RETURN, opcode)
	}
}
