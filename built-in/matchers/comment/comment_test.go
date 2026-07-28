package comments

import (
	"minimal/minimal-core/built-in/lexer"
	"testing"
)

func getLexer() (*lexer.Lexer, lexer.TokenType) {
	l := lexer.New()
	commentType := l.NewTokenType(
		lexer.TokenTypeMetadata{DisplayName: "a comment", DebugName: "Comment"},
	)

	identifierMatcher := NewCommentMatcher(commentType)
	l.AddMatcher(identifierMatcher)

	return l, commentType
}

func TestComments(t *testing.T) {
	source := "0| A comment"

	l, commentType := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: source[:1]},
		{Type: commentType, Value: source[1:]},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestNewLine(t *testing.T) {
	source := "| A comment\n| Second comment\r"

	l, commentType := getLexer()

	expected := []lexer.Token{
		{Type: commentType, Value: source[0:11]},
		{Type: lexer.UNKNOWN, Value: source[11:12]},
		{Type: commentType, Value: source[12:28]},
		{Type: lexer.UNKNOWN, Value: source[28:29]},
	}

	lexer.CheckTokens(t, l, expected, source)
}
