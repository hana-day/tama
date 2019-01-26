package tama

import (
	"github.com/hyusuk/tama/compiler"
	"github.com/hyusuk/tama/types"
)

func (s *State) OpenBase() *State {
	// set syntaxes
	for name, syntax := range compiler.DefaultSyntaxes {
		s.SetGlobal(name, syntax)
	}

	// set procedures
	s.RegisterFunc("+", 0, -1, genFnArith("+"))
	s.RegisterFunc("-", 1, -1, genFnArith("-"))
	s.RegisterFunc("*", 0, -1, genFnArith("*"))
	s.RegisterFunc("/", 1, -1, genFnArith("/"))
	s.RegisterFunc("cons", 2, 2, fnCons)
	s.RegisterFunc("car", 1, 1, fnCar)
	s.RegisterFunc("cdr", 1, 1, fnCdr)
	s.RegisterFunc("=", 2, -1, fnNumEq)
	return s
}

func fnCons(s *State, args []types.Object) (types.Object, error) {
	return types.Cons(args[0], args[1]), nil
}

func fnCar(s *State, args []types.Object) (types.Object, error) {
	if err := types.AssertType(types.TyPair, args[0]); err != nil {
		return nil, err
	}
	pair := args[0].(*types.Pair)
	return pair.Car(), nil
}

func fnCdr(s *State, args []types.Object) (types.Object, error) {
	if err := types.AssertType(types.TyPair, args[0]); err != nil {
		return nil, err
	}
	pair := args[0].(*types.Pair)
	return pair.Cdr(), nil
}

func fnNumEq(s *State, args []types.Object) (types.Object, error) {
	if err := types.AssertType(types.TyNumber, args...); err != nil {
		return nil, err
	}
	num := args[0].(types.Number)
	for _, arg := range args[1:] {
		num2 := arg.(types.Number)
		if num != num2 {
			return types.Boolean(false), nil
		}
	}
	return types.Boolean(true), nil
}

func genFnArith(name string) GoFunc {
	return func(s *State, args []types.Object) (types.Object, error) {
		if err := types.AssertType(types.TyNumber, args...); err != nil {
			return nil, err
		}
		result := args[0].(types.Number)
		for _, arg := range args[1:] {
			num := arg.(types.Number)
			switch name {
			case "+":
				result += num
			case "-":
				result -= num
			case "*":
				if num == 0 {
					return types.Number(0), nil
				}
				result *= num
			case "/":
				if num == 0 {
					return nil, types.NewInternalError("/: division by zero")
				}
				result /= num
			}
		}
		return result, nil
	}
}
