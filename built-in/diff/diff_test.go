package diff

import (
	"fmt"
	"testing"
)

func isEqual(a, b string) bool {
	return a == b
}

func TestDiff(t *testing.T) {
	a := []string{"a", "b", "c", "d"}
	b := []string{"b", "c", "a"}

	for _, line := range getDiff(a, b, isEqual) {
		switch line.Type {
		case Insert:
			fmt.Print("+ ")
		case Delete:
			fmt.Print("- ")
		default:
			fmt.Print("  ")
		}

		fmt.Println(line.Value)
	}
	
	t.Error()
}
