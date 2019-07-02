package mutation

import (
	"context"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

type Mutation func(vs *gatewayv1.VirtualService) error

//go:generate mockgen -destination mocks/mutator_mock.go -self_package github.com/solo-io/gloo/projects/gateway/pkg/api/v1 -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mutation Mutator

type Mutator interface {
	Create(input *v1.VirtualServiceInput, f Mutation) (*gatewayv1.VirtualService, error)
	Update(ref *core.ResourceRef, f Mutation) (*gatewayv1.VirtualService, error)
}

type mutator struct {
	ctx    context.Context
	client gatewayv1.VirtualServiceClient
}

func (m *mutator) Create(input *v1.VirtualServiceInput, f Mutation) (*gatewayv1.VirtualService, error) {
	virtualService := &gatewayv1.VirtualService{
		Metadata: core.Metadata{
			Namespace: input.GetRef().GetNamespace(),
			Name:      input.GetRef().GetName(),
		},
	}

	if err := f(virtualService); err != nil {
		return nil, err
	}
	return m.client.Write(virtualService, clients.WriteOpts{Ctx: m.ctx, OverwriteExisting: false})
}

func (m *mutator) Update(ref *core.ResourceRef, f Mutation) (*gatewayv1.VirtualService, error) {
	virtualService, err := m.client.Read(ref.GetNamespace(), ref.GetName(), clients.ReadOpts{Ctx: m.ctx})
	if err != nil {
		return nil, err
	}
	if err := f(virtualService); err != nil {
		return nil, err
	}
	return m.client.Write(virtualService, clients.WriteOpts{Ctx: m.ctx, OverwriteExisting: true})
}

func NewMutator(ctx context.Context, client gatewayv1.VirtualServiceClient) *mutator {
	return &mutator{
		ctx:    ctx,
		client: client,
	}
}
