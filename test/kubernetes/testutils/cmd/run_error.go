package cmd

import (
	"fmt"

	"github.com/solo-io/gloo/pkg/utils/cmdutils"
)

var _ cmdutils.RunError = &remoteRunError{}

type remoteRunError struct {
	command    []string // [Name Args...]
	output     []byte   // Captured Stdout / Stderr of the command
	inner      error    // Underlying error if any
	stackTrace error
}

func (e *remoteRunError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("command \"%s\" failed with error: %v", e.PrettyCommand(), e.inner)
}

// PrettyCommand pretty prints the command in a way that could be pasted
// into a shell
func (e *remoteRunError) PrettyCommand() string {
	if e == nil {
		return "RunError is nil"
	}

	if len(e.command) == 0 {
		return "no command args"
	}

	if len(e.command) == 1 {
		return e.command[0]
	}

	// The above cases should not happen, but we defend against it
	return cmdutils.PrettyCommand(e.command[0], e.command[1:]...)
}

func (e *remoteRunError) OutputString() string {
	if e == nil {
		return ""
	}
	return string(e.output)
}

// Cause mimics github.com/pkg/errors's Cause pattern for errors
func (e *remoteRunError) Cause() error {
	if e == nil {
		return nil
	}
	return e.stackTrace
}
