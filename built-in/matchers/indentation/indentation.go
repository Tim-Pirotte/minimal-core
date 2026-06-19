package indentation

// TODO Should this be merged with newline since it depends on the working of it

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
)

type IndentationMatcher struct {
    openBlockSymbol         byte
    whiteSpace              byte
    openBlock               lexer.TokenType
    closeBlock              lexer.TokenType
    endOfLine               lexer.TokenType
    spacesPerIndentation    uint
    currentIndentation      uint
    currentSpaceCount       uint
    correctlyIndentedBefore bool
}

// TODO option to force a certain whitespace count?
func NewIndentationMatcher(
    openBlockSymbol, whiteSpace byte,
    openBlock, closeBlock, endOfLine lexer.TokenType,
) *IndentationMatcher {
    return &IndentationMatcher{
        openBlockSymbol,
        whiteSpace,
        openBlock,
        closeBlock,
        endOfLine,
        0,
        0,
        0,
        false,
    }
}

// TODO handle not ok cases
func (i *IndentationMatcher) Match(t *lexer.LexerJob) uint {
    pos := uint(0)
    c, ok := t.Get(pos)

    if c == i.openBlockSymbol {
        pos++
        c, ok = t.Get(pos)

        for ; ok && c == ' '; c, ok = t.Get(pos) {
            pos++
        }

        if ok && !IsEOL(c) {
            return 0
        }
    }

    posBefore := pos

    for ; ok && (c == ' ' || IsEOL(c)); c, ok = t.Get(pos) {
		pos++

        if IsEOL(c) {
            posBefore = pos
        }
	}

    i.currentSpaceCount = pos - posBefore

    return pos
}

func (i *IndentationMatcher) Consume(t *lexer.LexerJob, length uint) {
    value, _ := t.GetRange(t.Position, length)
    isOpenBlock := false

    if c, _ := t.Get(0); c == i.openBlockSymbol {
        i.currentIndentation++
        i.correctlyIndentedBefore = false

        t.Emit(lexer.Token{
            Type:  i.openBlock,
            Value: value,
            Range: primitives.Range{Start: t.Position, Length: length},
        })

        isOpenBlock = true
    }

    level := uint(0)

    if i.currentSpaceCount != 0 {
        if i.spacesPerIndentation == 0 {
            i.spacesPerIndentation = i.currentSpaceCount
        }

        if i.correctlyIndentedBefore &&
           i.currentSpaceCount > i.currentIndentation * i.spacesPerIndentation {
            i.currentSpaceCount = i.currentIndentation * i.spacesPerIndentation
        }

        if i.currentSpaceCount % i.spacesPerIndentation != 0 {
            // TODO Error inconsistent indentation
            panic("Inconsistent indentation")
        }

        level = i.currentSpaceCount / i.spacesPerIndentation

        if level > i.currentIndentation {
            // TODO Error indented without a new block
            panic("Indent without a new block")
        }
    }

    i.correctlyIndentedBefore = true

    if !isOpenBlock && i.currentIndentation - level == 0 {
        t.Emit(lexer.Token{
            Type: i.endOfLine,
            Value: value,
            Range: primitives.Range{Start: t.Position, Length: length},
        })

        return
    }

    for range i.currentIndentation - level {
        t.Emit(lexer.Token{
            Type:  i.closeBlock,
            Value: value,
            Range: primitives.Range{Start: t.Position, Length: length},
        })
    }

    i.currentIndentation = level
}

func IsEOL(c byte) bool {
	// Do not inline
	// This is used to keep the EOL chars in sync with the comment matcher
	return c == '\n' || c == '\r'
}
