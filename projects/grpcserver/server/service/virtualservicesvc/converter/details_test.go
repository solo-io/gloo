package converter_test

import (
	"context"

	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/util"
	extauthapi "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	ratelimitapi "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/converter"
)

var virtualServiceConverter converter.VirtualServiceDetailsConverter

var _ = Describe("VirtualServiceDetailsConverter", func() {
	getVirtualService := func(pluginConfigs map[string]*types.Struct) *gatewayv1.VirtualService {
		return &gatewayv1.VirtualService{
			VirtualHost: &gloov1.VirtualHost{
				VirtualHostPlugins: &gloov1.VirtualHostPlugins{
					Extensions: &gloov1.Extensions{
						Configs: pluginConfigs,
					},
				},
			},
		}
	}

	getExtAuthConfig := func() *extauthapi.ExtAuthConfig {
		return &extauthapi.ExtAuthConfig{
			Vhost: "test-host",
			AuthConfig: &extauthapi.ExtAuthConfig_Oauth{
				Oauth: &extauthapi.ExtAuthConfig_OAuthConfig{
					ClientId: "test-client-id",
				},
			},
		}
	}

	getRateLimit := func() *ratelimitapi.IngressRateLimit {
		return &ratelimitapi.IngressRateLimit{
			AnonymousLimits: &ratelimitapi.RateLimit{
				Unit:            ratelimitapi.RateLimit_DAY,
				RequestsPerUnit: 10,
			},
		}
	}

	getInvalidStruct := func() *types.Struct {
		return &types.Struct{
			Fields: map[string]*types.Value{
				"foo": {
					Kind: &types.Value_StringValue{
						StringValue: "bar",
					},
				},
			},
		}
	}

	Describe("GetDetails", func() {
		BeforeEach(func() {
			virtualServiceConverter = converter.NewVirtualServiceDetailsConverter()
		})

		It("works", func() {
			extAuthConfig := getExtAuthConfig()
			extAuthStruct, err := util.MessageToStruct(extAuthConfig)
			Expect(err).NotTo(HaveOccurred())
			rateLimit := getRateLimit()
			rateLimitStruct, err := util.MessageToStruct(rateLimit)
			Expect(err).NotTo(HaveOccurred())

			for _, testCase := range []struct {
				desc            string
				configs         map[string]*types.Struct
				expectedPlugins *v1.Plugins
			}{
				{
					desc:            "for a nil config",
					configs:         nil,
					expectedPlugins: nil,
				},
				{
					desc: "for a valid extauth plugin",
					configs: map[string]*types.Struct{
						extauth.ExtensionName: extAuthStruct,
					},
					expectedPlugins: &v1.Plugins{
						ExtAuth: &v1.ExtAuthPlugin{
							Value: extAuthConfig,
						},
					},
				},
				{
					desc: "for a valid rate limit plugin",
					configs: map[string]*types.Struct{
						ratelimit.ExtensionName: rateLimitStruct,
					},
					expectedPlugins: &v1.Plugins{
						RateLimit: &v1.RateLimitPlugin{
							Value: rateLimit,
						},
					},
				},
				{
					desc: "for an invalid extauth plugin",
					configs: map[string]*types.Struct{
						extauth.ExtensionName: getInvalidStruct(),
					},
					expectedPlugins: &v1.Plugins{
						ExtAuth: &v1.ExtAuthPlugin{
							Error: converter.FailedToParseExtAuthConfig,
						},
					},
				},
				{
					desc: "for an invalid rate limit plugin",
					configs: map[string]*types.Struct{
						ratelimit.ExtensionName: getInvalidStruct(),
					},
					expectedPlugins: &v1.Plugins{
						RateLimit: &v1.RateLimitPlugin{
							Error: converter.FailedToParseRateLimitConfig,
						},
					},
				},
			} {
				virtualService := getVirtualService(testCase.configs)
				actual := virtualServiceConverter.GetDetails(context.Background(), virtualService)
				expected := &v1.VirtualServiceDetails{
					VirtualService: virtualService,
					Plugins:        testCase.expectedPlugins,
				}
				ExpectEqualProtoMessages(actual, expected, testCase.desc)
			}
		})
	})
})
