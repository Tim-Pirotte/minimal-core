package diff

import (
	"math/rand"
	"reflect"
	"testing"
)

func areStringsEqual(a, b string) bool {
	return a == b
}

func TestEmpty(t *testing.T) {
	a := []string{}
	b := []string{}

	expected := []DiffPart[string]{}

	actual := GetDiff(a, b, areStringsEqual)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nExpected: %v\nActual:   %v", formatDiff(expected), formatDiff(actual))
	}
}

func TestIdentical(t *testing.T) {
	a := []string{"A", "B"}
	b := []string{"A", "B"}

	expected := []DiffPart[string]{
		{Type: Equal, Value: "A"},
		{Type: Equal, Value: "B"},
	}

	actual := GetDiff(a, b, areStringsEqual)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nExpected: %v\nActual: %v", formatDiff(expected), formatDiff(actual))
	}
}

func TestInsert(t *testing.T) {
	a := []string{}
	b := []string{"A", "B"}

	expected := []DiffPart[string]{
		{Type: Insert, Value: "A"},
		{Type: Insert, Value: "B"},
	}

	actual := GetDiff(a, b, areStringsEqual)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nExpected: %v\nActual: %v", formatDiff(expected), formatDiff(actual))
	}
}

func TestDelete(t *testing.T) {
	a := []string{"A", "B"}
	b := []string{}

	expected := []DiffPart[string]{
		{Type: Delete, Value: "A"},
		{Type: Delete, Value: "B"},
	}

	actual := GetDiff(a, b, areStringsEqual)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nExpected: %v\nActual: %v", formatDiff(expected), formatDiff(actual))
	}
}

func TestReplacePrefix(t *testing.T) {
	a := []string{"A", "B", "C"}
	b := []string{"D", "B", "C"}

	expected := []DiffPart[string]{
		{Type: Delete, Value: "A"},
		{Type: Insert, Value: "D"},
		{Type: Equal, Value: "B"},
		{Type: Equal, Value: "C"},
	}

	actual := GetDiff(a, b, areStringsEqual)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nExpected: %v\nActual: %v", formatDiff(expected), formatDiff(actual))
	}
}

func TestReplaceMiddle(t *testing.T) {
	a := []string{"A", "B", "C"}
	b := []string{"A", "D", "C"}

	expected := []DiffPart[string]{
		{Type: Equal, Value: "A"},
		{Type: Delete, Value: "B"},
		{Type: Insert, Value: "D"},
		{Type: Equal, Value: "C"},
	}

	actual := GetDiff(a, b, areStringsEqual)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nExpected: %v\nActual: %v", formatDiff(expected), formatDiff(actual))
	}
}

func TestReplaceSuffix(t *testing.T) {
	a := []string{"A", "B", "C"}
	b := []string{"A", "B", "D"}

	expected := []DiffPart[string]{
		{Type: Equal, Value: "A"},
		{Type: Equal, Value: "B"},
		{Type: Delete, Value: "C"},
		{Type: Insert, Value: "D"},
	}

	actual := GetDiff(a, b, areStringsEqual)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nExpected: %v\nActual: %v", formatDiff(expected), formatDiff(actual))
	}
}

func TestComplex(t *testing.T) {
	a := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	b := []string{"K", "L", "C", "D", "M", "N", "G", "H", "O", "P"}

	expected := []DiffPart[string]{
		{Type: Delete, Value: "A"},
		{Type: Delete, Value: "B"},
		{Type: Insert, Value: "K"},
		{Type: Insert, Value: "L"},
		{Type: Equal,  Value: "C"},
		{Type: Equal,  Value: "D"},
		{Type: Delete, Value: "E"},
		{Type: Delete, Value: "F"},
		{Type: Insert, Value: "M"},
		{Type: Insert, Value: "N"},
		{Type: Equal,  Value: "G"},
		{Type: Equal,  Value: "H"},
		{Type: Delete, Value: "I"},
		{Type: Delete, Value: "J"},
		{Type: Insert, Value: "O"},
		{Type: Insert, Value: "P"},
	}

	actual := GetDiff(a, b, areStringsEqual)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nExpected: %v\nActual: %v", formatDiff(expected), formatDiff(actual))
	}
}

func formatDiff(parts []DiffPart[string]) []string {
	res := make([]string, len(parts))

	for i, p := range parts {
		prefix := "  "

		switch p.Type {
		case Insert:
			prefix = "+ "
		case Delete:
			prefix = "- "
		}

		res[i] = prefix + p.Value
	}

	return res
}

func areBytesEqual(a, b byte) bool {
	return a == b
}

func FuzzGetDiff(f *testing.F) {
	f.Add([]byte("abc"), []byte("abc"))
	f.Add([]byte(""), []byte("abc"))
	f.Add([]byte("abc"), []byte(""))
	f.Add([]byte(""), []byte(""))
	f.Add([]byte("abc"), []byte("adc"))

	f.Fuzz(func(t *testing.T, a []byte, b []byte) {
		GetDiff(a, b, areBytesEqual)
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
		GetDiff(src, dst, areBytesEqual)
	}
}

func BenchmarkInsert(b *testing.B) {
	src := []byte{}
	dst := generateRandomBytes(1_000_000)

	for b.Loop() {
		GetDiff(src, dst, areBytesEqual)
	}
}

func BenchmarkDelete(b *testing.B) {
	src := generateRandomBytes(1_000_000)
	dst := []byte{}

	for b.Loop() {
		GetDiff(src, dst, areBytesEqual)
	}
}

func BenchmarkRandom(b *testing.B) {
	src := generateRandomBytes(1_000_000)
	dst := generateRandomBytes(1_000_000)

	for b.Loop() {
		GetDiff(src, dst, areBytesEqual)
	}
}
