package prefix

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messaging"
	"strings"
)

type Rule struct {
    TokenTypes []lexer.TokenType
    Parser     RuleParser
}

type RuleParser interface {
    Parse(*lexer.Lexer, *ast.ASTSchema)
}

type trieNode struct {
    leaf     bool
    parser  RuleParser
    children map[lexer.TokenType]*trieNode
}

type PrefixParser struct {
    prefixes  *trieNode
    maxLength uint
}

func NewPrefixParser(m *messaging.Messenger, l *lexer.LexerScheme, prefixes []Rule) PrefixParser {
    root := &trieNode{false, nil, map[lexer.TokenType]*trieNode{}}
    maxLength := uint(0)

    for _, prefix := range prefixes {
        length := uint(len(prefix.TokenTypes))

        if length > maxLength {
            maxLength = length
        }

        node := root

        for _, tokenType := range prefix.TokenTypes {
            if _, ok := node.children[tokenType]; !ok {
                node.children[tokenType] = &trieNode{children: map[lexer.TokenType]*trieNode{}}
            }

            node = node.children[tokenType]
        }

        if node.leaf {
            logDuplicatePrefix(m, l, prefix.TokenTypes)
        }

        node.leaf = true
        node.parser = prefix.Parser
    }

    return PrefixParser{root, uint(maxLength)}
}

func (p *PrefixParser) Parse(l *lexer.Lexer, syntax *ast.ASTSchema) {
    var largestMatchParser RuleParser
    node := p.prefixes
    ok := true

    if node.leaf {
        largestMatchParser = node.parser
    }

    for pos := uint(0); ok && pos < p.maxLength; pos++ {
        tokenType := l.Peek(pos).Type
        node, ok = node.children[tokenType]

        if ok && node.leaf {
            largestMatchParser = node.parser
        }
    }

    if largestMatchParser != nil {
        largestMatchParser.Parse(l, syntax)
    }
}

func logDuplicatePrefix(m *messaging.Messenger, l *lexer.LexerScheme, prefix []lexer.TokenType) {
    sb := strings.Builder{}
    sb.WriteString("Prefix: [")

    for i, t := range prefix {
        metadata := l.GetTokenTypeMetadata(t)

        sb.WriteString(metadata.DebugName)

        if i + 1 != len(prefix) {
            sb.WriteString(", ")
        }
    }

    sb.WriteByte(']')

    m.Send(
        messaging.Message{
            Message: "Duplicate prefix in prefix parser",
            Severity: messaging.Error,
            Notes: []string{sb.String()},
        },
    )
}
