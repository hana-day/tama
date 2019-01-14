package types

import (
	"fmt"
)

type ValueType int

const (
	TyNumber ValueType = iota
	TyString
	TyClosure
	TyNil
	TySymbol
	TyPair
)

type Value interface {
	String() string
	Type() ValueType
}

type (
	Number  float64
	String  string
	Closure struct {
		IsGo bool

		// scheme closure only
		Proto *ClosureProto

		// go closure only
		Fn interface{}
	}
	NilType struct{}
	Symbol  struct {
		Name String
	}
	Pair struct {
		Car Value
		Cdr Value
	}
)

func (num Number) String() string {
	return fmt.Sprint(float64(num))
}

func (num Number) Type() ValueType { return TyNumber }

func (s String) String() string {
	return string(s)
}

func (s String) Type() ValueType {
	return TyString
}

type ClosureProto struct {
	Insts        []uint32
	Consts       []Value
	MaxStackSize int
}

func (cl *Closure) String() string {
	return "closure"
}

func (cl *Closure) Type() ValueType {
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

func (s *Symbol) String() string {
	return s.Name.String()
}

func (n *NilType) String() string {
	return "()"
}

func (n *NilType) Type() ValueType {
	return TyNil
}

var Nil = &NilType{}

func (p *Pair) String() string {
	return fmt.Sprintf("(%s . %s)", p.Car.String(), p.Cdr.String())
}

func (p *Pair) Type() ValueType {
	return TyPair
}

func (s *Symbol) Type() ValueType {
	return TySymbol
}

func Cons(car Value, cdr Value) *Pair {
	return &Pair{
		Car: car,
		Cdr: cdr,
	}
}

func Car(v Value) (Value, error) {
	p, ok := v.(*Pair)
	if !ok {
		return nil, fmt.Errorf("%v is not a pair", v)
	}
	return p.Car, nil
}

func Cdr(v Value) (Value, error) {
	p, ok := v.(*Pair)
	if !ok {
		return nil, fmt.Errorf("%v is not a pair", v)
	}
	return p.Cdr, nil
}

func List(args ...Value) Value {
	if len(args) == 0 {
		return Nil
	}
	return &Pair{
		Car: args[0],
		Cdr: List(args[1:len(args)]...),
	}
}
