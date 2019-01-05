package tama

import (
	"github.com/hyusuk/tama/compiler"
)

func runVM(s *State) {
	ci, _ := s.CallInfos.Top().(*CallInfo)
	cl := ci.Cl
	for _, inst := range cl.Insts {
		ra := ci.Base + compiler.GetArgA(inst)
		switch compiler.GetOpCode(inst) {
		case compiler.LOADK:
			bx := compiler.GetArgBx(inst)
			s.CallStack.Set(ra, cl.Consts[bx])
		case compiler.RETURN:
			b := compiler.GetArgB(inst)
			if b != 0 {
				s.CallStack.SetSp(ra + b - 1)
			}
		}
	}
}
