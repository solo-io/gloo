package export

import (
	"context"
	"os"
	"path/filepath"

	"github.com/solo-io/gloo/pkg/utils/cmdutils"
	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/errutils"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CollectEnvoyData queries the Admin API for the provided Envoy Deployment, and emits data to the EnvoyDataDir
func CollectEnvoyData(ctx context.Context, envoyDeployment metav1.ObjectMeta, envoyDataDir string) error {
	// 1. Open a port-forward to the Kubernetes Deployment, so that we can query the Envoy Admin API directly
	portForwarder, err := kubectl.NewCli().WithReceiver(os.Stdout).StartPortForward(ctx,
		portforward.WithDeployment(envoyDeployment.GetName(), envoyDeployment.GetNamespace()),
		portforward.WithRemotePort(int(defaults.EnvoyAdminPort)))
	if err != nil {
		return err
	}

	// 2. Close the port-forward when we're done accessing data
	defer func() {
		portForwarder.Close()
		portForwarder.WaitForStop()
	}()

	// 3. Create a CLI that connects to the Envoy Admin API
	adminCli := admincli.NewClient().
		WithReceiver(os.Stdout).
		WithCurlOptions(
			curl.WithHostPort(portForwarder.Address()),
		)

	// 4. Execute parallel requests, emitting output to defined files
	return errutils.AggregateConcurrent([]func() error{
		cmdutils.RunCommandOutputToFileFunc(
			adminCli.ConfigDumpCmd(ctx),
			filepath.Join(envoyDataDir, "config_dump.log"),
		),
		cmdutils.RunCommandOutputToFileFunc(
			adminCli.StatsCmd(ctx),
			filepath.Join(envoyDataDir, "stats.log"),
		),
	})
}
