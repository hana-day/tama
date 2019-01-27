package scanner

type Token int

const (
	ILLEGAL Token = iota
	EOF
	NUMBER
	LPAREN // "("
	RPAREN // ")"
	IDENT
	QUOTE // "'"
	TRUE  // "#t"
	FALSE // "#f"
	STRING
)
