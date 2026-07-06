package rawstring

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messaging"
)

const prefix = `r#"`

type RawStringMatcher struct {
    messenger *messaging.Messenger
	tokenType lexer.TokenType
}

func NewRawStringMatcher(messenger *messaging.Messenger, tt lexer.TokenType) *RawStringMatcher {
    return &RawStringMatcher{messenger, tt}
}

func (r *RawStringMatcher) New(_ *lexer.LexerJob) lexer.Matcher {
    return r
}

func (r *RawStringMatcher) Match(l *lexer.LexerJob) uint {
    if len(l.Data) - int(l.Position) < len(prefix) || l.GetNextN(uint(len(prefix))) != prefix {
        return 0
    }

    pos := uint(len(prefix))
    previousIsQuote := false

    for c, ok := l.Get(pos); ok; c, ok = l.Get(pos) {
        if previousIsQuote && c == '#' {
            return pos + 1
        }

        previousIsQuote = c == '"'

        pos++
    }

    r.sendUnclosedErr(l)

    return pos
}

func (r *RawStringMatcher) Consume(l *lexer.LexerJob, length uint) {
    l.Emit(lexer.Token{Type: r.tokenType, Value: l.GetNextN(length)})
}

func (r *RawStringMatcher) sendUnclosedErr(l *lexer.LexerJob) {
    r.messenger.Send(messaging.Message{
        Reference: "TODO",
        Message: `Raw string is not terminated with the "# sequence`,
        Severity: messaging.Error,
        Context: []messaging.Span{{Content: l.GetNextN(1), Note: "The raw string starts here"}},
        Notes: []string{"The remaining content will be interpreted as the raw string"},
    })
}
