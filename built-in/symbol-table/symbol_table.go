package symboltable

import "unique"

const (
    maxExpectedSymbols = 64
    expectedScopeDepth = 8
)

type SymbolType uint32

type SymbolTableScheme struct {
    lastSymbolType SymbolType
}

type SymbolTable struct {
    identifiers []unique.Handle[string]
    symbols     []SymbolData
    scopes      []uint32
}

type SymbolData struct {
    Type       SymbolType
    Reference  uint32
}

func NewScheme() *SymbolTableScheme {
    return &SymbolTableScheme{0}
}

func (s *SymbolTableScheme) NewSymbolType() SymbolType {
    s.lastSymbolType++

    return s.lastSymbolType
}

func New() *SymbolTable {
    return &SymbolTable{
        make([]unique.Handle[string], 0, maxExpectedSymbols),
        make([]SymbolData, 0, maxExpectedSymbols),
        make([]uint32, 1, expectedScopeDepth),
    }
}

func (s *SymbolTable) AddScope() {
    s.scopes = append(s.scopes, uint32(len(s.symbols)))
}

func (s *SymbolTable) RemoveScope() {
    if len(s.scopes) > 1 {
        scopeStart := s.scopes[len(s.scopes)-1]
        s.identifiers = s.identifiers[:scopeStart]
        s.symbols = s.symbols[:scopeStart]
        s.scopes = s.scopes[:len(s.scopes)-1]

        return
    }

    panic("attempt to remove the global scope")
}

// Returns a reference to the new symbol if ok or returns the already declared symbol of the current scope
func (s *SymbolTable) AddSymbol(identifier unique.Handle[string], data SymbolData) (*SymbolData, bool) {
    scopeStart := s.scopes[len(s.scopes)-1]

    for i := scopeStart; i < uint32(len(s.identifiers)); i++ {
        if s.identifiers[i] == identifier {
            return &s.symbols[i], false
        }
    }

    s.identifiers = append(s.identifiers, identifier)
    s.symbols = append(s.symbols, data)

    return &s.symbols[len(s.symbols) - 1], true
}

func (s *SymbolTable) Lookup(identifier unique.Handle[string]) (*SymbolData, bool) {
    for i := len(s.symbols) - 1; i >= 0; i-- {
        if s.identifiers[i] == identifier {
            return &s.symbols[i], true
        }
    }

    return nil, false
}
