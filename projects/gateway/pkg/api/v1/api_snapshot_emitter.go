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

	errs := make(chan error)
	var done sync.WaitGroup
	/* Create channel for Gateway */
	type gatewayListWithNamespace struct {
		list      GatewayList
		namespace string
	}
	gatewayChan := make(chan gatewayListWithNamespace)
	/* Create channel for VirtualService */
	type virtualServiceListWithNamespace struct {
		list      VirtualServiceList
		namespace string
	}
	virtualServiceChan := make(chan virtualServiceListWithNamespace)

	for _, namespace := range watchNamespaces {
		/* Setup watch for Gateway */
		gatewayNamespacesChan, gatewayErrs, err := c.gateway.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Gateway watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(opts.Ctx, errs, gatewayErrs, namespace+"-gateways")
		}(namespace)
		/* Setup watch for VirtualService */
		virtualServiceNamespacesChan, virtualServiceErrs, err := c.virtualService.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting VirtualService watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(opts.Ctx, errs, virtualServiceErrs, namespace+"-virtualServices")
		}(namespace)

		/* Watch for changes and update snapshot */
		go func(namespace string) {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case gatewayList := <-gatewayNamespacesChan:
					select {
					case <-opts.Ctx.Done():
						return
					case gatewayChan <- gatewayListWithNamespace{list: gatewayList, namespace: namespace}:
					}
				case virtualServiceList := <-virtualServiceNamespacesChan:
					select {
					case <-opts.Ctx.Done():
						return
					case virtualServiceChan <- virtualServiceListWithNamespace{list: virtualServiceList, namespace: namespace}:
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
			snapshots <- &sentSnapshot
		}
		for {
			select {
			case <-opts.Ctx.Done():
				close(snapshots)
				done.Wait()
				close(errs)
				return
			case gatewayNamespacedList := <-gatewayChan:
				namespace := gatewayNamespacedList.namespace
				gatewayList := gatewayNamespacedList.list

				newSnapshot := currentSnapshot.Clone()
				newSnapshot.Gateways.Clear(namespace)
				newSnapshot.Gateways.Add(gatewayList...)
				sync(newSnapshot)
			case virtualServiceNamespacedList := <-virtualServiceChan:
				namespace := virtualServiceNamespacedList.namespace
				virtualServiceList := virtualServiceNamespacedList.list

				newSnapshot := currentSnapshot.Clone()
				newSnapshot.VirtualServices.Clear(namespace)
				newSnapshot.VirtualServices.Add(virtualServiceList...)
				sync(newSnapshot)
			}
		}
	}()
	return snapshots, errs, nil
}
