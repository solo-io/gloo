package kubeutils

import (
	"context"
	"fmt"
	"io"

	"github.com/solo-io/gloo/pkg/cliutil/install"

	"github.com/solo-io/gloo/test/helpers"
)

type InstallRef struct {
	ClusterName string
	Namespace   string
}

func GetClusteredPreFailHandler(ctx context.Context, orchestrator Orchestrator, out io.Writer, installs []InstallRef) func() {
	return func() {
		for _, installRef := range installs {
			_ = orchestrator.SetClusterContext(ctx, installRef.ClusterName)

			kubeCli := &install.CmdKubectl{}
			kubeEvents, _ := kubeCli.KubectlOut(nil, "get", "events", "-n", installRef.Namespace)
			_, _ = fmt.Fprintf(out, string(kubeEvents))

			helpers.KubeDumpOnFail(out, installRef.Namespace)()
		}
	}
}
