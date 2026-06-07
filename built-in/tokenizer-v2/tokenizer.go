package tokenizerv2

import (
	"math"
	"minimal/minimal-core/built-in/primitives"
)

type TokenType uint

const (
	UNKNOWN TokenType = math.MaxUint - iota
	END
)

type Token struct {
	Type  TokenType
	Value string
	Range primitives.Range
}

type TokenTypeMetadata struct {
	DisplayName string
	DebugName   string
}

type Tokenizer struct {
	matchers      []Matcher
	lastTokenType TokenType
	tokenTypesMetadata map[TokenType]TokenTypeMetadata
}

type TokenizerJob struct {
	tokenizer   *Tokenizer
	Data        string
	Position    uint
	buffer      []Token
	read        uint
	write       uint
	spillage    []Token
	minSafePeek uint
	stopper     Stopper
	endReached  bool
}

type Matcher interface {
	Match(t *TokenizerJob) (length uint)
	Consume(t *TokenizerJob, length uint)
}

type Stopper interface {
	End(t *TokenizerJob) bool
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		[]Matcher{},
		TokenType(0),
		map[TokenType]TokenTypeMetadata{
			UNKNOWN: {"a character that is not a valid token", "UNKNOWN"},
			END: {"to the end", "END"},
		},
	}
}

func (t *Tokenizer) AddMatcher(matcher Matcher) {
	t.matchers = append(t.matchers, matcher)
}

func (t *Tokenizer) NewTokenType(metadata TokenTypeMetadata) TokenType {
	t.lastTokenType++
	t.tokenTypesMetadata[t.lastTokenType] = metadata

	return t.lastTokenType
}

func (t *Tokenizer) GetTokenTypeMetadata(tokenType TokenType) (TokenTypeMetadata, bool) {
	v, ok := t.tokenTypesMetadata[tokenType]

	return v, ok
}

func (t *Tokenizer) Tokenize(source string, stopper Stopper, minSafePeek uint) *TokenizerJob {
	capacity := minSafePeek

	job := &TokenizerJob{
		t,
		source,
		0,
		make([]Token, capacity),
		0,
		0,
		[]Token{},
		minSafePeek,
		stopper,
		false,
	}

	job.fillTokenBuffer()

	return job
}

func (t *TokenizerJob) Get(i uint) (byte, bool) {
	offset := t.Position + i

	if offset >= uint(len(t.Data)) {
		return 0, false
	}

	return t.Data[offset], true
}

func (t *TokenizerJob) GetRange(start, length uint) (string, bool) {
	if start + length > uint(len(t.Data)) {
		return "", false
	}

	return t.Data[start:start + length], true
}

func (t *TokenizerJob) Emit(token Token) {
	// TODO add spillage to fill buffer
	if t.write - t.read == uint(len(t.buffer)) {
		t.spillage = append(t.spillage, token)

		return
	}

	t.buffer[t.write] = token
	t.write++
}

func (t *TokenizerJob) Peek(n uint) Token {
	if n >= t.minSafePeek {
		panic("attempt to peek more tokens in advance than expected")
	}

	read := t.read + n

	if read >= t.write {
		return Token{END, "", primitives.Range{Start: uint(len(t.Data)), Length: 0}}
	}

	return t.buffer[read % uint(len(t.buffer))]
}

func (t *TokenizerJob) Advance() {
	t.read++

	if t.read + t.minSafePeek >= t.write {
		t.fillTokenBuffer()
	}
}

func (t *TokenizerJob) fillTokenBuffer() {
	for t.write - t.read != uint(len(t.buffer)) && !t.endReached {
		if t.stopper.End(t) {
			t.endReached = true

			return
		}

		largestLength := uint(0)
		var matcherWithLargestLength Matcher = nil

		for _, matcher := range t.tokenizer.matchers {
			length := matcher.Match(t)

			if length > largestLength {
				largestLength = length
				matcherWithLargestLength = matcher
			}
		}

		if matcherWithLargestLength != nil {
			matcherWithLargestLength.Consume(t, largestLength)
			t.Position += largestLength
		} else {
			t.Emit(Token{
				Type: UNKNOWN,
				Value: t.Data[t.Position:t.Position + 1],
				Range: primitives.Range{Start: t.Position, Length: 1},
			})

			t.Position++
		}
	}
}
