package strings

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messaging"
)

type StringMatcher struct {
    messenger *messaging.Messenger
    lexer     *lexer.Lexer
    strType lexer.TokenType
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
    // This matcher could violate max munch
    if c, _ := l.Get(0); c == '"' {
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
            l.Emit(lexer.Token{Type: s.strType, Value: l.GetNextN(pos + 1)})
            l.Position += pos + 1
            pos = 1

            level := 1

            for l.Position < uint(len(l.Data)) {
                switch c, _ := l.Get(0); c {
                // This forces every use of '{' as start of a token in an expression
                // to be properly closed by '}'
                case '{':
                    level++
                // There cannot be tokens starting with '}'
                // since the part after it would be ambiguous
                // with the end of a string interpolation
                case '}':
                    level--
                }

                if level == 0 {
                    break
                }

                largestLength := uint(0)
                var matcherWithLargestLength lexer.Matcher = nil

                for _, matcher := range l.Matchers {
                    length := matcher.Match(l)

                    if length > largestLength {
                        largestLength = length
                        matcherWithLargestLength = matcher
                    }
                }

                if matcherWithLargestLength != nil {
                    matcherWithLargestLength.Consume(l, largestLength)
                    l.Position += largestLength
                } else {
                    l.Emit(lexer.Token{Type: lexer.UNKNOWN, Value: l.GetNextN(1)})

                    l.Position++
                }
            }

            if level != 0 {
                // TODO Error: interpolation not properly closed
                panic("interpolation not closed")
            }
        case '\\':
            pos += 2
        default:
            pos++
        }
    }

    if !ok {
        s.sendUnclosedStrErr(l)

        pos--
    }

    l.Emit(lexer.Token{Type: s.strType, Value: l.GetNextN(pos + 1)})
    l.Position += pos
}

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
