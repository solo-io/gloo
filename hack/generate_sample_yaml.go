package main

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/rest"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-projects/projects/sqoop/pkg/api/v1"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("finished.")
}


// TODO (ilackarms): fix or delete
func run() error {
	//opts, err := setup.DefaultKubernetesConstructOpts()
	//if err != nil {
	//	return err
	//}
	//schemas := samples.Schemas()
	//resoverMapsClient, err := v1.NewResolverMapClient(opts.ResolverMaps)
	//if err != nil {
	//	return err
	//}
	//resolvermaps, err := resoverMapsClient.List(defaults.GlooSystem, clients.ListOpts{})
	//if err != nil {
	//	return err
	//}
	//err = resoverMapsClient.Delete(rm.Metadata.Namespace, rm.Metadata.Name, clients.DeleteOpts{})
	//_, err = resoverMapsClient.Write(rm, clients.WriteOpts{})
	//log.Printf("%v, %v", err, resolvermaps)
	//
	//schemaClient, err := v1.NewSchemaClient(opts.Schemas)
	//if err != nil {
	//	return err
	//}
	//for _, sch := range schemas {
	//	_, err := schemaClient.Write(sch, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true})
	//	if err != nil {
	//		log.Printf("failed writing schema %v", err)
	//	}
	//}
	return nil
}

var rm = &v1.ResolverMap{
	Types: map[string]*v1.TypeResolver{
		"Query": &v1.TypeResolver{
			Fields: map[string]*v1.FieldResolver{
				"pet": &v1.FieldResolver{
					Resolver: &v1.FieldResolver_GlooResolver{
						GlooResolver: &v1.GlooResolver{
							RequestTemplate: &v1.RequestTemplate{
								Headers: map[string]string{"Content-type": "application/json"},
								Body:    `{"id": {{.}} }`,
							},
							Action: &gloov1.RouteAction{
								Destination: &gloov1.RouteAction_Single{
									Single: &gloov1.Destination{
										Upstream: core.ResourceRef{Name: "petstir", Namespace: defaults.GlooSystem},
										DestinationSpec: &gloov1.DestinationSpec{
											DestinationType: &gloov1.DestinationSpec_Rest{
												Rest: &rest.DestinationSpec{
													FunctionName: "findPetById",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"pets": &v1.FieldResolver{
					Resolver: &v1.FieldResolver_GlooResolver{
						GlooResolver: &v1.GlooResolver{
							RequestTemplate: &v1.RequestTemplate{
								Headers: map[string]string{"Content-type": "application/json"},
							},
							Action: &gloov1.RouteAction{
								Destination: &gloov1.RouteAction_Single{
									Single: &gloov1.Destination{
										Upstream: core.ResourceRef{Name: "petstir", Namespace: defaults.GlooSystem},
										DestinationSpec: &gloov1.DestinationSpec{
											DestinationType: &gloov1.DestinationSpec_Rest{
												Rest: &rest.DestinationSpec{
													FunctionName: "findPetById",
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
		},
		"Mutation": &v1.TypeResolver{
			Fields: map[string]*v1.FieldResolver{
				"addPet": &v1.FieldResolver{
					Resolver: nil,
				},
			},
		},
		"Pet": &v1.TypeResolver{
			Fields: map[string]*v1.FieldResolver{
				"id": &v1.FieldResolver{
					Resolver: nil,
				},
				"name": &v1.FieldResolver{
					Resolver: nil,
				},
				"status": &v1.FieldResolver{
					Resolver: nil,
				},
			},
		},
	},
	Status: core.Status{
		State:               0,
		Reason:              "",
		ReportedBy:          "",
		SubresourceStatuses: map[string]*core.Status{},
	},
	Metadata: core.Metadata{
		Name:      "petstore",
		Namespace: defaults.GlooSystem,
		Annotations: map[string]string{
			"created_for": "petstore",
		},
	},
}
