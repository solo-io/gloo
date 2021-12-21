package graphql_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/solo-kit/test/matchers"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql"
)

var _ = Describe("Rest Resolver Test", func() {
	Context("Translates string extraction to correct value provider", func() {
		It("creates extractions only when necessary", func() {
			extractions := map[string]string{
				"no providers": "this must not have any providers",
				"one provider": "this must have one provider here: {$parent.id}, {$parent.id}",
			}
			vp, err := graphql.TranslateStringValueProviderMap(extractions)
			Expect(err).NotTo(HaveOccurred())
			Expect(vp).To(HaveLen(2))
			Expect(vp["no providers"].GetProviders()[graphql.ARBITRARY_PROVIDER_NAME].GetTypedProvider()).To(
				matchers.MatchProto(
					&v2.ValueProvider_TypedValueProvider{
						ValProvider: &v2.ValueProvider_TypedValueProvider_Value{
							Value: "this must not have any providers",
						},
					},
				))
			Expect(vp["one provider"].GetProviders()).To(HaveLen(1))
			Expect(vp["one provider"]).To(
				matchers.MatchProto(
					&v2.ValueProvider{
						Providers: map[string]*v2.ValueProvider_Provider{
							"parentid": {
								Provider: &v2.ValueProvider_Provider_GraphqlParent{
									GraphqlParent: &v2.ValueProvider_GraphQLParentExtraction{
										Path: []*v2.PathSegment{
											{Segment: &v2.PathSegment_Key{Key: "id"}},
										},
									},
								},
							},
						},
						ProviderTemplate: "this must have one provider here: {parentid}, {parentid}",
					},
				))
		})
	})
})
