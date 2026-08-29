package symboltable

import "unique"

type SymbolType uint32

type SymbolTable struct {
	symbols []Symbol
	scopes  []uint32
}

type Symbol struct {
	Identifier unique.Handle[string]
	Type       SymbolType
	Reference  uint32
}

func New() *SymbolTable {
	return &SymbolTable{[]Symbol{}, []uint32{0}}
}

func (s *SymbolTable) AddScope() {
	s.scopes = append(s.scopes, uint32(len(s.symbols)))
}

func (s *SymbolTable) RemoveScope() {
	if len(s.scopes) > 1 {
		scopeStart := s.scopes[len(s.scopes)-1]
		s.symbols = s.symbols[:scopeStart]
		s.scopes = s.scopes[:len(s.scopes)-1]

		return
	}

	// TODO error
}

// Returns a reference to the new symbol if ok or returns the already declared symbol of the current scope
func (s *SymbolTable) AddSymbol(newSymbol Symbol) (*Symbol, bool) {
	scopeStart := s.scopes[len(s.scopes)-1]

	for i := scopeStart; i < uint32(len(s.symbols)); i++ {
		if s.symbols[i].Identifier == newSymbol.Identifier {
			return &s.symbols[i], false
		}
	}

	s.symbols = append(s.symbols, newSymbol)

	return &s.symbols[len(s.symbols) - 1], true
}

func (s *SymbolTable) Lookup(identifier unique.Handle[string]) (*Symbol, bool) {
	for i := len(s.symbols) - 1; i >= 0; i-- {
		if s.symbols[i].Identifier == identifier {
			return &s.symbols[i], true
		}
	}

	return nil, false
}
