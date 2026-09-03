package indentation

import (
	"minimal/minimal-lang/built-in/lexer"
	"minimal/minimal-lang/built-in/messenger"
	"minimal/minimal-lang/built-in/substring"
	"strconv"
	"strings"
)

type IndentationMatcher struct {
    messenger               *messenger.Messenger
    openBlockSymbol         byte
    indentChar              byte
    openBlock               lexer.TokenType
    closeBlock              lexer.TokenType
    endOfLine               lexer.TokenType
    // Is a string so we can show where the number was derived from in error messages
    spacesPerLevel          string
    level                   uint32
    spaceCount              uint32
}

func NewIndentationMatcher(
    messenger *messenger.Messenger,
    openBlockSymbol, indentChar byte,
    openBlock, closeBlock, endOfLine lexer.TokenType,
    spacesPerLevel uint32,
) *IndentationMatcher {
    return &IndentationMatcher{
        messenger,
        openBlockSymbol,
        indentChar,
        openBlock,
        closeBlock,
        endOfLine,
        strings.Repeat(" ", int(spacesPerLevel)),
        0,
        0,
    }
}

func (i *IndentationMatcher) New(l *lexer.Lexer) lexer.Matcher {
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

    startIndent := uint32(0)

    for c, ok := l.Get(startIndent); ok && c == i.indentChar; c, ok = l.Get(startIndent) {
        startIndent++
    }

    if startIndent > 0 {
        i.sendPrefixIndentErr(l, startIndent)

        l.Position += startIndent
    }

    return m
}

func (i *IndentationMatcher) Match(l *lexer.Lexer) uint32 {
    pos := uint32(0)
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

func (i *IndentationMatcher) Consume(l *lexer.Lexer, length uint32) {
    isOpenBlock := false

    if c, _ := l.Get(0); c == i.openBlockSymbol {
        i.level++

        l.Emit(lexer.Token{
            Type:  i.openBlock,
            Value: l.GetNextN(length),
        })

        isOpenBlock = true
    }

    level := i.getIndentLevel(l, isOpenBlock, length)

    if !isOpenBlock && i.level - level == 0 {
        l.Emit(lexer.Token{
            Type: i.endOfLine,
            Value: l.GetNextN(length),
        })

        return
    }

    for range i.level - level {
        l.Emit(lexer.Token{
            Type:  i.closeBlock,
            Value: l.GetNextN(length),
        })
    }

    i.level = level
}

func IsEOL(c byte) bool {
	// This is used to keep the EOL chars in sync with the comment matcher
	return c == '\n' || c == '\r'
}

func (i *IndentationMatcher) getIndentLevel(l *lexer.Lexer, isOpenBlock bool, length uint32) uint32 {
    if i.spaceCount == 0 {
        return 0
    }

    if len(i.spacesPerLevel) == 0 {
        i.spacesPerLevel = l.Data[l.Position + length - i.spaceCount:l.Position + length]
    }

    if !isOpenBlock && i.spaceCount > uint32(len(i.spacesPerLevel)) * i.level {
        return i.level
    }

    if i.spaceCount % uint32(len(i.spacesPerLevel)) != 0 {
        i.sendInconsistentIndentErr(l, length)

        return i.level
    }

    level := i.spaceCount / uint32(len(i.spacesPerLevel))

    if level > i.level {
        i.sendMoreIndentErr(l, length)

        return i.level
    }

    return level
}

func (i *IndentationMatcher) sendPrefixIndentErr(l *lexer.Lexer, startIndent uint32) {
    i.messenger.Send(messenger.Message{
        Message: "Source code cannot start with indentation",
        Severity: messenger.Error,
        Context: []messenger.Span{{Content: l.GetNextN(startIndent)}},
        Notes: []string{"The indentation at the start will be skipped"},
    })
}

func (i *IndentationMatcher) sendInconsistentIndentErr(l *lexer.Lexer, length uint32) {
    context := l.Data[l.Position + length - i.spaceCount:l.Position + length]

    message := messenger.Message{
        Message: "Indentation is inconsistent",
        Severity: messenger.Error,
        Context: []messenger.Span{{Content: context}},
        AdditionalContext: []messenger.Span{},
        Notes: []string{
            "The indentation must be a multiple of " + strconv.Itoa(len(i.spacesPerLevel)),
            "The indentation of the incorrect line is " + strconv.Itoa(int(i.spaceCount)),
            "The indentation will be set to the current level",
        },
    }

    if substring.IsSubString(l.Data, i.spacesPerLevel) {
        message.AdditionalContext = append(message.AdditionalContext, messenger.Span{
            Content: i.spacesPerLevel,
            Note: "The indentation was derived here",
        })
    } else {
        message.Notes = append([]string{"The indentation was manually set"}, message.Notes...)
    }

    i.messenger.Send(message)
}

func (i *IndentationMatcher) sendMoreIndentErr(l *lexer.Lexer, length uint32) {
    context := l.Data[l.Position + length - i.spaceCount:l.Position + length]

    i.messenger.Send(messenger.Message{
        Message: "Got more indentation than expected",
        Severity: messenger.Error,
        Context: []messenger.Span{{Content: context}},
        AdditionalContext: []messenger.Span{},
        Notes: []string{
            "The indentation of the incorrect line is " + strconv.Itoa(int(i.spaceCount)),
            "The largest expected indentation is " + strconv.Itoa(int(i.level) * len(i.spacesPerLevel)),
            "The indentation will be set to the current level",
        },
    })
}
