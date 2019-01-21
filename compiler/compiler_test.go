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
	if opcode := GetOpCode(insts[0]); opcode != OP_LOADK {
		t.Fatalf("expected %d, but got %d", OP_LOADK, opcode)
	}
	if opcode := GetOpCode(insts[1]); opcode != OP_RETURN {
		t.Fatalf("expected %d, but got %d", OP_RETURN, opcode)
	}
}

func TestCompileSymbol(t *testing.T) {
	fs := newFuncState(nil)
	c := &Compiler{}
	c.compileSymbol(fs, &types.Symbol{Name: "test"})
	if len(fs.Proto.Insts) != 1 {
		t.Fatalf("expected %d, but got %d", 2, len(fs.Proto.Insts))
	}
	if opcode := GetOpCode(fs.Proto.Insts[0]); opcode != OP_GETGLOBAL {
		t.Fatalf("expected %d, but got %d", OP_GETGLOBAL, opcode)
	}
}

func TestNameStorage(t *testing.T) {
	ns := newNameStorage(0)
	if ns.Find("test") != -1 {
		t.Fatalf("unexpected index")
	}
	i := ns.Register("test")
	if i != 0 {
		t.Fatalf("expected %d, but got %d", 0, i)
	}
	i = ns.Find("test")
	if i != 0 {
		t.Fatalf("expected %d, but got %d", 0, i)
	}
	if ns.Capacity() != 2 {
		t.Fatalf("expected %d, but got %d", 2, ns.Capacity())
	}
	if ns.Len() != 1 {
		t.Fatalf("expected %d, but got %d", 1, ns.Len())
	}
}
