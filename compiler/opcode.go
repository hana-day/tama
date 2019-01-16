// This file defines tama's opcodes.
// The definition of opcodes is almost same with Lua 5.1.4's opcodes.
// See http://underpop.free.fr/l/lua/docs/a-no-frills-introduction-to-lua-5.1-vm-instructions.pdf

package compiler

const (
	RETURN int = iota
	LOADK
	GETGLOBAL
	SETGLOBAL
	MOVE
	CALL
)

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
