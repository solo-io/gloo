package extauth_test

import (
	envoycore "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	envoytransformation "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	glooTransformation "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
	extauthgloo "github.com/solo-io/gloo/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/test/matchers"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

var _ = Describe("ExtAuthzConfigGenerator", func() {

	var (
		extAuthzConfigGenerator extauthgloo.ExtAuthzConfigGenerator
		defaultSettings         *extauthv1.Settings
		namedSettings           map[string]*extauthv1.Settings
	)

	Context("GenerateListenerExtAuthzConfig", func() {

		Context("DefaultConfigGenerator (single ext_authz filter)", func() {

			JustBeforeEach(func() {
				extAuthzConfigGenerator = extauth.NewEnterpriseDefaultConfigGenerator(defaultSettings)
			})

			// The EnterpriseDefaultConfigGenerator just delegates generation to the OpenSource ConfigGenerator
			// Smoke test to ensure that basic settings are configured.

			When("settings point to a valid ext auth server", func() {

				var (
					upstream       *gloov1.Upstream
					expectedConfig *envoyauth.ExtAuthz
				)

				BeforeEach(func() {
					upstream = createExtAuthUpstream("extauth")
					usRef := upstream.Metadata.Ref()

					defaultSettings = &extauthv1.Settings{
						ExtauthzServerRef: usRef,
					}

					expectedConfig = &envoyauth.ExtAuthz{
						TransportApiVersion:       envoycore.ApiVersion_V3,
						MetadataContextNamespaces: []string{extauthgloo.JWTFilterName},
						Services: &envoyauth.ExtAuthz_GrpcService{
							GrpcService: &envoycore.GrpcService{
								Timeout: extauthgloo.DefaultTimeout,
								TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
									EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
										ClusterName: translator.UpstreamToClusterName(usRef),
									},
								},
							},
						},
					}
				})

				It("uses the expected defaults", func() {
					filters, err := extAuthzConfigGenerator.GenerateListenerExtAuthzConfig(nil, gloov1.UpstreamList{upstream})
					Expect(err).NotTo(HaveOccurred())
					Expect(filters).To(HaveLen(1))

					actualFilterConfig := filters[0]
					Expect(actualFilterConfig).To(matchers.MatchProto(expectedConfig))
				})

			})

		})

		Context("MultiConfigGenerator (multiple ext_authz filters)", func() {

			var (
				upstreams gloov1.UpstreamList
			)

			BeforeEach(func() {
				extAuthUpstreamDefault := createExtAuthUpstream("extauth-default")
				extAuthUpstreamNamed := createExtAuthUpstream("extauth-named")

				defaultSettings = createExtAuthSettings(extAuthUpstreamDefault)
				namedSettings = map[string]*extauthv1.Settings{
					"named": createExtAuthSettings(extAuthUpstreamNamed),
				}
				upstreams = gloov1.UpstreamList{
					extAuthUpstreamDefault,
					extAuthUpstreamNamed,
				}
			})

			JustBeforeEach(func() {
				extAuthzConfigGenerator = extauth.NewEnterpriseMultiConfigGenerator(defaultSettings, namedSettings)
			})

			When("defaultSettings are nil", func() {

				BeforeEach(func() {
					defaultSettings = nil
				})

				It("does not generate any configurations", func() {
					filters, err := extAuthzConfigGenerator.GenerateListenerExtAuthzConfig(nil, upstreams)
					Expect(err).NotTo(HaveOccurred())
					Expect(filters).To(HaveLen(0))
				})
			})

			When("valid defaultSettings and namedSettings are provided", func() {

				// Configurations are not always returned in the same order. Therefore, we explicitly look
				// for filters with a particular statPrefix
				getExtAuthzWithStatPrefix := func(configurations []*envoyauth.ExtAuthz, statPrefix string) *envoyauth.ExtAuthz {
					for _, config := range configurations {
						if config.StatPrefix == statPrefix {
							return config
						}
					}
					return nil
				}

				It("converts default extauth settings to ext_authz filter", func() {
					filters, err := extAuthzConfigGenerator.GenerateListenerExtAuthzConfig(nil, upstreams)
					Expect(err).NotTo(HaveOccurred())
					Expect(filters).To(HaveLen(2)) // 1 in default settings, 1 in named settings

					defaultFilterConfig := getExtAuthzWithStatPrefix(filters, "extauth-default")

					Expect(defaultFilterConfig.MetadataContextNamespaces).Should(ContainElement(extauth.SanitizeFilterName))
					Expect(defaultFilterConfig.FilterEnabledMetadata).To(matchers.MatchProto(&envoymatcher.MetadataMatcher{
						Filter: extauth.SanitizeFilterName,
						Path: []*envoymatcher.MetadataMatcher_PathSegment{
							{
								Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{
									Key: extauth.CustomAuthServerNameMetadataKey,
								},
							},
						},
						Value: &envoymatcher.ValueMatcher{
							MatchPattern: &envoymatcher.ValueMatcher_StringMatch{
								StringMatch: &envoymatcher.StringMatcher{
									MatchPattern: &envoymatcher.StringMatcher_Exact{
										Exact: extauth.DefaultAuthServiceName,
									},
								},
							},
						},
					}))
				})

				It("converts default extauth settings to ext_authz filter with namespaces added from listener ", func() {
					filters, err := extAuthzConfigGenerator.GenerateListenerExtAuthzConfig(createListener(), upstreams)
					Expect(err).NotTo(HaveOccurred())
					Expect(filters).To(HaveLen(2)) // 1 in default settings, 1 in named settings

					//MetaDataNameSpaces
					for _, filter := range filters {
						Expect(len(filter.MetadataContextNamespaces)).To(Equal(5))
						Expect(filter.MetadataContextNamespaces).To(ContainElements("envoy.filters.http.jwt_authn", "io.solo.filters.http.sanitize", "namespace1", "namespace2", "namespace3"))
					}

					defaultFilterConfig := getExtAuthzWithStatPrefix(filters, "extauth-default")

					Expect(defaultFilterConfig.MetadataContextNamespaces).Should(ContainElement(extauth.SanitizeFilterName))
					Expect(defaultFilterConfig.FilterEnabledMetadata).To(matchers.MatchProto(&envoymatcher.MetadataMatcher{
						Filter: extauth.SanitizeFilterName,
						Path: []*envoymatcher.MetadataMatcher_PathSegment{
							{
								Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{
									Key: extauth.CustomAuthServerNameMetadataKey,
								},
							},
						},
						Value: &envoymatcher.ValueMatcher{
							MatchPattern: &envoymatcher.ValueMatcher_StringMatch{
								StringMatch: &envoymatcher.StringMatcher{
									MatchPattern: &envoymatcher.StringMatcher_Exact{
										Exact: extauth.DefaultAuthServiceName,
									},
								},
							},
						},
					}))
				})

				It("converts named extauth settings to ext_authz filter", func() {
					filters, err := extAuthzConfigGenerator.GenerateListenerExtAuthzConfig(nil, upstreams)
					Expect(err).NotTo(HaveOccurred())
					Expect(filters).To(HaveLen(2)) // 1 in default settings, 1 in named settings

					namedFilterConfig := getExtAuthzWithStatPrefix(filters, "extauth-named")

					Expect(namedFilterConfig.MetadataContextNamespaces).Should(ContainElement(extauth.SanitizeFilterName))
					Expect(namedFilterConfig.FilterEnabledMetadata).To(matchers.MatchProto(&envoymatcher.MetadataMatcher{
						Filter: extauth.SanitizeFilterName,
						Path: []*envoymatcher.MetadataMatcher_PathSegment{
							{
								Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{
									Key: extauth.CustomAuthServerNameMetadataKey,
								},
							},
						},
						Value: &envoymatcher.ValueMatcher{
							MatchPattern: &envoymatcher.ValueMatcher_StringMatch{
								StringMatch: &envoymatcher.StringMatcher{
									MatchPattern: &envoymatcher.StringMatcher_Exact{
										Exact: "named",
									},
								},
							},
						},
					}))
				})

			})

		})

	})
})

func createExtAuthUpstream(name string) *gloov1.Upstream {
	return &gloov1.Upstream{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: "default",
		},
		UpstreamType: &gloov1.Upstream_Static{
			Static: &static.UpstreamSpec{
				Hosts: []*static.Host{{
					Addr: name,
					Port: 1234,
				}},
			},
		},
	}
}

func createExtAuthSettings(upstream *gloov1.Upstream) *extauthv1.Settings {
	return &extauthv1.Settings{
		ExtauthzServerRef: upstream.Metadata.Ref(),
		StatPrefix:        upstream.Metadata.Name,
	}
}

func createListener() *gloov1.HttpListener {
	return &gloov1.HttpListener{
		VirtualHosts: []*gloov1.VirtualHost{
			{
				Domains: []string{
					"*",
				},
				Options: &gloov1.VirtualHostOptions{
					StagedTransformations: &glooTransformation.TransformationStages{
						Early: &glooTransformation.RequestResponseTransformations{
							RequestTransforms: []*glooTransformation.RequestMatch{
								{
									RequestTransformation: &glooTransformation.Transformation{
										TransformationType: &glooTransformation.Transformation_TransformationTemplate{
											TransformationTemplate: &envoytransformation.TransformationTemplate{
												DynamicMetadataValues: []*envoytransformation.TransformationTemplate_DynamicMetadataValue{
													{
														Key:               "key1",
														MetadataNamespace: "namespace1",
														Value: &envoytransformation.InjaTemplate{
															Text: "test1",
														},
													},
													{
														Key:               "key2",
														MetadataNamespace: "namespace2",
														Value: &envoytransformation.InjaTemplate{
															Text: "test2",
														},
													},
													{
														Key:               "Key3",
														MetadataNamespace: "namespace3",
														Value: &envoytransformation.InjaTemplate{
															Text: "test3",
														},
													},
													{
														Key:               "Key4",
														MetadataNamespace: "namespace3", //duplicate to make sure dupes are removed
														Value: &envoytransformation.InjaTemplate{
															Text: "test4",
														},
													},
												},
											},
										},
									},
								},
							}, // reqTransform
						}, // regular
					}, // stagedTransform
				},
			},
		},
	}
}
