package aws_test

import (
	"context"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/ptypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/aws"

	envoyaws "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/aws"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	awsapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws"
	glooaws "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	. "github.com/solo-io/gloo/projects/gloo/pkg/plugins/aws"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("Plugin", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc

		awsPlugin   plugins.Plugin
		params      plugins.Params
		vhostParams plugins.VirtualHostParams

		// input resources
		upstream *v1.Upstream
		route    *v1.Route

		// output resources
		envoyCluster *envoy_config_cluster_v3.Cluster
		envoyRoute   *envoy_config_route_v3.Route
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		awsPlugin = aws.NewPlugin()

		upstreamName := "up"
		funcName := "foo"
		secretRef := &core.ResourceRef{
			Name:      "secret",
			Namespace: defaults.GlooSystem,
		}

		upstream = &v1.Upstream{
			Metadata: &core.Metadata{
				Name:      upstreamName,
				Namespace: defaults.GlooSystem,
			},
			UpstreamType: &v1.Upstream_Aws{
				Aws: &glooaws.UpstreamSpec{
					LambdaFunctions: []*glooaws.LambdaFunctionSpec{{
						LogicalName:        funcName,
						LambdaFunctionName: "foo",
						Qualifier:          "v1",
					}},
					Region:    "us-east1",
					SecretRef: secretRef,
				},
			},
		}
		route = &v1.Route{
			Action: &v1.Route_RouteAction{
				RouteAction: &v1.RouteAction{
					Destination: &v1.RouteAction_Single{
						Single: &v1.Destination{
							DestinationType: &v1.Destination_Upstream{
								Upstream: &core.ResourceRef{
									Namespace: defaults.GlooSystem,
									Name:      upstreamName,
								},
							},
							DestinationSpec: &v1.DestinationSpec{
								DestinationType: &v1.DestinationSpec_Aws{
									Aws: &awsapi.DestinationSpec{
										LogicalName: funcName,
									},
								},
							},
						},
					},
				},
			},
		}

		envoyCluster = &envoy_config_cluster_v3.Cluster{}
		envoyRoute = &envoy_config_route_v3.Route{
			Action: &envoy_config_route_v3.Route_Route{
				Route: &envoy_config_route_v3.RouteAction{
					ClusterSpecifier: &envoy_config_route_v3.RouteAction_Cluster{
						Cluster: upstreamName,
					},
				},
			},
		}

		params.Snapshot = &v1snap.ApiSnapshot{
			Secrets: v1.SecretList{{
				Metadata: &core.Metadata{
					Name:      secretRef.GetName(),
					Namespace: secretRef.GetNamespace(),
				},
				Kind: &v1.Secret_Aws{
					Aws: &v1.AwsSecret{
						AccessKey: "accessKeyValue",
						SecretKey: "secretKeyValue",
					},
				},
			}},
		}
		vhostParams = plugins.VirtualHostParams{Params: params}

		awsPlugin.Init(plugins.InitParams{
			Ctx: ctx,
		})
	})

	AfterEach(func() {
		cancel()
	})

	Context("ProcessRoute", func() {

		BeforeEach(func() {
			err := awsPlugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, envoyCluster)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should process route", func() {
			err := awsPlugin.(plugins.RoutePlugin).ProcessRoute(plugins.RouteParams{VirtualHostParams: vhostParams}, route, envoyRoute)
			Expect(err).NotTo(HaveOccurred())
			Expect(envoyRoute.TypedPerFilterConfig).To(HaveKey(FilterName))
		})

		It("should set route options defined by open source plugin", func() {
			route.Action.(*v1.Route_RouteAction).RouteAction.Destination.(*v1.RouteAction_Single).Single.DestinationSpec.DestinationType = &v1.DestinationSpec_Aws{
				Aws: &awsapi.DestinationSpec{
					UnwrapAsAlb:     false,
					LogicalName:     "foo",
					InvocationStyle: glooaws.DestinationSpec_ASYNC,
				},
			}
			err := awsPlugin.(plugins.RoutePlugin).ProcessRoute(plugins.RouteParams{VirtualHostParams: vhostParams}, route, envoyRoute)
			Expect(err).NotTo(HaveOccurred())
			Expect(envoyRoute.TypedPerFilterConfig).To(HaveKey(FilterName))
			cfg := &envoyaws.AWSLambdaPerRoute{}
			err = ptypes.UnmarshalAny(envoyRoute.TypedPerFilterConfig[FilterName], cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.GetAsync()).To(BeTrue())
			Expect(cfg.GetName()).To(Equal("foo"))
			Expect(cfg.GetUnwrapAsAlb()).To(Equal(false))
		})

		It("should set route transformer when unwrapAsApiGateway=True && unwrapAsAlb=False", func() {
			route.Action.(*v1.Route_RouteAction).RouteAction.Destination.(*v1.RouteAction_Single).Single.DestinationSpec.DestinationType = &v1.DestinationSpec_Aws{
				Aws: &awsapi.DestinationSpec{
					UnwrapAsApiGateway: true,
					UnwrapAsAlb:        false,
					LogicalName:        "foo",
				},
			}
			err := awsPlugin.(plugins.RoutePlugin).ProcessRoute(plugins.RouteParams{VirtualHostParams: vhostParams}, route, envoyRoute)
			Expect(err).NotTo(HaveOccurred())
			Expect(envoyRoute.TypedPerFilterConfig).To(HaveKey(FilterName))
			cfg := &envoyaws.AWSLambdaPerRoute{}
			err = ptypes.UnmarshalAny(envoyRoute.TypedPerFilterConfig[FilterName], cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.TransformerConfig).ToNot(BeNil())
		})

		It("should not set route transformer when unwrapAsApiGateway=True && unwrapAsAlb=True", func() {
			route.Action.(*v1.Route_RouteAction).RouteAction.Destination.(*v1.RouteAction_Single).Single.DestinationSpec.DestinationType = &v1.DestinationSpec_Aws{
				Aws: &awsapi.DestinationSpec{
					UnwrapAsApiGateway: true,
					UnwrapAsAlb:        true,
					LogicalName:        "foo",
				},
			}
			err := awsPlugin.(plugins.RoutePlugin).ProcessRoute(plugins.RouteParams{VirtualHostParams: vhostParams}, route, envoyRoute)
			Expect(err).NotTo(HaveOccurred())
			Expect(envoyRoute.TypedPerFilterConfig).To(HaveKey(FilterName))
			cfg := &envoyaws.AWSLambdaPerRoute{}
			err = ptypes.UnmarshalAny(envoyRoute.TypedPerFilterConfig[FilterName], cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.TransformerConfig).To(BeNil())
		})

		It("should not set route transformer when unwrapAsApiGateway=False", func() {
			route.Action.(*v1.Route_RouteAction).RouteAction.Destination.(*v1.RouteAction_Single).Single.DestinationSpec.DestinationType = &v1.DestinationSpec_Aws{
				Aws: &awsapi.DestinationSpec{
					UnwrapAsApiGateway: false,
					LogicalName:        "foo",
				},
			}
			err := awsPlugin.(plugins.RoutePlugin).ProcessRoute(plugins.RouteParams{VirtualHostParams: vhostParams}, route, envoyRoute)
			Expect(err).NotTo(HaveOccurred())
			Expect(envoyRoute.TypedPerFilterConfig).To(HaveKey(FilterName))
			cfg := &envoyaws.AWSLambdaPerRoute{}
			err = ptypes.UnmarshalAny(envoyRoute.TypedPerFilterConfig[FilterName], cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.TransformerConfig).To(BeNil())
		})

	})

})
