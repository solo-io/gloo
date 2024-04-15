package export

import (
	"context"
	"github.com/solo-io/gloo/pkg/utils/envoyutils/adminctl"
	"github.com/solo-io/gloo/pkg/utils/errutils"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"os"
	"path/filepath"
)

func CollectEnvoyData(ctx context.Context, namespace, envoyDataDir string) error {
	// 1. Open a port-forward to the Kubernetes Deployment, so that we can query the Envoy Admin API directly
	portForwarder, err := kubectl.NewCli(os.Stdout).StartPortForward(ctx,
		portforward.WithDeployment(kubeutils.GatewayProxyDeploymentName, namespace),
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
	adminCli := adminctl.NewCli(os.Stdout, portForwarder.Address())

	// 4. Execute parallel requests, emitting output to defined files
	return errutils.AggregateConcurrent([]func() error{
		RunCommandOutputToFile(
			adminCli.ConfigDumpCmd(ctx),
			filepath.Join(envoyDataDir, "config_dump.log"),
		),
		RunCommandOutputToFile(
			adminCli.StatsCmd(ctx),
			filepath.Join(envoyDataDir, "stats.log"),
		),
	})
}
