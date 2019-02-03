package types

type Vector []Object

func (v Vector) Type() ObjectType {
	return TyVector
}

func (v Vector) String() string {
	return "vector"
}
