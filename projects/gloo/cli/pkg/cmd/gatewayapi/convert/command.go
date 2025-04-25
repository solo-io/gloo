package convert

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"

	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"

	"github.com/spf13/cobra"
)

const (
	RandomSuffix = 4
	RandomSeed   = 1
)

func RootCmd(op *options.Options) *cobra.Command {
	opts := &Options{
		Options: op,
	}
	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert Gloo Edge APIs to Gateway API",
		Long:  "Convert Gloo Edge APIs to Gateway APIs by either providing Kubernetes YAML files or a Gloo Gateway input snapshot.",
		Example: `# This command converts Gloo Edge APIs to Kubernetes Gateway API YAML files and places them in the '--output-dir' directory, grouped by namespace.
# To generate Gateway API YAML files from a Gloo Gateway snapshot that is retrieved from a running 'gloo' pod. The 'output-dir' must not exist.
  glooctl gateway-api convert --gloo-control-plane deploy/gloo --output-dir ./_output

# To generate Gateway API YAML files from a single Kubernetes YAML file. The 'output-dir' must not exist.
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output

# To delete and recreate the content in the 'output-dir', add the 'delete-output-dir'' option.
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output --delete-output-dir

# To generate Gateway API YAML files from a single Kubernetes YAML file, but place all the output configurations in to the same file. 
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output --retain-input-folder-structure

# To load a bunch of '*.yaml' or '*.yml' files in nested directories. You can also use the '--retain-input-folder-structure' option to keep the original file structure, which can be helpful in CI/CD pipelines.
  glooctl gateway-api convert --input-dir ./gloo-configs --output-dir ./_output --retain-input-folder-structure

To download a Gloo Gateway snapshot from a running 'gloo' pod (verison 1.17+) and generate Gateway API YAML files from that snapshot. 
  kubectl -n gloo-system port-forward deploy/gloo 9091
  curl localhost:9091/snapshots/input > gg-input.json
  
  glooctl gateway-api convert --input-snapshot gg-input.json --output-dir ./_output

# To get the stats for each migration, such as the number of configuration files that were generated, add the '--print-stats' option. 
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output --print-stats

# To retain non-Gloo Gateway API YAML files, add  the '--include-unknown' option. 
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output --include-unknown`,
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

func run(opts *Options) error {

	foundFiles, err := findFiles(opts)
	if err != nil {
		return err
	}

	output := NewGatewayAPIOutput()
	var inputSnapshot *snapshot.Instance
	snapshotFile := opts.GlooSnapshotFile

	// the snapshot file comes from control plane
	if opts.ControlPlaneName != "" {
		snapshotFile = foundFiles[0]
	}

	if snapshotFile != "" {
		inputSnapshot, err = snapshot.FromGlooSnapshot(snapshotFile)
		if err != nil {
			return err
		}
	} else {
		// yaml files
		// snapshot file
		inputSnapshot, err = snapshot.FromYamlFiles(foundFiles)
		if err != nil {
			return err
		}
	}

	output.edgeCache = inputSnapshot

	fmt.Printf("Successfully loaded %d files\n", len(foundFiles))
	// preprocessing
	if err := output.PreProcess(true); err != nil {
		return err
	}

	// now we need to convert the easy stuff like route tables
	// preprocessing
	if err := output.Convert(); err != nil {
		return err
	}
	fmt.Printf("Processing complete, entering post processing...\n")

	// now we need to convert the easy stuff like route tables
	// preprocessing
	if err := output.PostProcess(opts); err != nil {
		return err
	}
	fmt.Printf("Post processing complete, writing to files...\n")

	// now we need to convert the easy stuff like route tables
	// preprocessing
	if err := output.Write(opts); err != nil {
		return err
	}

	if opts.Stats {
		printMetrics(output, len(foundFiles))
	}

	return nil
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

func findFiles(opts *Options) ([]string, error) {
	var files []string
	if opts.ControlPlaneName != "" {
		// we need to download the file to the output dir and add it to the files list
		if folderExists(opts.OutputDir) {
			if !opts.DeleteOutputDir {
				return nil, fmt.Errorf("output-dir %s already exists, not writing files", opts.OutputDir)
			}
			if err := os.RemoveAll(opts.OutputDir); err != nil {
				return nil, err
			}
		}
		if err := os.MkdirAll(opts.OutputDir, os.ModePerm); err != nil {
			return nil, err
		}

		cli, shutdownFunc, err := NewPortForwardedClient(context.Background(), kubectl.NewCli().WithKubeContext(opts.Top.KubeContext), opts.ControlPlaneName, opts.ControlPlaneNamespace)
		if err != nil {
			return nil, err
		}

		defer shutdownFunc()
		filePath := filepath.Join(opts.OutputDir, "gg-input.json")
		inputSnapshotFile := fileAtPath(filePath)
		err = cli.RequestPathCmd(context.Background(), "/snapshots/input").WithStdout(inputSnapshotFile).Run().Cause()
		if err != nil {
			return nil, fmt.Errorf("error getting gloo snapshot: %f\n", err)
		}
		files = append(files, filePath)

	} else if opts.InputDir != "" {
		fs, err := findYamlFiles(opts.InputDir)
		if err != nil {
			return nil, err
		}
		files = fs
	} else if opts.InputFile != "" {
		files = append(files, opts.InputFile)
	} else if opts.GlooSnapshotFile != "" {
		files = append(files, opts.GlooSnapshotFile)
	}

	return files, nil
}

func findYamlFiles(directory string) ([]string, error) {
	var files []string
	libRegEx, e := regexp.Compile("^.+\\.(yaml|yml)$")
	if e != nil {
		return nil, e
	}

	e = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			if !strings.Contains(info.Name(), "kustomization") {
				files = append(files, path)
			}
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	return files, nil
}
