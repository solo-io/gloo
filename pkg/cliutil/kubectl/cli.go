package kubectl

import (
	"bytes"
	"context"
	"errors"

	"github.com/solo-io/gloo/pkg/utils/cmdutils"

	"io"
	"os"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
)

// Cli is a utility for executing `kubectl` commands
type Cli struct {
	// receiver is the default destination for the kubectl stdout and stderr
	receiver io.Writer

	// kubeContext is the optional value of the context for a given kubernetes cluster
	// If it is not supplied, no context will be included in the command
	kubeContext string
}

// NewCli returns NewCliWithKubeContext with an empty Kubernetes context
func NewCli(receiver io.Writer) (*Cli, error) {
	return NewCliWithKubeContext(receiver, "")
}

// NewCliWithKubeContext returns an implementation of KubectlCmd, or an error if one of the provided receivers was nil
func NewCliWithKubeContext(receiver io.Writer, kubeContext string) (*Cli, error) {
	if receiver == nil {
		return nil, errors.New("receiver must not be nil")
	}

	return &Cli{
		receiver:    receiver,
		kubeContext: kubeContext,
	}, nil
}

func (c *Cli) Command(ctx context.Context, args ...string) cmdutils.Cmd {
	if c.kubeContext != "" {
		args = append([]string{"--context", c.kubeContext}, args...)
	}

	cmd := cmdutils.Command(ctx, "kubectl", args...)

	//disable DEBUG=1 from getting through to kube
	env := os.Environ()
	for i, pair := range env {
		if strings.HasPrefix(pair, "DEBUG") {
			env = append(env[:i], env[i+1:]...)
			break
		}
	}
	cmd.SetEnv(env...)

	// For convenience, we set the stdout and stderr to the receiver
	// This can still be overwritten by consumers who use the commands
	cmd.SetStdout(c.receiver)
	cmd.SetStderr(c.receiver)
	return cmd
}

func (c *Cli) ExecuteCommand(ctx context.Context, args ...string) error {
	return c.Command(ctx, args...).Run().Cause()
}

func (c *Cli) ApplyCmd(ctx context.Context, content []byte, extraArgs ...string) cmdutils.Cmd {
	args := append([]string{"apply", "-f", "-"}, extraArgs...)

	cmd := c.Command(ctx, args...)
	cmd.SetStdin(bytes.NewBuffer(content))
	return cmd
}

func (c *Cli) Apply(ctx context.Context, content []byte, extraArgs ...string) error {
	return c.ApplyCmd(ctx, content, extraArgs...).Run().Cause()
}

func (c *Cli) DeleteCmd(ctx context.Context, content []byte, extraArgs ...string) cmdutils.Cmd {
	args := append([]string{"delete", "-f", "-"}, extraArgs...)

	cmd := c.Command(ctx, args...)
	cmd.SetStdin(bytes.NewBuffer(content))
	return cmd
}

func (c *Cli) Delete(ctx context.Context, content []byte, extraArgs ...string) error {
	return c.DeleteCmd(ctx, content, extraArgs...).Run().Cause()
}

func (c *Cli) CopyCmd(ctx context.Context, from, to string) cmdutils.Cmd {
	return c.Command(ctx, "cp", from, to)
}

func (c *Cli) Copy(ctx context.Context, from, to string) error {
	return c.CopyCmd(ctx, from, to).Run().Cause()
}

func (c *Cli) StartPortForward(ctx context.Context, options ...portforward.Option) (portforward.PortForwarder, error) {
	options = append([]portforward.Option{
		// We define some default values, which users can then override
		portforward.WithWriters(c.receiver, c.receiver),
		portforward.WithKubeContext(c.kubeContext),
	}, options...)

	portForwarder := portforward.NewCliPortForwarder(options...)
	err := portForwarder.Start(
		ctx,
		retry.LastErrorOnly(true),
		retry.Delay(100*time.Millisecond),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(5),
	)
	return portForwarder, err
}
