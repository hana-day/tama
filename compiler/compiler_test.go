package compiler

import (
	"github.com/hyusuk/tama/types"
	"testing"
)

func TestCompileNumber(t *testing.T) {
	num := types.Number(1)
	objs := []types.Object{num}
	cl, err := Compile(objs)
	if err != nil {
		t.Fatal(err)
	}
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
	fs := newFuncState(nil)
	c := &Compiler{}
	c.compileSymbol(fs, &types.Symbol{Name: "test"})
	if len(fs.Proto.Insts) != 1 {
		t.Fatalf("expected %d, but got %d", 2, len(fs.Proto.Insts))
	}
	if opcode := GetOpCode(fs.Proto.Insts[0]); opcode != GETGLOBAL {
		t.Fatalf("expected %d, but got %d", GETGLOBAL, opcode)
	}
}
