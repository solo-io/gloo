package converter_test

import (
	"context"

	"github.com/gogo/protobuf/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth"
	ratelimitapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/ratelimit"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/util"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/converter"
)

var (
	virtualServiceConverter converter.VirtualServiceDetailsConverter
	mockCtrl                *gomock.Controller
	rawGetter               *mocks.MockRawGetter
)

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

	getExpectedRaw := func() *v1.Raw {
		return &v1.Raw{FileName: "fn", Content: "ct"}
	}

	Describe("List", func() {
		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			rawGetter = mocks.NewMockRawGetter(mockCtrl)
			virtualServiceConverter = converter.NewVirtualServiceDetailsConverter(rawGetter)
		})

		AfterEach(func() {
			mockCtrl.Finish()
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
				expectedRaw     *v1.Raw
			}{
				{
					desc:            "for a nil config",
					configs:         nil,
					expectedPlugins: nil,
					expectedRaw:     getExpectedRaw(),
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
					expectedRaw: getExpectedRaw(),
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
					expectedRaw: getExpectedRaw(),
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
					expectedRaw: getExpectedRaw(),
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
					expectedRaw: getExpectedRaw(),
				},
			} {
				virtualService := getVirtualService(testCase.configs)
				if testCase.expectedRaw != nil {
					rawGetter.EXPECT().
						GetRaw(context.Background(), virtualService, gatewayv1.VirtualServiceCrd).
						Return(testCase.expectedRaw)
				}
				actual := virtualServiceConverter.GetDetails(context.Background(), virtualService)
				expected := &v1.VirtualServiceDetails{
					VirtualService: virtualService,
					Plugins:        testCase.expectedPlugins,
					Raw:            testCase.expectedRaw,
				}
				ExpectEqualProtoMessages(actual, expected, testCase.desc)
			}
		})
	})
})
