package rest_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/solo-kit/test/matchers"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/resolvers/rest"
	resolver_utils "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/resolvers/utils"
)

var _ = Describe("Rest Resolver Test", func() {
	Context("Translates string extraction to correct value provider", func() {
		It("creates extractions only when necessary", func() {
			extractions := map[string]string{
				"no providers": "this must not have any providers",
				"one provider": "this must have one provider here: {$parent.id}, {$parent.id}",
			}
			vp, err := resolver_utils.TranslateStringValueProviderMap(extractions)
			Expect(err).NotTo(HaveOccurred())
			Expect(vp).To(HaveLen(2))
			Expect(vp["no providers"].GetProviders()[resolver_utils.ARBITRARY_PROVIDER_NAME].GetTypedProvider()).To(
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

	Context("Translates string setter to correct templated path", func() {
		It("translates single response with no interpolation", func() {
			templatedPath, err := rest.TranslateSetter("{$body[0][*].details.firstname}")
			Expect(err).NotTo(HaveOccurred())
			// This should be empty because we don't want string interpolation, we want the actual body.details.firstname object
			// whether it is a string, list, or object and we don't want it coerced into a string type
			Expect(templatedPath.PathTemplate).To(Equal(""))
			Expect(templatedPath.NamedPaths).To(HaveKeyWithValue("body[0][*]detailsfirstname", &v2.Path{
				Segments: []*v2.PathSegment{
					{
						Segment: &v2.PathSegment_Index{
							Index: 0,
						},
					},
					{
						Segment: &v2.PathSegment_All{
							All: true,
						},
					},
					{
						Segment: &v2.PathSegment_Key{
							Key: "details",
						},
					},
					{
						Segment: &v2.PathSegment_Key{
							Key: "firstname",
						},
					},
				}}))
		})

		It("translates templated response with interpolation", func() {
			templatedPath, err := rest.TranslateSetter("fullname: {$body.details.firstname} {$body.details.lastname}")
			Expect(err).NotTo(HaveOccurred())
			Expect(templatedPath.PathTemplate).To(Equal("fullname: {bodydetailsfirstname} {bodydetailslastname}"))
			Expect(templatedPath.NamedPaths).To(HaveLen(2))
			Expect(templatedPath.NamedPaths).To(HaveKeyWithValue("bodydetailsfirstname", &v2.Path{
				Segments: []*v2.PathSegment{
					{
						Segment: &v2.PathSegment_Key{
							Key: "details",
						},
					},
					{
						Segment: &v2.PathSegment_Key{
							Key: "firstname",
						},
					},
				}}))
			Expect(templatedPath.NamedPaths).To(HaveKeyWithValue("bodydetailslastname", &v2.Path{
				Segments: []*v2.PathSegment{
					{
						Segment: &v2.PathSegment_Key{
							Key: "details",
						},
					},
					{
						Segment: &v2.PathSegment_Key{
							Key: "lastname",
						},
					},
				}}))
		})

	})
})
