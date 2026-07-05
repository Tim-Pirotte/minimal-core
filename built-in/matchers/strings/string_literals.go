package strings

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messaging"
)

type StringMatcher struct {
    messenger *messaging.Messenger
    lexer     *lexer.Lexer
    tokenType lexer.TokenType
}

func NewStringMatcher(
    messenger *messaging.Messenger,
    lexer *lexer.Lexer,
    tt lexer.TokenType,
) *StringMatcher {
    return &StringMatcher{messenger, lexer, tt}
}

func (s *StringMatcher) New(_ *lexer.LexerJob) lexer.Matcher {
    return s
}

func (s *StringMatcher) Match(l *lexer.LexerJob) uint {
    firstChar, _ := l.Get(0)

    if firstChar == '"' {
        return 1
    } else {
        return 0
    }
}

func (s *StringMatcher) Consume(l *lexer.LexerJob, length uint) {
    pos := uint(1)
    c, ok := l.Get(pos)

    for ; ok && c != '"'; c, ok = l.Get(pos) {
        switch c {
        case '{':
            // TODO Lex and offset position
            if c, ok := l.Get(0); !ok || c != '}' {
                // TODO Error: interpolation not properly closed
            }
        case '\\':
            pos += 2
        default:
            pos++
        }
    }

    if !ok {
        s.sendUnclosedStrErr(l)

        l.Emit(lexer.Token{Type: s.tokenType, Value: l.GetNextN(pos)})
    }
}

// if s.Position >= uint(len(s.Data)) {
//     return true
// }

// switch s.Data[s.Position] {
// case '{':
//     i.nesting++
// case '}':
//     i.nesting--
// }

// return i.nesting == 0

func (s *StringMatcher) sendUnclosedStrErr(l *lexer.LexerJob) {
    s.messenger.Send(messaging.Message{
        Reference: "TODO",
        Message: "String is not terminated with a quote",
        Severity: messaging.Error,
        Context: []messaging.Span{{Content: l.GetNextN(1), Note: "The string starts here"}},
        Notes: []string{"The remaining content will be interpreted as the string"},
    })
}

func (s *StringMatcher) sendCloseBraceErr(l *lexer.LexerJob, pos uint) {
    context := l.Data[l.Position + pos:l.Position+pos + 1]

    s.messenger.Send(messaging.Message{
        Reference: "TODO",
        Message: "Closing brace does not have a matching opening brace",
        Severity: messaging.Error,
        Context: []messaging.Span{{Content: context}},
        Notes: []string{"The closing brace will be ignored"},
    })
}
