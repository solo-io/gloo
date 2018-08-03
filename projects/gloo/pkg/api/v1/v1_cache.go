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
	AttributeList []*Attribute
	RoleList []*Role
	UpstreamList []*Upstream
	VirtualServiceList []*VirtualService
}

func (s Snapshot) Clone() Snapshot {
	var attributeList []*Attribute
	for _, attribute := range s.AttributeList {
		attributeList = append(attributeList, proto.Clone(attribute).(*Attribute))
	}
	var roleList []*Role
	for _, role := range s.RoleList {
		roleList = append(roleList, proto.Clone(role).(*Role))
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
		AttributeList: attributeList,
		RoleList: roleList,
		UpstreamList: upstreamList,
		VirtualServiceList: virtualServiceList,
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, attribute := range snapshotForHashing.AttributeList {
		resources.UpdateMetadata(attribute, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		attribute.SetStatus(core.Status{})
	}
	for _, role := range snapshotForHashing.RoleList {
		resources.UpdateMetadata(role, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		role.SetStatus(core.Status{})
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
	Attribute() AttributeClient
	Role() RoleClient
	Upstream() UpstreamClient
	VirtualService() VirtualServiceClient
	Snapshots(namespace string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewCache(attributeClient AttributeClient, roleClient RoleClient, upstreamClient UpstreamClient, virtualServiceClient VirtualServiceClient) Cache {
	return &cache{
		attribute: attributeClient,
		role: roleClient,
		upstream: upstreamClient,
		virtualService: virtualServiceClient,
	}
}

type cache struct {
	attribute AttributeClient
	role RoleClient
	upstream UpstreamClient
	virtualService VirtualServiceClient
}

func (c *cache) Register() error {
	if err := c.attribute.Register(); err != nil {
		return err
	}
	if err := c.role.Register(); err != nil {
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

func (c *cache) Attribute() AttributeClient {
	return c.attribute
}

func (c *cache) Role() RoleClient {
	return c.role
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
	attributeChan, attributeErrs, err := c.attribute.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Attribute watch")
	}
	roleChan, roleErrs, err := c.role.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Role watch")
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
			case attributeList := <-attributeChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.AttributeList = attributeList
				sync(newSnapshot)
			case roleList := <-roleChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.RoleList = roleList
				sync(newSnapshot)
			case upstreamList := <-upstreamChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.UpstreamList = upstreamList
				sync(newSnapshot)
			case virtualServiceList := <-virtualServiceChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.VirtualServiceList = virtualServiceList
				sync(newSnapshot)
			case err := <-attributeErrs:
				errs <- err
			case err := <-roleErrs:
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
