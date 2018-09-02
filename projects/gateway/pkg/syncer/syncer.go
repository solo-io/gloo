package syncer

import (
	"context"
	"strings"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/propagator"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

type syncer struct {
	writeNamespace  string
	reporter        reporter.Reporter
	propagator      *propagator.Propagator
	writeErrs       chan error
	proxyReconciler gloov1.ProxyReconciler
}

func NewSyncer(writeNamespace string, proxyClient gloov1.ProxyClient, reporter reporter.Reporter, propagator *propagator.Propagator, writeErrs chan error) v1.ApiSyncer {
	return &syncer{
		writeNamespace:  writeNamespace,
		reporter:        reporter,
		propagator:      propagator,
		writeErrs:       writeErrs,
		proxyReconciler: gloov1.NewProxyReconciler(proxyClient),
	}
}

func (s *syncer) Sync(ctx context.Context, snap *v1.ApiSnapshot) error {
	ctx = contextutils.WithLogger(ctx, "syncer")

	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("Beginning translation loop for snapshot %v", snap.Hash())
	logger.Debugf("%v", snap)

	desired, resourceErrs := translate(s.writeNamespace, snap)
	if err := s.reporter.WriteReports(ctx, resourceErrs); err != nil {
		return errors.Wrapf(err, "writing reports")
	}
	if err := resourceErrs.Validate(); err != nil {
		logger.Warnf("gateway %v was rejected due to invalid config: %v\nxDS cache will not be updated.", err)
		return nil
	}
	// proxy was deleted / none desired
	if desired == nil {
		return s.proxyReconciler.Reconcile(s.writeNamespace, nil, nil, clients.ListOpts{})
	}
	if err := s.proxyReconciler.Reconcile(s.writeNamespace, gloov1.ProxyList{desired}, nil, clients.ListOpts{}); err != nil {
		return err
	}

	// start propagating for new set of resources
	return s.propagator.PropagateStatuses(snap, desired, clients.WatchOpts{Ctx: ctx})
}

func translate(namespace string, snap *v1.ApiSnapshot) (*gloov1.Proxy, reporter.ResourceErrors) {
	resourceErrs := make(reporter.ResourceErrors)
	resourceErrs.Initialize(snap.Gateways.List().AsInputResources()...)
	resourceErrs.Initialize(snap.VirtualServices.List().AsInputResources()...)
	if len(snap.Gateways.List()) == 0 {
		return nil, resourceErrs
	}
	if len(snap.VirtualServices.List()) == 0 {
		return nil, resourceErrs
	}
	validateGateways(snap.Gateways.List(), resourceErrs)
	validateVirtualServices(snap.VirtualServices.List(), resourceErrs)
	meta := core.Metadata{
		Name:        joinGatewayNames(snap.Gateways.List()),
		Namespace:   namespace,
		Annotations: map[string]string{"owner_ref": "gateway"},
	}
	var listeners []*gloov1.Listener
	for _, gateway := range snap.Gateways.List() {
		listener := desiredListener(gateway, snap.VirtualServices.List(), resourceErrs)
		listeners = append(listeners, listener)
	}
	return &gloov1.Proxy{
		Metadata:  meta,
		Listeners: listeners,
	}, resourceErrs
}

func joinGatewayNames(gateways v1.GatewayList) string {
	var names []string
	for _, gw := range gateways {
		names = append(names, gw.Metadata.Name)
	}
	return strings.Join(names, ".")
}

// TODO(ilackarms): implement validation func
func validateGateways(gateways v1.GatewayList, resourceErrs reporter.ResourceErrors) {

}

func validateVirtualServices(virtualServices v1.VirtualServiceList, resourceErrs reporter.ResourceErrors) {

}

func desiredListener(gateway *v1.Gateway, virtualServices v1.VirtualServiceList, resourceErrs reporter.ResourceErrors) *gloov1.Listener {
	if len(gateway.VirtualServices) == 0 {
		resourceErrs.AddError(gateway, errors.Errorf("must specify at least one virtual service on gateway"))
	}
	var (
		virtualHosts []*gloov1.VirtualHost
		sslConfigs   []*gloov1.SslConfig
	)

	// add all virtual services if empty
	if len(gateway.VirtualServices) == 0 {
		for _, virtualService := range virtualServices {
			gateway.VirtualServices = append(gateway.VirtualServices, core.ResourceRef{
				Name:      virtualService.GetMetadata().Name,
				Namespace: virtualService.GetMetadata().Namespace,
			})
		}
	}

	for _, ref := range gateway.VirtualServices {
		// virtual service must live in the same namespace as gateway
		virtualService, err := virtualServices.Find(ref.Strings())
		if err != nil {
			resourceErrs.AddError(gateway, err)
			continue
		}
		virtualHosts = append(virtualHosts, virtualService.VirtualHost)
	}
	return &gloov1.Listener{
		Name:        gateway.Metadata.Name,
		BindAddress: gateway.BindAddress,
		BindPort:    gateway.BindPort,
		ListenerType: &gloov1.Listener_HttpListener{
			HttpListener: &gloov1.HttpListener{
				VirtualHosts: virtualHosts,
			},
		},
		SslConfiguations: sslConfigs,
	}
}
