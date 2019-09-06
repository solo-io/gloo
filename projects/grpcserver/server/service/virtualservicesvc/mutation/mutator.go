package mutation

import (
	"context"

	"github.com/solo-io/solo-projects/pkg/license"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/svccodes"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type Mutation func(vs *gatewayv1.VirtualService) error

//go:generate mockgen -destination mocks/mutator_mock.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mutation Mutator

type Mutator interface {
	Create(ref *core.ResourceRef, f Mutation) (*gatewayv1.VirtualService, error)
	Update(ref *core.ResourceRef, f Mutation) (*gatewayv1.VirtualService, error)
}

type mutator struct {
	ctx           context.Context
	client        gatewayv1.VirtualServiceClient
	licenseClient license.Client
}

func (m *mutator) Create(ref *core.ResourceRef, f Mutation) (*gatewayv1.VirtualService, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(m.ctx, m.licenseClient); err != nil {
		return nil, err
	}
	virtualService := &gatewayv1.VirtualService{
		Metadata: core.Metadata{
			Namespace: ref.GetNamespace(),
			Name:      ref.GetName(),
		},
	}
	return m.mutateAndWrite(virtualService, f, false)
}

func (m *mutator) Update(ref *core.ResourceRef, f Mutation) (*gatewayv1.VirtualService, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(m.ctx, m.licenseClient); err != nil {
		return nil, err
	}
	virtualService, err := m.client.Read(ref.GetNamespace(), ref.GetName(), clients.ReadOpts{Ctx: m.ctx})
	if err != nil {
		return nil, err
	}
	return m.mutateAndWrite(virtualService, f, true)
}

func (m *mutator) mutateAndWrite(vs *gatewayv1.VirtualService, f Mutation, overwrite bool) (*gatewayv1.VirtualService, error) {
	if err := f(vs); err != nil {
		return nil, err
	}
	vs.Status = core.Status{}
	return m.client.Write(vs, clients.WriteOpts{Ctx: m.ctx, OverwriteExisting: overwrite})
}

func NewMutator(ctx context.Context, client gatewayv1.VirtualServiceClient, licenseClient license.Client) Mutator {
	return &mutator{
		ctx:           ctx,
		client:        client,
		licenseClient: licenseClient,
	}
}
