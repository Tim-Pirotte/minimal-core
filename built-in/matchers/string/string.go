package strings

import (
	"minimal/minimal-lang/built-in/lexer"
	"minimal/minimal-lang/built-in/messenger"
)

type StringMatcher struct {
    messenger *messenger.Messenger
    lexer     *lexer.LexerScheme
    strType lexer.TokenType
}

func NewStringMatcher(
    messenger *messenger.Messenger,
    lexer *lexer.LexerScheme,
    tt lexer.TokenType,
) *StringMatcher {
    return &StringMatcher{messenger, lexer, tt}
}

func (s *StringMatcher) New(_ *lexer.Lexer) lexer.Matcher {
    return s
}

func (s *StringMatcher) Match(l *lexer.Lexer) uint32 {
    // This matcher could violate max munch
    if c, _ := l.Get(0); c == '\'' {
        return 1
    } else {
        return 0
    }
}

func (s *StringMatcher) Consume(l *lexer.Lexer, length uint32) {
    pos := uint32(1)
    c, ok := l.Get(pos)

    for ; ok && c != '\''; c, ok = l.Get(pos) {
        switch c {
        case '{':
            stringValue := l.GetNextN(pos + 1)
            l.Emit(lexer.Token{Type: s.strType, Value: stringValue})
            l.Position += pos + 1

            level := 1

            for l.Position < uint32(len(l.Data)) {
                switch c, _ := l.Get(0); c {
                // This requires every use of '{' as start of a token in an expression
                // to be properly closed by '}'
                case '{':
                    level++
                case '}':
                    level--
                }

                if level == 0 {
                    break
                }

                largestLength := uint32(0)
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
                l.Emit(lexer.Token{Type: s.strType, Value: l.GetNextN(0)})

                s.sendUnclosedInterpolationErr(stringValue[len(stringValue) - 1:])

                return
            }

            pos = 1
        case '\\':
            pos += 2
        default:
            pos++
        }
    }

    if !ok {
        s.sendUnclosedStrErr(l)

        l.Emit(lexer.Token{Type: s.strType, Value: l.GetNextN(pos)})
        l.Position += pos - 1

        return
    }

    l.Emit(lexer.Token{Type: s.strType, Value: l.GetNextN(pos + 1)})
    l.Position += pos
}

func (s *StringMatcher) sendUnclosedStrErr(l *lexer.Lexer) {
    s.messenger.Send(messenger.Message{
        Message: "String is not terminated with '",
        Severity: messenger.Error,
        Context: []messenger.Span{{Content: l.GetNextN(1), Note: "The string starts here"}},
        Notes: []string{"The remaining content will be interpreted as the string"},
    })
}

func (s *StringMatcher) sendUnclosedInterpolationErr(interpolationStart string) {
    s.messenger.Send(messenger.Message{
        Message: "String interpolation is not terminated with }",
        Severity: messenger.Error,
        Context: []messenger.Span{{Content: interpolationStart, Note: "Interpolation starts here"}},
        Notes: []string{"The string will be closed at the end of the source code"},
    })
}
