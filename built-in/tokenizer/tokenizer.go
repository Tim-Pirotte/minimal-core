package tokenizer

import (
	"bufio"
	"fmt"
	"io"
	"minimal/minimal-core/domain"
)

const MaxSymbolLengthBytes = 16

type Tokenizer struct {
	CurrentToken domain.Token
	reader       bufio.Reader
	position     uint

	keywords *trieNode
	symbols  *trieNode
	checks   []checkRule

	lastTokenType domain.TokenType
}

type checkRule struct {
	check func(rune) bool
	tokenType domain.TokenType
}

func NewTokenizer(reader bufio.Reader) Tokenizer {
	return Tokenizer{
		domain.Token{}, 
		reader,
		0,
		&trieNode{children: [256]*trieNode{}},
		&trieNode{children: [256]*trieNode{}},
		[]checkRule{},
		domain.EOF,
	}
}

func (t *Tokenizer) AddSymbol(symbol string) domain.TokenType {
	if len(symbol) > MaxSymbolLengthBytes {
		// TODO log error
	}

	tokenType := t.newTokenType()
	err := updateTrie(t.symbols, symbol, tokenType)

	if err != nil {
		// TODO log error
		fmt.Println(err.Error())
	}
	
	return tokenType
}

func (t *Tokenizer) AddCheck(check func(rune) bool) domain.TokenType {
	tokenType := t.newTokenType()
	t.checks = append(t.checks, checkRule{check, tokenType})

	return tokenType
}

func (t *Tokenizer) newTokenType() domain.TokenType {
	t.lastTokenType++
	return t.lastTokenType
}

func (t *Tokenizer) Consume() {
	t.skipWhiteSpace()
    
	if _, err := t.reader.Peek(1); err == io.EOF {
        t.setToEof(0) // TODO
        return
    }

	bestTokenType := domain.UNKNOWN
	bestLength := uint(0)

	if token, length, ok := t.checkSymbols(); ok {
		bestTokenType = token
		bestLength = length
	}

	if bestLength == 0 {
		t.reader.ReadRune()
		t.setToUnknown(t.position, 1)
	} else {
		t.reader.Discard(int(bestLength))
		t.CurrentToken.Type = bestTokenType
		t.CurrentToken.Span.Start = t.position
		t.CurrentToken.Span.Length = bestLength
		t.position += bestLength
	}
}

func (t *Tokenizer) skipWhiteSpace() {
	currentByte, err := t.reader.Peek(1)

	for err == nil && (currentByte[0] == ' ' || currentByte[0] == '\t') {
		_, discardErr := t.reader.Discard(1)
		t.position++

		if discardErr != nil {
			panic(fmt.Sprint("Could not properly discard:", discardErr))
		}

		currentByte, err = t.reader.Peek(1)
	}
}

func (t *Tokenizer) checkSymbols() (domain.TokenType, uint, bool) {
    node := t.symbols
    var lastMatch domain.TokenType
    foundMatch := false
    matchLength := uint(0)
    
    peeked, err := t.reader.Peek(MaxSymbolLengthBytes) 

	if err == bufio.ErrBufferFull || err == bufio.ErrNegativeCount {
		panic(fmt.Sprint("Something went wrong while peeking: ", err))
	}

    for i := range peeked {
        nextNode := node.children[peeked[i]]

        if nextNode == nil {
            break
        }
        
        node = nextNode

        if node.leaf {
            lastMatch = node.token
            foundMatch = true
            matchLength = uint(i) + 1
        }
    }

    return lastMatch, uint(matchLength), foundMatch
}

func (t *Tokenizer) setToEof(start uint) {
	t.CurrentToken.Type = domain.EOF
	t.CurrentToken.Value = ""
	t.CurrentToken.Span.Start = start
	t.CurrentToken.Span.Length = 0
}

func (t *Tokenizer) setToUnknown(start uint, length uint) {
	t.CurrentToken.Type = domain.UNKNOWN
	t.CurrentToken.Value = ""
	t.CurrentToken.Span.Start = start
	t.CurrentToken.Span.Length = length
}

const asciiMax = 127

func isAlphanumericOrUnicode(char rune) bool {
	return 'a' >= char && char <= 'z' || 'A' >= char && char <= 'Z' || '0' >= char && char <= '9' || char > asciiMax
}
