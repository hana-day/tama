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
			return arr, fmt.Errorf("%v is not slicable object", p)
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
		return nil, fmt.Errorf("cdar: attempt to get car of %v", p.cdr)
	}
	return cdr.car, nil
}

func (p *Pair) Cddr() (Object, error) {
	cdr, ok := p.cdr.(*Pair)
	if !ok {
		return nil, fmt.Errorf("cddr: attempt to get cdr of %v", p.cdr)
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
		return nil, fmt.Errorf("cddar: attempt to get car of %v", cddr)
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
		return nil, fmt.Errorf("cdddr: attempt to get cdr of %v", cddr)
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
