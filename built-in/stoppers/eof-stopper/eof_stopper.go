package eofstopper

import (
	"minimal/minimal-core/built-in/primitives"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
)

type EOFStopper struct {
	EOF tokenizerv2.TokenType
}

func NewEOFStopper(tokenizer tokenizerv2.Tokenizer) *EOFStopper {
	return &EOFStopper{
		tokenizer.NewTokenType(
			tokenizerv2.TokenTypeMetadata{DisplayName: "the end of the file", DebugName: "EOF"},
		),
	}
}

func (e *EOFStopper) End(s *tokenizerv2.TokenizerState) bool {
	if s.Position >= uint(len(s.Data)) {
		// TODO decide if we actually want an EOF token
		s.Emit(tokenizerv2.Token{Type: e.EOF, Value: "", Range: primitives.Range{Start: s.Position, Length: 0}})
		
		return true
	}

	return false
}
