package mutation

import (
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/util"

	//TODO: (Graham) handle plugins correclty once the reorg happens "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

//go:generate mockgen -destination mocks/mutation_factory_mock.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mutation MutationFactory

type MutationFactory interface {
	// Deprecated
	ConfigureVirtualService(input *v1.VirtualServiceInput) Mutation
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

func (*mutationFactory) ConfigureVirtualService(input *v1.VirtualServiceInput) Mutation {
	return func(vs *gatewayv1.VirtualService) error {
		// Only set metadata if this is a new Virtual Service
		if vs.GetMetadata().Namespace == "" {
			vs.Metadata.Namespace = input.GetRef().GetNamespace()
			vs.Metadata.Name = input.GetRef().GetName()
		}

		// Convert external auth config into type expected for extensions
		var extAuthStruct *types.Struct
		var err error
		if input.GetExtAuthConfig() != nil {
			switch t := input.ExtAuthConfig.(type) {
			case *v1.VirtualServiceInput_BasicAuth:
				return errors.Errorf("Basic auth is not supported.")
			case *v1.VirtualServiceInput_Oauth:
				extAuthStruct, err = getOauthVhostExtensionStruct(t.Oauth)
				if err != nil {
					return err
				}
			case *v1.VirtualServiceInput_CustomAuth:
				extAuthStruct, err = getCustomAuthVhostExtensionStruct(t.CustomAuth)
				if err != nil {
					return err
				}
			}
		}

		// Attempt to set secret ref -- error if there is a different SSL strategy in place.
		if input.GetSecretRef() != nil {
			if vs.SslConfig == nil {
				vs.SslConfig = &gloov1.SslConfig{}
			}

			switch vs.SslConfig.GetSslSecrets().(type) {
			case *gloov1.SslConfig_SslFiles:
				return AlreadyConfiguredSslWithFiles
			case *gloov1.SslConfig_Sds:
				return AlreadyConfiguredSslWithSds
			case *gloov1.SslConfig_SecretRef:
				vs.SslConfig.SslSecrets = &gloov1.SslConfig_SecretRef{SecretRef: input.GetSecretRef()}
			default:
				vs.SslConfig.SslSecrets = &gloov1.SslConfig_SecretRef{SecretRef: input.GetSecretRef()}
			}
		}

		if vs.GetVirtualHost() == nil {
			vs.VirtualHost = &gatewayv1.VirtualHost{}
		}

		if extAuthStruct != nil || input.GetRateLimitConfig() != nil {
			if vs.GetVirtualHost().GetVirtualHostPlugins() == nil {
				vs.VirtualHost.VirtualHostPlugins = &gloov1.VirtualHostPlugins{}
			}
			if vs.GetVirtualHost().GetVirtualHostPlugins().GetExtensions() == nil {
				vs.VirtualHost.VirtualHostPlugins.Extensions = &gloov1.Extensions{}
			}
			if vs.GetVirtualHost().GetVirtualHostPlugins().GetExtensions().GetConfigs() == nil {
				vs.VirtualHost.VirtualHostPlugins.Extensions.Configs = make(map[string]*types.Struct)
			}

			if extAuthStruct != nil {
				vs.VirtualHost.VirtualHostPlugins.Extensions.Configs[extauth.ExtensionName] = extAuthStruct
			}
			if input.GetRateLimitConfig() != nil {
				delete(vs.VirtualHost.VirtualHostPlugins.Extensions.Configs, ratelimit.ExtensionName)
				vs.VirtualHost.VirtualHostPlugins.RatelimitBasic = input.GetRateLimitConfig()
			}
		}

		vs.DisplayName = input.GetDisplayName()
		vs.VirtualHost.Domains = input.GetDomains()
		vs.VirtualHost.Routes = input.GetRoutes()
		return nil
	}
}

// Only sets fields that are non-nil in the input to allow for delta-style updates.
func (*mutationFactory) ConfigureVirtualServiceV2(input *v1.VirtualServiceInputV2) Mutation {
	return func(vs *gatewayv1.VirtualService) error {
		// Only set metadata if this is a new Virtual Service
		if vs.GetMetadata().Namespace == "" {
			vs.Metadata.Namespace = input.GetRef().GetNamespace()
			vs.Metadata.Name = input.GetRef().GetName()
		}

		// Convert external auth config into type expected for extensions
		var extAuthStruct *types.Struct
		var err error
		if input.GetExtAuthConfig() != nil {
			if input.GetExtAuthConfig().GetConfig() != nil {
				switch t := input.ExtAuthConfig.Config.Value.(type) {
				case *v1.ExtAuthInput_Config_Oauth:
					extAuthStruct, err = getOauthVhostExtensionStruct(t.Oauth)
					if err != nil {
						return err
					}
				case *v1.ExtAuthInput_Config_CustomAuth:
					extAuthStruct, err = getCustomAuthVhostExtensionStruct(t.CustomAuth)
					if err != nil {
						return err
					}
				}
			}
		}

		if input.GetSslConfig() != nil {
			vs.SslConfig = input.GetSslConfig().GetValue()
		}

		if vs.GetVirtualHost() == nil {
			vs.VirtualHost = &gatewayv1.VirtualHost{}
		}

		if input.GetExtAuthConfig() != nil || input.GetRateLimitConfig() != nil {
			if vs.GetVirtualHost().GetVirtualHostPlugins() == nil {
				vs.VirtualHost.VirtualHostPlugins = &gloov1.VirtualHostPlugins{}
			}
			if vs.GetVirtualHost().GetVirtualHostPlugins().GetExtensions() == nil {
				vs.VirtualHost.VirtualHostPlugins.Extensions = &gloov1.Extensions{}
			}
			if vs.GetVirtualHost().GetVirtualHostPlugins().GetExtensions().GetConfigs() == nil {
				vs.VirtualHost.VirtualHostPlugins.Extensions.Configs = make(map[string]*types.Struct)
			}

			if input.GetExtAuthConfig() != nil {
				if extAuthStruct == nil {
					delete(vs.VirtualHost.VirtualHostPlugins.Extensions.Configs, extauth.ExtensionName)
				} else {
					vs.VirtualHost.VirtualHostPlugins.Extensions.Configs[extauth.ExtensionName] = extAuthStruct
				}
			}
			if input.GetRateLimitConfig() != nil {
				delete(vs.VirtualHost.VirtualHostPlugins.Extensions.Configs, ratelimit.ExtensionName)
				vs.VirtualHost.VirtualHostPlugins.RatelimitBasic = input.GetRateLimitConfig().GetValue()
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

func getOauthVhostExtensionStruct(oauth *extauthapi.OAuth) (*types.Struct, error) {
	return util.MessageToStruct(&extauthapi.VhostExtension{
		Configs: []*extauthapi.VhostExtension_AuthConfig{{
			AuthConfig: &extauthapi.VhostExtension_AuthConfig_Oauth{
				Oauth: oauth,
			},
		}},
	})
}

func getCustomAuthVhostExtensionStruct(customAuth *extauthapi.CustomAuth) (*types.Struct, error) {
	return util.MessageToStruct(&extauthapi.VhostExtension{
		Configs: []*extauthapi.VhostExtension_AuthConfig{{
			AuthConfig: &extauthapi.VhostExtension_AuthConfig_CustomAuth{
				CustomAuth: customAuth,
			},
		}},
	})
}
