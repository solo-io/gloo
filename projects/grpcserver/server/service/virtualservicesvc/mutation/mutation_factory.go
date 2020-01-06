package mutation

import (
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/ratelimit"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

//go:generate mockgen -destination mocks/mutation_factory_mock.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mutation MutationFactory

type MutationFactory interface {
	// Deprecated
	ConfigureVirtualServiceV2(input *v1.VirtualServiceInputV2) Mutation
	CreateRoute(input *v1.RouteInput) Mutation
	UpdateRoute(input *v1.RouteInput) Mutation
	DeleteRoute(index uint32) Mutation
	SwapRoutes(index1, index2 uint32) Mutation
	ShiftRoutes(fromIndex, toIndex uint32) Mutation
}

type mutationFactory struct{}

func NewMutationFactory() MutationFactory {
	return &mutationFactory{}
}

// Only sets fields that are non-nil in the input to allow for delta-style updates.
func (*mutationFactory) ConfigureVirtualServiceV2(input *v1.VirtualServiceInputV2) Mutation {
	return func(vs *gatewayv1.VirtualService) error {
		// Only set metadata if this is a new Virtual Service
		if vs.GetMetadata().Namespace == "" {
			vs.Metadata.Namespace = input.GetRef().GetNamespace()
			vs.Metadata.Name = input.GetRef().GetName()
		}

		if input.GetSslConfig() != nil {
			vs.SslConfig = input.GetSslConfig().GetValue()
		}

		if vs.GetVirtualHost() == nil {
			vs.VirtualHost = &gatewayv1.VirtualHost{}
		}

		if input.GetRateLimitConfig() != nil {
			if vs.GetVirtualHost().GetOptions() == nil {
				vs.VirtualHost.Options = &gloov1.VirtualHostOptions{}
			}
			if input.GetRateLimitConfig() != nil {
				if vs.VirtualHost.Options.GetExtensions().GetConfigs() != nil {
					delete(vs.VirtualHost.Options.Extensions.Configs, ratelimit.ExtensionName)
				}
				vs.VirtualHost.Options.RatelimitBasic = input.GetRateLimitConfig().GetValue()
			}
		}

		if input.GetDisplayName() != nil {
			vs.DisplayName = input.GetDisplayName().GetValue()
		}

		if input.GetDomains() != nil {
			vs.VirtualHost.Domains = input.GetDomains().GetValues()
		}

		if input.GetRoutes() != nil {
			vs.VirtualHost.Routes = input.GetRoutes().GetValues()
		}

		return nil
	}
}

func (*mutationFactory) CreateRoute(input *v1.RouteInput) Mutation {
	return func(vs *gatewayv1.VirtualService) error {
		if input.GetRoute() == nil {
			return NoRouteProvidedError
		}

		index := int(input.GetIndex())
		if index > len(vs.GetVirtualHost().GetRoutes()) {
			index = len(vs.GetVirtualHost().GetRoutes())
		}
		vs.VirtualHost.Routes = append(vs.VirtualHost.Routes, nil)
		copy(vs.VirtualHost.Routes[index+1:], vs.VirtualHost.Routes[index:])
		vs.VirtualHost.Routes[index] = input.GetRoute()
		return nil
	}
}

func (*mutationFactory) UpdateRoute(input *v1.RouteInput) Mutation {
	return func(vs *gatewayv1.VirtualService) error {
		if input.GetRoute() == nil {
			return NoRouteProvidedError
		}

		index := int(input.GetIndex())
		if index > len(vs.VirtualHost.Routes)-1 {
			return IndexOutOfBoundsError
		}
		vs.VirtualHost.Routes[index] = input.GetRoute()
		return nil
	}
}

func (*mutationFactory) DeleteRoute(index uint32) Mutation {
	return func(vs *gatewayv1.VirtualService) error {
		if int(index) > len(vs.VirtualHost.Routes)-1 {
			return IndexOutOfBoundsError
		}
		vs.VirtualHost.Routes = append(vs.VirtualHost.Routes[:index], vs.VirtualHost.Routes[index+1:]...)
		return nil
	}
}

func (*mutationFactory) SwapRoutes(index1, index2 uint32) Mutation {
	return func(vs *gatewayv1.VirtualService) error {
		if int(index1) > len(vs.VirtualHost.Routes)-1 || int(index2) > len(vs.VirtualHost.Routes)-1 {
			return IndexOutOfBoundsError
		}
		vs.VirtualHost.Routes[index1], vs.VirtualHost.Routes[index2] = vs.VirtualHost.Routes[index2], vs.VirtualHost.Routes[index1]
		return nil
	}
}

func (*mutationFactory) ShiftRoutes(fromIndex, toIndex uint32) Mutation {
	return func(vs *gatewayv1.VirtualService) error {
		if int(fromIndex) > len(vs.VirtualHost.Routes)-1 || int(toIndex) > len(vs.VirtualHost.Routes)-1 {
			return IndexOutOfBoundsError
		}

		if toIndex < fromIndex {
			// anchor on the fromIndex and swap until all updated
			for i := toIndex; i < fromIndex; i++ {
				vs.VirtualHost.Routes[fromIndex], vs.VirtualHost.Routes[i] = vs.VirtualHost.Routes[i], vs.VirtualHost.Routes[fromIndex]
			}
		} else {
			// anchor on the toIndex and swap until all updated
			for i := toIndex; i > fromIndex; i-- {
				vs.VirtualHost.Routes[fromIndex], vs.VirtualHost.Routes[i] = vs.VirtualHost.Routes[i], vs.VirtualHost.Routes[fromIndex]
			}
		}

		return nil
	}
}
