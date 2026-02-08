package domain

type TokenType uint

const (
	UNKNOWN TokenType = iota
	IGNORE
	// Add new tokens before this one.
	// The reason for this is that the value of EOF will be used as
	// the starting value to increment for new tokens
	EOF
)

type Token struct {
	Type  TokenType
	Value string
	Span  Span
}
