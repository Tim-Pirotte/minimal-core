package domain

type TokenType uint

const (
	UNKNOWN TokenType = iota
	IGNORE
	EOF // Add new tokens before this one
)

type Token struct {
	Type  TokenType
	Value string
	Span  Span
}
