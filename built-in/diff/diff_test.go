package diff

import (
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
			actual := getDiff(tt.a, tt.b, isEqual)
			
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("\nGetDiff() mismatch:\nExpected: %v\nActual:   %v", formatDiff(tt.expected), formatDiff(actual))
			}
		})
	}
}

// Helper to format output neatly if a test fails
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
