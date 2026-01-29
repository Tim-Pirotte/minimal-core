package tokenizer

import (
	"fmt"
	"minimal/minimal-core/domain"
)

func updateTrie(root *trieNode, text string, tokenType domain.TokenType) error {
	node := root

	for _, char := range text {
		if node.children[char] == nil {
			node.children[char] = &trieNode{children: map[rune]*trieNode{}}
		}

		node = node.children[char]
	}

	if node.leaf {
		return fmt.Errorf("%v has already been declared", text)
	}

	node.leaf = true
	node.token = tokenType

	return nil
}
