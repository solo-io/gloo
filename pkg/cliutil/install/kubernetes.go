package install

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/pkg/cliutil/kubectl"
)

// Deprecated: Prefer kubectl.Kubectl
func KubectlApply(manifest []byte, extraArgs ...string) error {
	return Kubectl(bytes.NewBuffer(manifest), append([]string{"apply", "-f", "-"}, extraArgs...)...)
}

// Deprecated: Prefer kubectl.Kubectl
func KubectlApplyOut(manifest []byte, extraArgs ...string) ([]byte, error) {
	return KubectlOut(bytes.NewBuffer(manifest), append([]string{"apply", "-f", "-"}, extraArgs...)...)
}

// Deprecated: Prefer kubectl.Kubectl
func KubectlDelete(manifest []byte, extraArgs ...string) error {
	return Kubectl(bytes.NewBuffer(manifest), append([]string{"delete", "-f", "-"}, extraArgs...)...)
}

// Deprecated: Prefer kubectl.Kubectl
type KubeCli interface {
	Kubectl(stdin io.Reader, args ...string) error
	KubectlOut(stdin io.Reader, args ...string) ([]byte, error)
}

type CmdKubectl struct{}

var _ KubeCli = &CmdKubectl{}

func (k *CmdKubectl) Kubectl(stdin io.Reader, args ...string) error {
	return Kubectl(stdin, args...)
}

func (k *CmdKubectl) KubectlOut(stdin io.Reader, args ...string) ([]byte, error) {
	return KubectlOut(stdin, args...)
}

var verbose bool

func SetVerbose(b bool) {
	verbose = b
}

// Deprecated: Prefer kubectl.Kubectl
func Kubectl(stdin io.Reader, args ...string) error {
	_, err := KubectlOut(stdin, args...)
	return err
}

// Deprecated: Prefer kubectl.Kubectl
func KubectlOut(stdin io.Reader, args ...string) ([]byte, error) {
	var outWriter, errWriter io.Writer

	if verbose {
		outWriter = os.Stdout
		errWriter = os.Stderr
	} else {
		// use logfile
		cliutil.Initialize()
		outWriter = cliutil.GetLogger()
		errWriter = cliutil.GetLogger()
	}

	cli, err := kubectl.NewKubectl(outWriter, errWriter)
	if err != nil {
		return nil, err
	}
	stdout, err := cli.ExecuteSafe(context.Background(), stdin, args...)

	return []byte(stdout), err
}
