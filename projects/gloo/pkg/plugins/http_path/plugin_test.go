package http_path_test

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	gogotypes "github.com/gogo/protobuf/types"
	wrapperspb "github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pbhttp_path "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/http_path"
	envoy_type_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/http_path"
)

var _ = Describe("Plugin", func() {

	var (
		p               *Plugin
		params          plugins.Params
		upstream        *v1.Upstream
		upstreamSpec    *v1static.UpstreamSpec
		out             *envoyapi.Cluster
		baseHealthCheck *envoy_api_v2_core.HealthCheck_HttpHealthCheck
	)

	BeforeEach(func() {
		p = NewPlugin()
		out = new(envoyapi.Cluster)
		baseHealthCheck = &envoy_api_v2_core.HealthCheck_HttpHealthCheck{
			UseHttp2:               true,
			Host:                   "foo",
			Path:                   "/health",
			ServiceName:            "svc",
			RequestHeadersToRemove: []string{"test"},
			RequestHeadersToAdd: []*envoy_api_v2_core.HeaderValueOption{
				&envoy_api_v2_core.HeaderValueOption{
					Header: &envoy_api_v2_core.HeaderValue{Key: "key", Value: "value"},
					Append: &wrapperspb.BoolValue{Value: true},
				},
			},
		}
		out.HealthChecks = []*envoy_api_v2_core.HealthCheck{
			{
				HealthChecker: &envoy_api_v2_core.HealthCheck_HttpHealthCheck_{

					HttpHealthCheck: baseHealthCheck,
				},
			},
		}

		p.Init(plugins.InitParams{})
		upstreamSpec = &v1static.UpstreamSpec{
			Hosts: []*v1static.Host{{
				Addr: "localhost",
				Port: 1234,
				HealthCheckConfig: &v1static.Host_HealthCheckConfig{
					Path: "/foo",
				},
			}},
		}
		upstream = &v1.Upstream{
			Metadata: core.Metadata{
				Name:      "extauth-server",
				Namespace: "default",
			},
			UpstreamType: &v1.Upstream_Static{
				Static: upstreamSpec,
			},
		}

	})

	It("should create a custom health check when static upstream has a path", func() {
		p.ProcessUpstream(params, upstream, out)
		Expect(out.GetHealthChecks()[0].GetCustomHealthCheck().GetName()).To(Equal(HealthCheckerName))
		typedcfg := out.GetHealthChecks()[0].GetCustomHealthCheck().GetTypedConfig()
		var out pbhttp_path.HttpPath
		gogotypes.UnmarshalAny(&gogotypes.Any{
			TypeUrl: typedcfg.TypeUrl,
			Value:   typedcfg.Value,
		}, &out)

		Expect(out.HttpHealthCheck.Path).To(Equal(baseHealthCheck.Path))
		Expect(out.HttpHealthCheck.Host).To(Equal(baseHealthCheck.Host))
		Expect(out.HttpHealthCheck.CodecClientType).To(Equal(envoy_type_v3.CodecClientType_HTTP2))
		Expect(out.HttpHealthCheck.GetServiceNameMatcher().GetPrefix()).To(Equal(baseHealthCheck.ServiceName))
		Expect(out.HttpHealthCheck.RequestHeadersToRemove).To(Equal(baseHealthCheck.RequestHeadersToRemove))
		Expect(out.HttpHealthCheck.RequestHeadersToAdd[0].Header.Key).To(Equal(baseHealthCheck.RequestHeadersToAdd[0].Header.Key))
		Expect(out.HttpHealthCheck.RequestHeadersToAdd[0].Header.Value).To(Equal(baseHealthCheck.RequestHeadersToAdd[0].Header.Value))
		Expect(out.HttpHealthCheck.RequestHeadersToAdd[0].Append.Value).To(Equal(baseHealthCheck.RequestHeadersToAdd[0].Append.Value))

	})

})
