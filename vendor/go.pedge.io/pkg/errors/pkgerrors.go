package pkgerrors // import "go.pedge.io/pkg/errors"

import (
	"bytes"
	"errors"
	"fmt"
)

// New creates a new error.
//
// The first argument is the message, the rest are key/value pairs.
func New(message string, keyValues ...interface{}) error {
	len := len(keyValues)
	if len == 0 {
		return errors.New(message)
	}
	buffer := bytes.NewBuffer(nil)
	if message != "" {
		_, _ = buffer.WriteString(message)
		_ = buffer.WriteByte(' ')
	}
	if len%2 != 0 {
		keyValues = append(keyValues, "MISSING")
		len++
	}
	for i := 0; i < len; i += 2 {
		if i != 0 {
			_ = buffer.WriteByte(' ')
		}
		_, _ = buffer.WriteString(fmt.Sprintf("%v", keyValues[i]))
		_ = buffer.WriteByte('=')
		_, _ = buffer.WriteString(fmt.Sprintf("%v", keyValues[i+1]))
	}
	return errors.New(buffer.String())
}
