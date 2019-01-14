package tama

import (
	"github.com/hyusuk/tama/compiler"
)

func runVM(s *State) {
	ci, _ := s.CallInfos.Top().(*CallInfo)
	cl := ci.Cl
	base := ci.Base
	for _, inst := range cl.Proto.Insts {
		ra := base + compiler.GetArgA(inst)
		switch compiler.GetOpCode(inst) {
		case compiler.LOADK:
			bx := compiler.GetArgBx(inst)
			s.CallStack.Set(ra, cl.Proto.Consts[bx])
		case compiler.GETGLOBAL:
			bx := compiler.GetArgBx(inst)
			s.CallStack.Set(ra, s.Global[cl.Proto.Consts[bx].String()])
		case compiler.MOVE:
			rb := base + compiler.GetArgB(inst)
			s.CallStack.Set(ra, s.CallStack.Get(rb))
		case compiler.CALL:
			b := compiler.GetArgB(inst)
			if b != 0 {
				s.CallStack.SetSp(ra + b)
			}
			if err := s.precall(ra); err != nil {
				panic(err)
			}
			base = s.Base
		case compiler.RETURN:
			b := compiler.GetArgB(inst)
			if b != 0 {
				s.CallStack.SetSp(ra + b - 1)
			}
		}
	}
}
