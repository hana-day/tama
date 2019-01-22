package types

type UpValue struct {
	Next   *UpValue
	Index  int
	Closed bool
	obj    Object
}

func (uv *UpValue) Value(callStack *Stack) Object {
	if uv.Closed || callStack == nil {
		return uv.obj
	}
	return callStack.Get(uv.Index).(Object)
}

func (uv *UpValue) Close(callStack *Stack) {
	value := uv.Value(callStack).(Object)
	uv.obj = value
	uv.Closed = true
}

type Closure struct {
	IsGo bool

	// scheme closure only
	Proto *ClosureProto

	// go closure only
	Fn interface{}

	UpVals []*UpValue
}

type ClosureProto struct {
	Insts   []uint32
	Consts  []Object
	Args    []*Symbol
	Protos  []*ClosureProto // function prototypes inside the function
	NUpVals int
}

func NewClosureProto() *ClosureProto {
	return &ClosureProto{
		Insts:   []uint32{},
		Consts:  []Object{},
		NUpVals: 0,
	}
}

func (cl *Closure) String() string {
	if cl.IsGo {
		return "closure (go func)"
	}
	return "closure"
}

func (cl *Closure) Type() ObjectType {
	return TyClosure
}

func NewScmClosure(proto *ClosureProto, nUpVals int) *Closure {
	return &Closure{
		Proto:  proto,
		IsGo:   false,
		UpVals: make([]*UpValue, nUpVals),
	}
}

func NewGoClosure(fn interface{}) *Closure {
	return &Closure{
		IsGo: true,
		Fn:   fn,
	}
}
