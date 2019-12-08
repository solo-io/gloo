package syncer

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana"

	"github.com/solo-io/gloo/pkg/utils/setuputils"

	"github.com/gogo/protobuf/types"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/errutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	v1 "github.com/solo-io/solo-projects/projects/observability/pkg/api/v1"
	"k8s.io/client-go/rest"
)

const (
	observability     = "observability"
	grafanaUrl        = "GRAFANA_URL"
	grafanaUsername   = "GRAFANA_USERNAME"
	grafanaPassword   = "GRAFANA_PASSWORD"
	grafanaApiKey     = "GRAFANA_API_KEY"
	dashboardTemplate = "/observability/dashboard-template.json"
)

func Main() error {
	return setuputils.Main(setuputils.SetupOpts{
		LoggerName:  observability,
		ExitOnError: true,
		SetupFunc:   Setup,
	})
}

func Setup(ctx context.Context, kubeCache kube.SharedCache, inMemoryCache memory.InMemoryResourceCache, settings *gloov1.Settings) error {
	var (
		cfg *rest.Config
	)
	params := bootstrap.NewConfigFactoryParams(settings, inMemoryCache, kubeCache, &cfg, nil)
	upstreamFactory, err := bootstrap.ConfigFactoryForSettings(params, gloov1.UpstreamCrd)
	if err != nil {
		return err
	}

	refreshRate, err := types.DurationFromProto(settings.RefreshRate)
	if err != nil {
		return err
	}

	writeNamespace := settings.DiscoveryNamespace
	if writeNamespace == "" {
		writeNamespace = gloodefaults.GlooSystem
	}
	watchNamespaces := settings.WatchNamespaces
	var writeNamespaceProvided bool
	for _, ns := range watchNamespaces {
		if ns == writeNamespace {
			writeNamespaceProvided = true
			break
		}
	}
	if !writeNamespaceProvided {
		watchNamespaces = append(watchNamespaces, writeNamespace)
	}
	opts := Opts{
		WriteNamespace:  writeNamespace,
		WatchNamespaces: watchNamespaces,
		Upstreams:       upstreamFactory,
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: refreshRate,
		},
		DevMode: true,
	}

	return RunObservability(opts)
}

func RunObservability(opts Opts) error {
	opts.WatchOpts = opts.WatchOpts.WithDefaults()
	opts.WatchOpts.Ctx = contextutils.WithLogger(opts.WatchOpts.Ctx, "observability")

	upstreamClient, err := gloov1.NewUpstreamClient(opts.Upstreams)
	if err != nil {
		return err
	}
	if err := upstreamClient.Register(); err != nil {
		return err
	}

	creds, err := buildRestCredentials(opts.WatchOpts.Ctx)
	if err != nil {
		return err
	}

	grafanaApiUrl, err := getGrafanaApiUrl()
	if err != nil {
		return err
	}

	httpClient := http.DefaultClient
	restClient := grafana.NewRestClient(grafanaApiUrl, httpClient, creds)

	dashboardClient := grafana.NewDashboardClient(restClient)
	snapshotClient := grafana.NewSnapshotClient(restClient)

	dashboardJsonTemplate, err := getDashboardJsonTemplate(opts.WatchOpts.Ctx)
	if err != nil {
		return err
	}

	dashSyncer := NewGrafanaDashboardSyncer(dashboardClient, snapshotClient, dashboardJsonTemplate)

	emitter := v1.NewDashboardsEmitter(upstreamClient)
	eventLoop := v1.NewDashboardsEventLoop(emitter, dashSyncer)
	writeErrs := make(chan error)
	eventLoopErrs, err := eventLoop.Run(opts.WatchNamespaces, opts.WatchOpts)
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(opts.WatchOpts.Ctx, writeErrs, eventLoopErrs, "event_loop")

	go func() {
		for {
			select {
			case err := <-writeErrs:
				contextutils.LoggerFrom(opts.WatchOpts.Ctx).Errorf("error: %v", err)
			case <-opts.WatchOpts.Ctx.Done():
				close(writeErrs)
				return
			}
		}
	}()
	return nil
}

func buildRestCredentials(ctx context.Context) (*grafana.Credentials, error) {
	var (
		username = os.Getenv(grafanaUsername)
		password = os.Getenv(grafanaPassword)
		apiKey   = os.Getenv(grafanaApiKey)
		logger   = contextutils.LoggerFrom(ctx)
	)

	if apiKey != "" {
		logger.Info("Using api key for authentication to grafana")
		return &grafana.Credentials{
			BasicAuth: nil,
			ApiKey:    apiKey,
		}, nil
	} else if username != "" && password != "" {
		logger.Info("Using basic auth for authentication to grafana")
		return &grafana.Credentials{
			BasicAuth: &grafana.BasicAuth{
				Username: username,
				Password: password,
			},
			ApiKey: "",
		}, nil
	} else {
		return nil, grafana.IncompleteGrafanaCredentials
	}
}

func getGrafanaApiUrl() (string, error) {
	url := os.Getenv(grafanaUrl)
	if url == "" {
		return "", NoGrafanaUrl(grafanaUrl)
	}

	return url, nil
}

func getDashboardJsonTemplate(ctx context.Context) (string, error) {
	bytes, err := ioutil.ReadFile(dashboardTemplate)
	if err != nil {
		contextutils.LoggerFrom(ctx).Warnf("Error reading file %s - %s", dashboardTemplate, err.Error())
		return "", nil
	}
	return string(bytes), nil
}
