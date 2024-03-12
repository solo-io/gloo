package helm

import (
	"context"
	"fmt"
	"os"

	"github.com/solo-io/gloo/test/setup/helpers"
	"github.com/solo-io/gloo/test/setup/kubernetes"
	"github.com/solo-io/gloo/test/setup/types"
	"go.uber.org/zap"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	k8s "k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Opts struct {
	Cluster    string
	Client     k8s.Interface
	Controller client.Client
	Logger     *zap.SugaredLogger
}

func (o Opts) InstallChart(ctx context.Context, chart *types.Chart, cluster *kubernetes.Cluster) error {
	if chart == nil {
		return nil
	}

	name := chart.Name

	flags := &genericclioptions.ConfigFlags{
		KubeConfig: strPtr(cluster.GetKubeConfig()),
		Context:    strPtr(cluster.GetKubeContext()),
		Namespace:  strPtr(chart.Namespace),
	}

	settings := cli.New()
	kubeClient := kube.New(flags)

	client, err := kubeClient.Factory.KubernetesClientSet()
	if err != nil {
		return err
	}

	registry, err := registry.NewClient()
	if err != nil {
		return err
	}

	timerFn := helpers.TimerFunc((fmt.Sprintf("[%s] %s helm installation", cluster.GetName(), name)))
	defer timerFn()

	storage := storage.Init(driver.NewSecrets(client.CoreV1().Secrets(chart.Namespace)))

	knownReleases := map[string]bool{}

	list, err := storage.ListReleases()
	if err != nil {
		return err
	}

	for _, release := range list {
		knownReleases[release.Name] = true
	}

	cfg := &action.Configuration{
		RegistryClient: registry,
		Releases:       storage,
		KubeClient:     kubeClient,
		Capabilities:   chartutil.DefaultCapabilities,
	}

	logFn := func(format string, args ...any) {}
	if err = cfg.Init(flags, chart.Namespace, os.Getenv("HELM_DRIVER"), logFn); err != nil {
		return err
	}

	if _, ok := knownReleases[name]; !ok {
		if err = installAction(ctx, name, chart, cfg, settings); err != nil {
			return err
		}
		return nil
	}

	if err = upgradeAction(ctx, name, chart, cfg, settings); err != nil {
		return err
	}

	return nil
}

// Pulls a chart from a repository and returns the path to the chart.
func pullRepository(repoName, repoURL, version string, settings *cli.EnvSettings) (string, error) {
	puller := action.NewPull()

	puller.RepoURL = repoURL
	puller.Version = version

	path, err := puller.LocateChart(repoName, settings)
	if err != nil {
		return "", err
	}
	return path, nil
}

func upgradeAction(
	ctx context.Context,
	name string,
	c *types.Chart,
	cfg *action.Configuration,
	settings *cli.EnvSettings,
) error {
	action := action.NewUpgrade(cfg)

	action.Namespace = c.Namespace

	var (
		err       error
		chartPath = c.Local
		chart     *chart.Chart
	)

	if chartPath == "" {
		chartPath, err = pullRepository(name, c.Remote, c.Version, settings)
		if err != nil {
			return err
		}
	}

	chart, err = loader.Load(chartPath)
	if err != nil {
		return err
	}

	_, err = action.RunWithContext(ctx, name, chart, c.Values)
	return err
}

func installAction(
	ctx context.Context,
	name string,
	c *types.Chart,
	cfg *action.Configuration,
	settings *cli.EnvSettings,
) error {
	action := action.NewInstall(cfg)

	action.CreateNamespace = true
	action.Namespace = c.Namespace
	action.ReleaseName = name

	var (
		err       error
		chartPath = c.Local
		chart     *chart.Chart
	)

	if chartPath == "" {
		chartPath, err = pullRepository(name, c.Remote, c.Version, settings)
		if err != nil {
			return err
		}
	}

	chart, err = loader.Load(chartPath)
	if err != nil {
		return err
	}

	_, err = action.RunWithContext(ctx, chart, c.Values)
	return err
}

func strPtr(str string) *string {
	return &str
}
