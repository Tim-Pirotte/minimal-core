package tokenizer

import (
	"minimal/minimal-core/domain"
)

type TokenizerConfig struct {
	source       []byte
	position     uint
	matchers     []Matcher
	lastTokenType domain.TokenType
}

type Matcher interface {
	Match(t *TokenizerConfig) (valid bool, length uint, tokenTypeToPushIfLargest domain.TokenType, tokenContent string)
}

func NewTokenizerConfig(source []byte) TokenizerConfig {
	return TokenizerConfig{source, 0, []Matcher{}, domain.EOF}
}

func (t *TokenizerConfig) AddMatcher(matcher Matcher) {
	t.matchers = append(t.matchers, matcher)
}

func (t *TokenizerConfig) NewTokenType() domain.TokenType {
	t.lastTokenType++
	return t.lastTokenType
}

// Retrieves a byte at position n relative to the current position from the source.
// Ok wil be set to false if it is out of bounds
func (t *TokenizerConfig) Get(n int) (byte, bool) {
	index := int(t.position) + n

	if 0 <= index && index < len(t.source) {
		return t.source[index], true
	} else {
		return 0, false
	}
}

func (t *TokenizerConfig) tokenize() []domain.Token {
	tokens := make([]domain.Token, 0)

	for t.position < uint(len(t.source)) {
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
			if tokenTypeToPush != domain.IGNORE {
				tokens = append(
					tokens, 
					domain.Token{
						Type: tokenTypeToPush, 
						Value: tokenContent, 
						Span: domain.Span{Start: t.position, Length: largestLength}},
				)
			}

			t.position += largestLength
		} else {
			tokens = append(
				tokens, 
				domain.Token{
					Type: domain.UNKNOWN, 
					Value: string(t.source[t.position]),
					Span: domain.Span{Start: t.position, Length: 1},
				},
			)

			t.position++
		}
	}

	return append(
		tokens, 
		domain.Token{Type: domain.EOF, Value: "", Span: domain.Span{Start: t.position, Length: 0}},
	)
}
