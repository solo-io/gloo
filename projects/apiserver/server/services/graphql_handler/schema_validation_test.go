package graphql_handler_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	graphql_v1alpha1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1alpha1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/graphql_handler"
)

var _ = Describe("schema validation", func() {

	Context("with schema definition", func() {
		It("accepts empty schema", func() {
			err := graphql_handler.ValidateSchemaDefinition(&rpc_edge_v1.ValidateSchemaDefinitionRequest{
				Input: &rpc_edge_v1.ValidateSchemaDefinitionRequest_SchemaDefinition{
					SchemaDefinition: "",
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("accepts schema without resolve directives", func() {
			err := graphql_handler.ValidateSchemaDefinition(&rpc_edge_v1.ValidateSchemaDefinitionRequest{
				Input: &rpc_edge_v1.ValidateSchemaDefinitionRequest_SchemaDefinition{
					SchemaDefinition: `
	type Query {
		productsForHome: [Product]
	}
`,
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("rejects schema with resolve directive", func() {
			// this should fail because there is a resolve directive with no corresponding resolutions map
			err := graphql_handler.ValidateSchemaDefinition(&rpc_edge_v1.ValidateSchemaDefinitionRequest{
				Input: &rpc_edge_v1.ValidateSchemaDefinitionRequest_SchemaDefinition{
					SchemaDefinition: `
	type Query {
		productsForHome: [Product] @resolve(name: "Query|productsForHome")
	}
`,
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("resolver Query|productsForHome is not defined"))
		})

		It("rejects schema with syntax errors", func() {
			err := graphql_handler.ValidateSchemaDefinition(&rpc_edge_v1.ValidateSchemaDefinitionRequest{
				Input: &rpc_edge_v1.ValidateSchemaDefinitionRequest_SchemaDefinition{
					SchemaDefinition: `
	type Query {
		productsForHome: asdfas[Product]
	}
`,
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unable to parse graphql schema"))
		})

		It("accepts schema with valid cacheControl", func() {
			err := graphql_handler.ValidateSchemaDefinition(&rpc_edge_v1.ValidateSchemaDefinitionRequest{
				Input: &rpc_edge_v1.ValidateSchemaDefinitionRequest_SchemaDefinition{
					SchemaDefinition: `
	type Query {
		productsForHome: [Product] @cacheControl(maxAge: 60, inheritMaxAge: false, scope: private)
	}
`,
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("rejects schema with invalid cacheControl", func() {
			err := graphql_handler.ValidateSchemaDefinition(&rpc_edge_v1.ValidateSchemaDefinitionRequest{
				Input: &rpc_edge_v1.ValidateSchemaDefinitionRequest_SchemaDefinition{
					SchemaDefinition: `
	type Query {
		productsForHome: [Product] @cacheControl(maxAge: 60, inheritMaxAge: false, scope: other)
	}
`,
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unimplemented cacheControl scope type other"))
		})
	})

	Context("with graphqlapi spec", func() {
		It("accepts empty schema", func() {
			err := graphql_handler.ValidateSchemaDefinition(&rpc_edge_v1.ValidateSchemaDefinitionRequest{
				Input: &rpc_edge_v1.ValidateSchemaDefinitionRequest_Spec{
					Spec: &graphql_v1alpha1.GraphQLApiSpec{
						Schema: &graphql_v1alpha1.GraphQLApiSpec_ExecutableSchema{
							ExecutableSchema: &graphql_v1alpha1.ExecutableSchema{
								SchemaDefinition: "",
							},
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("accepts schema with resolver name that's in map", func() {
			err := graphql_handler.ValidateSchemaDefinition(&rpc_edge_v1.ValidateSchemaDefinitionRequest{
				Input: &rpc_edge_v1.ValidateSchemaDefinitionRequest_Spec{
					Spec: &graphql_v1alpha1.GraphQLApiSpec{
						Schema: &graphql_v1alpha1.GraphQLApiSpec_ExecutableSchema{
							ExecutableSchema: &graphql_v1alpha1.ExecutableSchema{
								SchemaDefinition: `
	type Query {
		productsForHome: [Product] @resolve(name: "Query|productsForHome")
	}
`,
								Executor: &graphql_v1alpha1.Executor{
									Executor: &graphql_v1alpha1.Executor_Local_{
										Local: &graphql_v1alpha1.Executor_Local{
											Resolutions: map[string]*graphql_v1alpha1.Resolution{
												"Query|productsForHome": {
													Resolver: &graphql_v1alpha1.Resolution_RestResolver{
														RestResolver: &graphql_v1alpha1.RESTResolver{},
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
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("rejects schema with resolver name not in map", func() {
			err := graphql_handler.ValidateSchemaDefinition(&rpc_edge_v1.ValidateSchemaDefinitionRequest{
				Input: &rpc_edge_v1.ValidateSchemaDefinitionRequest_Spec{
					Spec: &graphql_v1alpha1.GraphQLApiSpec{
						Schema: &graphql_v1alpha1.GraphQLApiSpec_ExecutableSchema{
							ExecutableSchema: &graphql_v1alpha1.ExecutableSchema{
								SchemaDefinition: `
	type Query {
		productsForHome: [Product] @resolve(name: "Query|productsForHome")
	}
`,
								Executor: &graphql_v1alpha1.Executor{
									Executor: &graphql_v1alpha1.Executor_Local_{
										Local: &graphql_v1alpha1.Executor_Local{
											Resolutions: map[string]*graphql_v1alpha1.Resolution{
												"Query|productsForHome123": {
													Resolver: &graphql_v1alpha1.Resolution_RestResolver{
														RestResolver: &graphql_v1alpha1.RESTResolver{},
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
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("resolver Query|productsForHome is not defined"))
		})
	})

})
