package types

import (
	"fmt"
)

type ValueType int

const (
	TyNumber ValueType = iota
	TyString
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
