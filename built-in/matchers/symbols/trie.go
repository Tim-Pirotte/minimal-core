package symbols

import (
	"fmt"
	"minimal/minimal-core/built-in/tokenizer"
)

type trieNode struct {
	leaf bool
	token tokenizer.TokenType
	children [256]*trieNode
}

func updateTrie(root *trieNode, text string, tokenType tokenizer.TokenType) error {
	node := root

	for i := 0; i < len(text); i++ {
        char := text[i] 

        if node.children[char] == nil {
            node.children[char] = &trieNode{}
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
