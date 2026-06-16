package comments

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
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
	l, commentType := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "0", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: commentType, Value: "| A comment", Range: primitives.Range{Start: 1, Length: 11}},
	}

	lexer.CheckTokens(t, l, expected, "0| A comment")
}

func TestNewLine(t *testing.T) {
	l, commentType := getLexer()

	expected := []lexer.Token{
		{Type: commentType, Value: "| A comment", Range: primitives.Range{Start: 0, Length: 11}},
		{Type: lexer.UNKNOWN, Value: "\n", Range: primitives.Range{Start: 11, Length: 1}},
		{Type: commentType, Value: "| Second comment", Range: primitives.Range{Start: 12, Length: 16}},
		{Type: lexer.UNKNOWN, Value: "\r", Range: primitives.Range{Start: 28, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, "| A comment\n| Second comment\r")
}
