package tama

import (
	"fmt"
)

type TValueType int

const (
	TyNumber TValueType = iota
)

type TValue interface {
	String() string
	Type() TValueType
}

type TNumber float64

func (num TNumber) String() string {
	return fmt.Sprint(float64(num))
}

func (num TNumber) Type() TValueType { return TyNumber }
