package parser

import (
	"github.com/hyusuk/tama/token"
)

type Node interface {
}

type Expr interface {
	Node
	exprNode()
}

func (*Primitive) exprNode() {}

type (
	Primitive struct {
		Kind  token.Token // token.INT or ?
		Value string      // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
	}
)
