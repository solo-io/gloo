package graphql_test

import (
	"context"
	"fmt"
	"testing"

	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/gogo/protobuf/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/aws"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/mocks"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
)

var (
	newPrefix          = "/new/prefix/"
	route1PathMatch    = "/some-path"
	upstream1Name      = "us-1"
	upstream1Namespace = "us-ns-1"
)

var T *testing.T

func TestConverter(t *testing.T) {
	RegisterFailHandler(Fail)
	T = t
	RunSpecs(t, "Converter Test Suite")
}

var _ = Describe("Converter", func() {

	var (
		upstreamClient *mocks.MockUpstreamClient
		converter      *graphql.Converter
	)

	BeforeEach(func() {
		ctrl := gomock.NewController(T)
		defer ctrl.Finish()

		// Get mock upstream client
		upstreamClient = mocks.NewMockUpstreamClient(ctrl)

		// Set the return value. We need this for the Proto -> GraphQL conversion. The proto route holds a ResourceRef
		// to the Upstream and the converter tries to resolve the dependency by retrieving the correspondent resource.
		upstreamClient.EXPECT().
			Read(upstream1Namespace, upstream1Name, gomock.Any()).
			Return(&PROTO.upstream, nil).
			AnyTimes()

		// We only need the upstream client
		resolvers := graphql.NewResolvers(
			upstreamClient,
			gloov1.ArtifactClient(nil),
			gloov1.SettingsClient(nil),
			gloov1.SecretClient(nil),
			gatewayv1.VirtualServiceClient(nil),
			corev1.CoreV1Interface(nil),
		)

		converter = graphql.NewConverter(resolvers, context.TODO())
	})

	It("correctly converts a GraphQL Route to the proto format", func() {
		route, err := converter.ConvertInputRoute(GRAPHQL.inputRoute)
		Expect(err).NotTo(HaveOccurred())

		Expect(route.Matcher).To(BeEquivalentTo(PROTO.route.Matcher))
		Expect(route.Action).To(BeEquivalentTo(PROTO.route.Action))
		Expect(route.RoutePlugins).To(BeEquivalentTo(PROTO.route.RoutePlugins))

		// Finally, check the whole struct
		Expect(route).To(BeEquivalentTo(PROTO.route))
	})

	It("correctly converts a Route proto to the GraphQL format", func() {
		route, err := converter.ConvertOutputRoute(nil, PROTO.route)
		Expect(err).NotTo(HaveOccurred())

		// Check individual fields to get more info in case of failure
		Expect(route.Matcher).To(BeEquivalentTo(GRAPHQL.route.Matcher))
		Expect(route.Destination).To(BeEquivalentTo(GRAPHQL.route.Destination))
		Expect(route.Plugins).To(BeEquivalentTo(GRAPHQL.route.Plugins))

		// Finally, check the whole struct
		Expect(route).To(BeEquivalentTo(GRAPHQL.route))
	})

	It("correctly converts a GraphQL VirtualService to the proto format", func() {

		inputVs := models.InputVirtualService{
			Domains: []string{"domain1", "domain2"},
			Metadata: models.InputMetadata{
				Name:      "test-vs-1",
				Namespace: "test-ns-1",
			},
			SslConfig: &models.InputSslConfig{
				SecretRef: models.InputResourceRef{
					Name:      "secret-1",
					Namespace: "secret-ns-1",
				},
			},
			Routes:          []models.InputRoute{GRAPHQL.inputRoute},
			RateLimitConfig: GRAPHQL.inputRlConfig,
		}

		expected := &gatewayv1.VirtualService{
			Metadata: core.Metadata{
				Name:      "test-vs-1",
				Namespace: "test-ns-1",
			},
			SslConfig: &gloov1.SslConfig{
				SslSecrets: &gloov1.SslConfig_SecretRef{
					SecretRef: &core.ResourceRef{
						Name:      "secret-1",
						Namespace: "secret-ns-1",
					},
				},
			},
			VirtualHost: &gloov1.VirtualHost{
				Domains: []string{"domain1", "domain2"},
				Routes:  []*gloov1.Route{PROTO.route},
				VirtualHostPlugins: &gloov1.VirtualHostPlugins{
					Extensions: &gloov1.Extensions{
						Configs: PROTO.rlConfig,
					},
				},
			},
		}

		result, err := converter.ConvertInputVirtualService(inputVs)
		Expect(err).NotTo(HaveOccurred())

		// Check individual fields to get more info in case of failure
		Expect(result.Metadata).To(BeEquivalentTo(expected.Metadata))
		Expect(result.SslConfig).To(BeEquivalentTo(expected.SslConfig))
		Expect(result.VirtualHost.Domains).To(BeEquivalentTo(expected.VirtualHost.Domains))
		Expect(result.VirtualHost.Routes).To(BeEquivalentTo(expected.VirtualHost.Routes))
		Expect(result.VirtualHost.VirtualHostPlugins).To(BeEquivalentTo(expected.VirtualHost.VirtualHostPlugins))

		// Finally, check the whole struct
		Expect(result).To(BeEquivalentTo(expected))
	})

	It("correctly converts a VirtualService proto to the GraphQL format", func() {

		inputVs := &gatewayv1.VirtualService{
			Status: core.Status{
				State: core.Status_Accepted,
			},
			Metadata: core.Metadata{
				Name:      "test-vs-1",
				Namespace: "test-ns-1",
			},
			SslConfig: &gloov1.SslConfig{
				SslSecrets: &gloov1.SslConfig_SecretRef{
					SecretRef: &core.ResourceRef{
						Name:      "secret-1",
						Namespace: "secret-ns-1",
					},
				},
			},
			VirtualHost: &gloov1.VirtualHost{
				Domains: []string{"domain1", "domain2"},
				Routes:  []*gloov1.Route{PROTO.route},
				VirtualHostPlugins: &gloov1.VirtualHostPlugins{
					Extensions: &gloov1.Extensions{
						Configs: PROTO.rlConfig,
					},
				},
			},
		}

		expected := &models.VirtualService{
			Metadata: models.Metadata{
				Name:      "test-vs-1",
				Namespace: "test-ns-1",
				GUID:      "*v1.VirtualService test-ns-1 test-vs-1",
			},
			Status: models.Status{
				State: models.StateAccepted,
			},
			Domains: []string{"domain1", "domain2"},
			SslConfig: &models.SslConfig{
				SecretRef: models.ResourceRef{
					Name:      "secret-1",
					Namespace: "secret-ns-1",
				},
			},
			Routes:          []models.Route{GRAPHQL.route},
			RateLimitConfig: GRAPHQL.rlConfig,
		}

		result, err := converter.ConvertOutputVirtualService(inputVs)
		Expect(err).NotTo(HaveOccurred())

		// Check individual fields to get more info in case of failure
		Expect(result.Domains).To(BeEquivalentTo(expected.Domains))
		Expect(result.Metadata).To(BeEquivalentTo(expected.Metadata))
		Expect(result.SslConfig).To(BeEquivalentTo(expected.SslConfig))
		// prune the cyclic VirtualService field
		for i := range result.Routes {
			result.Routes[i].VirtualService = nil
		}
		Expect(result.Routes).To(BeEquivalentTo(expected.Routes))
		Expect(result.RateLimitConfig).To(BeEquivalentTo(expected.RateLimitConfig))

		// Finally, check the whole struct
		Expect(result).To(BeEquivalentTo(expected))
	})

})

var GRAPHQL = struct {
	route         models.Route
	rlConfig      *models.RateLimitConfig
	inputRoute    models.InputRoute
	inputRlConfig *models.InputRateLimitConfig
	upstream      models.Upstream
}{
	route: models.Route{
		VirtualService: nil,
		Destination: models.SingleDestination{
			DestinationSpec: &models.AwsDestinationSpec{
				LogicalName:            "aws-1",
				ResponseTransformation: true,
				InvocationStyle:        models.AwsLambdaInvocationStyleSync,
			},
			Upstream: models.Upstream{
				Metadata: models.Metadata{
					Name:            upstream1Name,
					Namespace:       upstream1Namespace,
					GUID:            fmt.Sprintf("*v1.Upstream %v %v", upstream1Namespace, upstream1Name),
					ResourceVersion: "123",
				},
				Status: models.Status{
					State: models.StateAccepted,
				},
				Spec: &models.AwsUpstreamSpec{
					Region: "us-east-1",
					SecretRef: models.ResourceRef{
						Name:      "aws-secret-1",
						Namespace: upstream1Namespace,
					},
					Functions: []models.AwsLambdaFunction{{
						LogicalName:  "aws-1",
						FunctionName: "lambda-1",
						Qualifier:    "lambda-qualifier",
					}},
				},
			},
		},
		Matcher: models.Matcher{
			PathMatch:     route1PathMatch,
			PathMatchType: models.PathMatchTypePrefix,
			Methods:       []string{"GET", "POST"},
			Headers: []models.KeyValueMatcher{
				{
					Name:    "ikvm-1",
					Value:   "match-this",
					IsRegex: false,
				},
			},
			QueryParameters: []models.KeyValueMatcher{
				{
					Name:    "ikvm-2",
					Value:   "match-.*-this",
					IsRegex: true,
				},
			},
		},
		Plugins: &models.RoutePlugins{
			PrefixRewrite: &newPrefix,
		},
	},
	inputRoute: models.InputRoute{
		Destination: models.InputDestination{
			SingleDestination: &models.InputSingleDestination{
				Upstream: models.InputResourceRef{
					Name:      upstream1Name,
					Namespace: upstream1Namespace,
				},
				DestinationSpec: &models.InputDestinationSpec{
					Aws: &models.InputAwsDestinationSpec{
						LogicalName:            "aws-1",
						ResponseTransformation: true,
						InvocationStyle:        models.AwsLambdaInvocationStyleSync,
					},
				},
			},
		},
		Matcher: models.InputMatcher{
			PathMatch:     route1PathMatch,
			PathMatchType: models.PathMatchTypePrefix,
			Methods:       []string{"GET", "POST"},
			Headers: []models.InputKeyValueMatcher{
				{
					Name:    "ikvm-1",
					Value:   "match-this",
					IsRegex: false,
				},
			},
			QueryParameters: []models.InputKeyValueMatcher{
				{
					Name:    "ikvm-2",
					Value:   "match-.*-this",
					IsRegex: true,
				},
			},
		},
		Plugins: &models.InputRoutePlugins{
			PrefixRewrite: &newPrefix,
		},
	},
	rlConfig: &models.RateLimitConfig{
		AuthorizedLimits: &models.RateLimit{
			RequestsPerUnit: 100,
			Unit:            models.TimeUnitSecond,
		},
		AnonymousLimits: &models.RateLimit{
			RequestsPerUnit: 50,
			Unit:            models.TimeUnitSecond,
		},
	},
	inputRlConfig: &models.InputRateLimitConfig{
		AuthorizedLimits: &models.InputRateLimit{
			RequestsPerUnit: 100,
			Unit:            models.TimeUnitSecond,
		},
		AnonymousLimits: &models.InputRateLimit{
			RequestsPerUnit: 50,
			Unit:            models.TimeUnitSecond,
		},
	},
}

var PROTO = struct {
	route    *gloov1.Route
	rlConfig map[string]*types.Struct
	upstream gloov1.Upstream
}{
	route: &gloov1.Route{
		Matcher: &gloov1.Matcher{
			PathSpecifier: &gloov1.Matcher_Prefix{
				Prefix: route1PathMatch,
			},
			Methods: []string{"GET", "POST"},
			Headers: []*gloov1.HeaderMatcher{
				{
					Name:  "ikvm-1",
					Value: "match-this",
					Regex: false,
				},
			},
			QueryParameters: []*gloov1.QueryParameterMatcher{
				{
					Name:  "ikvm-2",
					Value: "match-.*-this",
					Regex: true,
				},
			},
		},
		Action: &gloov1.Route_RouteAction{
			RouteAction: &gloov1.RouteAction{
				Destination: &gloov1.RouteAction_Single{
					Single: &gloov1.Destination{
						DestinationType: &gloov1.Destination_Upstream{
							Upstream: &core.ResourceRef{
								Name:      upstream1Name,
								Namespace: upstream1Namespace,
							},
						},
						DestinationSpec: &gloov1.DestinationSpec{
							DestinationType: &gloov1.DestinationSpec_Aws{
								Aws: &aws.DestinationSpec{
									LogicalName:            "aws-1",
									ResponseTransformation: true,
									InvocationStyle:        aws.DestinationSpec_SYNC,
								},
							},
						},
					},
				},
			},
		},
		RoutePlugins: &gloov1.RoutePlugins{
			PrefixRewrite: &transformation.PrefixRewrite{
				PrefixRewrite: newPrefix,
			},
		},
	},
	rlConfig: map[string]*types.Struct{
		ratelimit.ExtensionName: {
			Fields: map[string]*types.Value{
				"authorized_limits": {
					Kind: &types.Value_StructValue{
						StructValue: &types.Struct{
							Fields: map[string]*types.Value{
								"unit": {
									Kind: &types.Value_StringValue{StringValue: "SECOND"},
								},
								"requests_per_unit": {
									Kind: &types.Value_NumberValue{NumberValue: 100},
								},
							},
						},
					},
				},
				"anonymous_limits": {
					Kind: &types.Value_StructValue{
						StructValue: &types.Struct{
							Fields: map[string]*types.Value{
								"unit": {
									Kind: &types.Value_StringValue{StringValue: "SECOND"},
								},
								"requests_per_unit": {
									Kind: &types.Value_NumberValue{NumberValue: 50},
								},
							},
						},
					},
				},
			},
		},
	},
	upstream: gloov1.Upstream{
		Status: core.Status{
			State: core.Status_Accepted,
		},
		Metadata: core.Metadata{
			Name:            upstream1Name,
			Namespace:       upstream1Namespace,
			ResourceVersion: "123",
		},
		UpstreamSpec: &gloov1.UpstreamSpec{
			UpstreamType: &gloov1.UpstreamSpec_Aws{
				Aws: &aws.UpstreamSpec{
					Region: "us-east-1",
					SecretRef: core.ResourceRef{
						Name:      "aws-secret-1",
						Namespace: upstream1Namespace,
					},
					LambdaFunctions: []*aws.LambdaFunctionSpec{{
						LogicalName:        "aws-1",
						LambdaFunctionName: "lambda-1",
						Qualifier:          "lambda-qualifier",
					}},
				},
			},
		},
	},
}
