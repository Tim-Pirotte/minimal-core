package strings

import "minimal/minimal-core/built-in/lexer"

type StringMatcher struct {
	tokenType lexer.TokenType
}

func NewStringMatcher(tt lexer.TokenType) *StringMatcher {
	return &StringMatcher{tt}
}

func (s *StringMatcher) New(_ *lexer.LexerJob) lexer.Matcher {
	return s
}

func (s *StringMatcher) Match(l *lexer.LexerJob) uint {
	firstChar, _ := l.Get(0)

	if firstChar != '"' {
		return 0
	}

	pos := uint(1)
	c, ok := l.Get(pos)

	for ; ok && c != '"'; c, ok = l.Get(pos) {
		pos++
	}

	if !ok {
		// TODO error no closing and only consume up to the next newline
	}

	return pos
}

func (s *StringMatcher) Consume(l *lexer.LexerJob, length uint) {
	l.Emit(lexer.Token{Type: s.tokenType, Value: l.GetNextN(length)})
}
