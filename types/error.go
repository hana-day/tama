package types

import "fmt"

type ErrorType int

const (
	// syntaxError is thrown when the tama engine encounters invald syntaxes.
	ErrSyntax ErrorType = iota
	// internalError indicates an error that occured internally in tama engine.
	ErrInternal
	// typeError is thrown when a object is not of the expected type.
	ErrType
)

type Error struct {
	s       string
	errType ErrorType
}

func NewSyntaxError(s string, v ...interface{}) *Error {
	return &Error{s: fmt.Sprintf(s, v...), errType: ErrSyntax}
}

func NewInternalError(s string, v ...interface{}) *Error {
	return &Error{s: fmt.Sprintf(s, v...), errType: ErrInternal}
}

func NewTypeError(s string, v ...interface{}) *Error {
	return &Error{s: fmt.Sprintf(s, v...), errType: ErrType}
}

func (e *Error) Type() ObjectType {
	return TyError
}

func (e *Error) Error() string {
	return e.s
}
