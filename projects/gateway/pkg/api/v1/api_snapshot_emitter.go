package v1

import (
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type ApiEmitter interface {
	Register() error
	Gateway() GatewayClient
	VirtualService() VirtualServiceClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *ApiSnapshot, <-chan error, error)
}

func NewApiEmitter(gatewayClient GatewayClient, virtualServiceClient VirtualServiceClient) ApiEmitter {
	return &apiEmitter{
		gateway:        gatewayClient,
		virtualService: virtualServiceClient,
	}
}

type apiEmitter struct {
	gateway        GatewayClient
	virtualService VirtualServiceClient
}

func (c *apiEmitter) Register() error {
	if err := c.gateway.Register(); err != nil {
		return err
	}
	if err := c.virtualService.Register(); err != nil {
		return err
	}
	return nil
}

func (c *apiEmitter) Gateway() GatewayClient {
	return c.gateway
}

func (c *apiEmitter) VirtualService() VirtualServiceClient {
	return c.virtualService
}

func (c *apiEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *ApiSnapshot, <-chan error, error) {
	snapshots := make(chan *ApiSnapshot)
	errs := make(chan error)
	var done sync.WaitGroup

	currentSnapshot := ApiSnapshot{}

	sync := func(newSnapshot ApiSnapshot) {
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
		done.Add(1)
		go func() {
			defer done.Done()
			errutils.AggregateErrs(opts.Ctx, errs, gatewayErrs, namespace+"-gateways")
		}()

		done.Add(1)
		go func(namespace string, gatewayChan <-chan GatewayList) {
			defer done.Done()
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
		done.Add(1)
		go func() {
			defer done.Done()
			errutils.AggregateErrs(opts.Ctx, errs, virtualServiceErrs, namespace+"-virtualServices")
		}()

		done.Add(1)
		go func(namespace string, virtualServiceChan <-chan VirtualServiceList) {
			defer done.Done()
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
			done.Wait()
			close(snapshots)
			close(errs)
		}
	}()
	return snapshots, errs, nil
}
