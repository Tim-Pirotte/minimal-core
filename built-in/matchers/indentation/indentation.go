package indentation

// TODO Should this be merged with newline since it depends on the working of it

import (
	"minimal/minimal-core/built-in/lexer"
	eol "minimal/minimal-core/built-in/matchers/end-of-line"
	"minimal/minimal-core/built-in/primitives"
)

type IndentationMatcher struct {
    openBlockSymbol  byte
    whiteSpace       byte
    openBlock        lexer.TokenType
    closeBlock       lexer.TokenType
    handleOpenBlock  bool
    isEndOfLine      bool
    indentation      uint
    indentationCount uint
}

func NewIndentationMatcher(openBlockSymbol, whiteSpace byte, openBlock, closeBlock lexer.TokenType) *IndentationMatcher {
    return &IndentationMatcher{openBlockSymbol, whiteSpace, openBlock, closeBlock, false, false, 0, 0}
}

// TODO handle not ok cases
func (i *IndentationMatcher) Match(t *lexer.LexerJob) uint {
    c, ok := t.Get(0)

    if c == i.openBlockSymbol {
        pos := uint(1)
        c, ok := t.Get(pos)

        for ; ok && c == ' '; c, ok = t.Get(pos) {
            pos++
        }

        if eol.IsEOL(c) {
            i.handleOpenBlock = true

            return pos
        }

        return 0
    }

    if !i.isEndOfLine {
        i.isEndOfLine = eol.IsEOL(c)

        return 0
    }

    pos := uint(0)

    for ; ok && c == ' '; c, ok = t.Get(pos) {
        pos++
    }

    i.handleOpenBlock = false
    i.isEndOfLine = false

    return pos
}

func (i *IndentationMatcher) Consume(t *lexer.LexerJob, length uint) {
    if i.handleOpenBlock {
        i.indentationCount++

        openBlock, _ := t.GetRange(t.Position, 1)

        t.Emit(lexer.Token{
            Type:  i.openBlock,
            Value: openBlock,
            Range: primitives.Range{Start: t.Position, Length: 1},
        },
        )

        return
    }

    c, _ := t.Get(length)

    if eol.IsEOL(c) {
        return
    }

    if i.indentation == 0 {
        i.indentation = length
    }

    if length%i.indentation != 0 {
        // TODO Error inconsistent indentation
        panic("Inconsistent indentation")
    }

    level := length / i.indentation

    if level > i.indentationCount {
        // TODO Error indented without a new block
        panic("Indent without a new block")
    }

    closeBlock, _ := t.GetRange(t.Position, length)

    for range i.indentationCount - level {
        t.Emit(lexer.Token{
            Type:  i.closeBlock,
            Value: closeBlock,
            Range: primitives.Range{Start: t.Position, Length: length},
        })
    }
}
