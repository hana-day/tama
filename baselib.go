package tama

import (
	"fmt"
	"github.com/hyusuk/tama/types"
)

func (s *State) OpenBase() *State {
	s.RegisterFunc("+", fnAdd)
	s.RegisterFunc("cons", fnCons)
	s.RegisterFunc("car", fnCar)
	s.RegisterFunc("cdr", fnCdr)
	return s
}

func fnAdd(s *State, args []types.Object) (types.Object, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("+: invalid syntax")
	}
	var result types.Number = 0
	for i := 0; i < len(args); i++ {
		num, ok := args[i].(types.Number)
		if !ok {
			return nil, fmt.Errorf("+: invalid value %v", args[i])
		}
		result += num
	}
	return result, nil
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
	return pair.Car, nil
}

func fnCdr(s *State, args []types.Object) (types.Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("cdr: invalid syntax")
	}
	pair, ok := args[0].(*types.Pair)
	if !ok {
		return nil, fmt.Errorf("cdr: invalid value")
	}
	return pair.Cdr, nil
}
