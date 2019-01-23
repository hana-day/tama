package types

import (
	"fmt"
)

type Pair struct {
	car Object
	cdr Object
}

func (p *Pair) String() string {
	return fmt.Sprintf("(%s . %s)", p.car.String(), p.cdr.String())
}

func (p *Pair) Type() ObjectType {
	return TyPair
}

func (p *Pair) Slice() ([]Object, error) {
	arr := []Object{}
	pair := p
	var ok bool
	for pair.cdr.Type() != TyNil {
		arr = append(arr, pair.car)
		pair, ok = pair.cdr.(*Pair)
		if !ok {
			return arr, fmt.Errorf("%v is not slicable object", p.String())
		}
	}
	arr = append(arr, pair.car)
	return arr, nil
}

func (p *Pair) Len() int {
	len := 1
	cdr := p.cdr
	for cdr.Type() == TyPair {
		len++
		p := cdr.(*Pair)
		cdr = p.cdr
	}
	return len
}

func (p *Pair) Car() Object {
	return p.car
}

func (p *Pair) Cdr() Object {
	return p.cdr
}

func (p *Pair) Cdar() (Object, error) {
	cdr, ok := p.cdr.(*Pair)
	if !ok {
		return nil, fmt.Errorf("pair required")
	}
	return cdr.car, nil
}

func (p *Pair) Cddr() (Object, error) {
	cdr, ok := p.cdr.(*Pair)
	if !ok {
		return nil, fmt.Errorf("pair required")
	}
	return cdr.cdr, nil
}

func (p *Pair) Cddar() (Object, error) {
	cddr, err := p.Cddr()
	if err != nil {
		return nil, err
	}
	pair, ok := cddr.(*Pair)
	if !ok {
		return nil, fmt.Errorf("pair required")
	}
	return pair.car, nil
}

func (p *Pair) Cdddr() (Object, error) {
	cddr, err := p.Cddr()
	if err != nil {
		return nil, err
	}
	pair, ok := cddr.(*Pair)
	if !ok {
		return nil, fmt.Errorf("pair required")
	}
	return pair.cdr, nil
}

func (p *Pair) First() Object {
	return p.Car()
}

func (p *Pair) Second() (Object, error) {
	return p.Cdar()
}

func (p *Pair) Third() (Object, error) {
	return p.Cddar()
}

func (s *Symbol) Type() ObjectType {
	return TySymbol
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
