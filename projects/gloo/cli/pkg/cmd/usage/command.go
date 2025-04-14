package usage

import (
	"context"
	"fmt"
	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gatewayapi/convert"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

type Options struct {
	*options.Options
	ControlPlaneName      string
	ControlPlaneNamespace string
	GlooSnapshotFile      string
}

func RootCmd(op *options.Options) *cobra.Command {
	opts := &Options{
		Options: op,
	}
	cmd := &cobra.Command{
		Use:     "usage",
		Short:   "Scan Gloo for feature usage",
		Long:    "",
		Example: ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.validate(); err != nil {
				return err
			}

			return run(opts)
		},
	}
	opts.addToFlags(cmd.PersistentFlags())
	cmd.SilenceUsage = true
	return cmd
}

func (opts *Options) validate() error {

	return nil
}

func run(opts *Options) error {
	tempDir, err := os.MkdirTemp("", "tmp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir) // Clean up the directory when done
	// Create a temporary directory
	var filePath string

	if opts.GlooSnapshotFile == "" {

		filePath, err = LoadSnapshotFromGloo(opts, tempDir)
		if err != nil {
			return err
		}
	} else {
		filePath = opts.GlooSnapshotFile
	}

	// scan for gloo gateways

	output := convert.NewGatewayAPIOutput()

	inputSnapshot, err := snapshot.FromGlooSnapshot(filePath)
	if err != nil {
		return err
	}
	output.EdgeCache(inputSnapshot)

	// go through the edge snapshot
	usage, err := generateUsage(output)
	if err != nil {
		return err
	}
	usage.Print()

	return nil

}

func findGlooProxyPods() {

}

func LoadSnapshotFromGloo(opts *Options, tempDir string) (string, error) {

	cli, shutdownFunc, err := NewPortForwardedClient(context.Background(), kubectl.NewCli().WithKubeContext(opts.Top.KubeContext), opts.ControlPlaneName, opts.ControlPlaneNamespace)
	if err != nil {
		return "", err
	}
	defer shutdownFunc()
	filePath := filepath.Join(tempDir, "gg-input.json")
	inputSnapshotFile := fileAtPath(filePath)
	err = cli.RequestPathCmd(context.Background(), "/snapshots/input").WithStdout(inputSnapshotFile).Run().Cause()

	return filePath, nil
}

func (o *Options) addToFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.ControlPlaneName, "gloo-control-plane", "deploy/gloo", "Name of the Gloo control plane pod")
	flags.StringVarP(&o.ControlPlaneNamespace, "gloo-control-plane-namespace", "n", "gloo-system", "Namespace of the Gloo control plane pod")
	flags.StringVar(&o.GlooSnapshotFile, "input-snapshot", "", "Gloo input snapshot file location")
}

func NewPortForwardedClient(ctx context.Context, kubectlCli *kubectl.Cli, proxySelector, namespace string) (*admincli.Client, func(), error) {
	selector := portforward.WithResourceSelector(proxySelector, namespace)

	// 1. Open a port-forward to the Kubernetes Deployment, so that we can query the Envoy Admin API directly
	portForwarder, err := kubectlCli.StartPortForward(ctx,
		selector,
		portforward.WithRemotePort(9095))
	if err != nil {
		return nil, nil, err
	}

	// 2. Close the port-forward when we're done accessing data
	deferFunc := func() {
		portForwarder.Close()
		portForwarder.WaitForStop()
	}

	// 3. Create a CLI that connects to the Envoy Admin API
	adminCli := admincli.NewClient().
		WithCurlOptions(
			curl.WithHostPort(portForwarder.Address()),
		)

	return adminCli, deferFunc, err
}

func fileAtPath(path string) *os.File {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		fmt.Printf("unable to openfile: %f\n", err)
	}
	return f
}
