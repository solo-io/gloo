package upstreams

import (
	"context"
	"fmt"
	"time"

	"github.com/solo-io/go-utils/contextutils"

	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams/consul"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams/kubernetes"
	"golang.org/x/sync/errgroup"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/errutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	skkube "github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
)

const (
	sourceGloo           = "gloo"
	sourceKube           = "kube"
	sourceConsul         = "consul"
	notImplementedErrMsg = "this operation is not supported by this client"
)

var (
	// used in tests
	TimerOverride <-chan time.Time
)

func NewHybridUpstreamClient(
	upstreamClient v1.UpstreamClient,
	serviceClient skkube.ServiceClient,
	consulClient consul.ConsulWatcher,
	settings *v1.Settings) (v1.UpstreamClient, error) {

	clientMap := make(map[string]v1.UpstreamClient)

	if upstreamClient == nil {
		return nil, eris.New("required upstream client is nil")
	}
	clientMap[sourceGloo] = upstreamClient

	if serviceClient != nil {
		clientMap[sourceKube] = kubernetes.NewKubernetesUpstreamClient(serviceClient)
	}

	if consulClient != nil {
		clientMap[sourceConsul] = consul.NewConsulUpstreamClient(consulClient, settings.GetConsulDiscovery())
	}

	return &hybridUpstreamClient{
		clientMap: clientMap,
	}, nil
}

type hybridUpstreamClient struct {
	clientMap map[string]v1.UpstreamClient
}

func (c *hybridUpstreamClient) BaseClient() clients.ResourceClient {
	// We need this modified base client to build reporters, which require generic clients.ResourceClient instances
	return newHybridBaseClient(c.clientMap[sourceGloo].BaseClient())
}

func (c *hybridUpstreamClient) Register() error {
	var err *multierror.Error
	for _, client := range c.clientMap {
		err = multierror.Append(err, client.Register())
	}
	return err.ErrorOrNil()
}

func (c *hybridUpstreamClient) Read(namespace, name string, opts clients.ReadOpts) (*v1.Upstream, error) {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return nil, fmt.Errorf(notImplementedErrMsg)
}

func (c *hybridUpstreamClient) Write(resource *v1.Upstream, opts clients.WriteOpts) (*v1.Upstream, error) {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return nil, fmt.Errorf(notImplementedErrMsg)
}

func (c *hybridUpstreamClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return fmt.Errorf(notImplementedErrMsg)
}

func (rc *hybridUpstreamClient) ApplyStatus(statusClient resources.StatusClient, inputResource resources.InputResource, opts clients.ApplyStatusOpts) (*v1.Upstream, error) {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return nil, fmt.Errorf(notImplementedErrMsg)
}

func (c *hybridUpstreamClient) List(namespace string, opts clients.ListOpts) (v1.UpstreamList, error) {
	var (
		result v1.UpstreamList
		errs   *multierror.Error
	)

	for _, client := range c.clientMap {
		upstreams, err := client.List(namespace, opts)
		errs = multierror.Append(errs, err)

		result = append(result, upstreams...)
	}

	return result, errs.ErrorOrNil()
}

type upstreamsWithSource struct {
	source    string
	upstreams v1.UpstreamList
}

func (c *hybridUpstreamClient) Watch(namespace string, opts clients.WatchOpts) (<-chan v1.UpstreamList, <-chan error, error) {
	opts = opts.WithDefaults()
	ctx := contextutils.WithLogger(opts.Ctx, "hybrid upstream client")
	logger := contextutils.LoggerFrom(ctx)

	logger.Infow("Starting hybrid upstream client watch",
		"issue", "8539",
		"namespace", namespace,
		"clientSources", func() []string {
			var sources []string
			for source := range c.clientMap {
				sources = append(sources, source)
			}
			return sources
		}(),
		"refreshRate", opts.RefreshRate)

	var (
		eg                   = errgroup.Group{}
		collectErrsChan      = make(chan error)
		collectUpstreamsChan = make(chan *upstreamsWithSource)
	)

	// first thing, do a list of everything to get the current state
	current := &hybridUpstreamSnapshot{upstreamsBySource: map[string]v1.UpstreamList{}}
	for source, client := range c.clientMap {
		logger.Infow("Getting initial upstream list from source",
			"issue", "8539",
			"source", source)

		upstreams, err := client.List(namespace, clients.ListOpts{Ctx: opts.Ctx, Selector: opts.Selector})
		if err != nil {
			logger.Errorw("Failed to get initial upstream list from source",
				"issue", "8539",
				"source", source,
				"error", err.Error())
			return nil, nil, err
		}

		logger.Infow("Got initial upstream list from source",
			"issue", "8539",
			"source", source,
			"upstreamCount", len(upstreams))

		current.setUpstreams(source, upstreams)
	}

	logger.Infow("Starting watch loops for all sources",
		"issue", "8539",
		"sourceCount", len(c.clientMap))

	for source, client := range c.clientMap {
		upstreamsFromSourceChan, errsFromSourceChan, err := client.Watch(namespace, opts)
		if err != nil {
			logger.Errorw("Failed to start watch for source",
				"issue", "8539",
				"source", source,
				"error", err.Error())
			return nil, nil, err
		}

		// Copy before passing to goroutines
		sourceName := source

		logger.Infow("Started watch for source",
			"issue", "8539",
			"source", sourceName)

		eg.Go(func() error {
			logger.Infow("Starting error aggregation goroutine",
				"issue", "8539",
				"source", sourceName)
			errutils.AggregateErrs(ctx, collectErrsChan, errsFromSourceChan, sourceName)
			return nil
		})

		eg.Go(func() error {
			logger.Infow("Starting upstream aggregation goroutine",
				"issue", "8539",
				"source", sourceName)
			aggregateUpstreams(ctx, collectUpstreamsChan, upstreamsFromSourceChan, sourceName)
			return nil
		})
	}

	upstreamsOut := make(chan v1.UpstreamList, 1)

	logger.Infow("Created output channel",
		"issue", "8539",
		"channelBufferSize", 1)

	go func() {
		var previousHash uint64
		syncCount := 0
		channelFullCount := 0

		logger.Infow("Starting sync goroutine",
			"issue", "8539")

		// return success for the sync (ie if there still needs changes or there is a hash error it's a false)
		syncFunc := func() bool {
			currentHash, err := current.hash()
			if currentHash == previousHash && err == nil {
				logger.Infow("Hash unchanged, skipping sync",
					"issue", "8539",
					"hash", currentHash)
				return true
			}

			toSend := current.clone()
			upstreamList := toSend.toList()

			logger.Infow("Syncing upstream list",
				"issue", "8539",
				"previousHash", previousHash,
				"currentHash", currentHash,
				"upstreamCount", len(upstreamList),
				"syncCount", syncCount)

			// empty the channel if not empty, as we only care about the latest
			select {
			case oldList := <-upstreamsOut:
				logger.Infow("Drained old upstream list from channel",
					"issue", "8539",
					"oldListSize", len(oldList),
					"newListSize", len(upstreamList))
			default:
				logger.Infow("Output channel was empty",
					"issue", "8539")
			}

			select {
			case upstreamsOut <- upstreamList:
				logger.Infow("Successfully sent upstream list to channel",
					"issue", "8539",
					"upstreamCount", len(upstreamList),
					"hash", currentHash,
					"syncCount", syncCount)
				previousHash = currentHash
				syncCount++
			default:
				logger.Warnw("Failed to send upstream list - channel is full",
					"issue", "8539",
					"upstreamCount", len(upstreamList),
					"hash", currentHash,
					"syncCount", syncCount,
					"channelFullCount", channelFullCount)
				channelFullCount++
				contextutils.LoggerFrom(ctx).DPanic("sending to a buffered channel blocked")
				return false
			}
			return true
		}

		// First time - sync the current state
		logger.Infow("Performing initial sync",
			"issue", "8539")
		needsSync := syncFunc()

		logger.Infow("Initial sync complete",
			"issue", "8539",
			"needsSync", needsSync)

		timerC := TimerOverride
		if timerC == nil {
			timer := time.NewTicker(time.Second * 1)
			timerC = timer.C
			defer timer.Stop()
			logger.Infow("Started retry timer",
				"issue", "8539",
				"interval", "1s")
		}

		for {
			select {
			case <-ctx.Done():
				logger.Infow("Context cancelled, shutting down sync goroutine",
					"issue", "8539",
					"finalSyncCount", syncCount,
					"finalChannelFullCount", channelFullCount)
				close(upstreamsOut)
				_ = eg.Wait() // will never return an error
				close(collectUpstreamsChan)
				close(collectErrsChan)
				return
			case upstreamWithSource, ok := <-collectUpstreamsChan:
				if ok {
					logger.Infow("Received upstream update from source",
						"issue", "8539",
						"source", upstreamWithSource.source,
						"upstreamCount", len(upstreamWithSource.upstreams))
					needsSync = true
					current.setUpstreams(upstreamWithSource.source, upstreamWithSource.upstreams)
				} else {
					logger.Warnw("Upstream collection channel closed",
						"issue", "8539")
				}
			case <-timerC:
				if len(upstreamsOut) != 0 {
					logger.Infow("failed to push hybrid upstream list to "+
						"channel (must be full), retrying in 1s",
						"issue", "8539",
						"list hash", previousHash,
						"channelLength", len(upstreamsOut),
						"channelFullCount", channelFullCount)
				}
				if needsSync {
					logger.Infow("Timer triggered sync attempt",
						"issue", "8539",
						"syncCount", syncCount)
					needsSync = !syncFunc()
				}
			}
		}
	}()

	logger.Infow("Hybrid upstream client watch setup complete",
		"issue", "8539")

	return upstreamsOut, collectErrsChan, nil
}

// Redirects src to dest adding source information
func aggregateUpstreams(ctx context.Context, dest chan *upstreamsWithSource, src <-chan v1.UpstreamList, sourceName string) {
	logger := contextutils.LoggerFrom(ctx)
	logger.Infow("Starting upstream aggregation",
		"issue", "8539",
		"source", sourceName)

	upstreamUpdateCount := 0

	for {
		select {
		case upstreams, ok := <-src:
			if !ok {
				logger.Infow("Source channel closed, stopping aggregation",
					"issue", "8539",
					"source", sourceName,
					"totalUpdates", upstreamUpdateCount)
				return
			}

			logger.Infow("Received upstream update from source",
				"issue", "8539",
				"source", sourceName,
				"upstreamCount", len(upstreams),
				"updateNumber", upstreamUpdateCount)

			upstreamUpdateCount++

			select {
			case dest <- &upstreamsWithSource{
				source:    sourceName,
				upstreams: upstreams,
			}:
				logger.Infow("Successfully forwarded upstream update",
					"issue", "8539",
					"source", sourceName,
					"upstreamCount", len(upstreams))
			case <-ctx.Done():
				logger.Infow("Context cancelled during upstream forwarding",
					"issue", "8539",
					"source", sourceName,
					"totalUpdates", upstreamUpdateCount)
				return
			}
		case <-ctx.Done():
			logger.Infow("Context cancelled, stopping upstream aggregation",
				"issue", "8539",
				"source", sourceName,
				"totalUpdates", upstreamUpdateCount)
			return
		}
	}
}
