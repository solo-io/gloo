package v8go_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	. "github.com/solo-io/solo-kit/test/matchers"
	enterprisev1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/enterprise/graphql/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/v8go"
)

var _ = Describe("SchemaDiff", func() {
	var runner *v8go.StitchingScriptRunner
	var err error

	BeforeEach(func() {
		runner, err = v8go.GetStitchingScriptRunner()
		Expect(err).ToNot(HaveOccurred())
	})

	It("returns empty diff when both schemas are empty", func() {
		in := &enterprisev1.GraphQLInspectorDiffInput{
			OldSchema: "",
			NewSchema: "",
		}
		out, err := runner.RunSchemaDiff(in)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(MatchProto(&enterprisev1.GraphQLInspectorDiffOutput{
			Changes: []*enterprisev1.GraphQLInspectorDiffOutput_Change{},
		}))
	})

	It("returns empty diff when old and new schema are the same", func() {
		schema := `
type Query {
	productsForHome: [Product]
}

type Product {
    id: String
    title: String
}
`
		in := &enterprisev1.GraphQLInspectorDiffInput{
			OldSchema: schema,
			NewSchema: schema,
		}
		out, err := runner.RunSchemaDiff(in)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(MatchProto(&enterprisev1.GraphQLInspectorDiffOutput{
			Changes: []*enterprisev1.GraphQLInspectorDiffOutput_Change{},
		}))
	})

	It("returns error on invalid syntax", func() {
		in := &enterprisev1.GraphQLInspectorDiffInput{
			OldSchema: "blarg",
			NewSchema: "",
		}
		_, err := runner.RunSchemaDiff(in)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Syntax Error: Unexpected Name \"blarg\""))
	})

	It("returns non-breaking change when field is added", func() {
		oldSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
    id: String
    title: String
}
`
		newSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
    id: String
	author: String
    title: String
}
`
		in := &enterprisev1.GraphQLInspectorDiffInput{
			OldSchema: oldSchema,
			NewSchema: newSchema,
		}
		out, err := runner.RunSchemaDiff(in)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(MatchProto(&enterprisev1.GraphQLInspectorDiffOutput{
			Changes: []*enterprisev1.GraphQLInspectorDiffOutput_Change{
				{
					Message:    "Field 'author' was added to object type 'Product'",
					Path:       "Product.author",
					ChangeType: "FIELD_ADDED",
					Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
						Level:  enterprisev1.GraphQLInspectorDiffOutput_NON_BREAKING,
						Reason: "",
					},
				},
			},
		}))
	})

	It("returns breaking change when field is removed", func() {
		oldSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
    id: String
	author: String
    title: String
}
`
		newSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
    id: String
    title: String
}
`
		in := &enterprisev1.GraphQLInspectorDiffInput{
			OldSchema: oldSchema,
			NewSchema: newSchema,
		}
		out, err := runner.RunSchemaDiff(in)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(MatchProto(&enterprisev1.GraphQLInspectorDiffOutput{
			Changes: []*enterprisev1.GraphQLInspectorDiffOutput_Change{
				{
					Message:    "Field 'author' was removed from object type 'Product'",
					Path:       "Product.author",
					ChangeType: "FIELD_REMOVED",
					Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
						Level:  enterprisev1.GraphQLInspectorDiffOutput_BREAKING,
						Reason: "Removing a field is a breaking change. It is preferable to deprecate the field before removing it.",
					},
				},
			},
		}))
	})

	It("can return a mixture of change types", func() {
		oldSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
    id: String @deprecated(reason: "used to be deprecated")
	author: String
    title: String
}
`
		newSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
    id: String
    title: String @deprecated(reason: "no longer needed")
	price: Int
}
`
		in := &enterprisev1.GraphQLInspectorDiffInput{
			OldSchema: oldSchema,
			NewSchema: newSchema,
		}
		out, err := runner.RunSchemaDiff(in)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(MatchProto(&enterprisev1.GraphQLInspectorDiffOutput{
			Changes: []*enterprisev1.GraphQLInspectorDiffOutput_Change{
				{
					Message:    "Field 'author' was removed from object type 'Product'",
					Path:       "Product.author",
					ChangeType: "FIELD_REMOVED",
					Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
						Level:  enterprisev1.GraphQLInspectorDiffOutput_BREAKING,
						Reason: "Removing a field is a breaking change. It is preferable to deprecate the field before removing it.",
					},
				},
				{
					Message:    "Field 'Product.id' is no longer deprecated",
					Path:       "Product.id",
					ChangeType: "FIELD_DEPRECATION_REMOVED",
					Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
						Level:  enterprisev1.GraphQLInspectorDiffOutput_DANGEROUS,
						Reason: "",
					},
				},
				{
					Message:    "Deprecation reason was removed from field 'Product.id'",
					Path:       "Product.id",
					ChangeType: "FIELD_DEPRECATION_REASON_REMOVED",
					Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
						Level:  enterprisev1.GraphQLInspectorDiffOutput_NON_BREAKING,
						Reason: "",
					},
				},
				{
					Message:    "Field 'price' was added to object type 'Product'",
					Path:       "Product.price",
					ChangeType: "FIELD_ADDED",
					Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
						Level:  enterprisev1.GraphQLInspectorDiffOutput_NON_BREAKING,
						Reason: "",
					},
				},
				{
					Message:    "Field 'Product.title' is deprecated",
					Path:       "Product.title",
					ChangeType: "FIELD_DEPRECATION_ADDED",
					Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
						Level:  enterprisev1.GraphQLInspectorDiffOutput_NON_BREAKING,
						Reason: "",
					},
				},
				{
					Message:    "Field 'Product.title' has deprecation reason 'no longer needed'",
					Path:       "Product.title",
					ChangeType: "FIELD_DEPRECATION_REASON_ADDED",
					Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
						Level:  enterprisev1.GraphQLInspectorDiffOutput_NON_BREAKING,
						Reason: "",
					},
				},
			},
		}))
	})

	It("does not error on custom directives", func() {
		oldSchema := `
type Query {
	productsForHome: [Product] @cacheControl(maxAge: 30, inheritMaxAge: true, scope: private)
}

type Product {
    id: String
    title: String
}
`
		newSchema := `
type Query {
	productsForHome: [Product] @cacheControl(maxAge: 60, inheritMaxAge: false, scope: private)
}

type Product {
    id: String
	author: String
    title: String @resolve(name: "title")
}
`
		in := &enterprisev1.GraphQLInspectorDiffInput{
			OldSchema: oldSchema,
			NewSchema: newSchema,
		}
		out, err := runner.RunSchemaDiff(in)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(MatchProto(&enterprisev1.GraphQLInspectorDiffOutput{
			Changes: []*enterprisev1.GraphQLInspectorDiffOutput_Change{
				{
					Message:    "Field 'author' was added to object type 'Product'",
					Path:       "Product.author",
					ChangeType: "FIELD_ADDED",
					Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
						Level:  enterprisev1.GraphQLInspectorDiffOutput_NON_BREAKING,
						Reason: "",
					},
				},
			},
		}))
	})

	It("returns error on unknown directives", func() {
		schema := `
type Query {
	productsForHome: [Product] @myDirective(maxAge: 30, inheritMaxAge: true, scope: private)
}

type Product {
    id: String
    title: String
}
`
		in := &enterprisev1.GraphQLInspectorDiffInput{
			OldSchema: schema,
			NewSchema: schema,
		}
		_, err := runner.RunSchemaDiff(in)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Unknown directive \"@myDirective\""))
	})

	// These tests make sure we are passing through the diff rules correctly to graphql-inspector.
	Context("diff rules", func() {

		It("respects the RULE_DANGEROUS_TO_BREAKING rule", func() {
			oldSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
    id: String @deprecated(reason: "used to be deprecated")
}
`
			newSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
    id: String
}
`
			in := &enterprisev1.GraphQLInspectorDiffInput{
				OldSchema: oldSchema,
				NewSchema: newSchema,
				Rules: []gloov1.GraphqlOptions_SchemaChangeValidationOptions_ProcessingRule{
					gloov1.GraphqlOptions_SchemaChangeValidationOptions_RULE_DANGEROUS_TO_BREAKING,
				},
			}
			out, err := runner.RunSchemaDiff(in)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(MatchProto(&enterprisev1.GraphQLInspectorDiffOutput{
				Changes: []*enterprisev1.GraphQLInspectorDiffOutput_Change{
					{
						Message:    "Field 'Product.id' is no longer deprecated",
						Path:       "Product.id",
						ChangeType: "FIELD_DEPRECATION_REMOVED",
						Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
							// this is normally a dangerous change, but RULE_DANGEROUS_TO_BREAKING turns it into a breaking change
							Level:  enterprisev1.GraphQLInspectorDiffOutput_BREAKING,
							Reason: "",
						},
					},
					{
						Message:    "Deprecation reason was removed from field 'Product.id'",
						Path:       "Product.id",
						ChangeType: "FIELD_DEPRECATION_REASON_REMOVED",
						Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
							Level:  enterprisev1.GraphQLInspectorDiffOutput_NON_BREAKING,
							Reason: "",
						},
					},
				},
			}))
		})

		It("respects the RULE_DEPRECATED_FIELD_REMOVAL_DANGEROUS rule", func() {
			oldSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
    id: String @deprecated(reason: "no longer needed")
	author: String
}
`
			newSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
	author: String
}
`
			in := &enterprisev1.GraphQLInspectorDiffInput{
				OldSchema: oldSchema,
				NewSchema: newSchema,
				Rules: []gloov1.GraphqlOptions_SchemaChangeValidationOptions_ProcessingRule{
					gloov1.GraphqlOptions_SchemaChangeValidationOptions_RULE_DEPRECATED_FIELD_REMOVAL_DANGEROUS,
				},
			}
			out, err := runner.RunSchemaDiff(in)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(MatchProto(&enterprisev1.GraphQLInspectorDiffOutput{
				Changes: []*enterprisev1.GraphQLInspectorDiffOutput_Change{
					{
						Message:    "Field 'id' (deprecated) was removed from object type 'Product'",
						Path:       "Product.id",
						ChangeType: "FIELD_REMOVED",
						Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
							// this is normally a breaking change, but RULE_DEPRECATED_FIELD_REMOVAL_DANGEROUS turns it into a dangerous change
							Level:  enterprisev1.GraphQLInspectorDiffOutput_DANGEROUS,
							Reason: "Removing a deprecated field is a breaking change. Before removing it, you may want to look at the field's usage to see the impact of removing the field.",
						},
					},
				},
			}))
		})

		It("respects the RULE_IGNORE_DESCRIPTION_CHANGES rule", func() {
			oldSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
	"This is the old description"
	author: String
}
`
			newSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
	"This is the new description"
	author: String
}
`
			in := &enterprisev1.GraphQLInspectorDiffInput{
				OldSchema: oldSchema,
				NewSchema: newSchema,
				Rules: []gloov1.GraphqlOptions_SchemaChangeValidationOptions_ProcessingRule{
					gloov1.GraphqlOptions_SchemaChangeValidationOptions_RULE_IGNORE_DESCRIPTION_CHANGES,
				},
			}
			out, err := runner.RunSchemaDiff(in)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(MatchProto(&enterprisev1.GraphQLInspectorDiffOutput{
				// normally changing a field description is a non-breaking change, but with the above
				// rule it will be ignored and not shown in the list of changes
				Changes: []*enterprisev1.GraphQLInspectorDiffOutput_Change{},
			}))
		})

		It("can process multiple rules", func() {
			oldSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
    id: String @deprecated(reason: "used to be deprecated")
	getSomething(a: String): String
}
`
			newSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
	getSomething(a: String, b: String): String
}
`
			in := &enterprisev1.GraphQLInspectorDiffInput{
				OldSchema: oldSchema,
				NewSchema: newSchema,
				Rules: []gloov1.GraphqlOptions_SchemaChangeValidationOptions_ProcessingRule{
					gloov1.GraphqlOptions_SchemaChangeValidationOptions_RULE_DANGEROUS_TO_BREAKING,
					gloov1.GraphqlOptions_SchemaChangeValidationOptions_RULE_DEPRECATED_FIELD_REMOVAL_DANGEROUS,
				},
			}
			out, err := runner.RunSchemaDiff(in)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(MatchProto(&enterprisev1.GraphQLInspectorDiffOutput{
				Changes: []*enterprisev1.GraphQLInspectorDiffOutput_Change{
					{
						Message:    "Argument 'b: String' added to field 'Product.getSomething'",
						Path:       "Product.getSomething.b",
						ChangeType: "FIELD_ARGUMENT_ADDED",
						Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
							// this is normally a dangerous change, but RULE_DANGEROUS_TO_BREAKING turns it into a breaking change
							Level:  enterprisev1.GraphQLInspectorDiffOutput_BREAKING,
							Reason: "Adding a new argument to an existing field may involve a change in resolve function logic that potentially may cause some side effects.",
						},
					},
					{
						Message:    "Field 'id' (deprecated) was removed from object type 'Product'",
						Path:       "Product.id",
						ChangeType: "FIELD_REMOVED",
						Criticality: &enterprisev1.GraphQLInspectorDiffOutput_Criticality{
							// this is normally a breaking change, but RULE_DEPRECATED_FIELD_REMOVAL_DANGEROUS turns it into a dangerous change
							Level:  enterprisev1.GraphQLInspectorDiffOutput_DANGEROUS,
							Reason: "Removing a deprecated field is a breaking change. Before removing it, you may want to look at the field's usage to see the impact of removing the field.",
						},
					},
				},
			}))
		})

	})
})
