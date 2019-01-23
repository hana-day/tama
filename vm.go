package tama

import (
	"fmt"
	"github.com/hyusuk/tama/compiler"
	"github.com/hyusuk/tama/types"
)

func runVM(s *State, debug bool) error {
	nexeccalls := 1
reentry:
	if debug {
		fmt.Println("[Enter function]")
	}
	ci, _ := s.CallInfos.Top().(*types.CallInfo)
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
			obj := s.CallStack.Get(ra)
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
			newCl := types.NewScmClosure(proto, proto.NUpVals)
			s.CallStack.Set(ra, newCl)
			if debug {
				fmt.Printf("%-20s ; R[%d] = %v\n", compiler.DumpInst(inst), ra, newCl)
			}
			for i := 0; i < proto.NUpVals; i++ {
				inst = ci.Cl.Proto.Insts[ci.Pc]
				ci.Pc++
				b := compiler.GetArgB(inst)
				switch compiler.GetOpCode(inst) {
				case compiler.OP_MOVE:
					uv := s.findUpValue(base + b)
					newCl.UpVals[i] = uv
					if debug {
						fmt.Printf("%-20s ; Up[%d] = R[%d]\n", compiler.DumpInst(inst), i, base+b)
					}
				case compiler.OP_GETUPVAL:
					newCl.UpVals[i] = ci.Cl.UpVals[b]
					if debug {
						fmt.Printf("%-20s ; Up[%d] = Up[%d]\n", compiler.DumpInst(inst), i, b)
					}
				}
			}

		case compiler.OP_CALL:
			b := compiler.GetArgB(inst)
			s.CallStack.SetSp(ra + b - 1)
			if debug {
				cl := s.CallStack.Get(ra)
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
		case compiler.OP_GETUPVAL:
			b := compiler.GetArgB(inst)
			uv := ci.Cl.UpVals[b]
			v := uv.Value(s.CallStack)
			s.CallStack.Set(ra, v)
			if debug {
				fmt.Printf("%-20s ; R[%d] = %v\n", compiler.DumpInst(inst), ra, v)
			}
		case compiler.OP_SETUPVAL:
			b := compiler.GetArgB(inst)
			uv := ci.Cl.UpVals[b]
			v := s.CallStack.Get(ra)
			uv.Set(s.CallStack, v)
			if debug {
				fmt.Printf("%-20s ; Up[%d] = %v\n", compiler.DumpInst(inst), b, v)
			}
		case compiler.OP_CLOSE:
			s.closeUpValues(ra)
			if debug {
				fmt.Printf("%-20s ; close %d\n", compiler.DumpInst(inst), ra)
			}
		case compiler.OP_TEST:
			c := compiler.GetArgC(inst)
			v := s.CallStack.Get(ra)
			if (c == 0) == IsTruthy(v) {
				ci.Pc++
			}
			if debug {
				incPc := (c == 0) == IsTruthy(v)
				if incPc {
					fmt.Printf("%-20s ; pc += 1\n", compiler.DumpInst(inst))
				} else {
					fmt.Printf("%-20s ; pc += 0\n", compiler.DumpInst(inst))
				}
			}
		case compiler.OP_JMP:
			sbx := compiler.GetArgSbx(inst)
			ci.Pc += sbx
			if debug {
				fmt.Printf("%-20s ; pc += %d\n", compiler.DumpInst(inst), sbx)
			}
		}
	}
	return nil
}
