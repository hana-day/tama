package types

import (
	"fmt"
)

type ValueType int

const (
	TyNumber ValueType = iota
	TyString
	TyClosure
)

type Value interface {
	String() string
	Type() ValueType
}

type Number float64
type String string

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

type Closure struct {
	IsGo bool

	// scheme closure only
	Proto *ClosureProto

	// go closure only
	Fn interface{}
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
