package krtquery

import (
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"istio.io/istio/pkg/config/schema/gvr"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/kube/kubetypes"

	solov1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
)

type Queries struct {
	Gateways                      krt.Collection[*gwv1.Gateway]
	HTTPRouteAttachements         krt.Collection[HTTPRouteAttachement]
	VirtualHostOptionAttachements krt.Collection[VirtualHostOptionAttachement]
}

func New(client kube.Client) (Queries, error) {
	filter := kclient.Filter{}

	gatewayClient := kclient.NewDelayedInformer[*gwv1.Gateway](client, gvr.KubernetesGateway, kubetypes.StandardInformer, filter)
	Gateways := krt.WrapClient[*gwv1.Gateway](gatewayClient, krt.WithName("Gateways"))

	httprouteClient := kclient.NewDelayedInformer[*gwv1.HTTPRoute](client, gvr.HTTPRoute, kubetypes.StandardInformer, filter)
	HTTPRoutes := krt.WrapClient[*gwv1.HTTPRoute](httprouteClient, krt.WithName("HTTPRoutes"))

	// TODO idk if this way of writing GVK actually works
	virtualHostOptionClient := kclient.NewDelayedInformer[*solov1.VirtualHostOption](client, solov1.SchemeGroupVersion.WithResource("virtualhostoption"), kubetypes.StandardInformer, filter)
	VirtualHostOptions := krt.WrapClient[*solov1.VirtualHostOption](virtualHostOptionClient, krt.WithName("VirtualHostOptions"))

	return Queries{
		Gateways,
		HTTPRouteAttachements(Gateways, HTTPRoutes),
		VirtualHostOptionAttachements(Gateways, VirtualHostOptions),
	}, nil
}
