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

func (c *CommentMatcher) Match(t *lexer.LexerJob) uint {
	start, _ := t.Get(0)

	if start != prefix {
		return 0
	}

	pos := uint(1)

	for ch, ok := t.Get(pos); ok && !indentation.IsEOL(ch); ch, ok = t.Get(pos) {
		pos++
	}

	return pos
}

func (c *CommentMatcher) Consume(t *lexer.LexerJob, length uint) {
	comment, _ := t.GetRange(t.Position, length)

	t.Emit(lexer.Token{
		Type:  c.tokenType,
		Value: comment,
		Range: primitives.Range{Start: t.Position, Length: length}},
	)
}
