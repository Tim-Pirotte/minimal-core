package strings

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messaging"
)

var interpolation EnclosingSet = EnclosingSet{
    openSequence: "{",
    closingSequence: "}",
}

type StringMatcher struct {
    messenger *messaging.Messenger
    tokenType lexer.TokenType
    enclosingSets []EnclosingSet
}

// TODO assert length > 1
type EnclosingSet struct {
    openSequence    string
    closingSequence string
}

func NewStringMatcher(
    messenger *messaging.Messenger,
    tt lexer.TokenType,
    enclosingSets []EnclosingSet,
) *StringMatcher {
    enclosingSets = append(enclosingSets, interpolation)

    return &StringMatcher{messenger, tt, enclosingSets}
}

func (s *StringMatcher) New(_ *lexer.LexerJob) lexer.Matcher {
    return s
}

func (s *StringMatcher) Match(l *lexer.LexerJob) uint {
    nesting := []EnclosingSet{}

    firstChar, _ := l.Get(0)

    if firstChar != '"' {
        return 0
    }

    pos := uint(1)
    c, ok := l.Get(pos)

    for ; ok && (c != '"' || len(nesting) != 0); c, ok = l.Get(pos) {
        if len(nesting) > 0 {
            firstToClose := nesting[len(nesting) - 1].closingSequence
            maybeClose := l.Data[l.Position + pos:l.Position + pos + uint(len(firstToClose))]

            if maybeClose == firstToClose {
                nesting = nesting[:len(nesting) - 1]
                pos += uint(len(firstToClose))
            } else {
                found := false

                for _, s := range s.enclosingSets {
                    maybeOpen := l.Data[l.Position + pos:l.Position + pos + uint(len(s.openSequence))]

                    if maybeOpen == s.openSequence {
                        nesting = append(nesting, s)
                        pos += uint(len(s.openSequence))
                        found = false
                        break
                    }
                }

                if !found {
                    pos++
                }
            }
        } else if c == '\\' {
            pos += 2
        } else if c == '{' {
            nesting = append(nesting, interpolation)
            pos++
        } else {
            pos++
        }
    }

    if !ok {
        s.sendUnclosedStrErr(l)

        return pos
    }

    // To include the ending quote
    return pos + 1
}

func (s *StringMatcher) Consume(l *lexer.LexerJob, length uint) {
    l.Emit(lexer.Token{Type: s.tokenType, Value: l.GetNextN(length)})
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
