package cmdutils

import (
	"fmt"
	"strconv"
	"strings"
)

// RunError represents an error running a Cmd
type RunError struct {
	Command    []string // [Name Args...]
	Output     []byte   // Captured Stdout / Stderr of the command
	Inner      error    // Underlying error if any
	StackTrace error
}

var _ error = &RunError{}

func (e *RunError) Error() string {
	return fmt.Sprintf("command \"%s\" failed with error: %v", e.PrettyCommand(), e.Inner)
}

// PrettyCommand pretty prints the command in a way that could be pasted
// into a shell
func (e *RunError) PrettyCommand() string {
	return PrettyCommand(e.Command[0], e.Command[1:]...)
}

func (e *RunError) OutputString() string {
	if e == nil {
		return ""
	}
	return string(e.Output)
}

// Cause mimics github.com/pkg/errors's Cause pattern for errors
func (e *RunError) Cause() error {
	if e == nil {
		return nil
	}
	return e.Inner
}

// PrettyCommand takes arguments identical to Cmder.Command,
// it returns a pretty printed command that could be pasted into a shell
func PrettyCommand(name string, args ...string) string {
	var out strings.Builder
	out.WriteString(strconv.Quote(name))
	for _, arg := range args {
		out.WriteByte(' ')
		out.WriteString(strconv.Quote(arg))
	}
	return out.String()
}
