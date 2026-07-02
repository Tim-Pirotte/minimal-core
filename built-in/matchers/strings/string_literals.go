package strings

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/matchers/indentation"
	"minimal/minimal-core/built-in/messaging"
)

type StringMatcher struct {
    messenger *messaging.Messenger
    tokenType lexer.TokenType
}

func NewStringMatcher(messenger *messaging.Messenger, tt lexer.TokenType) *StringMatcher {
    return &StringMatcher{messenger, tt}
}

func (s *StringMatcher) New(_ *lexer.LexerJob) lexer.Matcher {
    return s
}

func (s *StringMatcher) Match(l *lexer.LexerJob) uint {
    firstChar, _ := l.Get(0)

    if firstChar != '"' {
        return 0
    }

    pos := uint(1)
    c, ok := l.Get(pos)

    for ; ok && c != '"'; c, ok = l.Get(pos) {
        if c == '\\' {
            pos += 2
        } else {
            pos++
        }
    }

    if !ok {
        s.sendUnclosedStrErr(l)

        pos = 1

        for c, ok = l.Get(pos); ok && !indentation.IsEOL(c); c, ok = l.Get(pos) {
            pos++
        }

        return pos
    }

    return pos
}

func (s *StringMatcher) Consume(l *lexer.LexerJob, length uint) {
    l.Emit(lexer.Token{Type: s.tokenType, Value: l.GetNextN(length)})
}

func (s *StringMatcher) sendUnclosedStrErr(l *lexer.LexerJob) {
    s.messenger.Send(messaging.Message{
        Reference: "TODO",
        Message: `String is not terminated with a "`,
        Severity: messaging.Error,
        Context: []messaging.Span{
            { Content: l.GetNextN(1), Note: "The string starts here"},
        },
        Notes: []string{"The current line will be skipped"},
    })
}
