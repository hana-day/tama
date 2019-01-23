package types

import (
	"fmt"
)

type ObjectType int

const (
	TyNumber ObjectType = iota
	TyString
	TyClosure
	TyNil
	TySymbol
	TyPair
	TyBoolean

	TyCallInfo // for internal use
)

type Object interface {
	String() string
	Type() ObjectType
}

type Slicable interface {
	Slice() ([]Object, error)
}

type SlicableObject interface {
	Object
	Slicable
}

type (
	Number float64
	String string
	Nil    struct{}
	Symbol struct {
		Name String
	}
	Boolean bool
)

func (num Number) String() string {
	return fmt.Sprint(float64(num))
}

func (num Number) Type() ObjectType { return TyNumber }

func (s String) String() string {
	return string(s)
}

func (s String) Type() ObjectType {
	return TyString
}

func (s *Symbol) String() string {
	return s.Name.String()
}

func (n *Nil) String() string {
	return "()"
}

func (n *Nil) Type() ObjectType {
	return TyNil
}

func (n *Nil) Slice() ([]Object, error) {
	return []Object{}, nil
}

var NilObject = &Nil{}

func (b Boolean) Type() ObjectType {
	return TyBoolean
}

func (b Boolean) String() string {
	if b {
		return "#t"
	}
	return "#f"
}
