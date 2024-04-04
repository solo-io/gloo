package cmdutils

import (
	"context"
	"io"
)

// Cmd abstracts over running a command somewhere
type Cmd interface {
	// Run executes the command (like os/exec.Cmd.Run)
	// It returns a *RunError if there is any error, nil otherwise
	Run() *RunError
	// SetEnv sets the Env variables for the Cmd
	// Each entry should be of the form "key=value"
	SetEnv(...string) Cmd

	SetStdin(reader io.Reader) Cmd
	SetStdout(io.Writer) Cmd
	SetStderr(io.Writer) Cmd
}

// Cmder abstracts over creating commands
type Cmder interface {
	Command(context.Context, string, ...string) Cmd
}
