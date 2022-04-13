package translation_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/solo-kit/test/matchers"
	enterprisev1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/enterprise/graphql/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/translation"
)

var _ = Describe("SchemaDiff", func() {

	BeforeEach(func() {
		err := os.Setenv(translation.GraphqlJsRootEnvVar, "../../../plugins/graphql/js/")
		Expect(err).NotTo(HaveOccurred())
		err = os.Setenv(translation.GraphqlProtoRootEnvVar, "../../../../..//ui/src/proto/")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.Unsetenv(translation.GraphqlProtoRootEnvVar)).NotTo(HaveOccurred())
		Expect(os.Unsetenv(translation.GraphqlJsRootEnvVar)).NotTo(HaveOccurred())
	})

	It("returns empty diff when both schemas are empty", func() {
		in := &enterprisev1.GraphQLInspectorDiffInput{
			OldSchema: "",
			NewSchema: "",
		}
		out, err := translation.GetSchemaDiff(in)
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
		out, err := translation.GetSchemaDiff(in)
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
		_, err := translation.GetSchemaDiff(in)
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
		out, err := translation.GetSchemaDiff(in)
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
		out, err := translation.GetSchemaDiff(in)
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
		out, err := translation.GetSchemaDiff(in)
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
		out, err := translation.GetSchemaDiff(in)
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
		_, err := translation.GetSchemaDiff(in)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Unknown directive \"@myDirective\""))
	})
})
