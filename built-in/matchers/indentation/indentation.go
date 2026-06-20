package indentation

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
)

type IndentationMatcher struct {
    openBlockSymbol         byte
    indentChar              byte
    openBlock               lexer.TokenType
    closeBlock              lexer.TokenType
    endOfLine               lexer.TokenType
    spacesPerLevel          uint
    level                   uint
    spaceCount              uint
}

// TODO option to force a certain whitespace count?
func NewIndentationMatcher(
    openBlockSymbol, indentChar byte,
    openBlock, closeBlock, endOfLine lexer.TokenType,
    spacesPerLevel uint,
) *IndentationMatcher {
    return &IndentationMatcher{
        openBlockSymbol,
        indentChar,
        openBlock,
        closeBlock,
        endOfLine,
        spacesPerLevel,
        0,
        0,
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
    encounteredEOL := false

    for ; ok && (c == i.indentChar || IsEOL(c)); c, ok = t.Get(pos) {
		pos++

        if IsEOL(c) {
            posBefore = pos
            encounteredEOL = true
        }
	}

    if !encounteredEOL {
        return 0
    }

    i.spaceCount = pos - posBefore

    return pos
}

func (i *IndentationMatcher) Consume(t *lexer.LexerJob, length uint) {
    value, _ := t.GetRange(t.Position, length)
    isOpenBlock := false

    if c, _ := t.Get(0); c == i.openBlockSymbol {
        i.level++

        t.Emit(lexer.Token{
            Type:  i.openBlock,
            Value: value,
            Range: primitives.Range{Start: t.Position, Length: length},
        })

        isOpenBlock = true
    }

    level := i.getIndentLevel(isOpenBlock)

    if !isOpenBlock && i.level - level == 0 {
        t.Emit(lexer.Token{
            Type: i.endOfLine,
            Value: value,
            Range: primitives.Range{Start: t.Position, Length: length},
        })

        return
    }

    for range i.level - level {
        t.Emit(lexer.Token{
            Type:  i.closeBlock,
            Value: value,
            Range: primitives.Range{Start: t.Position, Length: length},
        })
    }

    i.level = level
}

func IsEOL(c byte) bool {
	// This is used to keep the EOL chars in sync with the comment matcher
	return c == '\n' || c == '\r'
}

func (i *IndentationMatcher) getIndentLevel(isOpenBlock bool) uint {
    if i.spaceCount == 0 {
        return 0
    }

    if i.spacesPerLevel == 0 {
        i.spacesPerLevel = i.spaceCount
    }

    if !isOpenBlock && i.spaceCount > i.spacesPerLevel * i.level {
        return i.level
    }

    if i.spaceCount % i.spacesPerLevel != 0 {
        // TODO Error inconsistent indentation
        panic("Inconsistent indentation")
    }

    level := i.spaceCount / i.spacesPerLevel

    if level > i.level {
        // TODO Error indented without a new block
        panic("Indent without a new block")
    }

    return level
}
