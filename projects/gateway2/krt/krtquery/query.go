package krtquery

import (
	"context"

	"k8s.io/client-go/tools/clientcmd"
	gwapi "sigs.k8s.io/gateway-api/apis/v1beta1"

	"istio.io/istio/pkg/config/schema/gvr"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/kube/kubetypes"

	solov1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
)

type Queries interface{}

type queries struct {
	HTTPRouteAttachements         krt.Collection[HTTPRouteAttachement]
	VirtualHostOptionAttachements krt.Collection[VirtualHostOptionAttachement]
}

func New(ctx context.Context, cfg clientcmd.ClientConfig) (Queries, error) {
	client, err := kube.NewClient(cfg, "gloo-gateway")
	if err != nil {
		return nil, err
	}

	filter := kclient.Filter{}

	gatewayClient := kclient.NewDelayedInformer[*gwapi.Gateway](client, gvr.KubernetesGateway, kubetypes.StandardInformer, filter)
	Gateways := krt.WrapClient[*gwapi.Gateway](gatewayClient, krt.WithName("Gateways"))

	httprouteClient := kclient.NewDelayedInformer[*gwapi.HTTPRoute](client, gvr.HTTPRoute, kubetypes.StandardInformer, filter)
	HTTPRoutes := krt.WrapClient[*gwapi.HTTPRoute](httprouteClient, krt.WithName("HTTPRoutes"))

	// TODO idk if this way of writing GVK actually works
	virtualHostOptionClient := kclient.NewDelayedInformer[*solov1.VirtualHostOption](client, solov1.SchemeGroupVersion.WithResource("virtualhostoption"), kubetypes.StandardInformer, filter)
	VirtualHostOptions := krt.WrapClient[*solov1.VirtualHostOption](virtualHostOptionClient, krt.WithName("VirtualHostOptions"))

	// start informers (probably should move this elsewhere)
	client.RunAndWait(ctx.Done())

	return queries{
		HTTPRouteAttachements(Gateways, HTTPRoutes),
		VirtualHostOptionAttachements(Gateways, VirtualHostOptions),
	}, nil
}
