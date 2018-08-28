package v1

import (
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type Snapshot struct {
	Artifacts    ArtifactListsByNamespace
	EndpointList EndpointList
	ProxyList    ProxyList
	SecretList   SecretList
	UpstreamList UpstreamList
}

func (s Snapshot) Clone() Snapshot {
	return Snapshot{
		Artifacts:    s.Artifacts.Clone(),
		EndpointList: endpointList,
		ProxyList:    proxyList,
		SecretList:   secretList,
		UpstreamList: upstreamList,
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, artifact := range snapshotForHashing.Artifacts.List() {
		resources.UpdateMetadata(artifact, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
	}
	for _, endpoint := range snapshotForHashing.EndpointList {
		resources.UpdateMetadata(endpoint, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
	}
	for _, proxy := range snapshotForHashing.ProxyList {
		resources.UpdateMetadata(proxy, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		proxy.SetStatus(core.Status{})
	}
	for _, secret := range snapshotForHashing.SecretList {
		resources.UpdateMetadata(secret, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
	}
	for _, upstream := range snapshotForHashing.UpstreamList {
		resources.UpdateMetadata(upstream, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		upstream.SetStatus(core.Status{})
	}
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}

type Cache interface {
	Register() error
	Artifact() ArtifactClient
	Endpoint() EndpointClient
	Proxy() ProxyClient
	Secret() SecretClient
	Upstream() UpstreamClient
	Snapshots(namespace string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewCache(artifactClient ArtifactClient, endpointClient EndpointClient, proxyClient ProxyClient, secretClient SecretClient, upstreamClient UpstreamClient) Cache {
	return &cache{
		artifact: artifactClient,
		endpoint: endpointClient,
		proxy:    proxyClient,
		secret:   secretClient,
		upstream: upstreamClient,
	}
}

type cache struct {
	artifact ArtifactClient
	endpoint EndpointClient
	proxy    ProxyClient
	secret   SecretClient
	upstream UpstreamClient
}

func (c *cache) Register() error {
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

func (c *cache) Artifact() ArtifactClient {
	return c.artifact
}

func (c *cache) Endpoint() EndpointClient {
	return c.endpoint
}

func (c *cache) Proxy() ProxyClient {
	return c.proxy
}

func (c *cache) Secret() SecretClient {
	return c.secret
}

func (c *cache) Upstream() UpstreamClient {
	return c.upstream
}

func (c *cache) Snapshots(namespaces []string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error) {
	snapshots := make(chan *Snapshot)
	errs := make(chan error)

	currentSnapshot := Snapshot{}

	sync := func(newSnapshot Snapshot) {
		if currentSnapshot.Hash() == newSnapshot.Hash() {
			return
		}
		currentSnapshot = newSnapshot
		snapshots <- &currentSnapshot
	}

	for _, namespace := range namespaces {
		artifactChan, artifactErrs, err := c.artifact.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Artifact watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, artifactErrs, namespace+"-artifacts")
		go func() {
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
		}()
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
