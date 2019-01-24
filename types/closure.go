package types

import "fmt"

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

func (uv *UpValue) Set(callStack *Stack, value Object) {
	if uv.Closed {
		uv.obj = value
	} else {
		callStack.Set(uv.Index, value)
	}
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
	Fn     interface{}
	fnName string

	UpVals []*UpValue
}

type ArgMode int

const (
	// (lambda (x y) ...)
	FixedArgMode ArgMode = iota
	// (lambda args ...)
	VArgMode
	// (lambda (x y . rest) ...)
	RestArgMode
)

type ClosureProto struct {
	Insts   []uint32
	Consts  []Object
	Args    []*Symbol
	Protos  []*ClosureProto // function prototypes inside the function
	NUpVals int
	Mode    ArgMode
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
		return fmt.Sprintf("closure (go func '%s')", cl.fnName)
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

func NewGoClosure(name string, fn interface{}) *Closure {
	return &Closure{
		IsGo:   true,
		Fn:     fn,
		fnName: name,
	}
}

type CallInfo struct {
	FuncSp int // function sp
	Base   int // local sp
	Cl     *Closure
	Pc     int
}

func (ci *CallInfo) Type() ObjectType {
	return TyCallInfo
}

func (ci *CallInfo) String() string {
	return fmt.Sprintf("call info for %s", ci.Cl.String())
}
