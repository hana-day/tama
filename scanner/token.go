package scanner

type Token int

const (
	EOF Token = iota
	INT
	LPAREN // "("
	RPAREN // ")"
	IDENT
)
