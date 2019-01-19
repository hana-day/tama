package tama

import (
	"fmt"
	"github.com/hyusuk/tama/compiler"
	"github.com/hyusuk/tama/types"
)

func runVM(s *State) error {
	nexeccalls := 1
reentry:
	ci, _ := s.CallInfos.Top().(*CallInfo)
	cl := ci.Cl
	base := ci.Base
	for {
		inst := cl.Proto.Insts[ci.Pc]
		ci.Pc++
		ra := base + compiler.GetArgA(inst)
		switch compiler.GetOpCode(inst) {
		case compiler.LOADK:
			bx := compiler.GetArgBx(inst)
			s.CallStack.Set(ra, cl.Proto.Consts[bx])
		case compiler.GETGLOBAL:
			bx := compiler.GetArgBx(inst)
			s.CallStack.Set(ra, s.Global[cl.Proto.Consts[bx].String()])
		case compiler.SETGLOBAL:
			bx := compiler.GetArgBx(inst)
			obj := s.CallStack.Get(ra).(types.Object)
			s.Global[cl.Proto.Consts[bx].String()] = obj
		case compiler.MOVE:
			rb := base + compiler.GetArgB(inst)
			s.CallStack.Set(ra, s.CallStack.Get(rb))
		case compiler.CLOSURE:
			bx := compiler.GetArgBx(inst)
			proto := cl.Proto.Protos[bx]
			newCl := types.NewScmClosure()
			newCl.Proto = proto
			s.CallStack.Set(ra, newCl)
		case compiler.CALL:
			b := compiler.GetArgB(inst)
			s.CallStack.SetSp(ra + b)
			precalledCi, err := s.precall(ra)
			if err != nil {
				return err
			}
			if !precalledCi.Cl.IsGo {
				nexeccalls++
				goto reentry
			}
		case compiler.RETURN:
			b := compiler.GetArgB(inst)
			if b != 2 {
				return fmt.Errorf("invalid number of returns")
			}
			s.postcall(ra)
			nexeccalls--
			if nexeccalls == 0 {
				return nil
			} else {
				goto reentry
			}
		}
	}
	return nil
}
