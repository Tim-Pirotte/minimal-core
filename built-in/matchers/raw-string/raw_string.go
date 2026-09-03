package rawstring

import (
	"fmt"
	"minimal/minimal-lang/built-in/lexer"
	"minimal/minimal-lang/built-in/messenger"
	"strings"
)

type RawStringMatcher struct {
    messenger *messenger.Messenger
	tokenType lexer.TokenType
}

func NewRawStringMatcher(messenger *messenger.Messenger, tt lexer.TokenType) *RawStringMatcher {
    return &RawStringMatcher{messenger, tt}
}

func (r *RawStringMatcher) New(_ *lexer.Lexer) lexer.Matcher {
    return r
}

func (r *RawStringMatcher) Match(l *lexer.Lexer) uint32 {
    pos := uint32(0)
    c, ok := l.Get(pos)

    for ; ok && c == '-'; c, ok = l.Get(pos) {
        pos++
    }

    dashes := pos

    if !ok || c != '\'' || dashes == 0 {
        return 0
    }

    pos++

    consecutiveDashes := uint32(0)
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

func (r *RawStringMatcher) Consume(l *lexer.Lexer, length uint32) {
    l.Emit(lexer.Token{Type: r.tokenType, Value: l.GetNextN(length)})
}

func (r *RawStringMatcher) sendUnclosedErr(l *lexer.Lexer, dashes uint32) {
    r.messenger.Send(messenger.Message{
        Message: `Raw string is not terminated with '` + strings.Repeat("-", int(dashes)),

        Severity: messenger.Error,
        Context: []messenger.Span{{Content: l.GetNextN(dashes + 1), Note: "The raw string starts here"}},
        Notes: []string{
            "The amount of dashes in the string prefix must match with the suffix",
            "The remaining content will be interpreted as the raw string",
        },
        Suggestions: findPossibleStringEnd(l, dashes),
    })
}

func findPossibleStringEnd(l *lexer.Lexer, expectedDashes uint32) []messenger.Suggestion {
    pos := uint32(expectedDashes + 1)
    longestEndSequence := uint32(0)
    startOfLongestEndSequence := uint32(0)
    consecutiveDashes := uint32(0)
    startOfLastQuote := uint32(0)
    quote := false

    for c, ok := l.Get(pos); ok; c, ok = l.Get(pos) {
        switch c {
        case '-':
            consecutiveDashes++

            if quote && consecutiveDashes + 1 > longestEndSequence {
                longestEndSequence = consecutiveDashes + 1
                startOfLongestEndSequence = startOfLastQuote
            }
        case '\'':
            consecutiveDashes = 0
            quote = true
            startOfLastQuote = pos

            if longestEndSequence == 0 {
                longestEndSequence = 1
                startOfLongestEndSequence = startOfLastQuote
            }
        default:
            consecutiveDashes = 0
            quote = false
        }

        pos++
    }

    if longestEndSequence > 0 {
        missing := expectedDashes - (longestEndSequence - 1)
        plural := ""

        if missing > 1 {
            plural = "es"
        }

        return []messenger.Suggestion{
            {
                Suggestion: "This looks like an ending sequence",
                Replacements: []messenger.Replacement{
                    {
                        From: messenger.Span{
                            Content: l.GetNextN(pos)[startOfLongestEndSequence:startOfLongestEndSequence + longestEndSequence],
                            Note: fmt.Sprintf("Missing %d dash%s", missing, plural),
                        },
                        To: messenger.Span{
                            Content: "'" + strings.Repeat("-", int(expectedDashes)),
                        },
                    },
                },
            },
        }
    }

    return nil
}
