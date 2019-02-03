package tama

import (
	"fmt"
	"github.com/hyusuk/tama/compiler"
	"github.com/hyusuk/tama/types"
)

func runVM(s *State, debug bool) error {
	nexeccalls := 1
	nuated := false            // true if came back by using the continuation
	var nuatedObj types.Object // argument of the continuation
reentry:
	if debug {
		fmt.Println("[Enter function]")
	}
	ci := s.CallInfos.Top().(*types.CallInfo)
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
			k := cl.Proto.Consts[bx].String()
			v, ok := s.GetGlobal(k)
			if !ok {
				return types.NewInternalError("unbound symbol '%s'", k)
			}
			s.CallStack.Set(ra, v)
			if debug {
				fmt.Printf("%-20s ; R[%d] = %v with key %s\n", compiler.DumpInst(inst), ra, v, k)
			}
		case compiler.OP_SETGLOBAL:
			bx := compiler.GetArgBx(inst)
			obj := s.CallStack.Get(ra)
			s.SetGlobal(cl.Proto.Consts[bx].String(), obj)
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
		case compiler.OP_RETURN:
			b := compiler.GetArgB(inst)
			if b != 2 {
				return types.NewInternalError("invalid number of returns")
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
			if (c == 0) == types.IsTruthy(v) {
				ci.Pc++
			}
			if debug {
				incPc := (c == 0) == types.IsTruthy(v)
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
		case compiler.OP_LOADUNDEF:
			rb := base + compiler.GetArgB(inst)
			for r := ra; r <= rb; r++ {
				s.CallStack.Set(r, types.UndefinedObject)
			}
			if debug {
				fmt.Printf("%-20s ; R[%d] ... R[%d] = undefined\n", compiler.DumpInst(inst), ra, rb)
			}
		case compiler.OP_CALL, compiler.OP_TAILCALL:
			b := compiler.GetArgB(inst)
			s.CallStack.SetSp(ra + b - 1)
			obj := s.CallStack.Get(ra)
			if debug {
				fmt.Printf("%-20s ; R[%d] = %v(R[%d]...R[%d])\n", compiler.DumpInst(inst), ra, obj, ra+1, ra+b-1)
			}
			switch o := obj.(type) {
			case *types.Closure:
				curCi, err := s.precall(ra)
				if err != nil {
					return err
				}

				switch compiler.GetOpCode(inst) {
				case compiler.OP_CALL:
					if !curCi.Cl.IsGo {
						nexeccalls++
						goto reentry
					}
				case compiler.OP_TAILCALL:
					// precalled scheme closure
					if !curCi.Cl.IsGo {
						nargs := s.CallStack.Sp() - ra

						// pop current call info
						_ = s.CallInfos.Pop()
						prevCi := s.CallInfos.Top().(*types.CallInfo)

						// place the current closure and arguments to the previous closure sp
						s.CallStack.Set(prevCi.FuncSp, curCi.Cl)
						for i := 0; i < nargs; i++ {
							s.CallStack.Set(prevCi.FuncSp+i+1, s.CallStack.Get(curCi.FuncSp+i+1))
						}
						s.CallStack.SetSp(prevCi.FuncSp + nargs)

						// set information of tailcalling function to the previous call info
						prevCi.Cl = curCi.Cl
						prevCi.Pc = 0

						goto reentry
					}
				}
			case *types.Continuation:
				cont := o
				s.CallStack.Restore(cont.CallStack)
				s.CallInfos.Restore(cont.CallInfos)
				ci := s.CallInfos.Top().(*types.CallInfo)
				ci.Pc = cont.Pc
				nexeccalls = cont.NExecCalls
				nuated = true
				nuatedObj = s.CallStack.Get(ra + b - 1)
				goto reentry
			}
		case compiler.OP_CALLCC:
			if nuated {
				s.CallStack.Set(ra, nuatedObj)
				nuated = false

				if debug {
					fmt.Printf("%-20s ; R[%d] = %v\n", compiler.DumpInst(inst), ra, nuatedObj)
				}
			} else {
				callinfos := s.CallInfos.Store(s.CallInfos.Sp())
				callstack := s.CallStack.Store(s.CallStack.Sp())
				pc := ci.Pc - 1
				cont := types.NewContinuation(callinfos, callstack, pc, nexeccalls)
				// set the current continuation as an argument
				s.CallStack.Set(ra+1, cont)
				s.CallStack.SetSp(ra + 1)

				if debug {
					fmt.Printf("%-20s ; R[%d](%v)\n", compiler.DumpInst(inst), ra, cont)
				}

				// same with OP_CALL
				precalledCi, err := s.precall(ra)
				if err != nil {
					return err
				}
				if !precalledCi.Cl.IsGo {
					nexeccalls++
					goto reentry
				}
			}
		}
	}
	return nil
}
