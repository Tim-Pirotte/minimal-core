package tokenizerv2

import (
	"math"
	"minimal/minimal-core/built-in/primitives"
)

type TokenType uint

const (
	UNKNOWN TokenType = math.MaxUint - iota
)

type Token struct {
	Type  TokenType
	Value string
	Range primitives.Range
}

type TokenTypeMetadata struct {
	DisplayName string
	DebugName   string
}

type Tokenizer struct {
	matchers      []Matcher
	lastTokenType TokenType
	tokenTypesMetadata map[TokenType]TokenTypeMetadata
}

type Matcher interface {
	Match(s *TokenizerState) (length uint)
	Consume(s *TokenizerState, length uint)
}

type Stopper interface {
	End(s *TokenizerState) bool
}

func NewTokenizer() Tokenizer {
	return Tokenizer{
		[]Matcher{}, 
		TokenType(0),
		map[TokenType]TokenTypeMetadata{UNKNOWN: {"a character that is not a valid token", "UNKNOWN"}},
	}
}

func (t *Tokenizer) AddMatcher(matcher Matcher) {
	t.matchers = append(t.matchers, matcher)
}

func (t *Tokenizer) NewTokenType(metadata TokenTypeMetadata) TokenType {
	t.lastTokenType++
	t.tokenTypesMetadata[t.lastTokenType] = metadata

	return t.lastTokenType
}

type TokenizerState struct {
	Data     string
	Position uint
	tokens   []Token
}

func (t *TokenizerState) Get(i uint) (byte, bool) {
	offset := t.Position + i
	
	if offset >= uint(len(t.Data)) {
		return 0, false
	}

	return t.Data[offset], true
}

func (t *TokenizerState) GetRange(start, length uint) (string, bool) {
	if start + length > uint(len(t.Data)) {
		return "", false
	}

	return t.Data[start:start + length], true
}

func (t *TokenizerState) Emit(token Token) {
	t.tokens = append(t.tokens, token)
}

func (t *Tokenizer) Tokenize(source string, stopper Stopper) []Token {
	s := &TokenizerState{source, 0, []Token{}}

	for !stopper.End(s) {
		largestLength := uint(0)
		var matcherWithLargestLength Matcher = nil

		for _, matcher := range t.matchers {
			length := matcher.Match(s)

			if length > largestLength {
				largestLength = length
				matcherWithLargestLength = matcher
			}
		}

		if matcherWithLargestLength != nil {
			matcherWithLargestLength.Consume(s, largestLength)	
			s.Position += largestLength
		} else {
			s.Emit(Token{
				Type: UNKNOWN, 
				Value: string(s.Data[s.Position]),
				Range: primitives.Range{Start: s.Position, Length: 1},
			})

			s.Position++
		}
	}

	return s.tokens
}

func (t *Tokenizer) GetTokenTypeMetadata(tokenType TokenType) (TokenTypeMetadata, bool) {
	v, ok := t.tokenTypesMetadata[tokenType]
	
	return v, ok
}
