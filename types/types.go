package types

import (
	"fmt"
)

type ValueType int

const (
	TyNumber ValueType = iota
)

type Value interface {
	String() string
	Type() ValueType
}

type Number float64

func (num Number) String() string {
	return fmt.Sprint(float64(num))
}

func (num Number) Type() ValueType { return TyNumber }
