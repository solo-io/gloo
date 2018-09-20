package samples

import (
	"os"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/aws"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/azure"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/kubernetes"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/rest"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/static"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/transformation"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
)


func AddSampleData(opts bootstrap.Opts, vsClient gatewayv1.VirtualServiceClient) error {
	upstreamClient, err := v1.NewUpstreamClient(opts.Upstreams)
	if err != nil {
		return err
	}
	secretClient, err := v1.NewSecretClient(opts.Secrets)
	if err != nil {
		return err
	}
	virtualServices, upstreams, secrets := VirtualServices(), Upstreams(), Secrets()
	for _, item := range virtualServices {
		if _, err := vsClient.Write(item, clients.WriteOpts{}); err != nil && !errors.IsExist(err) {
			return err
		}
	}
	for _, item := range upstreams {
		if _, err := upstreamClient.Write(item, clients.WriteOpts{}); err != nil && !errors.IsExist(err) {
			return err
		}
	}
	for _, item := range secrets {
		if _, err := secretClient.Write(item, clients.WriteOpts{}); err != nil && !errors.IsExist(err) {
			return err
		}
	}
	return nil
}


func MakeMetadata(name, namespace string) core.Metadata {
	return core.Metadata{
		Name:      name,
		Namespace: namespace,
	}
}

func Secrets() v1.SecretList {
	return v1.SecretList{
		{
			Metadata: MakeMetadata("some-secret", defaults.GlooSystem),
			Kind: &v1.Secret_Aws{
				Aws: &v1.AwsSecret{
					AccessKey: "TODO_REMOVE" + os.Getenv("AWS_ACCESS_KEY"),
					SecretKey: "TODO_REMOVE" + os.Getenv("AWS_SECRET_KEY"),
				},
			},
		},
		{
			Metadata: MakeMetadata("my-precious", defaults.GlooSystem),
			Kind: &v1.Secret_Azure{
				Azure: &v1.AzureSecret{
					ApiKeys: map[string]string{"TO": "DO"},
				},
			},
		},
		{
			Metadata: MakeMetadata("ssl-secret", defaults.GlooSystem),
			Kind: &v1.Secret_Tls{
				Tls: &v1.TlsSecret{
					// TODO(ilackarms): ssl secret
				},
			},
		},
	}
}

func Artifacts() v1.ArtifactList {
	return v1.ArtifactList{
		{
			Metadata: MakeMetadata("artifact", "default"),
			Data:     "// TODO(ilackarms)",
		},
	}
}

func Upstreams() v1.UpstreamList {
	return v1.UpstreamList{
		{
			Metadata: MakeMetadata("my-aws-account-pls-donthack", "default"),
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Aws{
					Aws: &aws.UpstreamSpec{
						Region: "us-east-1",
						SecretRef: core.ResourceRef{
							Namespace: defaults.GlooSystem,
							Name:      "some-secret",
						},
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
			Metadata: MakeMetadata("kube-1", defaults.GlooSystem),
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
			Metadata: MakeMetadata("kube-2", defaults.GlooSystem),
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Kube{
					Kube: &kubernetes.UpstreamSpec{
						ServiceName:      "palmer-eldritch",
						ServiceNamespace: "pkd",
						ServicePort:      8080,
						ServiceSpec: &plugins.ServiceSpec{
							PluginType: &plugins.ServiceSpec_Rest{
								Rest: &rest.ServiceSpec{
									Transformations: map[string]*transformation.TransformationTemplate{
										"my-rest-function": {
											// TODO(ilackarms/yuval-k)
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
			Metadata: MakeMetadata("azure", defaults.GlooSystem),
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Azure{
					Azure: &azure.UpstreamSpec{
						FunctionAppName: "one-cloud-to-rule-them-all",
						SecretRef: core.ResourceRef{
							Name:      "my-precious",
							Namespace: defaults.GlooSystem,
						},
						Functions: []*azure.UpstreamSpec_FunctionSpec{
							{
								FunctionName: "CreateRing",
								AuthLevel:    azure.UpstreamSpec_FunctionSpec_Anonymous,
							},
							{
								FunctionName: "DestroyRing",
								AuthLevel:    azure.UpstreamSpec_FunctionSpec_Function,
							},
							{
								FunctionName: "TransportRing",
								AuthLevel:    azure.UpstreamSpec_FunctionSpec_Admin,
							},
						},
					},
				},
			},
		},
		{
			Metadata: MakeMetadata("static-1", defaults.GlooSystem),
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Static{
					Static: &static.UpstreamSpec{
						Hosts: []*static.Host{
							{Addr: "127.0.0.1", Port: 8000},
						},
					},
				},
			},
		},
	}
}

func VirtualServices() gatewayv1.VirtualServiceList {
	meta1 := MakeMetadata("virtualservice1", defaults.GlooSystem)
	meta2 := MakeMetadata("virtualservice2", defaults.GlooSystem)
	return gatewayv1.VirtualServiceList{
		{
			Metadata: meta1,
			SslConfig: &v1.SslConfig{SslSecrets: &v1.SslConfig_SecretRef{
				SecretRef: &core.ResourceRef{Name: "ssl-secret", Namespace: "gloo-sytsem"},
			}},
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
													Upstream: core.ResourceRef{
														Name:      "my-aws-account-pls-donthack",
														Namespace: "default",
													},
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
													Upstream: core.ResourceRef{
														Name:      "azure",
														Namespace: defaults.GlooSystem,
													},
													DestinationSpec: &v1.DestinationSpec{
														DestinationType: &v1.DestinationSpec_Azure{
															Azure: &azure.DestinationSpec{
																FunctionName: "CreateRing",
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
													Upstream: core.ResourceRef{
														Name:      "my-aws-account-pls-donthack",
														Namespace: "default",
													},
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
													Upstream: core.ResourceRef{
														Name:      "azure",
														Namespace: defaults.GlooSystem,
													},
													DestinationSpec: &v1.DestinationSpec{
														DestinationType: &v1.DestinationSpec_Azure{
															Azure: &azure.DestinationSpec{
																FunctionName: "CreateRing",
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
								"/pizza",
							},
							Methods: []string{"GET"},
						},
						Action: &v1.Route_RouteAction{
							RouteAction: &v1.RouteAction{
								Destination: &v1.RouteAction_Single{
									Single: &v1.Destination{
										Upstream: core.ResourceRef{
											Name:      "static-1",
											Namespace: defaults.GlooSystem,
										},
									},
								},
							},
						},
					},
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
										Upstream: core.ResourceRef{
											Name:      "kube-1",
											Namespace: defaults.GlooSystem,
										},
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

func ResolverMaps() sqoopv1.ResolverMapList {
	return sqoopv1.ResolverMapList{
		{
			Metadata: MakeMetadata("resolvermap1", defaults.GlooSystem),
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
			Metadata: MakeMetadata("resolvermap2", defaults.GlooSystem),
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

func Schemas() sqoopv1.SchemaList {
	return sqoopv1.SchemaList{
		{
			Metadata: MakeMetadata("petstore", defaults.GlooSystem),
			InlineSchema: `# The query type, represents all of the entry points into our object graph
type Query {
    pets: [Pet]
    pet(id: Int!): Pet
}

type Mutation {
    addPet(pet: InputPet!): Pet
}

type Pet{
    id: ID!
    name: String!
    status: Status!
}

input InputPet{
    id: ID!
    name: String!
    tag: String
}

enum Status {
    pending
    available
}
`,
		},
	}
}
