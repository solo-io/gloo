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
	Gateways        GatewaysByNamespace
	VirtualServices VirtualServicesByNamespace
}

func (s Snapshot) Clone() Snapshot {
	return Snapshot{
		Gateways:        s.Gateways.Clone(),
		VirtualServices: s.VirtualServices.Clone(),
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, gateway := range snapshotForHashing.Gateways.List() {
		resources.UpdateMetadata(gateway, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		gateway.SetStatus(core.Status{})
	}
	for _, virtualService := range snapshotForHashing.VirtualServices.List() {
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
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
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

func (c *cache) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error) {
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

	for _, namespace := range watchNamespaces {
		gatewayChan, gatewayErrs, err := c.gateway.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Gateway watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, gatewayErrs, namespace+"-gateways")
		go func(namespace string, gatewayChan <-chan GatewayList) {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case gatewayList := <-gatewayChan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.Gateways.Clear(namespace)
					newSnapshot.Gateways.Add(gatewayList...)
					sync(newSnapshot)
				}
			}
		}(namespace, gatewayChan)
		virtualServiceChan, virtualServiceErrs, err := c.virtualService.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting VirtualService watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, virtualServiceErrs, namespace+"-virtualServices")
		go func(namespace string, virtualServiceChan <-chan VirtualServiceList) {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case virtualServiceList := <-virtualServiceChan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.VirtualServices.Clear(namespace)
					newSnapshot.VirtualServices.Add(virtualServiceList...)
					sync(newSnapshot)
				}
			}
		}(namespace, virtualServiceChan)
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
