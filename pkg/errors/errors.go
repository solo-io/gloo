package errors

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type alreadyExistsErr struct {
	meta core.Metadata
}

func (err *alreadyExistsErr) Error() string {
	return fmt.Sprintf("already exists: %v", err.meta)
}

func NewAlreadyExistsErr(meta core.Metadata) *alreadyExistsErr {
	return &alreadyExistsErr{meta: meta}
}

func IsAlreadyExists(err error) bool {
	switch err.(type) {
	case *alreadyExistsErr:
		return true
	}
	return false
}

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

func Errors(msgs []string) error {
	return errors.Errorf(strings.Join(msgs, "\n"))
}
