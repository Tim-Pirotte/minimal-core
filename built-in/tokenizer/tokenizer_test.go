package tokenizer

import (
	"bufio"
	"minimal/minimal-core/domain"
	"strings"
	"testing"
)

func TestLexEmpty(t *testing.T) {
	expected := []domain.TokenType{domain.EOF}
	r := bufio.NewReader(strings.NewReader(""))

	tokenizer := NewTokenizer(*r)

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

func TestLexSymbols(t *testing.T) {
	expected := []domain.TokenType{
		domain.UNKNOWN, 
		2, 
		domain.UNKNOWN, 
		3, 
		domain.UNKNOWN, 
		domain.EOF,
	}

	r := bufio.NewReader(strings.NewReader("5 + 4 - 3"))
	
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
	expected := []domain.TokenType{
		domain.UNKNOWN, 
		2, 
		domain.UNKNOWN, 
		3, 
		domain.UNKNOWN, 
		domain.EOF,
	}

	r := bufio.NewReader(strings.NewReader("+ +-+-+- 4 -+-+-+ -"))
	
	tokenizer := NewTokenizer(*r)

	tokenizer.AddSymbol("+-+-+-")
	tokenizer.AddSymbol("-+-+-+")

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

func TestLexUnicodeSymbols(t *testing.T) {
	expected := []domain.TokenType{
		domain.UNKNOWN, 
		2, 
		domain.UNKNOWN, 
		3, 
		domain.UNKNOWN, 
		domain.EOF,
	}

	r := bufio.NewReader(strings.NewReader("5 ❄️ 4 🔥 3"))
	
	tokenizer := NewTokenizer(*r)

	tokenizer.AddSymbol("❄️")
	tokenizer.AddSymbol("🔥")

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

func TestLexZeroWidthJoinerSymbols(t *testing.T) {
	expected := []domain.TokenType{
		domain.UNKNOWN, 
		2, 
		domain.UNKNOWN, 
		3, 
		domain.UNKNOWN, 
		domain.EOF,
	}

	r := bufio.NewReader(strings.NewReader("5 🐻‍❄️ 4 🐈‍⬛ 3"))
	
	tokenizer := NewTokenizer(*r)

	tokenizer.AddSymbol("🐻‍❄️")
	tokenizer.AddSymbol("🐈‍⬛")

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
