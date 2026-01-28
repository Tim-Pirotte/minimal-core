package tokenizer

import (
	"io"
	"minimal/minimal-core/domain"
)

type Tokenizer struct {
	CurrentToken domain.Token
	reader       io.Reader
	position     uint

	keywords *node
	symbols  *node
	checks   []checkRule

	lastTokenType domain.TokenType
}

type node struct {
	token domain.TokenType
	children map[rune]*node
}

type checkRule struct {
	check func(rune) bool
	tokenType domain.TokenType
}

func NewTokenizer(reader io.Reader) Tokenizer {
	return Tokenizer{
		domain.Token{}, 
		reader,
		0,
		&node{children: map[rune]*node{}},
		&node{children: map[rune]*node{}},
		[]checkRule{},
		0,
	}
}

func (t *Tokenizer) newTokenType() domain.TokenType {
	t.lastTokenType++
	return t.lastTokenType
}

func (t *Tokenizer) AddKeyword(keyword string) domain.TokenType {
	return t.newTokenType()
}

func (t *Tokenizer) AddSymbol(symbol string) domain.TokenType {
	return t.newTokenType()
}

func (t *Tokenizer) AddCheck(check func(byte) bool) domain.TokenType {
	return t.newTokenType()
}

func (t *Tokenizer) Consume() {

}
