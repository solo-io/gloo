package test_fixtures

import (
	"fmt"
	"net"
	"net/http"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

const (
	virtualHostNameFmt = "vhost-%d"
	routeNameFmt       = "route-%d"
	proxyNameFmt       = "proxy-%d"
	authConfigNameFmt  = "ac-%d"
)

// AuthConfigSlice returns a list of AuthConfig resources which will be used in benchmark testing
// Because we are testing efficiency, we create many simple AuthConfigs, since the complexity of the AuthConfigs
// themselves is not important
func AuthConfigSlice(namespace string, numberOfAuthConfigs int) extauth.AuthConfigList {
	var authConfigs extauth.AuthConfigList

	for i := 0; i < numberOfAuthConfigs; i++ {
		authConfig := BasicAuthConfig(fmt.Sprintf(authConfigNameFmt, i), namespace)
		authConfigs = append(authConfigs, authConfig)
	}

	return authConfigs
}

// BasicAuthConfig returns an AuthConfig with a single definition for a Jwt AuthService
func BasicAuthConfig(name, namespace string) *extauth.AuthConfig {
	return &extauth.AuthConfig{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Configs: []*extauth.AuthConfig_Config{
			{
				Name: &wrappers.StringValue{Value: "jwt"},
				AuthConfig: &extauth.AuthConfig_Config_Jwt{
					Jwt: &empty.Empty{},
				},
			},
		},
		BooleanExpr: nil,
	}
}

// InvalidAuthConfig returns an AuthConfig with invalid configuration that the Control Plane should reject
func InvalidAuthConfig(name, namespace string) *extauth.AuthConfig {
	return &extauth.AuthConfig{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Configs: []*extauth.AuthConfig_Config{
			{
				Name: &wrappers.StringValue{Value: "passthrough"},
				AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
					PassThroughAuth: &extauth.PassThroughAuth{
						Protocol: &extauth.PassThroughAuth_Grpc{
							Grpc: &extauth.PassThroughGrpc{
								Address: "", // Empty address is invalid
							},
						},
					},
				},
			},
		},
		BooleanExpr: nil,
	}
}

// ComplexAuthConfig returns an AuthConfig with multiple configurations defined
// This is used to validate that AuthConfigs can be shuffled and returned in a consistent order
func ComplexAuthConfig(name, namespace string) *extauth.AuthConfig {
	return &extauth.AuthConfig{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		BooleanExpr: &wrappers.StringValue{Value: "oauth2 || passthrough || jwt"},
		Configs: []*extauth.AuthConfig_Config{
			{
				Name: &wrappers.StringValue{Value: "oauth2"},
				AuthConfig: &extauth.AuthConfig_Config_Oauth2{
					Oauth2: &extauth.OAuth2{
						OauthType: &extauth.OAuth2_AccessTokenValidation{
							AccessTokenValidation: &extauth.AccessTokenValidation{
								UserinfoUrl: fmt.Sprintf("fake-url-%s", name),
							},
						},
					},
				},
			},
			{
				Name: &wrappers.StringValue{Value: "passthrough"},
				AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
					PassThroughAuth: &extauth.PassThroughAuth{
						Protocol: &extauth.PassThroughAuth_Grpc{
							Grpc: &extauth.PassThroughGrpc{
								Address: fmt.Sprintf("passthrough-%s", name),
							},
						},
					},
				},
			},
			{
				Name: &wrappers.StringValue{Value: "jwt"},
				AuthConfig: &extauth.AuthConfig_Config_Jwt{
					Jwt: &empty.Empty{},
				},
			},
		},
	}
}

type ResourceFrequency struct {
	AuthConfigs          int
	Proxies              int
	VirtualHostsPerProxy int
	RoutesPerVirtualHost int
}

func ProxySlice(namespace string, frequency ResourceFrequency) v1.ProxyList {
	var proxies v1.ProxyList

	for i := 0; i < frequency.Proxies; i++ {
		proxy := BasicProxy(fmt.Sprintf(proxyNameFmt, i), namespace, frequency.VirtualHostsPerProxy, frequency.RoutesPerVirtualHost)
		proxies = append(proxies, proxy)
	}

	return proxies
}

func BasicProxy(name, namespace string, numberOfVirtualServices, numberOfRoutes int) *v1.Proxy {
	var virtualHosts []*v1.VirtualHost

	for i := 0; i < numberOfVirtualServices; i++ {
		vhostName := fmt.Sprintf(virtualHostNameFmt, i)
		vhost := &v1.VirtualHost{
			Name:    vhostName,
			Domains: []string{vhostName},
			Routes:  []*v1.Route{},
			Options: nil,
		}

		for j := 0; j < numberOfRoutes; j++ {
			route := &v1.Route{
				Name: fmt.Sprintf("%s-%s-%s", name, vhostName, fmt.Sprintf(routeNameFmt, j)),
				Action: &v1.Route_DirectResponseAction{
					DirectResponseAction: &v1.DirectResponseAction{
						Status: http.StatusOK,
					},
				},
				Options: &v1.RouteOptions{
					Extauth: &extauth.ExtAuthExtension{
						Spec: &extauth.ExtAuthExtension_ConfigRef{
							ConfigRef: &core.ResourceRef{
								Name:      fmt.Sprintf(authConfigNameFmt, j),
								Namespace: namespace,
							},
						},
					},
				},
			}
			vhost.Routes = append(vhost.Routes, route)
		}
		virtualHosts = append(virtualHosts, vhost)
	}

	return &v1.Proxy{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Listeners: []*v1.Listener{{
			Name:        "listener",
			BindAddress: net.IPv4zero.String(),
			BindPort:    defaults.HttpPort,
			ListenerType: &v1.Listener_HttpListener{
				HttpListener: &v1.HttpListener{
					VirtualHosts: virtualHosts,
				},
			},
		}},
	}

}

func ProxyWithExtAuthExtensionOnRoute(name, namespace string, extAuthExtension *extauth.ExtAuthExtension) *v1.Proxy {
	return &v1.Proxy{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Listeners: []*v1.Listener{{
			Name:        "listener",
			BindAddress: net.IPv4zero.String(),
			BindPort:    defaults.HttpPort,
			ListenerType: &v1.Listener_HttpListener{
				HttpListener: &v1.HttpListener{
					VirtualHosts: []*v1.VirtualHost{{
						Name:    "virtual-host",
						Domains: []string{"virtual-host"},
						Routes: []*v1.Route{{
							Name: "route",
							Action: &v1.Route_DirectResponseAction{
								DirectResponseAction: &v1.DirectResponseAction{
									Status: http.StatusOK,
								},
							},
							Options: &v1.RouteOptions{
								Extauth: extAuthExtension,
							},
						}},
					}},
				},
			},
		}},
	}
}
