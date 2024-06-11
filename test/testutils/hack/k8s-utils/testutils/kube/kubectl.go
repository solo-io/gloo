package kube

import (
	"context"
	"fmt"
	"io"

	"github.com/solo-io/gloo/test/testutils/hack/go-utils/testutils/kubectl"

	. "github.com/onsi/gomega"
)

type KubectlParams struct {
	KubectlCmdParams kubectl.Params
	KubeContext      string
	Logger           io.Writer
}

func mustExecute(ctx context.Context, params KubectlParams) string {
	data, err := execute(ctx, params)
	Expect(err).NotTo(HaveOccurred())
	return data
}

func execute(ctx context.Context, params KubectlParams) (string, error) {
	params.KubectlCmdParams.Args = append([]string{"--context", params.KubeContext}, params.KubectlCmdParams.Args...)
	fmt.Fprintf(params.Logger, "Executing: kubectl %v \n", params.KubectlCmdParams.Args)
	p := kubectl.NewParams()
	p.Stdin = params.KubectlCmdParams.Stdin
	p.Stdout = params.KubectlCmdParams.Stdout
	p.Stderr = params.KubectlCmdParams.Stderr
	p.Env = params.KubectlCmdParams.Env
	p.Args = params.KubectlCmdParams.Args
	readerChan, err := kubectl.KubectlOutChan(ctx, p)
	if err != nil {
		return "", err
	}
	select {
	case <-ctx.Done():
		return "", nil
	case reader := <-readerChan:
		data, err := io.ReadAll(reader)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(params.Logger, "<kubectl %v> output: %v\n", params.KubectlCmdParams.Args, string(data))
		return string(data), nil
	}
}
