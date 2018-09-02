package v1

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

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
	return &apiEmitter{
		artifact: artifactClient,
		endpoint: endpointClient,
		proxy:    proxyClient,
		secret:   secretClient,
		upstream: upstreamClient,
	}
}

type apiEmitter struct {
	artifact ArtifactClient
	endpoint EndpointClient
	proxy    ProxyClient
	secret   SecretClient
	upstream UpstreamClient
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
	snapshots := make(chan *ApiSnapshot)
	errs := make(chan error)

	currentSnapshot := ApiSnapshot{}

	sync := func(newSnapshot ApiSnapshot) {
		if currentSnapshot.Hash() == newSnapshot.Hash() {
			return
		}
		currentSnapshot = newSnapshot
		snapshots <- &currentSnapshot
	}

	for _, namespace := range watchNamespaces {
		artifactChan, artifactErrs, err := c.artifact.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Artifact watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, artifactErrs, namespace+"-artifacts")
		go func(namespace string, artifactChan <-chan ArtifactList) {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case artifactList := <-artifactChan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.Artifacts.Clear(namespace)
					newSnapshot.Artifacts.Add(artifactList...)
					sync(newSnapshot)
				}
			}
		}(namespace, artifactChan)
		endpointChan, endpointErrs, err := c.endpoint.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Endpoint watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, endpointErrs, namespace+"-endpoints")
		go func(namespace string, endpointChan <-chan EndpointList) {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case endpointList := <-endpointChan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.Endpoints.Clear(namespace)
					newSnapshot.Endpoints.Add(endpointList...)
					sync(newSnapshot)
				}
			}
		}(namespace, endpointChan)
		proxyChan, proxyErrs, err := c.proxy.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Proxy watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, proxyErrs, namespace+"-proxies")
		go func(namespace string, proxyChan <-chan ProxyList) {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case proxyList := <-proxyChan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.Proxies.Clear(namespace)
					newSnapshot.Proxies.Add(proxyList...)
					sync(newSnapshot)
				}
			}
		}(namespace, proxyChan)
		secretChan, secretErrs, err := c.secret.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Secret watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, secretErrs, namespace+"-secrets")
		go func(namespace string, secretChan <-chan SecretList) {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case secretList := <-secretChan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.Secrets.Clear(namespace)
					newSnapshot.Secrets.Add(secretList...)
					sync(newSnapshot)
				}
			}
		}(namespace, secretChan)
		upstreamChan, upstreamErrs, err := c.upstream.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Upstream watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, upstreamErrs, namespace+"-upstreams")
		go func(namespace string, upstreamChan <-chan UpstreamList) {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case upstreamList := <-upstreamChan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.Upstreams.Clear(namespace)
					newSnapshot.Upstreams.Add(upstreamList...)
					sync(newSnapshot)
				}
			}
		}(namespace, upstreamChan)
	}

	go func() {
		select {
		case <-opts.Ctx.Done():
			close(snapshots)
			close(errs)
		}
	}()
	return snapshots, errs, nil
}
