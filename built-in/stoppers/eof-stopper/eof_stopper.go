package eofstopper

import (
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
)

type EOFStopper struct {}

func NewEOFStopper() *EOFStopper {
	return &EOFStopper{}
}

func (e *EOFStopper) End(s *tokenizerv2.TokenizerJob) bool {
	return s.Position >= uint(len(s.Data))
}
