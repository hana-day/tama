package types

import (
	"fmt"
)

type ObjectType int

const (
	TyNumber ObjectType = iota
	TyString
	TyClosure
	TyNil
	TySymbol
	TyPair
)

type Object interface {
	String() string
	Type() ObjectType
}

type Arrayable interface {
	Array() []Object
}

type ArrayableObject interface {
	Object
	Arrayable
}

type (
	Number  float64
	String  string
	Closure struct {
		IsGo bool

		// scheme closure only
		Proto *ClosureProto

		// go closure only
		Fn interface{}
	}
	Nil    struct{}
	Symbol struct {
		Name String
	}
	Pair struct {
		Car Object
		Cdr Object
	}
)

func (num Number) String() string {
	return fmt.Sprint(float64(num))
}

func (num Number) Type() ObjectType { return TyNumber }

func (s String) String() string {
	return string(s)
}

func (s String) Type() ObjectType {
	return TyString
}

type LocVar struct {
	Name  String
	Index int
}

type ClosureProto struct {
	Insts   []uint32
	Consts  []Object
	Args    []*Symbol
	LocVars map[String]*LocVar
	Protos  []*ClosureProto // function prototypes inside the function
}

func (cl *Closure) String() string {
	return "closure"
}

func (cl *Closure) Type() ObjectType {
	return TyClosure
}

func NewScmClosure() *Closure {
	return &Closure{
		IsGo: false,
	}
}

func NewGoClosure() *Closure {
	return &Closure{
		IsGo: true,
	}
}

func (s *Symbol) String() string {
	return s.Name.String()
}

func (n *Nil) String() string {
	return "()"
}

func (n *Nil) Type() ObjectType {
	return TyNil
}

func (n *Nil) Array() []Object {
	return []Object{}
}

var NilObject = &Nil{}

func (p *Pair) String() string {
	return fmt.Sprintf("(%s . %s)", p.Car.String(), p.Cdr.String())
}

func (p *Pair) Type() ObjectType {
	return TyPair
}

func (p *Pair) Array() []Object {
	arr := []Object{}
	pair := p
	for pair.Cdr.Type() != TyNil {
		arr = append(arr, pair.Car)
		pair = pair.Cdr.(*Pair)
	}
	arr = append(arr, pair.Car)
	return arr
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
