package mutation

import (
	"context"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type Mutation func(upstream *gloov1.Upstream) error

//go:generate mockgen -destination mocks/mutator_mock.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/mutation Mutator

type Mutator interface {
	Create(ctx context.Context, ref *core.ResourceRef, f Mutation) (*gloov1.Upstream, error)
	Update(ctx context.Context, ref *core.ResourceRef, f Mutation) (*gloov1.Upstream, error)
}

type mutator struct {
	client gloov1.UpstreamClient
}

var _ Mutator = &mutator{}

func NewMutator(client gloov1.UpstreamClient) Mutator {
	return &mutator{client: client}
}

func (m *mutator) Create(ctx context.Context, ref *core.ResourceRef, f Mutation) (*gloov1.Upstream, error) {
	upstream := &gloov1.Upstream{
		Metadata: core.Metadata{
			Namespace: ref.GetNamespace(),
			Name:      ref.GetName(),
		},
	}
	return m.mutateAndWrite(ctx, upstream, f, false)
}

func (m *mutator) Update(ctx context.Context, ref *core.ResourceRef, f Mutation) (*gloov1.Upstream, error) {
	virtualService, err := m.client.Read(ref.GetNamespace(), ref.GetName(), clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return m.mutateAndWrite(ctx, virtualService, f, true)
}

func (m *mutator) mutateAndWrite(ctx context.Context, upstream *gloov1.Upstream, f Mutation, overwrite bool) (*gloov1.Upstream, error) {
	if err := f(upstream); err != nil {
		return nil, err
	}
	upstream.Status = core.Status{}
	return m.client.Write(upstream, clients.WriteOpts{Ctx: ctx, OverwriteExisting: overwrite})
}
