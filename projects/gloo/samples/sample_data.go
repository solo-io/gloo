package samples

import (
	"fmt"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/aws"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/azure"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/kubernetes"
)

func MakeMetadata(kind, namespace string, i int) core.Metadata {
	return core.Metadata{
		Name:      fmt.Sprintf("%v-%v", kind, i),
		Namespace: namespace,
	}
}

func Upstreams() v1.UpstreamList {
	return v1.UpstreamList{
		{
			Metadata: MakeMetadata("upstream", "some-namespace", 1),
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
			Metadata: MakeMetadata("upstream", "some-namespace", 2),
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
			Metadata: MakeMetadata("upstream", "some-namespace", 3),
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Kube{
					Kube: &kubernetes.UpstreamSpec{
						ServiceName:      "palmer-eldritch",
						ServiceNamespace: "pkd",
						ServicePort:      8080,
						ServiceSpec:      &plugins.ServiceSpec{},
					},
				},
			},
		},
		{
			Metadata: MakeMetadata("upstream", "some-namespace", 4),
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Azure{
					Azure: &azure.UpstreamSpec{
						FunctionAppName: "one-cloud-to-rule-them-all",
						SecretRef:       "my-precious",
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
	meta1 := MakeMetadata("virtualservice", "some-namespace", 1)
	meta2 := MakeMetadata("virtualservice", "some-namespace", 2)
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
