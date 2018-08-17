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
	GatewayList        GatewayList
	VirtualServiceList VirtualServiceList
}

func (s Snapshot) Clone() Snapshot {
	var gatewayList []*Gateway
	for _, gateway := range s.GatewayList {
		gatewayList = append(gatewayList, proto.Clone(gateway).(*Gateway))
	}
	var virtualServiceList []*VirtualService
	for _, virtualService := range s.VirtualServiceList {
		virtualServiceList = append(virtualServiceList, proto.Clone(virtualService).(*VirtualService))
	}
	return Snapshot{
		GatewayList:        gatewayList,
		VirtualServiceList: virtualServiceList,
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, gateway := range snapshotForHashing.GatewayList {
		resources.UpdateMetadata(gateway, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		gateway.SetStatus(core.Status{})
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
	Gateway() GatewayClient
	VirtualService() VirtualServiceClient
	Snapshots(namespace string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewCache(gatewayClient GatewayClient, virtualServiceClient VirtualServiceClient) Cache {
	return &cache{
		gateway:        gatewayClient,
		virtualService: virtualServiceClient,
	}
}

type cache struct {
	gateway        GatewayClient
	virtualService VirtualServiceClient
}

func (c *cache) Register() error {
	if err := c.gateway.Register(); err != nil {
		return err
	}
	if err := c.virtualService.Register(); err != nil {
		return err
	}
	return nil
}

func (c *cache) Gateway() GatewayClient {
	return c.gateway
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
	gatewayChan, gatewayErrs, err := c.gateway.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Gateway watch")
	}
	virtualServiceChan, virtualServiceErrs, err := c.virtualService.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting VirtualService watch")
	}

	go func() {
		for {
			select {
			case gatewayList := <-gatewayChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.GatewayList = gatewayList
				sync(newSnapshot)
			case virtualServiceList := <-virtualServiceChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.VirtualServiceList = virtualServiceList
				sync(newSnapshot)
			case err := <-gatewayErrs:
				errs <- err
			case err := <-virtualServiceErrs:
				errs <- err
			}
		}
	}()
	return snapshots, errs, nil
}
