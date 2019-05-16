package syncer

import (
	"context"
	"net/http"
	"time"

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

const observability = "observability"

func Main() error {
	start := time.Now()
	check.CallCheck(observability, version.Version, start)
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
	upstreamFactory, err := bootstrap.ConfigFactoryForSettings(
		settings,
		inMemoryCache,
		kubeCache,
		gloov1.UpstreamCrd,
		&cfg,
	)
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

	emitter := v1.NewDashboardsEmitter(upstreamClient)
	dashSyncer, err := NewGrafanaDashboardSyncer(http.DefaultClient)
	if err != nil {
		return err
	}
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
