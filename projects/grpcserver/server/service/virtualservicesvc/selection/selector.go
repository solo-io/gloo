package selection

import (
	"context"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	sk_errors "github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/kube"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mocks/selector_mock.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/selection VirtualServiceSelector

type VirtualServiceSelector interface {
	SelectOrCreate(ctx context.Context, ref *core.ResourceRef) (*gatewayv1.VirtualService, error)
}

type virtualServiceSelector struct {
	virtualServiceClient gatewayv1.VirtualServiceClient
	namespaceClient      kube.NamespaceClient
	podNamespace         string
}

var _ VirtualServiceSelector = &virtualServiceSelector{}

func NewVirtualServiceSelector(virtualServiceClient gatewayv1.VirtualServiceClient, namespaceClient kube.NamespaceClient, podNamespace string) VirtualServiceSelector {
	return &virtualServiceSelector{
		virtualServiceClient: virtualServiceClient,
		namespaceClient:      namespaceClient,
		podNamespace:         podNamespace,
	}
}

func (s *virtualServiceSelector) SelectOrCreate(ctx context.Context, ref *core.ResourceRef) (*gatewayv1.VirtualService, error) {
	// Read or create virtual service as specified
	if ref.GetNamespace() != "" && ref.GetName() != "" {
		found, err := s.virtualServiceClient.Read(ref.GetNamespace(), ref.GetName(), clients.ReadOpts{Ctx: ctx})
		if err != nil && !sk_errors.IsNotExist(err) {
			return nil, err
		}
		if found != nil {
			return found, nil
		}
	}

	// Create a new default virtual service with the given name
	if ref.GetNamespace() != "" || ref.GetName() != "" {
		return s.create(ctx, ref)
	}

	// Look for an existing virtual service with * domain
	namespaces, err := s.namespaceClient.ListNamespaces()
	if err != nil {
		return nil, err
	}
	for _, ns := range namespaces {
		allVirtualServices, err := s.virtualServiceClient.List(ns, clients.ListOpts{Ctx: ctx})
		if err != nil {
			return nil, err
		}
		for _, vs := range allVirtualServices {
			for _, domain := range vs.GetVirtualHost().GetDomains() {
				if domain == "*" {
					contextutils.LoggerFrom(ctx).Infow("Selected virtual service with domain *", zap.Any("virtualService", vs))
					return vs, nil
				}
			}
		}
	}

	// Create a new default virtual service
	return s.create(ctx, ref)
}

func (s *virtualServiceSelector) create(ctx context.Context, ref *core.ResourceRef) (*gatewayv1.VirtualService, error) {
	virtualService := &gatewayv1.VirtualService{
		Metadata: core.Metadata{
			Namespace: ref.GetNamespace(),
			Name:      ref.GetName(),
		},
		VirtualHost: &gatewayv1.VirtualHost{
			Domains: []string{"*"},
		},
	}
	if virtualService.GetMetadata().Namespace == "" {
		virtualService.Metadata.Namespace = s.podNamespace
	}
	if virtualService.GetMetadata().Name == "" {
		virtualService.Metadata.Name = "default"
	}

	written, err := s.virtualServiceClient.Write(virtualService, clients.WriteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	contextutils.LoggerFrom(ctx).Infow("Created new default virtual service", zap.Any("virtualService", virtualService))
	return written, nil

}
