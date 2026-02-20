package tokenizer

type Tokenizer struct {
	tokens   []Token
	position int
}

func NewTokenizer(config TokenizerConfig, source []byte) Tokenizer {
	return Tokenizer{
		config.tokenize(source),
		0,
	}
}

// Returns the current token
func (t *Tokenizer) CurrentToken() Token {
	return t.tokens[t.position]
}

// Goes to the next token or logs an error if the EOF token is consumed
func (t *Tokenizer) Consume() {
	if t.position+1 == len(t.tokens) {
		// TODO log error
	}

	t.position++
}

// Look ahead n tokens with 0 being the current token.
// Returns EOF if n goes out of bounds
func (t *Tokenizer) Peek(n int) Token {
	if t.position+n >= len(t.tokens) {
		return t.tokens[len(t.tokens)-1]
	}

	return t.tokens[t.position+n]
}
