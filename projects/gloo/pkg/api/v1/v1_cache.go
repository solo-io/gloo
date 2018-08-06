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
	ArtifactList       []*Artifact
	AttributeList      []*Attribute
	EndpointList       []*Endpoint
	RoleList           []*Role
	SecretList         []*Secret
	UpstreamList       []*Upstream
	VirtualServiceList []*VirtualService
}

func (s Snapshot) Clone() Snapshot {
	var artifactList []*Artifact
	for _, artifact := range s.ArtifactList {
		artifactList = append(artifactList, proto.Clone(artifact).(*Artifact))
	}
	var attributeList []*Attribute
	for _, attribute := range s.AttributeList {
		attributeList = append(attributeList, proto.Clone(attribute).(*Attribute))
	}
	var endpointList []*Endpoint
	for _, endpoint := range s.EndpointList {
		endpointList = append(endpointList, proto.Clone(endpoint).(*Endpoint))
	}
	var roleList []*Role
	for _, role := range s.RoleList {
		roleList = append(roleList, proto.Clone(role).(*Role))
	}
	var secretList []*Secret
	for _, secret := range s.SecretList {
		secretList = append(secretList, proto.Clone(secret).(*Secret))
	}
	var upstreamList []*Upstream
	for _, upstream := range s.UpstreamList {
		upstreamList = append(upstreamList, proto.Clone(upstream).(*Upstream))
	}
	var virtualServiceList []*VirtualService
	for _, virtualService := range s.VirtualServiceList {
		virtualServiceList = append(virtualServiceList, proto.Clone(virtualService).(*VirtualService))
	}
	return Snapshot{
		ArtifactList:       artifactList,
		AttributeList:      attributeList,
		EndpointList:       endpointList,
		RoleList:           roleList,
		SecretList:         secretList,
		UpstreamList:       upstreamList,
		VirtualServiceList: virtualServiceList,
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, artifact := range snapshotForHashing.ArtifactList {
		resources.UpdateMetadata(artifact, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
	}
	for _, attribute := range snapshotForHashing.AttributeList {
		resources.UpdateMetadata(attribute, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		attribute.SetStatus(core.Status{})
	}
	for _, endpoint := range snapshotForHashing.EndpointList {
		resources.UpdateMetadata(endpoint, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
	}
	for _, role := range snapshotForHashing.RoleList {
		resources.UpdateMetadata(role, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		role.SetStatus(core.Status{})
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
	for _, virtualService := range snapshotForHashing.VirtualServiceList {
		resources.UpdateMetadata(virtualService, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		virtualService.SetStatus(core.Status{})
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
	Attribute() AttributeClient
	Endpoint() EndpointClient
	Role() RoleClient
	Secret() SecretClient
	Upstream() UpstreamClient
	VirtualService() VirtualServiceClient
	Snapshots(namespace string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewCache(artifactClient ArtifactClient, attributeClient AttributeClient, endpointClient EndpointClient, roleClient RoleClient, secretClient SecretClient, upstreamClient UpstreamClient, virtualServiceClient VirtualServiceClient) Cache {
	return &cache{
		artifact:       artifactClient,
		attribute:      attributeClient,
		endpoint:       endpointClient,
		role:           roleClient,
		secret:         secretClient,
		upstream:       upstreamClient,
		virtualService: virtualServiceClient,
	}
}

type cache struct {
	artifact       ArtifactClient
	attribute      AttributeClient
	endpoint       EndpointClient
	role           RoleClient
	secret         SecretClient
	upstream       UpstreamClient
	virtualService VirtualServiceClient
}

func (c *cache) Register() error {
	if err := c.artifact.Register(); err != nil {
		return err
	}
	if err := c.attribute.Register(); err != nil {
		return err
	}
	if err := c.endpoint.Register(); err != nil {
		return err
	}
	if err := c.role.Register(); err != nil {
		return err
	}
	if err := c.secret.Register(); err != nil {
		return err
	}
	if err := c.upstream.Register(); err != nil {
		return err
	}
	if err := c.virtualService.Register(); err != nil {
		return err
	}
	return nil
}

func (c *cache) Artifact() ArtifactClient {
	return c.artifact
}

func (c *cache) Attribute() AttributeClient {
	return c.attribute
}

func (c *cache) Endpoint() EndpointClient {
	return c.endpoint
}

func (c *cache) Role() RoleClient {
	return c.role
}

func (c *cache) Secret() SecretClient {
	return c.secret
}

func (c *cache) Upstream() UpstreamClient {
	return c.upstream
}

func (c *cache) VirtualService() VirtualServiceClient {
	return c.virtualService
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
	attributeChan, attributeErrs, err := c.attribute.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Attribute watch")
	}
	endpointChan, endpointErrs, err := c.endpoint.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Endpoint watch")
	}
	roleChan, roleErrs, err := c.role.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Role watch")
	}
	secretChan, secretErrs, err := c.secret.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Secret watch")
	}
	upstreamChan, upstreamErrs, err := c.upstream.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Upstream watch")
	}
	virtualServiceChan, virtualServiceErrs, err := c.virtualService.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting VirtualService watch")
	}

	go func() {
		for {
			select {
			case artifactList := <-artifactChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.ArtifactList = artifactList
				sync(newSnapshot)
			case attributeList := <-attributeChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.AttributeList = attributeList
				sync(newSnapshot)
			case endpointList := <-endpointChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.EndpointList = endpointList
				sync(newSnapshot)
			case roleList := <-roleChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.RoleList = roleList
				sync(newSnapshot)
			case secretList := <-secretChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.SecretList = secretList
				sync(newSnapshot)
			case upstreamList := <-upstreamChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.UpstreamList = upstreamList
				sync(newSnapshot)
			case virtualServiceList := <-virtualServiceChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.VirtualServiceList = virtualServiceList
				sync(newSnapshot)
			case err := <-artifactErrs:
				errs <- err
			case err := <-attributeErrs:
				errs <- err
			case err := <-endpointErrs:
				errs <- err
			case err := <-roleErrs:
				errs <- err
			case err := <-secretErrs:
				errs <- err
			case err := <-upstreamErrs:
				errs <- err
			case err := <-virtualServiceErrs:
				errs <- err
			}
		}
	}()
	return snapshots, errs, nil
}
