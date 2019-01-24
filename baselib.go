package tama

import (
	"fmt"
	"github.com/hyusuk/tama/compiler"
	"github.com/hyusuk/tama/types"
)

func (s *State) OpenBase() *State {
	// set syntaxes
	for name, syntax := range compiler.DefaultSyntaxes {
		s.SetGlobal(name, syntax)
	}

	// set procedures
	s.RegisterFunc("+", genFnArith("+"))
	s.RegisterFunc("-", genFnArith("-"))
	s.RegisterFunc("*", genFnArith("*"))
	s.RegisterFunc("/", genFnArith("/"))
	s.RegisterFunc("cons", fnCons)
	s.RegisterFunc("car", fnCar)
	s.RegisterFunc("cdr", fnCdr)
	s.RegisterFunc("=", fnNumEq)
	return s
}

func fnCons(s *State, args []types.Object) (types.Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("cons: invalid syntax")
	}
	return types.Cons(args[0], args[1]), nil
}

func fnCar(s *State, args []types.Object) (types.Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("car: invalid syntax")
	}
	pair, ok := args[0].(*types.Pair)
	if !ok {
		return nil, fmt.Errorf("car: invalid value")
	}
	return pair.Car(), nil
}

func fnCdr(s *State, args []types.Object) (types.Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("cdr: invalid syntax")
	}
	pair, ok := args[0].(*types.Pair)
	if !ok {
		return nil, fmt.Errorf("cdr: invalid value")
	}
	return pair.Cdr(), nil
}

func fnNumEq(s *State, args []types.Object) (types.Object, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("=: insufficient number of arguments")
	}
	num, ok := args[0].(types.Number)
	if !ok {
		return nil, fmt.Errorf("=: non-numerical argument")
	}
	for _, arg := range args[1:] {
		num2, ok := arg.(types.Number)
		if !ok {
			return nil, fmt.Errorf("=: non-numerical argument")
		}
		if num != num2 {
			return types.Boolean(false), nil
		}
	}
	return types.Boolean(true), nil
}

func genFnArith(name string) GoFunc {
	return func(s *State, args []types.Object) (types.Object, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("%s: insufficient number of arguments", name)
		}
		result, ok := args[0].(types.Number)
		if !ok {
			return nil, fmt.Errorf("%s: non-numerical argument", name)
		}
		for _, arg := range args[1:] {
			num, ok := arg.(types.Number)
			if !ok {
				return nil, fmt.Errorf("%s: non-numerical argument", name)
			}
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
					return nil, fmt.Errorf("/: division by zero")
				}
				result /= num
			}
		}
		return result, nil
	}
}
