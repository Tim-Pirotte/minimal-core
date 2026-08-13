package numbers

import (
	"minimal/minimal-core/built-in/lexer"
	"testing"
)

func getLexer() (*lexer.LexerScheme, lexer.TokenType) {
	l := lexer.New(1)
	numberType := l.NewTokenType(
		lexer.TokenTypeMetadata{DisplayName: "a number", DebugName: "Number"},
	)

	identifierMatcher := NewNumberMatcher(numberType)
	l.AddMatcher(identifierMatcher)

	return l, numberType
}

func TestNumber(t *testing.T) {
	source := "0123456789"

	l, numberType := getLexer()

	expected := []lexer.Token{
		{Type: numberType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMixed(t *testing.T) {
	source := "/0:12345 6789"

	l, numberType := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: source[:1]},
		{Type: numberType, Value: source[1:2]},
		{Type: lexer.UNKNOWN, Value: source[2:3]},
		{Type: numberType, Value: source[3:8]},
		{Type: lexer.UNKNOWN, Value: source[8:9]},
		{Type: numberType, Value: source[9:13]},
	}

	// / and : are next to 0 and 9 respectively in the ASCII table
	lexer.CheckTokens(t, l, expected, source)
}
