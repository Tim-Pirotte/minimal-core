package symboltable

// go test -bench="." -benchtime="5s"  -cpuprofile="cpu.prof" -memprofile="mem.prof" -v
// go tool pprof -http=:8080 cpu.prof

import (
	"strconv"
	"testing"
	"unique"
)

func BenchmarkCreation(b *testing.B) {
    scopes := 50
    scopeDepth := 5
    symbols := 7

    identifiers := make([]string, scopeDepth * symbols)

    for i := range len(identifiers) {
        identifiers[i] = strconv.Itoa(i)
    }

    b.ResetTimer()
    for b.Loop() {
        s := New()

        for range scopes {
            for i := range scopeDepth {
                s.AddScope()

                for j := range symbols {
                    s.AddSymbol(
                        unique.Make(identifiers[j * scopeDepth + i]),
                        SymbolData{Type: 0, Reference: 0},
                    )
                }
            }

            for range scopeDepth {
                s.RemoveScope()
            }
        }
    }
}

// func BenchmarkBestCaseLookup(b *testing.B) {

// }

// func BenchmarkWorstCaseLookup(b *testing.B) {

// }

// func BenchmarkAverageCaseLookup(b *testing.B) {

// }
