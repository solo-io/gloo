package convert

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
		Long:  "Convert Gloo Edge APIs to Gateway API by either providing kubernetes yaml files or a Gloo Gateway input snapshot",
		Example: `# This command converts Gloo Edge APIs to Gloo Gateway API yaml and places them in the '--output-dir' directory arranged by namespace.
# To generate gateway api by getting snapshot directly from running Gloo pod. The 'output-dir'' must not exist
  glooctl gateway-api convert --gloo-control-plane deploy/gloo --output-dir ./_output

# To generate gateway api by a single kubernetes yaml file. The 'output-dir'' must not exist
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output

# To delete and recreate the output dir add 'delete-output-dir'
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output --delete-output-dir

# To generate gateway api by a single kubernetes yaml file but place all the output configurations in the same file.
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output --retain-input-folder-structure

# To load a bunch of *.yaml or *.yml files in nested directories. You can also retain the same file structure which is helpful in CI/CD pipelines.
  glooctl gateway-api convert --input-dir ./gloo-configs --output-dir ./_output --retain-input-folder-structure

# Configuration can also be pulled an translated directly from a running Gloo Gateway Instance (version 1.17+).
# Port forward to the running gloo instance to grab its snapshot
  kubectl -n gloo-system port-forward deploy/gloo 9091
  curl localhost:9091/snapshots/input > gg-input.json
  
  glooctl gateway-api convert --input-snapshot gg-input.json --output-dir ./_output

# For all commands you can print stats about the migration such as number of configs
  glooctl gateway-api convert --input-file gloo-yamls.yaml --output-dir ./_output --print-stats

# If the yaml files contain non Gloo API yaml they can be retained by adding '--include-unknown'
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

	filesMetrics.Add(float64(len(foundFiles)))

	output := NewGatewayAPIOutput()
	var isSnapshotFile bool
	if opts.GlooSnapshotFile != "" {
		isSnapshotFile = true
	}

	if err := output.Load(foundFiles, isSnapshotFile); err != nil {
		return err
	}

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
		printMetrics(output)
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
