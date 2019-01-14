package compiler

import (
	"github.com/hyusuk/tama/types"
	"testing"
)

func TestCompileNumber(t *testing.T) {
	num := types.Number(1)
	objs := []types.Object{num}
	cl, _ := Compile(objs)
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

func TestCompileSymbol(t *testing.T) {
	c := &Compiler{
		Proto: newClosureProto(),
	}
	c.compileSymbol(&types.Symbol{
		Name: "test",
	})
	if len(c.Proto.Insts) != 1 {
		t.Fatalf("expected %d, but got %d", 2, len(c.Proto.Insts))
	}
	if opcode := GetOpCode(c.Proto.Insts[0]); opcode != GETGLOBAL {
		t.Fatalf("expected %d, but got %d", GETGLOBAL, opcode)
	}
}
