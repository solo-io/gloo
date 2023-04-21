package utils

import (
	"context"

	"github.com/solo-io/gloo/pkg/utils/statusutils"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"k8s.io/client-go/kubernetes"
)

type TestClients struct {
	UpstreamClient       gloov1.UpstreamClient
	VirtualServiceClient gatewayv1.VirtualServiceClient
	StatusClient         resources.StatusClient
}

func NewTestClients(ctx context.Context, ns string) (*TestClients, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, err
	}
	cache := kube.NewKubeCache(ctx)

	// upstreams
	upstreamClientFactory := &factory.KubeResourceClientFactory{
		Crd:         gloov1.UpstreamCrd,
		Cfg:         cfg,
		SharedCache: cache,
	}
	upstreamClient, err := gloov1.NewUpstreamClient(ctx, upstreamClientFactory)
	if err != nil {
		return nil, err
	}
	err = upstreamClient.Register()
	if err != nil {
		return nil, err
	}

	// virtual services
	virtualServiceClientFactory := &factory.KubeResourceClientFactory{
		Crd:         gatewayv1.VirtualServiceCrd,
		Cfg:         cfg,
		SharedCache: cache,
	}
	virtualServiceClient, err := gatewayv1.NewVirtualServiceClient(ctx, virtualServiceClientFactory)
	if err != nil {
		return nil, err
	}
	err = virtualServiceClient.Register()
	if err != nil {
		return nil, err
	}

	statusClient := statusutils.GetStatusClientFromEnvOrDefault(ns)

	return &TestClients{
		UpstreamClient:       upstreamClient,
		VirtualServiceClient: virtualServiceClient,
		StatusClient:         statusClient,
	}, nil
}

func NewKubeClient() (kubernetes.Interface, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
}
