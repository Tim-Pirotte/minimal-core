package tokenizer

import (
	"bufio"
	"minimal/minimal-core/domain"
	"strings"
	"testing"
)

func TestLexEmpty(t *testing.T) {

}

func TestLexSymbols(t *testing.T) {
	expected := []domain.TokenType{0, 2, 0, 3, 0, 1}

	input := "5 + 4 - 3"
	r := bufio.NewReader(strings.NewReader(input))
	
	tokenizer := NewTokenizer(*r)

	tokenizer.AddSymbol("+")
	tokenizer.AddSymbol("-")

	maxIter := 20
	currentIter := 0
	actual := make([]domain.TokenType, 0)

	for tokenizer.CurrentToken.Type != domain.EOF && currentIter <= maxIter {
		tokenizer.Consume()
		actual = append(actual, tokenizer.CurrentToken.Type)
		currentIter++
	}

	if currentIter == maxIter {
		t.Error("Infinite loop")
	} else {
		for i, e := range expected {
			if actual[i] != e {
				t.Error("Expected", e, "at index", i, "but got", actual[i])
			}
		}
	}
}

func TestLexMultiCharSymbols(t *testing.T) {

}

func TestLexUnicodeSymbols(t *testing.T) {

}
