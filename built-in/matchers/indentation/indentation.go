package indentation

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messaging"
	"minimal/minimal-core/built-in/primitives"
)

type IndentationMatcher struct {
    messenger               *messaging.Messenger
    openBlockSymbol         byte
    indentChar              byte
    openBlock               lexer.TokenType
    closeBlock              lexer.TokenType
    endOfLine               lexer.TokenType
    spacesPerLevel          uint
    level                   uint
    spaceCount              uint
}

func NewIndentationMatcher(
    messenger *messaging.Messenger,
    openBlockSymbol, indentChar byte,
    openBlock, closeBlock, endOfLine lexer.TokenType,
    spacesPerLevel uint,
) *IndentationMatcher {
    return &IndentationMatcher{
        messenger,
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

func (i *IndentationMatcher) New(l *lexer.LexerJob) lexer.Matcher {
	m := &IndentationMatcher{
        i.messenger,
        i.openBlockSymbol,
        i.indentChar,
        i.openBlock,
        i.closeBlock,
        i.endOfLine,
        i.spacesPerLevel,
        0,
        0,
    }

    startIndent := uint(0)

    for c, ok := l.Get(startIndent); ok && c == i.indentChar; c, ok = l.Get(startIndent) {
        startIndent++
    }

    if startIndent > 0 {
        context, _ := l.GetNextN(startIndent)
        i.messenger.Send(getCannotStartWithIndentMessage(context))

        l.Position += startIndent
    }

    return m
}

func (i *IndentationMatcher) Match(l *lexer.LexerJob) uint {
    pos := uint(0)
    c, ok := l.Get(pos)

    if c == i.openBlockSymbol {
        pos++

        for c, ok = l.Get(pos); ok && c == ' '; c, ok = l.Get(pos) {
            pos++
        }

        if !ok {
            return pos
        }

        if !IsEOL(c) {
            return 0
        }
    }

    posBefore := pos
    encounteredEOL := false

    for ; ok && (c == i.indentChar || IsEOL(c)); c, ok = l.Get(pos) {
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

func (i *IndentationMatcher) Consume(l *lexer.LexerJob, length uint) {
    value, _ := l.GetNextN(length)
    isOpenBlock := false

    if c, _ := l.Get(0); c == i.openBlockSymbol {
        i.level++

        l.Emit(lexer.Token{
            Type:  i.openBlock,
            Value: value,
            Range: primitives.Range{Start: l.Position, Length: length},
        })

        isOpenBlock = true
    }

    level := i.getIndentLevel(l, isOpenBlock, length)

    if !isOpenBlock && i.level - level == 0 {
        l.Emit(lexer.Token{
            Type: i.endOfLine,
            Value: value,
            Range: primitives.Range{Start: l.Position, Length: length},
        })

        return
    }

    for range i.level - level {
        l.Emit(lexer.Token{
            Type:  i.closeBlock,
            Value: value,
            Range: primitives.Range{Start: l.Position, Length: length},
        })
    }

    i.level = level
}

func IsEOL(c byte) bool {
	// This is used to keep the EOL chars in sync with the comment matcher
	return c == '\n' || c == '\r'
}

func (i *IndentationMatcher) getIndentLevel(l *lexer.LexerJob, isOpenBlock bool, length uint) uint {
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
        context := l.Data[l.Position + length - i.spaceCount:l.Position + length]

        // TODO Show the place at which the indentation was decided
        i.messenger.Send(messaging.Message{
            Reference: "TODO",
            Message: "Indentation is inconsistent",
            Severity: messaging.Error,
            Context: []messaging.Span{{Content: context}},
            Notes: []string{"The indentation will be set to the current level"},
        })

        return i.level
    }

    level := i.spaceCount / i.spacesPerLevel

    if level > i.level {
        // TODO Error indented without a new block
        panic("Indent without a new block")
    }

    return level
}

func getCannotStartWithIndentMessage(context string) messaging.Message {
    return messaging.Message{
        Reference: "TODO",
        Message: "Source code cannot start with indentation",
        Severity: messaging.Error,
        Context: []messaging.Span{{Content: context}},
        Notes: []string{"The indentation at the start will be skipped"},
    }
}
