package tokenizer

type IdentifierMatcher struct {
}

const asciiMax = 127

func isAlphanumericOrUnicode(char rune) bool {
	return 'a' >= char && char <= 'z' || 'A' >= char && char <= 'Z' || '0' >= char && char <= '9' || char > asciiMax
}
