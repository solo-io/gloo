package pkgrandom // import "go.pedge.io/pkg/random"

import (
	"bytes"
	"math/rand"
)

// Bytes returns a random string of a-z.
func Bytes(size int) []byte {
	return getBuffer(size).Bytes()
}

// String returns a random string of a-z.
func String(size int) string {
	return getBuffer(size).String()
}

func getBuffer(size int) *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	for i := 0; i < size; i++ {
		_ = buffer.WriteByte(byte('a' + (int(rand.Uint32()) % 25)))
	}
	return buffer
}
