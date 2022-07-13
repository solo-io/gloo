package validation_test

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/validation"
)

var _ = Describe("SchemaValidation", func() {

	var (
		settings *gloo_v1.Settings
	)

	BeforeEach(func() {
		settings = &gloo_v1.Settings{
			Spec: gloo_v1.SettingsSpec{
				GraphqlOptions: &gloo_v1.GraphqlOptions{
					SchemaChangeValidationOptions: &gloo_v1.GraphqlOptions_SchemaChangeValidationOptions{
						RejectBreakingChanges: &wrappers.BoolValue{
							Value: true,
						},
					},
				},
			},
		}
	})

	It("allows schema update with non-breaking and dangerous changes", func() {
		oldSchema := `
type Query {
	productsForHome: [Product]
}

type Product {
    id: String @deprecated(reason: "used to be deprecated")
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
	author: String
}
`
		// adding a new field (author) is a non-breaking change
		// un-deprecating a field (id) is a dangerous change
		// both of these changes should be accepted without error
		err := validation.ValidateSchemaUpdate(oldSchema, newSchema, settings)
		Expect(err).NotTo(HaveOccurred())
	})
	It("rejects schema update with breaking changes if rejectBreakingChanges is true", func() {
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
		// removing a field (author) is a breaking change and is rejected
		err := validation.ValidateSchemaUpdate(oldSchema, newSchema, settings)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Field 'author' was removed from object type 'Product'"))
	})
	It("allows schema update with breaking changes if rejectBreakingChanges is false", func() {
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
		// removing a field (author) is a breaking change, but is not rejected
		settings.Spec.GraphqlOptions.SchemaChangeValidationOptions.RejectBreakingChanges.Value = false
		err := validation.ValidateSchemaUpdate(oldSchema, newSchema, settings)
		Expect(err).NotTo(HaveOccurred())
	})
	It("rejects schema update with dangerous changes if rejectBreakingChanges and dangerousToBreaking are true", func() {
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
		// un-deprecating a field is normally a dangerous change, but we treat it as breaking and reject it
		settings.Spec.GraphqlOptions.SchemaChangeValidationOptions.ProcessingRules = []gloo_v1.GraphqlOptions_SchemaChangeValidationOptions_ProcessingRule{
			gloo_v1.GraphqlOptions_SchemaChangeValidationOptions_RULE_DANGEROUS_TO_BREAKING,
		}
		err := validation.ValidateSchemaUpdate(oldSchema, newSchema, settings)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Field 'Product.id' is no longer deprecated"))
	})
})
