package diff

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestGetDiff(t *testing.T) {
	isEqual := func(x, y string) bool {
		return x == y
	}

	tests := []struct {
		name     string
		a        []string
		b        []string
		expected []DiffPart[string]
	}{
		{
			name:     "Empty inputs",
			a:        []string{},
			b:        []string{},
			expected: []DiffPart[string]{},
		},
		{
			name: "Identical slices",
			a:    []string{"apple", "banana"},
			b:    []string{"apple", "banana"},
			expected: []DiffPart[string]{
				{Type: Equal, Value: "apple"},
				{Type: Equal, Value: "banana"},
			},
		},
		{
			name: "Pure insertions",
			a:    []string{},
			b:    []string{"apple", "banana"},
			expected: []DiffPart[string]{
				{Type: Insert, Value: "apple"},
				{Type: Insert, Value: "banana"},
			},
		},
		{
			name: "Pure deletions",
			a:    []string{"apple", "banana"},
			b:    []string{},
			expected: []DiffPart[string]{
				{Type: Delete, Value: "apple"},
				{Type: Delete, Value: "banana"},
			},
		},
		{
			name: "Mixed changes",
			a:    []string{"A", "B", "C"},
			b:    []string{"A", "D", "C"},
			expected: []DiffPart[string]{
				{Type: Equal, Value: "A"},
				{Type: Delete, Value: "B"},
				{Type: Insert, Value: "D"},
				{Type: Equal, Value: "C"},
			},
		},
		{
			name: "Complex reordering and edits",
			a:    []string{"X", "A", "B", "Y"},
			b:    []string{"A", "B", "Z"},
			expected: []DiffPart[string]{
				{Type: Delete, Value: "X"},
				{Type: Equal, Value: "A"},
				{Type: Equal, Value: "B"},
				{Type: Delete, Value: "Y"},
				{Type: Insert, Value: "Z"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetDiff(tt.a, tt.b, isEqual)
			
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("\nGetDiff() mismatch:\nExpected: %v\nActual:   %v", formatDiff(tt.expected), formatDiff(actual))
			}
		})
	}
}

func formatDiff(parts []DiffPart[string]) []string {
	res := make([]string, len(parts))
	for i, p := range parts {
		prefix := " "
		switch p.Type {
		case Insert:
			prefix = "+"
		case Delete:
			prefix = "-"
		}
		res[i] = prefix + p.Value
	}
	return res
}

func byteEqual(a, b byte) bool {
	return a == b
}

func FuzzGetDiff(f *testing.F) {
	f.Add([]byte("abc"), []byte("abc"))
	f.Add([]byte(""), []byte("abc"))
	f.Add([]byte("abc"), []byte(""))
	f.Add([]byte(""), []byte(""))

	f.Fuzz(func(t *testing.T, a []byte, b []byte) {
		GetDiff(a, b, byteEqual)
	})
}

func generateRandomBytes(count int) []byte {
	rng := rand.New(rand.NewSource(42))
	b := make([]byte, count)
	_, _ = rng.Read(b)

	return b
}

func BenchmarkIdentical(b *testing.B) {
	src := generateRandomBytes(1_000_000)
	dst := make([]byte, len(src))
	copy(dst, src)

	for b.Loop() {
		GetDiff(src, dst, byteEqual)
	}
}

func BenchmarkInsert(b *testing.B) {
	src := []byte{}
	dst := generateRandomBytes(1_000_000)

	for b.Loop() {
		GetDiff(src, dst, byteEqual)
	}
}

func BenchmarkDelete(b *testing.B) {
	src := generateRandomBytes(1_000_000)
	dst := []byte{}

	for b.Loop() {
		GetDiff(src, dst, byteEqual)
	}
}

func BenchmarkRandom(b *testing.B) {
	src := generateRandomBytes(1_000_000)
	dst := generateRandomBytes(1_000_000)

	for b.Loop() {
		GetDiff(src, dst, byteEqual)
	}
}
