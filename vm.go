package tama

import (
	"fmt"
	"github.com/hyusuk/tama/compiler"
	"github.com/hyusuk/tama/types"
)

func runVM(s *State, debug bool) error {
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
		case compiler.OP_LOADK:
			bx := compiler.GetArgBx(inst)
			s.CallStack.Set(ra, cl.Proto.Consts[bx])
			if debug {
				fmt.Printf("%-20s ; R[%d] = %v\n", compiler.DumpInst(inst), ra, cl.Proto.Consts[bx])
			}
		case compiler.OP_GETGLOBAL:
			bx := compiler.GetArgBx(inst)
			s.CallStack.Set(ra, s.Global[cl.Proto.Consts[bx].String()])
			if debug {
				fmt.Printf("%-20s ; R[%d] = %v\n", compiler.DumpInst(inst), ra, s.Global[cl.Proto.Consts[bx].String()])
			}
		case compiler.OP_SETGLOBAL:
			bx := compiler.GetArgBx(inst)
			obj := s.CallStack.Get(ra).(types.Object)
			s.Global[cl.Proto.Consts[bx].String()] = obj
			if debug {
				fmt.Printf("%-20s ; Gbl[%v] = %v\n", compiler.DumpInst(inst), cl.Proto.Consts[bx].String(), obj)
			}
		case compiler.OP_MOVE:
			rb := base + compiler.GetArgB(inst)
			s.CallStack.Set(ra, s.CallStack.Get(rb))
			if debug {
				fmt.Printf("%-20s ; R[%d] = R[%d]\n", compiler.DumpInst(inst), ra, rb)
			}
		case compiler.OP_CLOSURE:
			bx := compiler.GetArgBx(inst)
			proto := cl.Proto.Protos[bx]
			newCl := types.NewScmClosure()
			newCl.Proto = proto
			s.CallStack.Set(ra, newCl)
			if debug {
				fmt.Printf("%-20s ; R[%d] = %v\n", compiler.DumpInst(inst), ra, newCl)
			}
		case compiler.OP_CALL:
			b := compiler.GetArgB(inst)
			s.CallStack.SetSp(ra + b - 1)
			if debug {
				cl := s.CallStack.Get(ra).(*types.Closure)
				fmt.Printf("%-20s ; R[%d] = %v(R[%d]...R[%d])\n", compiler.DumpInst(inst), ra, cl, ra+1, ra+b-1)
			}
			precalledCi, err := s.precall(ra)
			if err != nil {
				return err
			}
			if !precalledCi.Cl.IsGo {
				nexeccalls++
				goto reentry
			}
		case compiler.OP_RETURN:
			b := compiler.GetArgB(inst)
			if b != 2 {
				return fmt.Errorf("vm: invalid number of returns")
			}
			if debug {
				fmt.Printf("%-20s ; return R[%d]\n", compiler.DumpInst(inst), ra)
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
