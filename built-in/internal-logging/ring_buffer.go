package logging

import (
	"io"
	"sync"
)

type ringBuffer struct {
	mutex     sync.RWMutex
	data   []byte
	cursor int
	full   bool
}

func newRingBuffer(maxBytes uint) *ringBuffer {
	return &ringBuffer{
		data: make([]byte, maxBytes),
	}
}

func (r *ringBuffer) Write(toWrite []byte) (originalLength int, err error) {
	if len(r.data) == 0 {
        return len(toWrite), nil 
    }

	r.mutex.Lock()
	defer r.mutex.Unlock()

	originalLength = len(toWrite)

	if len(toWrite) > len(r.data) {
		toWrite = toWrite[:len(r.data)]
	}

	remainingSpace := len(r.data) - r.cursor

	if len(toWrite) <= remainingSpace {
		copy(r.data[r.cursor:], toWrite)
	} else {
		copy(r.data[r.cursor:], toWrite[:remainingSpace])
		copy(r.data[0:], toWrite[remainingSpace:])
		r.full = true
	}

	r.cursor = (r.cursor + len(toWrite)) % len(r.data)

	return originalLength, nil
}

func (r *ringBuffer) writeTo(w io.Writer) (int64, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    var total int64

	if r.full {
		n, err := w.Write(r.data[r.cursor:])
		total += int64(n)

		if err != nil {
			return total, err
		}
	}

	n, err := w.Write(r.data[:r.cursor])
	total += int64(n)

	return total, err
}
