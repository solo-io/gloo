package internalerrors

import "github.com/pkg/errors"

// internalError is used to identify errors caused by a malfunctioning component.
// Mostly used by the translator to differentiate errors from user config
// and errors caused by internal failures
type internalError struct {
	Err error
}

func New(err error, format string, args ...interface{}) *internalError {
	return &internalError{
		Err: errors.Wrapf(err, format, args...),
	}
}

func (err *internalError) Error() string {
	return err.Err.Error()
}

func IsInternal(err error) bool {
	_, ok := err.(*internalError)
	return ok
}
