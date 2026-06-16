package comments

import (
	"minimal/minimal-core/built-in/lexer"
	"testing"
)

func getLexer() (*lexer.Lexer, lexer.TokenType) {
	l := lexer.NewLexer()
	commentType := l.NewTokenType(
		lexer.TokenTypeMetadata{DisplayName: "a comment", DebugName: "Comment"},
	)

	identifierMatcher := NewCommentMatcher(commentType)
	l.AddMatcher(identifierMatcher)

	return l, commentType
}

func TestComments(t *testing.T) {

}
