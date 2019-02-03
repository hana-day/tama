package types

type Continuation struct {
	CallInfos  *Stack
	CallStack  *Stack
	Pc         int
	NExecCalls int
}

func NewContinuation(callinfos, callstack *Stack, pc, nexeccalls int) *Continuation {
	return &Continuation{
		CallInfos:  callinfos,
		CallStack:  callstack,
		Pc:         pc,
		NExecCalls: nexeccalls,
	}
}

func (cont *Continuation) Type() ObjectType {
	return TyContinuation
}

func (cont *Continuation) String() string {
	return "continuation"
}
