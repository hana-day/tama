// This file defines tama's opcodes.
// The definition of opcodes is almost same with Lua 5.1.4's opcodes.
// See http://underpop.free.fr/l/lua/docs/a-no-frills-introduction-to-lua-5.1-vm-instructions.pdf

package compiler

const (
	RETURN int = iota
	LOADK
)

func getOpCode(inst uint32) int {
	return int(inst >> 26)
}

func getArgA(inst uint32) int {
	return int(inst>>18) & 0xff
}

func getArgB(inst uint32) int {
	return int(inst & 0x1ff)
}

func getArgBx(inst uint32) int {
	return int(inst & 0x3ffff)
}

func getArgC(inst uint32) int {
	return int(inst>>9) & 0x1ff
}

func setOpCode(inst *uint32, opcode int) {
	*inst = (*inst & 0x3ffffff) | uint32(opcode<<26)
}

func setArgA(inst *uint32, arg int) {
	*inst = (*inst & 0xfc03ffff) | uint32((arg&0xff)<<18)
}

func setArgB(inst *uint32, arg int) {
	*inst = (*inst & 0xfffffe00) | uint32(arg&0x1ff)
}

func setArgC(inst *uint32, arg int) {
	*inst = (*inst & 0xfffc01ff) | uint32((arg&0x1ff)<<9)
}

func setArgBx(inst *uint32, arg int) {
	*inst = (*inst & 0xfffc0000) | uint32(arg&0x3ffff)
}

func createABC(op int, a int, b int, c int) uint32 {
	var inst uint32 = 0
	setOpCode(&inst, op)
	setArgA(&inst, a)
	setArgB(&inst, b)
	setArgC(&inst, c)
	return inst
}

func createABx(op int, a int, bx int) uint32 {
	var inst uint32 = 0
	setOpCode(&inst, op)
	setArgA(&inst, a)
	setArgBx(&inst, bx)
	return inst
}
