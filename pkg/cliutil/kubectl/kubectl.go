package kubectl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/threadsafe"
)

// Kubectl is a utility for executing `kubectl` commands
type Kubectl struct {
	// receiver is destination for the kubectl stdout and stderr
	receiver io.Writer

	// kubeContext is the optional value of the context for a given kubernetes cluster
	// If it is not supplied, no context will be included in the command
	kubeContext string
}

// NewKubectl returns NewKubectlWithKubeContext with an empty Kubernetes context
func NewKubectl(receiver io.Writer) (*Kubectl, error) {
	return NewKubectlWithKubeContext(receiver, "")
}

// NewKubectlWithKubeContext returns an implementation of Kubectl, or an error if one of the provided receivers was nil
func NewKubectlWithKubeContext(receiver io.Writer, kubeContext string) (*Kubectl, error) {
	if receiver == nil {
		return nil, errors.New("receiver must not be nil")
	}

	return &Kubectl{
		receiver:    receiver,
		kubeContext: kubeContext,
	}, nil
}

func (k *Kubectl) GetReceiver() io.Writer {
	return k.receiver
}

func (k *Kubectl) Execute(ctx context.Context, in io.Reader, args ...string) (stdOut string, stdErr string, executeErr error) {
	if k.kubeContext != "" {
		args = append([]string{"--context", k.kubeContext}, args...)
	}

	cmd := createKubectlCommand(ctx, args...)
	if in != nil {
		cmd.Stdin = in
	}

	var stdout, stderr threadsafe.Buffer
	cmd.Stdout = io.MultiWriter(&stdout, k.receiver)
	cmd.Stderr = io.MultiWriter(&stderr, k.receiver)

	_, _ = fmt.Fprintf(k.receiver, "Executing: %s \n", strings.Join(cmd.Args, " "))
	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func (k *Kubectl) ExecuteSafe(ctx context.Context, in io.Reader, args ...string) (stdOut string, executeErr error) {
	stdout, stderr, err := k.Execute(ctx, in, args...)
	if err != nil {
		return stdout, eris.Wrapf(err, "failed to execute command: %s", stderr)
	}
	if stderr != "" {
		return stdout, eris.Errorf("failed to execute command: %s", stderr)
	}

	return stdout, nil
}

func (k *Kubectl) Apply(ctx context.Context, content []byte, extraArgs ...string) error {
	args := append([]string{"apply", "-f", "-"}, extraArgs...)
	buf := bytes.NewBuffer(content)

	_, err := k.ExecuteSafe(ctx, buf, args...)
	return err
}

func (k *Kubectl) Delete(ctx context.Context, content []byte, extraArgs ...string) error {
	args := append([]string{"delete", "-f", "-"}, extraArgs...)
	buf := bytes.NewBuffer(content)

	_, err := k.ExecuteSafe(ctx, buf, args...)
	return err
}

func createKubectlCommand(ctx context.Context, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "kubectl", args...)
	cmd.Env = os.Environ()
	// disable DEBUG=1 from getting through to kube
	for i, pair := range cmd.Env {
		if strings.HasPrefix(pair, "DEBUG") {
			cmd.Env = append(cmd.Env[:i], cmd.Env[i+1:]...)
			break
		}
	}
	return cmd
}
