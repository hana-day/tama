package tama

func runVM(s *State) {
	ci, _ := s.CallInfos.Top().(*CallInfo)
	cl := ci.Cl
	base := ci.Base
	for _, inst := range cl.Proto.Insts {
		ra := base + GetArgA(inst)
		switch GetOpCode(inst) {
		case LOADK:
			bx := GetArgBx(inst)
			s.CallStack.Set(ra, cl.Proto.Consts[bx])
		case GETGLOBAL:
			bx := GetArgBx(inst)
			s.CallStack.Set(ra, s.Global[cl.Proto.Consts[bx].String()])
		case MOVE:
			rb := base + GetArgB(inst)
			s.CallStack.Set(ra, s.CallStack.Get(rb))
		case CALL:
			b := GetArgB(inst)
			if b != 0 {
				s.CallStack.SetSp(ra + b)
			}
			s.precall(ra)
			base = s.Base
		case RETURN:
			b := GetArgB(inst)
			s.CallStack.dump()
			if b != 0 {
				s.CallStack.SetSp(ra + b - 1)
			}
		}
	}
}
