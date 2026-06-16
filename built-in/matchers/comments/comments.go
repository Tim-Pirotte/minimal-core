package comments

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
)

const prefix = '|'

type CommentMatcher struct {
	tokenType lexer.TokenType
}

func NewCommentMatcher(tt lexer.TokenType) *CommentMatcher {
	return &CommentMatcher{tt}
}

func (n *CommentMatcher) Match(t *lexer.LexerJob) uint {
	start, _ := t.Get(0)

	if start != prefix {
		return 0
	}

	pos := uint(1)

	// Keep these EOL sequences with the EOL matcher
	for c, ok := t.Get(pos); ok && c != '\n' && c != '\r'; c, ok = t.Get(pos) {
		pos++
	}

	return pos
}

func (i *CommentMatcher) Consume(t *lexer.LexerJob, length uint) {
	comment, _ := t.GetRange(t.Position, length)

	t.Emit(lexer.Token{
		Type:  i.tokenType,
		Value: comment,
		Range: primitives.Range{Start: t.Position, Length: length}},
	)
}
