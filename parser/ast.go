package parser

import (
	"github.com/hyusuk/tama/scanner"
)

type Node interface {
}

type Expr interface {
	Node
	exprNode()
}

func (*Primitive) exprNode() {}
func (*Ident) exprNode()     {}
func (*CallExpr) exprNode()  {}

type (
	Primitive struct {
		Kind  scanner.Token // token.INT or ?
		Value string        // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
	}

	Ident struct {
		Name string
	}

	CallExpr struct {
		Func Expr
		Args []Expr
	}
)
