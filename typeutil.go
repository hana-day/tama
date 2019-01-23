package tama

import (
	"github.com/hyusuk/tama/types"
)

func IsTruthy(obj types.Object) bool {
	b, ok := obj.(types.Boolean)
	if !ok || bool(b) {
		return true
	}
	return false
}

func IsFalse(obj types.Object) bool {
	return !IsTruthy(obj)
}
