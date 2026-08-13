package rawstring

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messaging"
	"strings"
)

type RawStringMatcher struct {
    messenger *messaging.Messenger
	tokenType lexer.TokenType
}

func NewRawStringMatcher(messenger *messaging.Messenger, tt lexer.TokenType) *RawStringMatcher {
    return &RawStringMatcher{messenger, tt}
}

func (r *RawStringMatcher) New(_ *lexer.Lexer) lexer.Matcher {
    return r
}

func (r *RawStringMatcher) Match(l *lexer.Lexer) uint {
    pos := uint(0)
    c, ok := l.Get(pos)

    for ; ok && c == '-'; c, ok = l.Get(pos) {
        pos++
    }

    dashes := pos

    if !ok || c != '\'' || dashes == 0 {
        return 0
    }

    pos++

    consecutiveDashes := uint(0)
    quote := false

    for ; ok; c, ok = l.Get(pos) {
        switch c {
        case '-':
            consecutiveDashes++

            if quote && consecutiveDashes == dashes {
                return pos + 1
            }
        case '\'':
            consecutiveDashes = 0
            quote = true
        default:
            consecutiveDashes = 0
            quote = false
        }

        pos++
    }

    r.sendUnclosedErr(l, dashes)

    return pos
}

func (r *RawStringMatcher) Consume(l *lexer.Lexer, length uint) {
    l.Emit(lexer.Token{Type: r.tokenType, Value: l.GetNextN(length)})
}

func (r *RawStringMatcher) sendUnclosedErr(l *lexer.Lexer, dashes uint) {
    r.messenger.Send(messaging.Message{
        Message: `Raw string is not terminated with '` + strings.Repeat("-", int(dashes)),

        Severity: messaging.Error,
        Context: []messaging.Span{{Content: l.GetNextN(dashes + 1), Note: "The raw string starts here"}},
        Notes: []string{
            "The amount of dashes in the string prefix must match with the suffix",
            "The remaining content will be interpreted as the raw string",
        },
    })
}

// TODO rescan the string and keep track of the longest ending sequence
