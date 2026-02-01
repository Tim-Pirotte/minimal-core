package tokenizer

import (
	"minimal/minimal-core/domain"
)

type TokenizerConfig struct {
	matchers     []Matcher
	lastTokenType domain.TokenType
}

type Matcher interface {
	Match(s *Source) (length uint, tokenTypeToPushIfLargest domain.TokenType, tokenContent string)
}

func NewTokenizerConfig() TokenizerConfig {
	return TokenizerConfig{[]Matcher{}, domain.EOF}
}

func (t *TokenizerConfig) AddMatcher(matcher Matcher) {
	t.matchers = append(t.matchers, matcher)
}

func (t *TokenizerConfig) NewTokenType() domain.TokenType {
	t.lastTokenType++
	return t.lastTokenType
}

type Source struct {
	source []byte
	position uint
}

// Retrieves a byte at position n relative to the current position from the source.
// Ok wil be set to false if it is out of bounds
func (s *Source) Get(n int) (byte, bool) {
	index := int(s.position) + n

	if 0 <= index && index < len(s.source) {
		return s.source[index], true
	} else {
		return 0, false
	}
}

func (t *TokenizerConfig) tokenize(source []byte) []domain.Token {
	s := Source{source, 0}
	tokens := make([]domain.Token, 0)

	for s.position < uint(len(s.source)) {
		foundMatch := false
		largestLength := uint(0)
		var tokenTypeToPush domain.TokenType
		var tokenContent string

		for _, matcher := range t.matchers {
			length, tt, content := matcher.Match(&s)

			if length > largestLength {
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
						Span: domain.Span{Start: s.position, Length: largestLength}},
				)
			}

			s.position += largestLength
		} else {
			tokens = append(
				tokens, 
				domain.Token{
					Type: domain.UNKNOWN, 
					Value: string(s.source[s.position]),
					Span: domain.Span{Start: s.position, Length: 1},
				},
			)

			s.position++
		}
	}

	return append(
		tokens, 
		domain.Token{Type: domain.EOF, Value: "", Span: domain.Span{Start: s.position, Length: 0}},
	)
}
