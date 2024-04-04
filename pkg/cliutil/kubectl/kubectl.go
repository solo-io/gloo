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
	// receiver is destination for the kubectl stdout and stderr
	receiver io.Writer

	// kubeContext is the optional value of the context for a given kubernetes cluster
	// If it is not supplied, no context will be included in the command
	kubeContext string
}

// NewKubectl returns NewKubectlWithKubeContext with an empty Kubernetes context
func NewKubectl(receiver io.Writer) (*Cli, error) {
	return NewKubectlWithKubeContext(receiver, "")
}

// NewKubectlWithKubeContext returns an implementation of KubectlCmd, or an error if one of the provided receivers was nil
func NewKubectlWithKubeContext(receiver io.Writer, kubeContext string) (*Cli, error) {
	if receiver == nil {
		return nil, errors.New("receiver must not be nil")
	}

	return &Cli{
		receiver:    receiver,
		kubeContext: kubeContext,
	}, nil
}

func (k *Cli) Command(ctx context.Context, args ...string) cmdutils.Cmd {
	if k.kubeContext != "" {
		args = append([]string{"--context", k.kubeContext}, args...)
	}

	cmd := cmdutils.CommandContext(ctx, "kubectl", args...)
	env := os.Environ()
	// disable DEBUG=1 from getting through to kube
	for i, pair := range env {
		if strings.HasPrefix(pair, "DEBUG") {
			env = append(env[:i], env[i+1:]...)
			break
		}
	}

	cmd.SetEnv(env...)

	// For convenience, we set the stdout and stderr to the receiver
	// This can still be overwritten by consumers who use the commands
	cmd.SetStdout(k.receiver)
	cmd.SetStderr(k.receiver)
	return cmd
}

func (k *Cli) ExecuteCommand(ctx context.Context, args ...string) error {
	return k.Command(ctx, args...).Run().Cause()
}

func (k *Cli) ApplyCmd(ctx context.Context, content []byte, extraArgs ...string) cmdutils.Cmd {
	args := append([]string{"apply", "-f", "-"}, extraArgs...)

	cmd := k.Command(ctx, args...)
	cmd.SetStdin(bytes.NewBuffer(content))
	return cmd
}

func (k *Cli) Apply(ctx context.Context, content []byte, extraArgs ...string) error {
	return k.ApplyCmd(ctx, content, extraArgs...).Run().Cause()
}

func (k *Cli) DeleteCmd(ctx context.Context, content []byte, extraArgs ...string) cmdutils.Cmd {
	args := append([]string{"delete", "-f", "-"}, extraArgs...)

	cmd := k.Command(ctx, args...)
	cmd.SetStdin(bytes.NewBuffer(content))
	return cmd
}

func (k *Cli) Delete(ctx context.Context, content []byte, extraArgs ...string) error {
	return k.DeleteCmd(ctx, content, extraArgs...).Run().Cause()
}

func (k *Cli) CopyCmd(ctx context.Context, from, to string) cmdutils.Cmd {
	return k.Command(ctx, "cp", from, to)
}

func (k *Cli) Copy(ctx context.Context, from, to string) error {
	return k.CopyCmd(ctx, from, to).Run().Cause()
}

func (k *Cli) StartPortForward(ctx context.Context, options ...portforward.Option) (portforward.PortForwarder, error) {
	options = append([]portforward.Option{
		// We define some default values, which users can then override
		portforward.WithWriters(k.receiver, k.receiver),
		portforward.WithKubeContext(k.kubeContext),
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
