package cmd

import (
	"context"
	"io"

	"github.com/solo-io/gloo/pkg/utils/cmdutils"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/k8s-utils/testutils/kube"
)

func RemoteCurlCmd(ctx context.Context, receiver io.Writer, kubeContext string, podExecOpt kubectl.PodExecOptions, curlOpts ...curl.Option) cmdutils.Cmd {
	f := NewRemoteCmdFactory(ctx, RemoteCmderParams{
		Receiver:      receiver,
		KubeContext:   kubeContext,
		Image:         kube.DefaultCurlImage,
		FromContainer: podExecOpt.Container,
		FromNamespace: podExecOpt.Namespace,
		FromPod:       podExecOpt.Name,
	})
	return f.Command(ctx, "curl", curl.BuildArgs(curlOpts...)...)
}
