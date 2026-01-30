package tokenizer

import (
	"minimal/minimal-core/domain"
)

type Tokenizer struct {
	tokens       []domain.Token
	position     int
}

func NewTokenizer(config TokenizerConfig) Tokenizer {
	return Tokenizer{
		config.tokenize(),
		0,
	}
}

// Returns the current token
func (t *Tokenizer) CurrentToken() domain.Token {
	return t.tokens[t.position]
}

// Goes to the next token or logs an error if the EOF token is consumed
func (t *Tokenizer) Consume() {
	if t.position + 1 == len(t.tokens) {
		// TODO log error
	}

	t.position++
}

// Look ahead n tokens with 0 being the current token.
// Returns EOF if n goes out of bounds
func (t *Tokenizer) Peek(n int) domain.Token {
	if t.position + n >= len(t.tokens) {
		return t.tokens[len(t.tokens) - 1]
	}

	return t.tokens[t.position + n]
}
