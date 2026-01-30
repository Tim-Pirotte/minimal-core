package tokenizer

import (
	"bufio"
	"io"
	"minimal/minimal-core/domain"
)

type TokenizerConfig struct {
	reader   *bufio.Reader
	position uint
	matchers     []Matcher
	lastTokenType domain.TokenType
}

type Matcher interface {
	Match(t *TokenizerConfig) (valid bool, length uint, tokenTypeToPushIfLargest domain.TokenType, tokenContent string)
}

func NewTokenizerConfig(reader *bufio.Reader) TokenizerConfig {
	return TokenizerConfig{reader, 0, []Matcher{}, domain.EOF}
}

func (t *TokenizerConfig) AddMatcher(matcher Matcher) {
	t.matchers = append(t.matchers, matcher)
}

func (t *TokenizerConfig) NewTokenType() domain.TokenType {
	t.lastTokenType++
	return t.lastTokenType
}

func (t *TokenizerConfig) tokenize() []domain.Token {
	tokens := make([]domain.Token, 0)

	for _, err := t.reader.Peek(1); err != io.EOF; {
		if err != nil {
			// TODO log error
		}

		foundMatch := false
		largestLength := uint(0)
		var tokenTypeToPush domain.TokenType
		var tokenContent string

		for _, matcher := range t.matchers {
			valid, length, tt, content := matcher.Match(t)

			if valid && length > largestLength {
				foundMatch = true
				largestLength = length
				tokenTypeToPush = tt
				tokenContent = content
			}
		}

		if foundMatch {
			tokens = append(
				tokens, 
				domain.Token{
					Type: tokenTypeToPush, 
					Value: tokenContent, 
					Span: domain.Span{Start: t.position, Length: largestLength}},
			)
		} else {
			byteToSkip, err := t.reader.ReadByte()

			if err != nil {
				// TODO log error
			}

			tokens = append(
				tokens, 
				domain.Token{Type: domain.UNKNOWN, Value: string(byteToSkip), Span: domain.Span{Start: t.position, Length: 1}},
			)
		}
	}

	return tokens
}
