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
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

type syncer struct {
	namespace       string
	reporter        reporter.Reporter
	proxyReconciler gloov1.ProxyReconciler
}

func NewSyncer(namespace string, proxyClient gloov1.ProxyClient, reporter reporter.Reporter) v1.Syncer {
	return &syncer{
		namespace:       namespace,
		reporter:        reporter,
		proxyReconciler: gloov1.NewProxyReconciler(proxyClient),
	}
}

func metaForSnap(snap *v1.Snapshot) core.Metadata {}

func (s *syncer) Sync(ctx context.Context, snap *v1.Snapshot) error {
	ctx = contextutils.WithLogger(ctx, "gateway.syncer")
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("Beginning translation loop for snapshot %v", snap.Hash())
	logger.Debugf("%v", snap)

	desired, resourceErrs := translate(s.namespace, snap)
	if err := s.reporter.WriteReports(ctx, resourceErrs); err != nil {
		return errors.Wrapf(err, "writing reports")
	}
	if err := resourceErrs.Validate(); err != nil {
		logger.Warnf("gateway %v was rejected due to invalid config: %v\nxDS cache will not be updated.", err)
		return nil
	}
	// proxy was deleted / none desired
	if desired == nil {
		return s.proxyReconciler.Reconcile(s.namespace, nil, nil, clients.ListOpts{})
	}
	return s.proxyReconciler.Reconcile(s.namespace, nil, nil, clients.ListOpts{})
}

func translate(namespace string, snap *v1.Snapshot) (*gloov1.Proxy, reporter.ResourceErrors) {
	resourceErrs := make(reporter.ResourceErrors)
	resourceErrs.Initialize(snap.GatewayList.AsInputResources()...)
	resourceErrs.Initialize(snap.VirtualServiceList.AsInputResources()...)
	if len(snap.GatewayList) == 0 {
		return nil, resourceErrs
	}
	if len(snap.VirtualServiceList) == 0 {
		return nil, resourceErrs
	}
	validateGateways(snap.GatewayList, resourceErrs)
	validateVirtualServices(snap.VirtualServiceList, resourceErrs)
	meta := core.Metadata{
		Name:        joinGatewayNames(snap.GatewayList),
		Namespace:   namespace,
		Annotations: map[string]string{"owner_ref": "gateway"},
	}
	var listeners []*gloov1.Listener
	for _, gateway := range snap.GatewayList {
		listener := desiredListener(gateway, snap.VirtualServiceList, resourceErrs)
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

func desiredListener(gateway *v1.Gateway, virtualSerivces v1.VirtualServiceList, resourceErrs reporter.ResourceErrors) *gloov1.Listener {
	if len(gateway.VirtualServices) == 0 {
		resourceErrs.AddError(gateway, errors.Errorf("must specify at least one virtual service on gateway"))
	}
	var (
		virtualHosts []*gloov1.VirtualHost
		sslConfigs   []*gloov1.SslConfig
	)
	for _, name := range gateway.VirtualServices {
		// virtual service must live in the same namespace as gateway
		virtualService, err := virtualSerivces.Find(gateway.Metadata.Namespace, name)
		if err != nil {
			resourceErrs.AddError(gateway, err)
			continue
		}

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
