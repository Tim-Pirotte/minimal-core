package identifiers

import (
	"minimal/minimal-core/built-in/lexer"
	"strings"
	"testing"
)

func getLexer() (*lexer.LexerScheme, lexer.TokenType) {
	l := lexer.NewScheme()
	identifierType := l.NewTokenType(
		lexer.TokenTypeMetadata{NounPhrase: "an identifier", DebugName: "Identifier"},
	)

	identifierMatcher := NewIdentifierMatcher(identifierType)
	l.AddMatcher(identifierMatcher)

	return l, identifierType
}

func TestIdentifier(t *testing.T) {
	source := "identifier1"

	l, identifierType := getLexer()

	expected := []lexer.Token{
		{Type: identifierType, Value: source[:11]},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMultipleIdentifiers(t *testing.T) {
	source := "identifier1 identifier2"

	l, identifierType := getLexer()

	expected := []lexer.Token{
		{Type: identifierType, Value: source[:11]},
		{Type: lexer.UNKNOWN, Value: source[11:12]},
		{Type: identifierType, Value: source[12:23]},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestUnicode(t *testing.T) {
	source := "👾"

	l, identifierType := getLexer()

	expected := []lexer.Token{
		{Type: identifierType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestZeroWidthJoiner(t *testing.T) {
	source := "🐻‍❄️"

	l, identifierType := getLexer()

	expected := []lexer.Token{
		{Type: identifierType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestStartingWithNumber(t *testing.T) {
	source := "1identifier"

	l, identifierType := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: source[:1]},
		{Type: identifierType, Value: source[1:11]},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestBounds(t *testing.T) {
	source := "`az{@AZ[\u007f\u0080\U0010ffff"

	l, identifierType := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: source[:1]},
		{Type: identifierType, Value: source[1:3]},
		{Type: lexer.UNKNOWN, Value: source[3:4]},
		{Type: lexer.UNKNOWN, Value: source[4:5]},
		{Type: identifierType, Value: source[5:7]},
		{Type: lexer.UNKNOWN, Value: source[7:8]},
		{Type: lexer.UNKNOWN, Value: source[8:9]},
		{Type: identifierType, Value: source[9:15]},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func FuzzLexIdentifier(f *testing.F) {
	l, identifierType := getLexer()

	f.Add("identifier")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 0 && '0' <= input[0] && input[0] <= '9' {
			input = input[1:]
		}

		var cleanedInput strings.Builder

		for _, c := range []byte(input) {
			if isAlphaOrUnicode(c) {
				cleanedInput.WriteString(string([]byte{c}))
			}
		}

		input = cleanedInput.String()

		if input == "" {
			return
		}

		expected := []lexer.Token{{Type: identifierType, Value: input}}

		lexer.CheckTokens(t, l, expected, input)
	})
}
