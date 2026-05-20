package tokenizerv2

import (
	"minimal/minimal-core/built-in/primitives"
)

type TokenType uint

const (
	UNKNOWN TokenType = iota
	EOF
)

type Token struct {
	Type  TokenType
	Value string
	Range primitives.Range
}

type Tokenizer struct {
	matchers     []Matcher
	lastTokenType TokenType
}

type Matcher interface {
	Match(s *TokenizerState) (length uint)
	Consume(s *TokenizerState, length uint)
}

func NewTokenizer() Tokenizer {
	return Tokenizer{
		[]Matcher{}, 
		// We start from the highest value for TokenType and decrement due to the constant TokenTypes declared in types.go
		^TokenType(0),
	}
}

func (t *Tokenizer) AddMatcher(matcher Matcher) {
	t.matchers = append(t.matchers, matcher)
}

func (t *Tokenizer) NewTokenType() TokenType {
	t.lastTokenType--

	return t.lastTokenType
}

type TokenizerState struct {
	data     string
	Position uint
	tokens   []Token
}

func (t *TokenizerState) Get(i uint) (byte, bool) {
	offset := t.Position + i
	
	if offset >= uint(len(t.data)) {
		return 0, false
	}

	return t.data[offset], true
}

func (t *TokenizerState) GetRange(start, length uint) (string, bool) {
	if start + length > uint(len(t.data)) {
		return "", false
	}

	return t.data[start:start + length], true
}

func (t *TokenizerState) Emit(token Token) {
	t.tokens = append(t.tokens, token)
}

func (t *Tokenizer) Tokenize(source string) []Token {
	s := TokenizerState{source, 0, []Token{}}

	for s.Position < uint(len(s.data)) {
		largestLength := uint(0)
		var matcherWithLargestLength Matcher = nil

		for _, matcher := range t.matchers {
			length := matcher.Match(&s)

			if length > largestLength {
				largestLength = length
				matcherWithLargestLength = matcher
			}
		}

		if matcherWithLargestLength != nil {
			matcherWithLargestLength.Consume(&s, largestLength)	
			s.Position += largestLength
		} else {
			s.Emit(Token{
				Type: UNKNOWN, 
				Value: string(s.data[s.Position]),
				Range: primitives.Range{Start: s.Position, Length: 1},
			})

			s.Position++
		}
	}

	s.Emit(Token{Type: EOF, Value: "", Range: primitives.Range{Start: s.Position, Length: 0}})

	return s.tokens
}
