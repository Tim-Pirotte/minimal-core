package tokenizerv2

import (
	usermessaging "minimal/minimal-core/built-in/user-messaging"
)

type TokenType uint

const (
	UNKNOWN TokenType = iota
	EOF
)

type Token struct {
	Type  TokenType
	Value string
	Span  usermessaging.Span
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
	// We start from the highest value for TokenType and decrement due to the constant TokenTypes declared in types.go
	return Tokenizer{[]Matcher{}, ^TokenType(0)}
}

func (t *Tokenizer) AddMatcher(matcher Matcher) {
	t.matchers = append(t.matchers, matcher)
}

func (t *Tokenizer) NewTokenType() TokenType {
	t.lastTokenType--

	return t.lastTokenType
}

type TokenizerState struct {
	Data []byte
	Position uint
	tokens []Token
}

func (t *TokenizerState) Emit(token Token) {
	t.tokens = append(t.tokens, token)
}

func (t *Tokenizer) Tokenize(source []byte) []Token {
	s := TokenizerState{source, 0, []Token{}}

	for s.Position < uint(len(s.Data)) {
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
				Value: string(s.Data[s.Position]),
				Span: usermessaging.Span{Start: s.Position, Length: 1},
			})

			s.Position++
		}
	}

	s.Emit(Token{Type: EOF, Value: "", Span: usermessaging.Span{Start: s.Position, Length: 0}})

	return s.tokens
}
