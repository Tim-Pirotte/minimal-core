package tokenizer

import (
	"bufio"
	"fmt"
	"minimal/minimal-core/domain"
)

type Tokenizer struct {
	CurrentToken domain.Token
	reader       bufio.Reader
	position     uint

	keywords *trieNode
	symbols  *trieNode
	checks   []checkRule

	lastTokenType domain.TokenType
}

type trieNode struct {
	leaf bool
	token domain.TokenType
	children map[rune]*trieNode
}

type checkRule struct {
	check func(rune) bool
	tokenType domain.TokenType
}

func NewTokenizer(reader bufio.Reader) Tokenizer {
	return Tokenizer{
		domain.Token{}, 
		reader,
		0,
		&trieNode{children: map[rune]*trieNode{}},
		&trieNode{children: map[rune]*trieNode{}},
		[]checkRule{},
		0,
	}
}

func (t *Tokenizer) newTokenType() domain.TokenType {
	t.lastTokenType++
	return t.lastTokenType
}

func (t *Tokenizer) AddKeyword(keyword string) domain.TokenType {
	tokenType := t.newTokenType()
	err := updateTrie(t.keywords, keyword, tokenType)

	if err != nil {
		// TODO log error
		fmt.Println(err.Error())
	}

	return tokenType
}

func (t *Tokenizer) AddSymbol(symbol string) domain.TokenType {
	tokenType := t.newTokenType()
	err := updateTrie(t.symbols, symbol, tokenType)

	if err != nil {
		// TODO log error
		fmt.Println(err.Error())
	}
	
	return tokenType
}

func (t *Tokenizer) AddCheck(check func(rune) bool) domain.TokenType {
	tokenType := t.newTokenType()
	t.checks = append(t.checks, checkRule{check, tokenType})

	return tokenType
}

func (t *Tokenizer) Consume() {
	t.skipWhiteSpace()
}

func (t *Tokenizer) setToEof(start uint) {
	t.CurrentToken.Type = domain.EOF
	t.CurrentToken.Value = ""
	t.CurrentToken.Span.Start = start
	t.CurrentToken.Span.Length = 0
}

func (t *Tokenizer) skipWhiteSpace() {
	currentByte, err := t.reader.Peek(1)

	for err == nil && (currentByte[0] == ' ' || currentByte[0] == '\t') {
		_, discardErr := t.reader.Discard(1)

		if discardErr != nil {
			panic(fmt.Sprint("Could not properly discard:", discardErr))
		}

		currentByte, err = t.reader.Peek(1)
	}
}

const asciiMax = 127

func isAlphanumericOrUnicode(char rune) bool {
	return 'a' >= char && char <= 'z' || 'A' >= char && char <= 'Z' || '0' >= char && char <= '9' || char > asciiMax
}
