package v1

import (
	"sync"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

var (
	MSnapshotIn  = stats.Int64("snapemitter/snapin", "The number of snapshots in", "1")
	MSnapshotOut = stats.Int64("snapemitter/snapout", "The number of snapshots out", "1")

	//	KeyNamespace, _ = tag.NewKey("namespace")

	SnapshotInView = &view.View{
		Name:        "snapemitter/snapin",
		Measure:     MSnapshotIn,
		Description: "The number of snapshots updates coming in",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{
			//		KeyNamespace,
		},
	}
	SnapshotOutView = &view.View{
		Name:        "snapemitter/snapout",
		Measure:     MSnapshotOut,
		Description: "The number of snapshots updates going out",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{
			//			KeyNamespace,
		},
	}
)

func init() {
	view.Register(SnapshotInView, SnapshotOutView)
}

type ApiEmitter interface {
	Register() error
	Artifact() ArtifactClient
	Endpoint() EndpointClient
	Proxy() ProxyClient
	Secret() SecretClient
	Upstream() UpstreamClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *ApiSnapshot, <-chan error, error)
}

func NewApiEmitter(artifactClient ArtifactClient, endpointClient EndpointClient, proxyClient ProxyClient, secretClient SecretClient, upstreamClient UpstreamClient) ApiEmitter {
	return NewApiEmitterWithEmit(artifactClient, endpointClient, proxyClient, secretClient, upstreamClient, make(chan struct{}))
}

func NewApiEmitterWithEmit(artifactClient ArtifactClient, endpointClient EndpointClient, proxyClient ProxyClient, secretClient SecretClient, upstreamClient UpstreamClient, emit <-chan struct{}) ApiEmitter {
	return &apiEmitter{
		artifact:  artifactClient,
		endpoint:  endpointClient,
		proxy:     proxyClient,
		secret:    secretClient,
		upstream:  upstreamClient,
		forceEmit: emit,
	}
}

type apiEmitter struct {
	forceEmit <-chan struct{}
	artifact  ArtifactClient
	endpoint  EndpointClient
	proxy     ProxyClient
	secret    SecretClient
	upstream  UpstreamClient
}

func (c *apiEmitter) Register() error {
	if err := c.artifact.Register(); err != nil {
		return err
	}
	if err := c.endpoint.Register(); err != nil {
		return err
	}
	if err := c.proxy.Register(); err != nil {
		return err
	}
	if err := c.secret.Register(); err != nil {
		return err
	}
	if err := c.upstream.Register(); err != nil {
		return err
	}
	return nil
}

func (c *apiEmitter) Artifact() ArtifactClient {
	return c.artifact
}

func (c *apiEmitter) Endpoint() EndpointClient {
	return c.endpoint
}

func (c *apiEmitter) Proxy() ProxyClient {
	return c.proxy
}

func (c *apiEmitter) Secret() SecretClient {
	return c.secret
}

func (c *apiEmitter) Upstream() UpstreamClient {
	return c.upstream
}

func (c *apiEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *ApiSnapshot, <-chan error, error) {

	ctx := opts.Ctx

	errs := make(chan error)
	var done sync.WaitGroup
	/* Create channel for Artifact */
	type artifactListWithNamespace struct {
		list      ArtifactList
		namespace string
	}
	artifactChan := make(chan artifactListWithNamespace)
	/* Create channel for Endpoint */
	type endpointListWithNamespace struct {
		list      EndpointList
		namespace string
	}
	endpointChan := make(chan endpointListWithNamespace)
	/* Create channel for Proxy */
	type proxyListWithNamespace struct {
		list      ProxyList
		namespace string
	}
	proxyChan := make(chan proxyListWithNamespace)
	/* Create channel for Secret */
	type secretListWithNamespace struct {
		list      SecretList
		namespace string
	}
	secretChan := make(chan secretListWithNamespace)
	/* Create channel for Upstream */
	type upstreamListWithNamespace struct {
		list      UpstreamList
		namespace string
	}
	upstreamChan := make(chan upstreamListWithNamespace)

	for _, namespace := range watchNamespaces {
		/* Setup watch for Artifact */
		artifactNamespacesChan, artifactErrs, err := c.artifact.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Artifact watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, artifactErrs, namespace+"-artifacts")
		}(namespace)
		/* Setup watch for Endpoint */
		endpointNamespacesChan, endpointErrs, err := c.endpoint.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Endpoint watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, endpointErrs, namespace+"-endpoints")
		}(namespace)
		/* Setup watch for Proxy */
		proxyNamespacesChan, proxyErrs, err := c.proxy.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Proxy watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, proxyErrs, namespace+"-proxies")
		}(namespace)
		/* Setup watch for Secret */
		secretNamespacesChan, secretErrs, err := c.secret.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Secret watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, secretErrs, namespace+"-secrets")
		}(namespace)
		/* Setup watch for Upstream */
		upstreamNamespacesChan, upstreamErrs, err := c.upstream.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Upstream watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, upstreamErrs, namespace+"-upstreams")
		}(namespace)

		/* Watch for changes and update snapshot */
		go func(namespace string) {
			for {
				select {
				case <-ctx.Done():
					return
				case artifactList := <-artifactNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case artifactChan <- artifactListWithNamespace{list: artifactList, namespace: namespace}:
					}
				case endpointList := <-endpointNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case endpointChan <- endpointListWithNamespace{list: endpointList, namespace: namespace}:
					}
				case proxyList := <-proxyNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case proxyChan <- proxyListWithNamespace{list: proxyList, namespace: namespace}:
					}
				case secretList := <-secretNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case secretChan <- secretListWithNamespace{list: secretList, namespace: namespace}:
					}
				case upstreamList := <-upstreamNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case upstreamChan <- upstreamListWithNamespace{list: upstreamList, namespace: namespace}:
					}
				}
			}
		}(namespace)
	}

	snapshots := make(chan *ApiSnapshot)

	go func() {
		currentSnapshot := ApiSnapshot{}
		sync := func(newSnapshot ApiSnapshot) {

			if currentSnapshot.Hash() == newSnapshot.Hash() {
				return
			}
			currentSnapshot = newSnapshot
			sentSnapshot := currentSnapshot.Clone()

			stats.Record(ctx, MSnapshotOut.M(1))
			snapshots <- &sentSnapshot
		}
		for {
			select {
			case <-ctx.Done():
				close(snapshots)
				done.Wait()
				close(errs)
				return
			case <-c.forceEmit:
				sentSnapshot := currentSnapshot.Clone()
				snapshots <- &sentSnapshot
			case artifactNamespacedList := <-artifactChan:
				namespace := artifactNamespacedList.namespace
				artifactList := artifactNamespacedList.list

				newSnapshot := currentSnapshot.Clone()
				newSnapshot.Artifacts.Clear(namespace)
				newSnapshot.Artifacts.Add(artifactList...)
				sync(newSnapshot)
			case endpointNamespacedList := <-endpointChan:
				namespace := endpointNamespacedList.namespace
				endpointList := endpointNamespacedList.list

				newSnapshot := currentSnapshot.Clone()
				newSnapshot.Endpoints.Clear(namespace)
				newSnapshot.Endpoints.Add(endpointList...)
				sync(newSnapshot)
			case proxyNamespacedList := <-proxyChan:
				namespace := proxyNamespacedList.namespace
				proxyList := proxyNamespacedList.list

				newSnapshot := currentSnapshot.Clone()
				newSnapshot.Proxies.Clear(namespace)
				newSnapshot.Proxies.Add(proxyList...)
				sync(newSnapshot)
			case secretNamespacedList := <-secretChan:
				namespace := secretNamespacedList.namespace
				secretList := secretNamespacedList.list

				newSnapshot := currentSnapshot.Clone()
				newSnapshot.Secrets.Clear(namespace)
				newSnapshot.Secrets.Add(secretList...)
				sync(newSnapshot)
			case upstreamNamespacedList := <-upstreamChan:
				namespace := upstreamNamespacedList.namespace
				upstreamList := upstreamNamespacedList.list

				newSnapshot := currentSnapshot.Clone()
				newSnapshot.Upstreams.Clear(namespace)
				newSnapshot.Upstreams.Add(upstreamList...)
				sync(newSnapshot)
			}
			// if we got here its because a new entry in the channel
			stats.Record(ctx, MSnapshotIn.M(1))
		}
	}()
	return snapshots, errs, nil
}
