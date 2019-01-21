// This file defines tama's opcodes.
// The definition of opcodes is almost same with Lua 5.1.4's opcodes.
// See http://underpop.free.fr/l/lua/docs/a-no-frills-introduction-to-lua-5.1-vm-instructions.pdf

package compiler

import "fmt"

const (
	OP_RETURN int = iota
	OP_LOADK
	OP_GETGLOBAL
	OP_SETGLOBAL
	OP_MOVE
	OP_CLOSURE
	OP_CALL
	OP_GETUPVAL
)

type opType int

const (
	opTypeABC opType = iota
	opTypeABx
	opTypeASbx
)

type opProp struct {
	Name string
	Type opType
}

var opProps = []opProp{
	opProp{"RETURN", opTypeABC},
	opProp{"LOADK", opTypeABx},
	opProp{"GETGLOBAL", opTypeABx},
	opProp{"SETGLOBAL", opTypeABx},
	opProp{"MOVE", opTypeABC},
	opProp{"CLOSURE", opTypeABx},
	opProp{"CALL", opTypeABC},
	opProp{"GETUPVAL", opTypeABC},
}

func GetOpType(inst uint32) opType {
	return opProps[GetOpCode(inst)].Type
}

func GetOpName(inst uint32) string {
	return opProps[GetOpCode(inst)].Name
}

func GetOpCode(inst uint32) int {
	return int(inst >> 26)
}

func GetArgA(inst uint32) int {
	return int(inst>>18) & 0xff
}

func GetArgB(inst uint32) int {
	return int(inst & 0x1ff)
}

func GetArgBx(inst uint32) int {
	return int(inst & 0x3ffff)
}

func GetArgC(inst uint32) int {
	return int(inst>>9) & 0x1ff
}

func SetOpCode(inst *uint32, opcode int) {
	*inst = (*inst & 0x3ffffff) | uint32(opcode<<26)
}

func SetArgA(inst *uint32, arg int) {
	*inst = (*inst & 0xfc03ffff) | uint32((arg&0xff)<<18)
}

func SetArgB(inst *uint32, arg int) {
	*inst = (*inst & 0xfffffe00) | uint32(arg&0x1ff)
}

func SetArgC(inst *uint32, arg int) {
	*inst = (*inst & 0xfffc01ff) | uint32((arg&0x1ff)<<9)
}

func SetArgBx(inst *uint32, arg int) {
	*inst = (*inst & 0xfffc0000) | uint32(arg&0x3ffff)
}

func CreateABC(op int, a int, b int, c int) uint32 {
	var inst uint32 = 0
	SetOpCode(&inst, op)
	SetArgA(&inst, a)
	SetArgB(&inst, b)
	SetArgC(&inst, c)
	return inst
}

func CreateABx(op int, a int, bx int) uint32 {
	var inst uint32 = 0
	SetOpCode(&inst, op)
	SetArgA(&inst, a)
	SetArgBx(&inst, bx)
	return inst
}

func dumpABC(inst uint32) string {
	opname := GetOpName(inst)
	a := GetArgA(inst)
	b := GetArgB(inst)
	c := GetArgC(inst)
	return fmt.Sprintf("%s %d %d %d", opname, a, b, c)
}

func dumpABx(inst uint32) string {
	opname := GetOpName(inst)
	a := GetArgA(inst)
	bx := GetArgBx(inst)
	return fmt.Sprintf("%s %d %d", opname, a, bx)
}

func DumpInst(inst uint32) string {
	switch GetOpType(inst) {
	case opTypeABC:
		return dumpABC(inst)
	case opTypeABx:
		return dumpABx(inst)
	default:
		panic("unsupported optype")
	}
}
