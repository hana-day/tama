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
	cl, _ := Compile(exprs)
	insts := cl.Proto.Insts
	if len(insts) != 2 {
		t.Fatalf("expected %d, but got %d", 2, len(insts))
	}
	if opcode := GetOpCode(insts[0]); opcode != LOADK {
		t.Fatalf("expected %d, but got %d", LOADK, opcode)
	}
	if opcode := GetOpCode(insts[1]); opcode != RETURN {
		t.Fatalf("expected %d, but got %d", RETURN, opcode)
	}
}

func TestCompileIdent(t *testing.T) {
	c := &Compiler{
		Proto: newClosureProto(),
	}
	c.compileIdent(&parser.Ident{
		Name: "test",
	})
	if len(c.Proto.Insts) != 2 {
		t.Fatalf("expected %d, but got %d", 2, len(c.Proto.Insts))
	}
	if opcode := GetOpCode(c.Proto.Insts[0]); opcode != LOADK {
		t.Fatalf("expected %d, but got %d", LOADK, opcode)
	}
	if opcode := GetOpCode(c.Proto.Insts[1]); opcode != GETGLOBAL {
		t.Fatalf("expected %d, but got %d", GETGLOBAL, opcode)
	}
}
