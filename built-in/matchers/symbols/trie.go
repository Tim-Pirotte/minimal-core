package symbols

import (
	"fmt"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
)

const byteValueCount = 256

type trieNode struct {
	leaf bool
	token tokenizerv2.TokenType
	children [byteValueCount]*trieNode
}

func updateTrie(root *trieNode, text string, tokenType tokenizerv2.TokenType) error {
	node := root

	for _, char := range []byte(text) {
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
