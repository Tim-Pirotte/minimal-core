package comments

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/matchers/indentation"
	"minimal/minimal-core/built-in/primitives"
)

const prefix = '|'

type CommentMatcher struct {
	tokenType lexer.TokenType
}

func NewCommentMatcher(tt lexer.TokenType) *CommentMatcher {
	return &CommentMatcher{tt}
}

func (c *CommentMatcher) New(_ *lexer.LexerJob) lexer.Matcher {
	return c
}

func (c *CommentMatcher) Match(l *lexer.LexerJob) uint {
	start, _ := l.Get(0)

	if start != prefix {
		return 0
	}

	pos := uint(1)

	for ch, ok := l.Get(pos); ok && !indentation.IsEOL(ch); ch, ok = l.Get(pos) {
		pos++
	}

	return pos
}

func (c *CommentMatcher) Consume(l *lexer.LexerJob, length uint) {
	l.Emit(lexer.Token{
		Type:  c.tokenType,
		Value: l.GetNextN(length),
		Range: primitives.Range{Start: l.Position, Length: length}},
	)
}
