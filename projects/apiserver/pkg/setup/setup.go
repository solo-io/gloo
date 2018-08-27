package setup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	apiserver "github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/graph"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/aws"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/azure"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/kubernetes"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
)

func Setup(port int) error {
	// TODO (ilackarms): pass in the factory with cli flags
	inputFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
		Cache: memory.NewInMemoryResourceCache(),
	})
	upstreams, err := v1.NewUpstreamClient(inputFactory)
	if err != nil {
		return err
	}
	virtualServices, err := gatewayv1.NewVirtualServiceClient(inputFactory)
	if err != nil {
		return err
	}
	resolverMaps, err := sqoopv1.NewResolverMapClient(inputFactory)
	if err != nil {
		return err
	}

	err = addSampleData(upstreams, virtualServices, resolverMaps)
	if err != nil {
		return err
	}

	http.Handle("/", handler.Playground("Solo-ApiServer", "/query"))
	http.Handle("/query", handler.GraphQL(graph.NewExecutableSchema(graph.Config{
		Resolvers: apiserver.NewResolvers(upstreams, virtualServices, resolverMaps),
	}),
		handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
			rc := graphql.GetResolverContext(ctx)
			fmt.Println("Entered", rc.Object, rc.Field.Name)
			res, err = next(ctx)
			fmt.Println("Left", rc.Object, rc.Field.Name, "=>", res, err)
			return res, err
		}),
	))

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func addSampleData(usClient v1.UpstreamClient,
	vsClient gatewayv1.VirtualServiceClient,
	rmClient sqoopv1.ResolverMapClient) error {
	upstreams, virtualServices, resolverMaps := sampleData()
	for _, us := range upstreams {
		_, err := usClient.Write(us, clients.WriteOpts{})
		if err != nil {
			return err
		}
	}
	for _, vs := range virtualServices {
		_, err := vsClient.Write(vs, clients.WriteOpts{})
		if err != nil {
			return err
		}
	}
	for _, rm := range resolverMaps {
		_, err := rmClient.Write(rm, clients.WriteOpts{})
		if err != nil {
			return err
		}
	}
	return nil
}

func sampleData() (v1.UpstreamList, gatewayv1.VirtualServiceList, sqoopv1.ResolverMapList) {
	return sampleUpstreams(), sampleVirtualServices(), sampleResolverMaps()
}

func sampleUpstreams() v1.UpstreamList {
	return v1.UpstreamList{
		{
			Metadata: makeMetadata("upstream", "some-namespace", 1),
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Aws{
					Aws: &aws.UpstreamSpec{
						Region:    "us-east-1",
						SecretRef: "some-secret",
						LambdaFunctions: []*aws.LambdaFunctionSpec{
							{
								LogicalName:        "my_func_v1",
								LambdaFunctionName: "my_func",
								Qualifier:          "v1",
							},
							{
								LogicalName:        "my_func_v2",
								LambdaFunctionName: "my_func",
								Qualifier:          "$LATEST",
							},
						},
					},
				},
			},
		},
		{
			Metadata: makeMetadata("upstream", "some-namespace", 2),
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Kube{
					Kube: &kubernetes.UpstreamSpec{
						ServiceName:      "perky-pat",
						ServiceNamespace: "default",
						ServicePort:      8080,
					},
				},
			},
		},
		{
			Metadata: makeMetadata("upstream", "some-namespace", 3),
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Kube{
					Kube: &kubernetes.UpstreamSpec{
						ServiceName:      "palmer-eldritch",
						ServiceNamespace: "pkd",
						ServicePort:      8080,
						ServiceSpec: &plugins.ServiceSpec{
							PluginType: &plugins.ServiceSpec_Empty{
								Empty: "eventually this will be replaced with an actual plugin (gRPC or Swagger)",
							},
						},
					},
				},
			},
		},
	}
}

func sampleVirtualServices() gatewayv1.VirtualServiceList {
	meta1 := makeMetadata("virtualservice", "some-namespace", 1)
	meta2 := makeMetadata("virtualservice", "some-namespace", 2)
	return gatewayv1.VirtualServiceList{
		{
			Metadata:  meta1,
			SslConfig: &v1.SslConfig{SslSecrets: &v1.SslConfig_SecretRef{SecretRef: "some-secret"}},
			VirtualHost: &v1.VirtualHost{
				Name:    meta1.Name,
				Domains: []string{"sqoop.didoop.com"},
				Routes: []*v1.Route{
					{
						Matcher: &v1.Matcher{
							PathSpecifier: &v1.Matcher_Prefix{
								"/makemeapizza",
							},
							Headers: []*v1.HeaderMatcher{
								{
									Name:  "x-custom-header",
									Value: "*",
									Regex: true,
								},
								{
									Name:  "x-special-header",
									Value: "dinosaurs",
									Regex: false,
								},
							},
							QueryParameters: []*v1.QueryParameterMatcher{
								{
									Name:  "favorite_day",
									Value: "friday",
								},
								{
									Name:  "best_xmen",
									Value: "professor_x*",
									Regex: true,
								},
							},
							Methods: []string{"GET", "POST"},
						},
						Action: &v1.Route_RouteAction{
							RouteAction: &v1.RouteAction{
								Destination: &v1.RouteAction_Multi{
									Multi: &v1.MultiDestination{
										Destinations: []*v1.WeightedDestination{
											{
												Destination: &v1.Destination{
													UpstreamName: "my-aws-account-pls-donthack",
													DestinationSpec: &v1.DestinationSpec{
														DestinationType: &v1.DestinationSpec_Aws{
															Aws: &aws.DestinationSpec{
																LogicalName:     "my_func_v1",
																InvocationStyle: aws.DestinationSpec_ASYNC,
															},
														},
													},
												},
												Weight: 1,
											},
											{
												Destination: &v1.Destination{
													UpstreamName: "my-azure-account-pls-donthack",
													DestinationSpec: &v1.DestinationSpec{
														DestinationType: &v1.DestinationSpec_Azure{
															Azure: &azure.DestinationSpec{
																FunctionName: "my_other_func_v1",
															},
														},
													},
												},
												Weight: 2,
											},
										},
									},
								},
							},
						},
					},
					{
						Matcher: &v1.Matcher{
							PathSpecifier: &v1.Matcher_Prefix{
								"/makemeasalad",
							},
							Methods: []string{"GET", "POST"},
						},
						Action: &v1.Route_RouteAction{
							RouteAction: &v1.RouteAction{
								Destination: &v1.RouteAction_Multi{
									Multi: &v1.MultiDestination{
										Destinations: []*v1.WeightedDestination{
											{
												Destination: &v1.Destination{
													UpstreamName: "my-aws-account-pls-donthack",
													DestinationSpec: &v1.DestinationSpec{
														DestinationType: &v1.DestinationSpec_Aws{
															Aws: &aws.DestinationSpec{
																LogicalName:     "my_func_v1",
																InvocationStyle: aws.DestinationSpec_ASYNC,
															},
														},
													},
												},
												Weight: 12,
											},
											{
												Destination: &v1.Destination{
													UpstreamName: "my-azure-account-pls-donthack",
													DestinationSpec: &v1.DestinationSpec{
														DestinationType: &v1.DestinationSpec_Azure{
															Azure: &azure.DestinationSpec{
																FunctionName: "my_other_func_v1",
															},
														},
													},
												},
												Weight: 25,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Metadata: meta2,
			VirtualHost: &v1.VirtualHost{
				Name:    meta2.Name,
				Domains: []string{"*"},
				Routes: []*v1.Route{
					{
						Matcher: &v1.Matcher{
							PathSpecifier: &v1.Matcher_Prefix{
								"/frenchfries",
							},
							Methods: []string{"GET", "POST", "PATCH", "PUT", "OPTIONS"},
						},
						Action: &v1.Route_RouteAction{
							RouteAction: &v1.RouteAction{
								Destination: &v1.RouteAction_Single{
									Single: &v1.Destination{
										UpstreamName: "my-kube-service",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func sampleResolverMaps() sqoopv1.ResolverMapList {
	return sqoopv1.ResolverMapList{
		{
			Metadata: makeMetadata("resolvermap", "some-namespace", 1),
			Types: map[string]*sqoopv1.TypeResolver{
				"Foo": {
					Fields: map[string]*sqoopv1.FieldResolver{
						"field1": {Resolver: &sqoopv1.FieldResolver_GlooResolver{}},
						"field2": {Resolver: &sqoopv1.FieldResolver_TemplateResolver{}},
					},
				},
				"Bar": {
					Fields: map[string]*sqoopv1.FieldResolver{
						"field1": {Resolver: &sqoopv1.FieldResolver_GlooResolver{}},
						"field2": {Resolver: &sqoopv1.FieldResolver_TemplateResolver{}},
					},
				},
			},
		},
		{
			Metadata: makeMetadata("resolvermap", "some-namespace", 2),
			Types: map[string]*sqoopv1.TypeResolver{
				"Baz": {
					Fields: map[string]*sqoopv1.FieldResolver{
						"field1": {Resolver: &sqoopv1.FieldResolver_GlooResolver{}},
						"field2": {Resolver: &sqoopv1.FieldResolver_TemplateResolver{}},
					},
				},
				"Qux": {
					Fields: map[string]*sqoopv1.FieldResolver{
						"field1": {Resolver: &sqoopv1.FieldResolver_GlooResolver{}},
						"field2": {Resolver: &sqoopv1.FieldResolver_TemplateResolver{}},
					},
				},
			},
		},
	}
}

func makeMetadata(kind, namespace string, i int) core.Metadata {
	return core.Metadata{
		Name:      fmt.Sprintf("%v-%v", kind, i),
		Namespace: namespace,
	}
}
