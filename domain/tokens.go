package domain

type TokenType uint

const EOF TokenType = 0

type Token struct {
	Type  TokenType
	Value string
	Span  Span
}
