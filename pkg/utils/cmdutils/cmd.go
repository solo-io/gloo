package cmdutils

import (
	"context"
	"io"
)

// Cmd abstracts over running a command somewhere, this is useful for testing
type Cmd interface {
	// Run executes the command (like os/exec.Cmd.Run), it should return
	// a *RunError if there is any error
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
	// command, args..., just like os/exec.Cmd
	Command(string, ...string) Cmd
	CommandContext(context.Context, string, ...string) Cmd
}
