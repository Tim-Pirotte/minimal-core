package symboltable

import (
	"strconv"
	"testing"
	"unique"
)

func TestEmpty(t *testing.T) {
    s := New()

    if len(s.identifiers) != 0 {
        t.Errorf("Expected len(s.identifiers) to be 0 but got %d", len(s.identifiers))
    }

    if len(s.symbols) != 0 {
        t.Errorf("Expected len(s.symbols) to be 0 but got %d", len(s.symbols))
    }

    if len(s.scopes) != 1 {
        t.Errorf("Expected len(s.scopes) to be 1 but got %d", len(s.scopes))
	}
}

func TestGlobal(t *testing.T) {
    scheme := NewScheme()
    symbolType := scheme.NewSymbolType()

    s := New()

    identifier := unique.Make("global")
    data := SymbolData{Type: symbolType, Reference: 7}

    s.AddSymbol(identifier, data)

    if l, ok := s.Lookup(identifier); ok {
        if !(*l == data) {
            t.Errorf("Expected the lookup to be %+v but got %+v", data, *l)
        }
    } else {
        t.Error("Expected the lookup to succeed")
    }
}

func TestScopes(t *testing.T) {
    s := New()

    for range 5 {
        for range 3 {
            s.AddScope()
        }

        if len(s.scopes) != 4 {
            t.Errorf("Expected len(s.scopes) to be 4 but got %d", len(s.scopes))
        }

        for range 3 {
            s.RemoveScope()
        }

        if len(s.scopes) != 1 {
            t.Errorf("Expected len(s.scopes) to be 1 but got %d", len(s.scopes))
        }
    }
}

func TestSymbols(t *testing.T) {
    s := New()

    for i := range 3 {
        s.AddScope()

        for j := range 2 {
            s.AddSymbol(unique.Make(strconv.Itoa(i * 2 + j)), SymbolData{})
        }
    }

    for i := range 3 * 2 {
        if _, ok := s.Lookup(unique.Make(strconv.Itoa(i))); !ok {
            t.Errorf("Could not find identifier %d", i)
        }
    }
}

func TestDuplicateSymbol(t *testing.T) {
    s := New()

    if symbol, ok := s.AddSymbol(unique.Make("duplicate"), SymbolData{Reference: 1}); ok {
        if symbol.Reference != 1 {
            t.Errorf("Expected symbol.Reference to be 1 but got %d", symbol.Reference)
        }
    } else {
        t.Errorf("Expected the first symbol to be added")
    }

    if symbol, ok := s.AddSymbol(unique.Make("duplicate"), SymbolData{Reference: 2}); !ok {
        if symbol.Reference != 1 {
            t.Errorf("Expected symbol.Reference to be 1 but got %d", symbol.Reference)
        }
    } else {
        t.Errorf("Expected the second symbol to be a duplicate")
    }
}

func TestRemoveTooManyScopes(t *testing.T) {
    s := New()

    defer func() {
        if r := recover(); r != nil {
            if r != "attempt to remove the global scope" {
                t.Errorf("Expected the panic to be 'attempt to remove the global scope' but got '%+v'", r)
            }
        } else {
            t.Errorf("Expected RemoveScope to panic")
        }
    }()

    for range 3 {
        s.AddScope()
    }

    for range 4 {
        s.RemoveScope()
    }
}
