package tokenizer

import "fmt"

func (t *Tokenizer) skipWhiteSpace() {
	currentByte, err := t.reader.Peek(1)

	for err == nil && (currentByte[0] == ' ' || currentByte[0] == '\t') {
		_, discardErr := t.reader.Discard(1)
		t.position++

		if discardErr != nil {
			panic(fmt.Sprint("Could not properly discard:", discardErr))
		}

		currentByte, err = t.reader.Peek(1)
	}
}