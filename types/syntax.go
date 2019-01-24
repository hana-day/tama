package types

import "fmt"

type Syntax struct {
	Name String
	Fn   interface{}
}

func NewSyntax(name string, fn interface{}) *Syntax {
	return &Syntax{String(name), fn}
}

func (s *Syntax) Type() ObjectType {
	return TySyntax
}

func (s *Syntax) String() string {
	return fmt.Sprintf("syntax %s", s.Name)
}
