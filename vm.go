package tama

func runVM(s *State) {
	ci, _ := s.CallInfos.Top().(*CallInfo)
	cl := ci.Cl
	for _, inst := range cl.Proto.Insts {
		ra := ci.Base + GetArgA(inst)
		switch GetOpCode(inst) {
		case LOADK:
			bx := GetArgBx(inst)
			s.CallStack.Set(ra, cl.Proto.Consts[bx])
		case GETGLOBAL:
			bx := GetArgBx(inst)
			s.CallStack.Set(ra, s.Global[cl.Proto.Consts[bx].String()])
		case RETURN:
			b := GetArgB(inst)
			if b != 0 {
				s.CallStack.SetSp(ra + b - 1)
			}
		}
	}
}
