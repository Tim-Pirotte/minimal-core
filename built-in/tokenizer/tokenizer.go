package tokenizer

import (
	"bufio"
	"fmt"
	"io"
	"minimal/minimal-core/domain"
)

const MaxSymbolLength = 16

type Tokenizer struct {
	CurrentToken domain.Token
	reader       bufio.Reader
	position     uint

	keywords *trieNode
	symbols  *trieNode
	checks   []checkRule

	lastTokenType domain.TokenType
}

type trieNode struct {
	leaf bool
	token domain.TokenType
	children [256]*trieNode
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

func (t *Tokenizer) newTokenType() domain.TokenType {
	t.lastTokenType++
	return t.lastTokenType
}

func (t *Tokenizer) AddKeyword(keyword string) domain.TokenType {
	tokenType := t.newTokenType()
	err := updateTrie(t.keywords, keyword, tokenType)

	if err != nil {
		// TODO log error
		fmt.Println(err.Error())
	}

	return tokenType
}

func (t *Tokenizer) AddSymbol(symbol string) domain.TokenType {
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

func (t *Tokenizer) Consume() {
	t.skipWhiteSpace()

	_, err := t.reader.Peek(1)
    
	if err == io.EOF {
        t.setToEof(0) // TODO
        return
    }

	ok := t.checkSymbols()

	if ok {
		return
	}

	t.reader.ReadRune()
	t.setToUnknown(0, 0) // TODO
}

func (t *Tokenizer) skipWhiteSpace() {
	currentByte, err := t.reader.Peek(1)

	for err == nil && (currentByte[0] == ' ' || currentByte[0] == '\t') {
		_, discardErr := t.reader.Discard(1)

		if discardErr != nil {
			panic(fmt.Sprint("Could not properly discard:", discardErr))
		}

		currentByte, err = t.reader.Peek(1)
	}
}

func (t *Tokenizer) checkSymbols() bool {
    node := t.symbols
    var lastMatch domain.TokenType
    foundMatch := false
    matchLen := 0
    
    peeked, _ := t.reader.Peek(MaxSymbolLength) 

    for i := range peeked {
        nextNode := node.children[peeked[i]]

        if nextNode == nil {
            break
        }
        
        node = nextNode
		
        if node.leaf {
            lastMatch = node.token
            foundMatch = true
            matchLen = i + 1
        }
    }

    if foundMatch {
        t.reader.Discard(matchLen)
        t.CurrentToken.Type = lastMatch
		t.CurrentToken.Span.Start = 0 // TODO
		t.CurrentToken.Span.Length = 0 // TODO

        return true
    }

    return false
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
