package types

import (
	"fmt"
)

type Pair struct {
	Car Object
	Cdr Object
}

func (p *Pair) String() string {
	return fmt.Sprintf("(%s . %s)", p.Car.String(), p.Cdr.String())
}

func (p *Pair) Type() ObjectType {
	return TyPair
}

func (p *Pair) Slice() ([]Object, error) {
	arr := []Object{}
	pair := p
	var ok bool
	for pair.Cdr.Type() != TyNil {
		arr = append(arr, pair.Car)
		pair, ok = pair.Cdr.(*Pair)
		if !ok {
			return arr, fmt.Errorf("%v is not slicable object", p.String())
		}
	}
	arr = append(arr, pair.Car)
	return arr, nil
}

func (p *Pair) Len() int {
	len := 1
	cdr := p.Cdr
	for cdr.Type() == TyPair {
		len++
		p := cdr.(*Pair)
		cdr = p.Cdr
	}
	return len
}

func (s *Symbol) Type() ObjectType {
	return TySymbol
}

func Cons(car Object, cdr Object) *Pair {
	return &Pair{
		Car: car,
		Cdr: cdr,
	}
}

func Car(v Object) (Object, error) {
	p, ok := v.(*Pair)
	if !ok {
		return nil, fmt.Errorf("car: %v is not a pair", v)
	}
	return p.Car, nil
}

func Cdr(v Object) (Object, error) {
	p, ok := v.(*Pair)
	if !ok {
		return nil, fmt.Errorf("cdr: %v is not a pair", v)
	}
	return p.Cdr, nil
}

func Cdar(v Object) (Object, error) {
	cdr, err := Cdr(v)
	if err != nil {
		return nil, err
	}
	return Car(cdr)
}

func Cddr(v Object) (Object, error) {
	cdr, err := Cdr(v)
	if err != nil {
		return nil, err
	}
	return Cdr(cdr)
}

func List(args ...Object) Object {
	if len(args) == 0 {
		return NilObject
	}
	return &Pair{
		Car: args[0],
		Cdr: List(args[1:len(args)]...),
	}
}
