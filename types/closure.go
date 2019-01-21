package types

type Closure struct {
	IsGo bool

	// scheme closure only
	Proto *ClosureProto

	// go closure only
	Fn interface{}
}

type ClosureProto struct {
	Insts  []uint32
	Consts []Object
	Args   []*Symbol
	Protos []*ClosureProto // function prototypes inside the function
}

func NewClosureProto() *ClosureProto {
	return &ClosureProto{
		Insts:  []uint32{},
		Consts: []Object{},
	}
}

func (cl *Closure) String() string {
	return "closure"
}

func (cl *Closure) Type() ObjectType {
	return TyClosure
}

func NewScmClosure() *Closure {
	return &Closure{
		IsGo: false,
	}
}

func NewGoClosure() *Closure {
	return &Closure{
		IsGo: true,
	}
}
