package domain

type TokenType uint

const (
	UNKNOWN TokenType = iota
	EOF
)

type Token struct {
	Type  TokenType
	Value string
	Span  Span
}
