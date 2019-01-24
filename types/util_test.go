package types

import (
	"testing"
)

func TestIs(t *testing.T) {
	testcases := []struct {
		obj    Object
		expect bool
		fn     func(Object) bool
	}{
		// test IsTruthy
		{Boolean(true), true, IsTruthy},
		{Number(0), true, IsTruthy},
		{Number(1), true, IsTruthy},
		{NilObject, true, IsTruthy},
		{Boolean(false), false, IsTruthy},

		// test IsFalse
		{Boolean(false), true, IsFalse},
		{Number(0), false, IsFalse},
		{Boolean(true), false, IsFalse},

		// test  IsNull
		{NilObject, true, IsNull},
		{Boolean(false), false, IsNull},

		// test IsPair
		{Cons(Number(1), Number(2)), true, IsPair},
		{Number(1), false, IsPair},

		// test IsNumber
		{Number(1), true, IsNumber},
		{Boolean(true), false, IsNumber},

		// test IsString
		{String("a"), true, IsString},
		{Number(1), false, IsString},

		// test IsSymbol
		{NewSymbol("a"), true, IsSymbol},
		{String("a"), false, IsSymbol},

		// test IsClosure
		{NewScmClosure(nil, 0), true, IsClosure},
		{NewGoClosure("a", nil), true, IsClosure},
		{Number(1), false, IsClosure},

		// test IsList
		{NilObject, true, IsList},
		{List(Number(1), Number(2)), true, IsList},
		{Cons(Number(1), NilObject), true, IsList},
		{Number(1), false, IsList},
		{Cons(Number(1), Number(2)), false, IsList},
	}
	for i, tc := range testcases {
		actual := tc.fn(tc.obj)
		if actual != tc.expect {
			t.Fatalf("case %d: expected %t, but got %t", i, tc.expect, actual)
		}
	}
}
