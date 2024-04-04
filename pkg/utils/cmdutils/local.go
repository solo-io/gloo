package cmdutils

import (
	"context"
	"github.com/pkg/errors"
	"github.com/solo-io/go-utils/threadsafe"
	"io"
	"os/exec"
)

// LocalCmd wraps os/exec.Cmd, implementing the kind/pkg/exec.Cmd interface
type LocalCmd struct {
	*exec.Cmd
}

var _ Cmd = &LocalCmd{}

// LocalCmder is a factory for LocalCmd, implementing Cmder
type LocalCmder struct{}

var _ Cmder = &LocalCmder{}

// Command returns a new exec.Cmd backed by Cmd
func (c *LocalCmder) Command(name string, arg ...string) Cmd {
	return &LocalCmd{
		Cmd: exec.Command(name, arg...),
	}
}

// CommandContext is like Command but includes a context
func (c *LocalCmder) CommandContext(ctx context.Context, name string, arg ...string) Cmd {
	return &LocalCmd{
		Cmd: exec.CommandContext(ctx, name, arg...),
	}
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
