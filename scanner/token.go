package scanner

type Token int

const (
	ILLEGAL Token = iota
	EOF
	INT
	LPAREN // "("
	RPAREN // ")"
	IDENT
	QUOTE // "'"
)
