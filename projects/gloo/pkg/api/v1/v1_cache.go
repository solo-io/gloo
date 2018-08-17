package v1

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
)

type Snapshot struct {
	ArtifactList ArtifactList
	EndpointList EndpointList
	ProxyList    ProxyList
	SecretList   SecretList
	UpstreamList UpstreamList
}

func (s Snapshot) Clone() Snapshot {
	var artifactList []*Artifact
	for _, artifact := range s.ArtifactList {
		artifactList = append(artifactList, proto.Clone(artifact).(*Artifact))
	}
	var endpointList []*Endpoint
	for _, endpoint := range s.EndpointList {
		endpointList = append(endpointList, proto.Clone(endpoint).(*Endpoint))
	}
	var proxyList []*Proxy
	for _, proxy := range s.ProxyList {
		proxyList = append(proxyList, proto.Clone(proxy).(*Proxy))
	}
	var secretList []*Secret
	for _, secret := range s.SecretList {
		secretList = append(secretList, proto.Clone(secret).(*Secret))
	}
	var upstreamList []*Upstream
	for _, upstream := range s.UpstreamList {
		upstreamList = append(upstreamList, proto.Clone(upstream).(*Upstream))
	}
	return Snapshot{
		ArtifactList: artifactList,
		EndpointList: endpointList,
		ProxyList:    proxyList,
		SecretList:   secretList,
		UpstreamList: upstreamList,
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, artifact := range snapshotForHashing.ArtifactList {
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

func (c *cache) Snapshots(namespace string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error) {
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
	artifactChan, artifactErrs, err := c.artifact.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Artifact watch")
	}
	endpointChan, endpointErrs, err := c.endpoint.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Endpoint watch")
	}
	proxyChan, proxyErrs, err := c.proxy.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Proxy watch")
	}
	secretChan, secretErrs, err := c.secret.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Secret watch")
	}
	upstreamChan, upstreamErrs, err := c.upstream.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Upstream watch")
	}

	go func() {
		for {
			select {
			case artifactList := <-artifactChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.ArtifactList = artifactList
				sync(newSnapshot)
			case endpointList := <-endpointChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.EndpointList = endpointList
				sync(newSnapshot)
			case proxyList := <-proxyChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.ProxyList = proxyList
				sync(newSnapshot)
			case secretList := <-secretChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.SecretList = secretList
				sync(newSnapshot)
			case upstreamList := <-upstreamChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.UpstreamList = upstreamList
				sync(newSnapshot)
			case err := <-artifactErrs:
				errs <- err
			case err := <-endpointErrs:
				errs <- err
			case err := <-proxyErrs:
				errs <- err
			case err := <-secretErrs:
				errs <- err
			case err := <-upstreamErrs:
				errs <- err
			}
		}
	}()
	return snapshots, errs, nil
}
