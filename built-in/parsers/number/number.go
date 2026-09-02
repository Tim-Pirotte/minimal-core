package numberparser

import (
	"fmt"
	"math"
	"math/bits"
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/parsers/pratt"
)

// TODO this could be ambiguous with units (e.g. 0 x with x defined as a type)

type NumberParser struct {
    baseParsers    map[byte]baseParser
    integerType    lexer.TokenType
    identifierType lexer.TokenType
    integer        ast.NodeType
    float          ast.NodeType
    integers       []IntegerLiteral
    floats         []FloatLiteral
}

type BaseParserData struct {
    parser   BaseParser
    baseChar byte
    base     uint
}

type BaseParser interface {
    Parse(*NumberParser, *lexer.Lexer, base) []ast.Node
}

type base uint

type baseParser struct {
    baseParser BaseParser
    base       base
}

type IntegerLiteral struct {
    parts []uint
}

type FloatLiteral struct {
    IntegerPart    IntegerLiteral
    FractionalPart IntegerLiteral
}

func New(
    baseParsers []BaseParserData,
    integerType, identifierType lexer.TokenType,
    integer, float ast.NodeType,
) NumberParser {
    np := NumberParser{
        integerType: integerType,
        identifierType: identifierType,
        integer: integer,
        float: float,
    }

    for _, bp := range baseParsers {
        if bp.base == 0 {
            // TODO error
        }

        maxSafe := math.MaxUint / bp.base

        // TODO is this correct
        if bp.base - 1 >= maxSafe {
            // TODO error
        }

        if '0' <= bp.baseChar && bp.baseChar <= '9' {
            // TODO error
        }

        if _, ok := np.baseParsers[bp.baseChar]; ok {
            // TODO error
        }

        np.baseParsers[bp.baseChar] = baseParser{bp.parser, base(bp.base)}
    }

    return np
}

func (n *NumberParser) AddInteger(integer IntegerLiteral) uint32 {
    reference := len(n.integers)
    n.integers = append(n.integers, integer)

    return uint32(reference)
}

func (n *NumberParser) AddFloat(float FloatLiteral) uint32 {
    reference := len(n.floats)
    n.floats = append(n.floats, float)

    return uint32(reference)
}

func (n *NumberParser) GetTokenType() lexer.TokenType {
    return n.integerType
}

func (n *NumberParser) ParsePrefix(_ *pratt.PrattParser, l *lexer.Lexer, minBindingPower uint) []ast.Node {
    prefix := l.Peek(0).Value

    l.Advance()

    if prefix != "0" {
        decimal := IntegerLiteral{}

        for _, c := range []byte(prefix) {
            decimal.AddDigit(base(10), uint(c - '0'))
        }

        for t := l.Peek(0); t.Type == n.integerType || t.Type == n.identifierType; t = l.Peek(0) {
            l.Advance()

            if t.Type == n.integerType {
                for _, c := range []byte(prefix) {
                    decimal.AddDigit(base(10), uint(c - '0'))
                }
            } else {
                for _, c := range []byte(t.Value) {
                    if !('0' <= c && c <= '9') {
                        // TODO error message

                        // TODO What should we now return?
                    }

                    decimal.AddDigit(base(10), uint(c - '0'))
                }
            }
        }

        reference := n.AddInteger(decimal)

        return []ast.Node{{Type: n.integer, Reference: reference}}
    }

    base := l.Peek(0)

    if base.Type != n.identifierType {
        reference := n.AddInteger(IntegerLiteral{})

        return []ast.Node{{Type: n.integer, Reference: reference}}
    }

    if bp, ok := n.baseParsers[base.Value[0]]; ok {
        return bp.baseParser.Parse(n, l, bp.base)
    }

    // Invalid base
    l.Advance()

    for t := l.Peek(0).Type; t == n.integerType || t == n.identifierType; t = l.Peek(0).Type {
        l.Advance()
    }

    // TODO what to return?
    return []ast.Node{}
}

func (i *IntegerLiteral) AddDigit(base base, digit uint) {
    if digit >= uint(base) {
        // TODO
        panic(fmt.Sprintf("digit %d is not valid for base %d", digit, base))
    }

    carry := digit

    for j := range i.parts {
        // I think that hi can only reach MAX - 1 so we can ignore the carry below
        hi, lo := bits.Mul(i.parts[j], uint(base))
        lo, carryOut := bits.Add(lo, carry, 0)
        hi, _ = bits.Add(hi, 0, carryOut)

        i.parts[j] = lo
        carry = hi
    }

    if carry > 0 {
		i.parts = append(i.parts, carry)
	}
}
