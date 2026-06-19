package indentation

// TODO Should this be merged with newline since it depends on the working of it

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
)

type IndentationMatcher struct {
    openBlockSymbol   byte
    whiteSpace        byte
    openBlock         lexer.TokenType
    closeBlock        lexer.TokenType
    endOfLine         lexer.TokenType
    indentation       uint
    indentationCount  uint
    currentSpaceCount uint
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
    }
}

// TODO handle not ok cases
func (i *IndentationMatcher) Match(t *lexer.LexerJob) uint {
    c, ok := t.Get(0)

    // TODO merge the duplicate code
    if c == i.openBlockSymbol {
        pos := uint(1)
        c, ok := t.Get(pos)

        for ; ok && c == ' '; c, ok = t.Get(pos) {
            pos++
        }

        posBefore := pos

        for ; ok && IsEOL(c); c, ok = t.Get(pos) {
		    pos++
	    }

        if pos == posBefore {
            return 0
        }

        posBefore = pos

        for ; ok && c == ' '; c, ok = t.Get(pos) {
            pos++
        }

        i.currentSpaceCount = pos - posBefore

        return pos
    }

    pos := uint(0)

    for ; ok && IsEOL(c); c, ok = t.Get(pos) {
		pos++
	}

    posBefore := pos

    for ; ok && c == ' '; c, ok = t.Get(pos) {
        pos++
    }

    i.currentSpaceCount = pos - posBefore

    return pos
}

func (i *IndentationMatcher) Consume(t *lexer.LexerJob, length uint) {
    value, _ := t.GetRange(t.Position, length)
    isOpenBlock := false

    if c, _ := t.Get(0); c == i.openBlockSymbol {
        i.indentationCount++

        t.Emit(lexer.Token{
                Type:  i.openBlock,
                Value: value,
                Range: primitives.Range{Start: t.Position, Length: length},
            },
        )

        isOpenBlock = true
    }

    level := uint(0)

    if i.currentSpaceCount != 0 {
        if i.indentation == 0 {
            i.indentation = i.currentSpaceCount
        }

        if i.currentSpaceCount % i.indentation != 0 {
            // TODO Error inconsistent indentation
            panic("Inconsistent indentation")
        }

        level = i.currentSpaceCount / i.indentation

        if level > i.indentationCount {
            // TODO Error indented without a new block
            panic("Indent without a new block")
        }
    }

    if !isOpenBlock && i.indentationCount - level == 0 {
        t.Emit(lexer.Token{
            Type: i.endOfLine,
            Value: value,
            Range: primitives.Range{Start: t.Position, Length: length},
        })

        return
    }

    for range i.indentationCount - level {
        t.Emit(lexer.Token{
            Type:  i.closeBlock,
            Value: value,
            Range: primitives.Range{Start: t.Position, Length: length},
        })
    }

    i.indentationCount = level
}

func IsEOL(c byte) bool {
	// Do not inline
	// This is used to keep the EOL chars in sync with the comment matcher
	return c == '\n' || c == '\r'
}
