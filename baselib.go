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
	s.RegisterFunc("+", 0, -1, fnAdd)
	s.RegisterFunc("-", 1, -1, fnSub)
	s.RegisterFunc("*", 0, -1, fnMul)
	s.RegisterFunc("/", 1, -1, fnDiv)
	s.RegisterFunc("cons", 2, 2, fnCons)
	s.RegisterFunc("car", 1, 1, fnCar)
	s.RegisterFunc("cdr", 1, 1, fnCdr)
	s.RegisterFunc("=", 2, -1, fnNumEq)
	s.RegisterFunc("<", 2, -1, genFnComp("<"))
	s.RegisterFunc(">", 2, -1, genFnComp(">"))
	s.RegisterFunc("<=", 2, -1, genFnComp("<="))
	s.RegisterFunc(">=", 2, -1, genFnComp(">="))
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

func fnAdd(s *State, args []types.Object) (types.Object, error) {
	if err := types.AssertType(types.TyNumber, args...); err != nil {
		return nil, err
	}
	var result types.Number = 0
	for _, arg := range args {
		num := arg.(types.Number)
		result += num
	}
	return result, nil
}

func fnSub(s *State, args []types.Object) (types.Object, error) {
	if err := types.AssertType(types.TyNumber, args...); err != nil {
		return nil, err
	}
	num := args[0].(types.Number)
	if len(args) == 1 {
		return -1 * num, nil
	}
	result := num
	for _, arg := range args[1:] {
		num := arg.(types.Number)
		result -= num
	}
	return result, nil
}

func fnMul(s *State, args []types.Object) (types.Object, error) {
	if err := types.AssertType(types.TyNumber, args...); err != nil {
		return nil, err
	}
	result := types.Number(1)
	for _, arg := range args {
		num := arg.(types.Number)
		if num == 0 {
			return types.Number(0), nil
		}
		result *= num
	}
	return result, nil
}

func fnDiv(s *State, args []types.Object) (types.Object, error) {
	if err := types.AssertType(types.TyNumber, args...); err != nil {
		return nil, err
	}
	num := args[0].(types.Number)
	if len(args) == 1 {
		return 1 / num, nil
	}
	result := num
	for _, arg := range args[1:] {
		num := arg.(types.Number)
		if num == 0 {
			return nil, types.NewInternalError("division by zero")
		}
		result /= num
	}
	return result, nil
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

func genFnComp(name string) GoFunc {
	return func(s *State, args []types.Object) (types.Object, error) {
		if err := types.AssertType(types.TyNumber, args...); err != nil {
			return nil, err
		}
		prev := args[0].(types.Number)
		var yes bool
		for _, arg := range args[1:] {
			next := arg.(types.Number)
			switch name {
			case "<":
				yes = prev < next
			case ">":
				yes = prev > next
			case "<=":
				yes = prev <= next
			case ">=":
				yes = prev >= next
			}
			prev = next
			if !yes {
				return types.Boolean(false), nil
			}
		}
		return types.Boolean(true), nil
	}
}
