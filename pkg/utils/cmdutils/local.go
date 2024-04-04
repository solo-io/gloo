package cmdutils

import (
	"context"
	"io"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/solo-io/go-utils/threadsafe"
)

var (
	_            Cmd   = &LocalCmd{}
	_            Cmder = &LocalCmder{}
	DefaultCmder       = &LocalCmder{}
)

// Command is a convenience wrapper over DefaultCmder.Command
func Command(ctx context.Context, command string, args ...string) Cmd {
	return DefaultCmder.Command(ctx, command, args...)
}

// LocalCmder is a factory for LocalCmd, implementing Cmder
type LocalCmder struct{}

// Command is like Command but includes a context
func (c *LocalCmder) Command(ctx context.Context, name string, arg ...string) Cmd {
	return &LocalCmd{
		Cmd: exec.CommandContext(ctx, name, arg...),
	}
}

// LocalCmd wraps os/exec.Cmd, implementing the kind/pkg/exec.Cmd interface
type LocalCmd struct {
	*exec.Cmd
}

// SetWorkingDirectory sets the working directory
func (cmd *LocalCmd) SetWorkingDirectory(dir string) Cmd {
	cmd.Dir = dir
	return cmd
}

// SetEnv sets env
func (cmd *LocalCmd) SetEnv(env ...string) Cmd {
	cmd.Env = env
	return cmd
}

// SetStdin sets stdin
func (cmd *LocalCmd) SetStdin(r io.Reader) Cmd {
	cmd.Stdin = r
	return cmd
}

// SetStdout set stdout
func (cmd *LocalCmd) SetStdout(w io.Writer) Cmd {
	cmd.Stdout = w
	return cmd
}

// SetStderr sets stderr
func (cmd *LocalCmd) SetStderr(w io.Writer) Cmd {
	cmd.Stderr = w
	return cmd
}

// Run runs the command
// If the returned error is non-nil, it should be of type *RunError
func (cmd *LocalCmd) Run() *RunError {
	var combinedOutput threadsafe.Buffer

	cmd.Stdout = io.MultiWriter(cmd.Stdout, &combinedOutput)
	cmd.Stderr = io.MultiWriter(cmd.Stderr, &combinedOutput)

	if err := cmd.Cmd.Run(); err != nil {
		return &RunError{
			Command:    cmd.Args,
			Output:     combinedOutput.Bytes(),
			Inner:      err,
			StackTrace: errors.WithStack(err),
		}
	}
	return nil
}
