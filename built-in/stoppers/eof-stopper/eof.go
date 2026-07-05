package eofstopper

import (
	"minimal/minimal-core/built-in/lexer"
)

type EOFStopper struct {}

func NewEOFStopper() *EOFStopper {
	return &EOFStopper{}
}

func (e *EOFStopper) End(s *lexer.LexerJob) bool {
	return s.Position >= uint(len(s.Data))
}
