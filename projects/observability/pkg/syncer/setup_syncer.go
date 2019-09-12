package syncer

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana"

	"github.com/solo-io/gloo/pkg/utils/setuputils"
	"github.com/solo-io/gloo/pkg/version"
	check "github.com/solo-io/go-checkpoint"

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
	observability    = "observability"
	GRAFANA_USERNAME = "GRAFANA_USERNAME"
	GRAFANA_PASSWORD = "GRAFANA_PASSWORD"
)

func Main() error {
	check.NewUsageClient().Start(observability, version.Version)
	return setuputils.Main(setuputils.SetupOpts{
		LoggingPrefix: observability,
		ExitOnError:   false,
		SetupFunc:     Setup,
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
	opts.WatchOpts.Ctx = contextutils.WithLogger(opts.WatchOpts.Ctx, "gateway")

	upstreamClient, err := gloov1.NewUpstreamClient(opts.Upstreams)
	if err != nil {
		return err
	}
	if err := upstreamClient.Register(); err != nil {
		return err
	}

	username := os.Getenv(GRAFANA_USERNAME)
	password := os.Getenv(GRAFANA_PASSWORD)
	if username == "" || password == "" {
		contextutils.LoggerFrom(opts.WatchOpts.Ctx).Fatalf("grafana username and password cannot be empty")
	}
	basicAuthString := fmt.Sprintf("%s:%s", username, password)
	apiUrl := fmt.Sprintf("%s:%s", SERVICE_LINK, SERVICE_PORT)

	httpClient := http.DefaultClient
	restClient := grafana.NewRestClient(apiUrl, basicAuthString, httpClient)

	dashboardClient := grafana.NewDashboardClient(restClient)
	snapshotClient := grafana.NewSnapshotClient(restClient)

	dashSyncer := NewGrafanaDashboardSyncer(dashboardClient, snapshotClient)

	emitter := v1.NewDashboardsEmitter(upstreamClient)
	eventLoop := v1.NewDashboardsEventLoop(emitter, dashSyncer)
	writeErrs := make(chan error)
	eventLoopErrs, err := eventLoop.Run(opts.WatchNamespaces, opts.WatchOpts)
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(opts.WatchOpts.Ctx, writeErrs, eventLoopErrs, "event_loop")

	logger := contextutils.LoggerFrom(opts.WatchOpts.Ctx)

	go func() {
		for {
			select {
			case err := <-writeErrs:
				logger.Errorf("error: %v", err)
			case <-opts.WatchOpts.Ctx.Done():
				close(writeErrs)
				return
			}
		}
	}()
	return nil
}
