package comments

import (
	"minimal/minimal-lang/built-in/lexer"
	"minimal/minimal-lang/built-in/matchers/indentation"
)

const prefix = '|'

type CommentMatcher struct {
	tokenType lexer.TokenType
}

func NewCommentMatcher(tt lexer.TokenType) *CommentMatcher {
	return &CommentMatcher{tt}
}

func (c *CommentMatcher) New(_ *lexer.Lexer) lexer.Matcher {
	return c
}

func (c *CommentMatcher) Match(l *lexer.Lexer) uint32 {
	start, _ := l.Get(0)

	if start != prefix {
		return 0
	}

	pos := uint32(1)

	for ch, ok := l.Get(pos); ok && !indentation.IsEOL(ch); ch, ok = l.Get(pos) {
		pos++
	}

	return pos
}

func (c *CommentMatcher) Consume(l *lexer.Lexer, length uint32) {
	l.Emit(lexer.Token{Type:  c.tokenType, Value: l.GetNextN(length)})
}
