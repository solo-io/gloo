package cmd

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/test/setup/defaults"
	"github.com/solo-io/gloo/test/setup/types"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"

	"github.com/solo-io/gloo/test/setup/config"
	"github.com/solo-io/gloo/test/setup/helm"
	"github.com/solo-io/gloo/test/setup/helpers"
	"github.com/solo-io/gloo/test/setup/istio"
	"github.com/solo-io/gloo/test/setup/kind"
	"github.com/solo-io/gloo/test/setup/kubernetes"
)

type options struct {
	configPath  string
	printConfig bool

	stdin io.Reader
}

func (o *options) AddToFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&o.configPath, "config", "c", "<empty>", "Path to the configuration file (Use '-' for stdin).")
	flags.BoolVar(&o.printConfig, "print", false, "Print the configuration file.")
}

func setupCommand(ctx context.Context) *cobra.Command {
	var (
		opts = &options{}
		cmd  = &cobra.Command{
			Use:          "setup",
			Short:        "Setup test environment",
			SilenceUsage: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				opts.stdin = cmd.InOrStdin()

				return setupFunc(ctx, opts)
			},
		}
	)
	opts.AddToFlags(cmd.PersistentFlags())

	return cmd
}

func setupFunc(ctx context.Context, opts *options) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	config, err := config.Load(ctx, opts.configPath, opts.stdin)
	if err != nil {
		return errors.Wrap(err, "failed to load configuration")
	}

	if opts.printConfig {
		return config.Print()
	}

	var istioctlBinary string
	if os.Getenv(defaults.IstioctlVersionEnv) != "" {
		istioctlBinary, err = istio.DownloadIstio(ctx, os.Getenv(defaults.IstioctlVersionEnv))
		if err != nil {
			return fmt.Errorf("failed to download istio: %w", err)
		}
		contextutils.LoggerFrom(ctx).Infof("Using Istio binary '%s'", istioctlBinary)
	}

	clusters := kubernetes.NewClusterMapping()

	kubeconfigDir := filepath.Join(os.Getenv("HOME"), ".kube", "kind")

	if _, err := os.Stat(kubeconfigDir); errors.Is(err, fs.ErrNotExist) {
		if err := os.MkdirAll(kubeconfigDir, 0700); err != nil {
			return errors.Wrap(err, "failed to create directory for kubeconfig files")
		}
	}

	for err := range runParallel(ctx, initClusters(ctx, kubeconfigDir, clusters, config.Clusters...)...) {
		if err != nil {
			return errors.Wrap(err, "failed to init clusters")
		}
	}

	funcs, err := configureClusters(ctx, istioctlBinary, clusters, config)
	if err != nil {
		return errors.Wrap(err, "failed to create cluster processing functions")
	}

	for err := range runParallel(ctx, funcs...) {
		if err != nil {
			return errors.Wrap(err, "failed to process clusters")
		}
	}

	helpers.ListTimers()

	return nil
}

// RunParallel will run the provided functions in parallel and return a channel of errors.
func runParallel(ctx context.Context, funcs ...func() error) <-chan error {
	var (
		errs = make(chan error)
		sync = &sync.WaitGroup{}
	)

	for _, f := range funcs {
		f := f

		sync.Add(1)

		go func(fn func() error) {
			defer sync.Done()

			errs <- fn()
		}(f)
	}

	go func() {
		defer close(errs)

		sync.Wait()
	}()

	return errs
}

// InitClusters will create objects that represent the clusters in the environment. Will create clusters if kind configuration is provided.
func initClusters(
	ctx context.Context,
	dir string,
	clusters *kubernetes.ClusterMapping,
	set ...*types.Cluster,
) (out []func() error) {
	for _, cluster := range set {
		cluster := cluster

		out = append(out, func() error {
			var (
				err        error
				kubeConfig string
			)

			if cluster.KindConfig != nil {
				if err := kind.Get(cluster.KindConfig); err != nil {
					if err = kind.Create(cluster.KindConfig); err != nil {
						return errors.Wrap(err, "failed to create kind cluster")
					}
				}

				for _, image := range cluster.GetImages() {
					if err = kind.LoadImage(image, cluster.Name); err != nil {
						return errors.Wrap(err, "failed to load image")
					}
				}
			}

			// TODO: Support loading images built from source to remote registries (e.g: AWS, GCP, etc ...)

			kubeContext, kubeConfig, client, controller, err := kubernetes.NewClient(ctx, dir, cluster.Name)
			if err != nil {
				return errors.Wrapf(err, "using %s failed to create kubernetes client", os.Getenv("KUBECONFIG"))
			}

			err = clusters.UpdateClusterMapping(ctx, cluster.Name, kubeContext, kubeConfig, client, controller)
			if err != nil {
				return errors.Wrap(err, "failed to update cluster mapping")
			}

			return nil
		})
	}

	return out
}

// configureClusters will process the clusters in the configuration -- installing the gloo control plane along
// with any configured helm charts for that cluster.
func configureClusters(
	ctx context.Context,
	istioctlBinary string,
	clusters *kubernetes.ClusterMapping,
	config *types.Config,
) (out []func() error, err error) {
	clusterMap := map[string]*kubernetes.Cluster{}
	for _, cluster := range config.GetClusters() {
		clusterMap[cluster.Name] = clusters.GetWorkload(cluster.Name)
	}

	for _, cluster := range config.GetClusters() {
		cluster := cluster

		out = append(out,
			func() error {
				info := clusters.GetWorkload(cluster.Name)

				// install istio and charts on cluster
				if err := processCluster(ctx, istioctlBinary, cluster, info); err != nil {
					return errors.Wrap(err, "failed to process cluster")
				}

				// install helper test applications on the cluster
				if err = processApps(ctx, cluster, info); err != nil {
					return errors.Wrap(err, "failed to process apps")
				}

				return nil
			})
	}

	return out, nil
}

// processCluster will process the cluster in the configuration async.
func processCluster(
	ctx context.Context,
	istioctlBinary string,
	clusterInfo *types.Cluster,
	cluster *kubernetes.Cluster,
) (err error) {
	if cluster == nil || clusterInfo == nil {
		return nil
	}

	var (
		name       = cluster.GetName()
		controller = cluster.GetController()
		client     = cluster.GetKubernetes()
	)

	installer := helm.Opts{
		Cluster:    name,
		Client:     client,
		Controller: controller,
		Logger:     contextutils.LoggerFrom(ctx),
	}

	if err = processCRDs(ctx, cluster); err != nil {
		return errors.Wrap(err, "failed to process crds")
	}

	// Prioritized charts
	if err = processCharts(ctx, installer, cluster, clusterInfo.GetPrioritizedCharts()...); err != nil {
		return errors.Wrap(err, "failed to process charts")
	}

	if err := processComponents(ctx, istioctlBinary, clusterInfo, cluster); err != nil {
		return err
	}

	// Unprioritized charts
	if err = processCharts(ctx, installer, cluster, clusterInfo.GetUnprioritizedCharts()...); err != nil {
		return errors.Wrap(err, "failed to process charts")
	}

	return err
}

func processCRDs(ctx context.Context, cluster *kubernetes.Cluster) error {
	// TODO: make this conditional
	// Install gateway apis

	cmd := exec.Command("kubectl", "apply", "--context", cluster.GetKubeContext(), "--kubeconfig", cluster.GetKubeConfig(), "-f", "https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.0.0/standard-install.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("istioctl install failed: %w", err)
	}

	return nil
}

// processComponents will process the cluster istio and app portion of the configuration.
func processComponents(
	ctx context.Context,
	istioctlBinary string,
	cluster *types.Cluster,
	info *kubernetes.Cluster,
) error {
	var err error
	for _, ns := range cluster.Namespaces {
		ns := ns

		if err = kubernetes.CreateOrUpdate[*corev1.Namespace](ctx, ns, info.GetController()); err != nil {
			return errors.Wrap(err, "failed to create namespace")
		}
	}

	for _, operator := range cluster.IstioOperators {
		operator := operator

		if err = istio.Install(ctx, istioctlBinary, cluster.Name, operator, info); err != nil {
			return errors.Wrap(err, "failed to install istio")
		}
	}

	return err
}

func processApps(ctx context.Context, cluster *types.Cluster, info *kubernetes.Cluster) error {
	var err error
	for _, app := range cluster.Apps {
		app := app

		if err = kubernetes.DeployApplication(ctx, cluster.Name, app.Versions, app, info); err != nil {
			contextutils.LoggerFrom(ctx).Errorf("failed to deploy test apps: %v", err)

			return errors.Wrap(err, "failed to deploy test apps")
		}
	}
	defer func() {
		for _, ns := range cluster.Namespaces {
			ns := ns

			if err = kubernetes.RolloutStatus(ns.Name, info); err != nil {
				return
			}
		}
	}()
	return err
}

func processCharts(
	ctx context.Context,
	installer helm.Opts,
	cluster *kubernetes.Cluster,
	charts ...*types.Chart,
) error {
	for _, chart := range charts {
		contextutils.LoggerFrom(ctx).Infof("installing chart %s", chart.Name)

		if err := installer.InstallChart(ctx, chart, cluster); err != nil {
			return errors.Wrap(err, "failed to install chart")
		}

		if err := kubernetes.RolloutStatus(chart.Namespace, cluster); err != nil {
			return errors.Wrap(err, "failed to wait for rollout")
		}
	}
	return nil
}
