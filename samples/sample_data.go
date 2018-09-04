package samples

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/aws"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/azure"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/kubernetes"
	"os"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/rest"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/transformation"
)

func MakeMetadata(name, namespace string) core.Metadata {
	return core.Metadata{
		Name:      name,
		Namespace: namespace,
	}
}

func Secrets() v1.SecretList {
	return v1.SecretList{
		{
			Metadata: MakeMetadata("some-secret", "gloo-system"),
			Data: map[string]string{
				"access_key": "TODO_REMOVE"+os.Getenv("AWS_ACCESS_KEY"),
				"secret_key": "TODO_REMOVE"+os.Getenv("AWS_SECRET_KEY"),
			},
		},
		{
			Metadata: MakeMetadata("my-precious", "gloo-system"),
			Data: map[string]string{
				// TODO(ilackarms): azure secrets
			},
		},
		{
			Metadata: MakeMetadata("ssl-secret", "gloo-system"),
			Data: map[string]string{
				// TODO(ilackarms): ssl secret
			},
		},
	}
}

func Artifacts() v1.ArtifactList {
	return v1.ArtifactList{
		{
			Metadata: MakeMetadata("artifact", "default"),
			Data: map[string]string{
				// TODO(ilackarms)
			},
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
							Namespace: "gloo-system",
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
			Metadata: MakeMetadata("kube-1", "gloo-system"),
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
			Metadata: MakeMetadata("kube-2", "gloo-system"),
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Kube{
					Kube: &kubernetes.UpstreamSpec{
						ServiceName:      "palmer-eldritch",
						ServiceNamespace: "pkd",
						ServicePort:      8080,
						ServiceSpec: &plugins.ServiceSpec{
							PluginType: &plugins.ServiceSpec_Rest{
								Rest: &rest.ServiceSpec{
									Transformation: map[string]*transformation.TransformationTemplate{
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
			Metadata: MakeMetadata("azure", "gloo-system"),
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Azure{
					Azure: &azure.UpstreamSpec{
						FunctionAppName: "one-cloud-to-rule-them-all",
						SecretRef: core.ResourceRef{
							Name:      "my-precious",
							Namespace: "gloo-system",
						},
						Functions: []*azure.UpstreamSpec_FunctionSpec{
							{
								FunctionName: "CreateRing",
								AuthLevel:    "dwarf_lvl",
							},
							{
								FunctionName: "DestroyRing",
								AuthLevel:    "hobbit_lvl",
							},
							{
								FunctionName: "TransportRing",
								AuthLevel:    "hobbit_lvl",
							},
						},
					},
				},
			},
		},
	}
}

func VirtualServices() gatewayv1.VirtualServiceList {
	meta1 := MakeMetadata("virtualservice1", "gloo-system")
	meta2 := MakeMetadata("virtualservice2", "gloo-system")
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
														Namespace: "gloo-system",
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
														Namespace: "gloo-system",
													},
													DestinationSpec: &v1.DestinationSpec{
														DestinationType: &v1.DestinationSpec_Azure{
															Azure: &azure.DestinationSpec{
																FunctionName: "my_func_v2",
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
										Upstream: core.ResourceRef{
											Name:      "kube-1",
											Namespace: "gloo-system",
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
