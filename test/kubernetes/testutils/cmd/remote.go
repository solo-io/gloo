package cmd

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/solo-io/go-utils/threadsafe"
	"github.com/solo-io/k8s-utils/testutils/kube"

	"github.com/solo-io/gloo/pkg/utils/cmdutils"
)

var (
	_            cmdutils.Cmd   = &RemoteCmd{}
	_            cmdutils.Cmder = &RemoteCmder{}
	defaultCmder                = &RemoteCmder{}
)

// Command is a convenience wrapper over defaultCmder.Command
func Command(ctx context.Context, command string, args ...string) cmdutils.Cmd {
	return defaultCmder.Command(ctx, command, args...).
		WithStdout(io.Discard).
		WithStderr(io.Discard)
}

// RemoteCmderParams define some values to pass to the created RemoteCmd
type RemoteCmderParams struct {
	Receiver      io.Writer
	KubeContext   string
	Image         kube.Image
	FromContainer string
	FromNamespace string
	FromPod       string
}

// RemoteCmder is a factory for RemoteCmd, implementing Cmder
type RemoteCmder struct {
	RemoteCmderParams
}

func NewRemoteCmdFactory(ctx context.Context, params RemoteCmderParams) *RemoteCmder {
	r := &RemoteCmder{
		params,
	}
	return r
}

// Command returns a Cmd which includes the running process's `Environment`
func (c *RemoteCmder) Command(ctx context.Context, name string, arg ...string) cmdutils.Cmd {
	p := kube.EphemeralPodParams{
		Logger:        c.Receiver,
		KubeContext:   c.KubeContext,
		Image:         c.Image,
		FromContainer: c.FromContainer,
		FromNamespace: c.FromNamespace,
		FromPod:       c.FromPod,
		ExecCmdPath:   name,
		Args:          arg,
	}
	cmd := &RemoteCmd{
		p,
	}

	// By default, assign the env variables for the command
	// Consumers of this Cmd can then override it, if they want
	return cmd.WithEnv(os.Environ()...)
}

// RemoteCmd wraps os/exec.Cmd, implementing the cmdutils.Cmd interface
type RemoteCmd struct {
	kube.EphemeralPodParams
}

// WithEnv sets env
func (cmd *RemoteCmd) WithEnv(env ...string) cmdutils.Cmd {
	//disable DEBUG=1 from getting through to command
	for i, pair := range env {
		if strings.HasPrefix(pair, "DEBUG") {
			env = append(env[:i], env[i+1:]...)
			break
		}
	}

	cmd.Env = env
	return cmd
}

// WithStdin sets stdin
func (cmd *RemoteCmd) WithStdin(r io.Reader) cmdutils.Cmd {
	cmd.Stdin = r
	return cmd
}

// WithStdout set stdout
func (cmd *RemoteCmd) WithStdout(w io.Writer) cmdutils.Cmd {
	cmd.Stdout = w
	return cmd
}

// WithStderr sets stderr
func (cmd *RemoteCmd) WithStderr(w io.Writer) cmdutils.Cmd {
	cmd.Stderr = w
	return cmd
}

// Run runs the command
// If the returned error is non-nil, it should be of type *RunError
func (cmd *RemoteCmd) Run() cmdutils.RunError {
	var combinedOutput threadsafe.Buffer

	cmd.Stdout = io.MultiWriter(cmd.Stdout, &combinedOutput)
	cmd.Stderr = io.MultiWriter(cmd.Stderr, &combinedOutput)

	_, err := kube.ExecFromEphemeralPod(context.Background(), cmd.EphemeralPodParams)
	if err != nil {
		return &remoteRunError{
			command:    cmd.Args,
			output:     combinedOutput.Bytes(),
			inner:      err,
			stackTrace: errors.WithStack(err),
		}
	}
	return (*remoteRunError)(nil)
}
