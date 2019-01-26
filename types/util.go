package types

func IsTruthy(obj Object) bool {
	b, ok := obj.(Boolean)
	if !ok || bool(b) {
		return true
	}
	return false
}

func IsFalse(obj Object) bool {
	return !IsTruthy(obj)
}

func IsNull(obj Object) bool {
	return obj.Type() == TyNil
}

func IsPair(obj Object) bool {
	return obj.Type() == TyPair
}

func IsNumber(obj Object) bool {
	return obj.Type() == TyNumber
}

func IsString(obj Object) bool {
	return obj.Type() == TyString
}

func IsSymbol(obj Object) bool {
	return obj.Type() == TySymbol
}

func IsClosure(obj Object) bool {
	return obj.Type() == TyClosure
}

func IsList(obj Object) bool {
	if IsNull(obj) {
		return true
	}
	o := obj
	pair, ok := o.(*Pair)
	for ok {
		o = pair.Cdr()
		pair, ok = o.(*Pair)
	}
	return IsNull(o)
}

func Cons(car Object, cdr Object) *Pair {
	return &Pair{
		car: car,
		cdr: cdr,
	}
}

func List(args ...Object) Object {
	if len(args) == 0 {
		return NilObject
	}
	return &Pair{
		car: args[0],
		cdr: List(args[1:len(args)]...),
	}
}

func AssertType(typ ObjectType, objs ...Object) error {
	for _, obj := range objs {
		if obj.Type() != typ {
			return NewTypeError("%s required, but got %v", typeProps[typ].name, obj)
		}
	}
	return nil
}
