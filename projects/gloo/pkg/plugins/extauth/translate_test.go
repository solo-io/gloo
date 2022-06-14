package extauth_test

import (
	"context"
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"

	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	static_plugin_gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("Translate", func() {

	var (
		params           plugins.Params
		virtualHost      *v1.VirtualHost
		upstream         *v1.Upstream
		secret           *v1.Secret
		route            *v1.Route
		authConfig       *extauth.AuthConfig
		authConfigRef    *core.ResourceRef
		extAuthExtension *extauth.ExtAuthExtension
		clientSecret     *extauth.OauthSecret
		apiKeySecret     *extauth.ApiKeySecret
	)

	BeforeEach(func() {

		upstream = &v1.Upstream{
			Metadata: &core.Metadata{
				Name:      "extauth",
				Namespace: "default",
			},
			UpstreamType: &v1.Upstream_Static{
				Static: &static_plugin_gloo.UpstreamSpec{
					Hosts: []*static_plugin_gloo.Host{{
						Addr: "test",
						Port: 1234,
					}},
				},
			},
		}
		route = &v1.Route{
			Matchers: []*matchers.Matcher{{
				PathSpecifier: &matchers.Matcher_Prefix{
					Prefix: "/",
				},
			}},
			Action: &v1.Route_RouteAction{
				RouteAction: &v1.RouteAction{
					Destination: &v1.RouteAction_Single{
						Single: &v1.Destination{
							DestinationType: &v1.Destination_Upstream{
								Upstream: upstream.Metadata.Ref(),
							},
						},
					},
				},
			},
		}

		apiKeySecret = &extauth.ApiKeySecret{
			ApiKey: "apiKey1",
		}

		clientSecret = &extauth.OauthSecret{
			ClientSecret: "1234",
		}

		secret = &v1.Secret{
			Metadata: &core.Metadata{
				Name:      "secret",
				Namespace: "default",
			},
			Kind: &v1.Secret_Oauth{
				Oauth: clientSecret,
			},
		}
		secretRef := secret.Metadata.Ref()
		authConfig = &extauth.AuthConfig{
			Metadata: &core.Metadata{
				Name:      "oauth",
				Namespace: "gloo-system",
			},
			Configs: []*extauth.AuthConfig_Config{{
				AuthConfig: &extauth.AuthConfig_Config_Oauth2{
					Oauth2: &extauth.OAuth2{
						OauthType: &extauth.OAuth2_OidcAuthorizationCode{

							OidcAuthorizationCode: &extauth.OidcAuthorizationCode{
								ClientSecretRef:          secretRef,
								ClientId:                 "ClientId",
								IssuerUrl:                "IssuerUrl",
								AuthEndpointQueryParams:  map[string]string{"test": "additional_auth_query_params"},
								TokenEndpointQueryParams: map[string]string{"test": "additional_token_query_params"},
								AppUrl:                   "AppUrl",
								CallbackPath:             "CallbackPath",
							},
						},
					},
				},
			}},
		}
		authConfigRef = authConfig.Metadata.Ref()
		extAuthExtension = &extauth.ExtAuthExtension{
			Spec: &extauth.ExtAuthExtension_ConfigRef{
				ConfigRef: authConfigRef,
			},
		}

		params.Snapshot = &v1snap.ApiSnapshot{
			Upstreams:   v1.UpstreamList{upstream},
			AuthConfigs: extauth.AuthConfigList{authConfig},
		}
	})

	JustBeforeEach(func() {

		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			Options: &v1.VirtualHostOptions{
				Extauth: extAuthExtension,
			},
			Routes: []*v1.Route{route},
		}

		proxy := &v1.Proxy{
			Metadata: &core.Metadata{
				Name:      "secret",
				Namespace: "default",
			},
			Listeners: []*v1.Listener{{
				Name: "default",
				ListenerType: &v1.Listener_HttpListener{
					HttpListener: &v1.HttpListener{
						VirtualHosts: []*v1.VirtualHost{virtualHost},
					},
				},
			}},
		}

		params.Snapshot.Proxies = v1.ProxyList{proxy}
		params.Snapshot.Secrets = v1.SecretList{secret}
	})

	It("should translate oauth config for extauth server", func() {
		translated, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
		Expect(err).NotTo(HaveOccurred())
		Expect(translated.AuthConfigRefName).To(Equal(authConfigRef.Key()))
		Expect(translated.Configs).To(HaveLen(1))
		actual := translated.Configs[0].GetOauth2()
		expected := authConfig.Configs[0].GetOauth2()
		Expect(actual.GetOidcAuthorizationCode().IssuerUrl).To(Equal(expected.GetOidcAuthorizationCode().IssuerUrl))
		Expect(actual.GetOidcAuthorizationCode().AuthEndpointQueryParams).To(Equal(expected.GetOidcAuthorizationCode().AuthEndpointQueryParams))
		Expect(actual.GetOidcAuthorizationCode().TokenEndpointQueryParams).To(Equal(expected.GetOidcAuthorizationCode().TokenEndpointQueryParams))
		Expect(actual.GetOidcAuthorizationCode().ClientId).To(Equal(expected.GetOidcAuthorizationCode().ClientId))
		Expect(actual.GetOidcAuthorizationCode().ClientSecret).To(Equal(clientSecret.ClientSecret))
		Expect(actual.GetOidcAuthorizationCode().AppUrl).To(Equal(expected.GetOidcAuthorizationCode().AppUrl))
		Expect(actual.GetOidcAuthorizationCode().CallbackPath).To(Equal(expected.GetOidcAuthorizationCode().CallbackPath))
	})

	It("will fail if the oidc auth proto has a new top level field", func() {
		// This test is important as it checks whether the oidc auth code proto have a new top level field.
		// This should happen very rarely, and should be used as an indication that the `translateOidcAuthorizationCode` function
		// most likely needs to change.

		Expect(reflect.TypeOf(extauth.ExtAuthConfig_OidcAuthorizationCodeConfig{}).NumField()).To(
			Equal(21),
			"wrong number of fields found",
		)
	})

	Context("with api key extauth", func() {
		BeforeEach(func() {

			secret = &v1.Secret{
				Metadata: &core.Metadata{
					Name:      "secretName",
					Namespace: "default",
					Labels:    map[string]string{"team": "infrastructure"},
				},
				Kind: &v1.Secret_ApiKey{
					ApiKey: apiKeySecret,
				},
			}
			secretRef := secret.Metadata.Ref()

			authConfig = &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      "apikey",
					Namespace: "gloo-system",
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_ApiKeyAuth{
						ApiKeyAuth: &extauth.ApiKeyAuth{
							HeaderName:       "x-api-key",
							ApiKeySecretRefs: []*core.ResourceRef{secretRef},
						},
					},
				}},
			}
			authConfigRef = authConfig.Metadata.Ref()
			extAuthExtension = &extauth.ExtAuthExtension{
				Spec: &extauth.ExtAuthExtension_ConfigRef{
					ConfigRef: authConfigRef,
				},
			}

			params.Snapshot = &v1snap.ApiSnapshot{
				AuthConfigs: extauth.AuthConfigList{authConfig},
			}
		})

		Context("secret is malformed", func() {
			It("returns expected error when secret is not of API key type", func() {
				secret.Kind = &v1.Secret_Aws{}
				_, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(NonApiKeySecretError(secret).Error())))
			})

			It("returns expected error when the secret does not contain an API key", func() {
				secret.GetApiKey().ApiKey = ""
				_, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(EmptyApiKeyError(secret).Error())))
			})
		})

		Context("with api key extauth, secret ref matching", func() {
			It("should translate api keys config for extauth server - matching secret ref", func() {
				translated, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
				Expect(err).NotTo(HaveOccurred())
				Expect(translated.AuthConfigRefName).To(Equal(authConfigRef.Key()))
				Expect(translated.Configs).To(HaveLen(1))
				actual := translated.Configs[0].GetApiKeyAuth()
				Expect(actual).To(Equal(&extauth.ExtAuthConfig_ApiKeyAuthConfig{
					HeaderName: "x-api-key",
					ValidApiKeys: map[string]*extauth.ExtAuthConfig_ApiKeyAuthConfig_KeyMetadata{
						"apiKey1": {
							Username: "secretName",
						},
					},
				}))
			})

			It("should translate api keys config for extauth server - mismatching secret ref", func() {
				secret.Metadata.Name = "mismatchName"
				_, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("list did not find secret"))
			})
		})

		Describe("API keys with metadata", func() {

			BeforeEach(func() {
				secret = &v1.Secret{
					Metadata: &core.Metadata{
						Name:      "secretName",
						Namespace: "default",
						Labels:    map[string]string{"team": "infrastructure"},
					},
					Kind: &v1.Secret_ApiKey{
						ApiKey: &extauth.ApiKeySecret{
							ApiKey: "apiKey1",
							Metadata: map[string]string{
								"user-id": "123",
							},
						},
					},
				}
				secretRef := secret.Metadata.Ref()
				authConfig = &extauth.AuthConfig{
					Metadata: &core.Metadata{
						Name:      "apikey",
						Namespace: "gloo-system",
					},
					Configs: []*extauth.AuthConfig_Config{{
						AuthConfig: &extauth.AuthConfig_Config_ApiKeyAuth{
							ApiKeyAuth: &extauth.ApiKeyAuth{
								HeaderName:       "x-api-key",
								ApiKeySecretRefs: []*core.ResourceRef{secretRef},
								HeadersFromMetadata: map[string]*extauth.ApiKeyAuth_SecretKey{
									"x-user-id": {
										Name:     "user-id",
										Required: true,
									},
								},
							},
						},
					}},
				}
				authConfigRef = authConfig.Metadata.Ref()
				extAuthExtension = &extauth.ExtAuthExtension{
					Spec: &extauth.ExtAuthExtension_ConfigRef{
						ConfigRef: authConfigRef,
					},
				}

				params.Snapshot = &v1snap.ApiSnapshot{
					AuthConfigs: extauth.AuthConfigList{authConfig},
				}
			})

			It("should fail if required metadata is missing on the secret", func() {
				secret.GetApiKey().Metadata = nil

				_, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(MissingRequiredMetadataError("user-id", secret).Error())))
			})

			It("should include secret metadata in the API key metadata", func() {
				translated, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
				Expect(err).NotTo(HaveOccurred())
				Expect(translated.AuthConfigRefName).To(Equal(authConfigRef.Key()))
				Expect(translated.Configs).To(HaveLen(1))
				actual := translated.Configs[0].GetApiKeyAuth()
				Expect(actual).To(Equal(&extauth.ExtAuthConfig_ApiKeyAuthConfig{
					HeaderName: "x-api-key",
					ValidApiKeys: map[string]*extauth.ExtAuthConfig_ApiKeyAuthConfig_KeyMetadata{
						"apiKey1": {
							Username: "secretName",
							Metadata: map[string]string{
								"user-id": "123",
							},
						},
					},
					HeadersFromKeyMetadata: map[string]string{
						"x-user-id": "user-id",
					},
				}))
			})

			When("metadata is not required", func() {
				It("does not fail if the secret does not contain the metadata", func() {
					secret.GetApiKey().Metadata = nil
					authConfig.GetConfigs()[0].GetApiKeyAuth().GetHeadersFromMetadata()["x-user-id"].Required = false

					translated, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
					Expect(err).NotTo(HaveOccurred())
					Expect(translated.AuthConfigRefName).To(Equal(authConfigRef.Key()))
					Expect(translated.Configs).To(HaveLen(1))
					actual := translated.Configs[0].GetApiKeyAuth()
					Expect(actual).To(Equal(&extauth.ExtAuthConfig_ApiKeyAuthConfig{
						HeaderName: "x-api-key",
						ValidApiKeys: map[string]*extauth.ExtAuthConfig_ApiKeyAuthConfig_KeyMetadata{
							"apiKey1": {
								Username: "secretName",
							},
						},
						HeadersFromKeyMetadata: map[string]string{
							"x-user-id": "user-id",
						},
					}))
				})
			})
		})

		Context("with api key ext auth, label matching", func() {
			BeforeEach(func() {
				authConfig = &extauth.AuthConfig{
					Metadata: &core.Metadata{
						Name:      "apikey",
						Namespace: "gloo-system",
					},
					Configs: []*extauth.AuthConfig_Config{{
						AuthConfig: &extauth.AuthConfig_Config_ApiKeyAuth{
							ApiKeyAuth: &extauth.ApiKeyAuth{
								LabelSelector: map[string]string{"team": "infrastructure"},
							},
						},
					}},
				}
				authConfigRef = authConfig.Metadata.Ref()
				extAuthExtension = &extauth.ExtAuthExtension{
					Spec: &extauth.ExtAuthExtension_ConfigRef{
						ConfigRef: authConfigRef,
					},
				}

				params.Snapshot = &v1snap.ApiSnapshot{
					AuthConfigs: extauth.AuthConfigList{authConfig},
				}
			})

			It("should translate api keys config for extauth server - matching label", func() {
				translated, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
				Expect(err).NotTo(HaveOccurred())
				Expect(translated.AuthConfigRefName).To(Equal(authConfigRef.Key()))
				Expect(translated.Configs).To(HaveLen(1))
				actual := translated.Configs[0].GetApiKeyAuth()
				Expect(actual.ValidApiKeys).To(Equal(map[string]*extauth.ExtAuthConfig_ApiKeyAuthConfig_KeyMetadata{
					"apiKey1": {
						Username: "secretName",
					},
				}))
			})

			Context("should translate apikeys config for extauth server", func() {

				It("should not error - mismatched labels", func() {
					secret.Metadata.Labels = map[string]string{"missingLabel": "missingValue"}
					_, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should not error - empty labels", func() {
					secret.Metadata.Labels = map[string]string{}
					_, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should not error - nil labels", func() {
					secret.Metadata.Labels = nil
					_, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
					Expect(err).NotTo(HaveOccurred())
				})

			})

		})
	})

	Context("with OPA extauth", func() {
		BeforeEach(func() {
			authConfig = &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      "oauth",
					Namespace: "gloo-system",
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_OpaAuth{
						OpaAuth: &extauth.OpaAuth{
							Modules: []*core.ResourceRef{{Namespace: "namespace", Name: "name"}},
							Query:   "true",
						},
					},
				}},
			}
			authConfigRef = authConfig.Metadata.Ref()
			extAuthExtension = &extauth.ExtAuthExtension{
				Spec: &extauth.ExtAuthExtension_ConfigRef{
					ConfigRef: authConfigRef,
				},
			}

			params.Snapshot = &v1snap.ApiSnapshot{
				Upstreams:   v1.UpstreamList{upstream},
				AuthConfigs: extauth.AuthConfigList{authConfig},
			}

			params.Snapshot.Artifacts = v1.ArtifactList{
				{

					Metadata: &core.Metadata{
						Name:      "name",
						Namespace: "namespace",
					},
					Data: map[string]string{"module.rego": "package foo"},
				},
			}
		})

		It("should translate OPA config without options specified", func() {
			translated, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
			Expect(err).NotTo(HaveOccurred())
			Expect(translated.AuthConfigRefName).To(Equal(authConfigRef.Key()))
			Expect(translated.Configs).To(HaveLen(1))
			actual := translated.Configs[0].GetOpaAuth()
			expected := authConfig.Configs[0].GetOpaAuth()
			Expect(actual.Query).To(Equal(expected.Query))
			data := params.Snapshot.Artifacts[0].Data
			Expect(actual.Modules).To(Equal(data))
			Expect(actual.Options).To(Equal(expected.Options))
		})

		It("Should translate OPA config with options specified", func() {
			// Specify additional options in Opa Config.
			opaAuth := authConfig.Configs[0].GetOpaAuth()
			opaAuth.Options = &extauth.OpaAuthOptions{
				FastInputConversion: true,
			}

			translated, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
			Expect(err).NotTo(HaveOccurred())
			Expect(translated.AuthConfigRefName).To(Equal(authConfigRef.Key()))
			Expect(translated.Configs).To(HaveLen(1))
			Expect(translated.Configs[0].GetOpaAuth().GetOptions().GetFastInputConversion()).To(Equal(true))
			actual := translated.Configs[0].GetOpaAuth()
			expected := authConfig.Configs[0].GetOpaAuth()
			Expect(actual.Query).To(Equal(expected.Query))
			data := params.Snapshot.Artifacts[0].Data
			Expect(actual.Modules).To(Equal(data))
			Expect(actual.Options).To(Equal(expected.Options))
		})
	})

	Context("with AccessTokenValidation extauth", func() {
		BeforeEach(func() {
			authConfig = &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      "oauth",
					Namespace: "gloo-system",
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_Oauth2{
						Oauth2: &extauth.OAuth2{
							OauthType: &extauth.OAuth2_AccessTokenValidation{
								AccessTokenValidation: &extauth.AccessTokenValidation{
									ValidationType: &extauth.AccessTokenValidation_Introspection{
										Introspection: &extauth.IntrospectionValidation{
											IntrospectionUrl:    "introspection-url",
											ClientId:            "client-id",
											ClientSecretRef:     secret.Metadata.Ref(),
											UserIdAttributeName: "sub",
										},
									},
									CacheTimeout: ptypes.DurationProto(time.Minute),
									UserinfoUrl:  "user-info-url",
									ScopeValidation: &extauth.AccessTokenValidation_RequiredScopes{
										RequiredScopes: &extauth.AccessTokenValidation_ScopeList{
											Scope: []string{"foo", "bar"},
										},
									},
								},
							},
						},
					},
				}},
			}
			authConfigRef = authConfig.Metadata.Ref()
			extAuthExtension = &extauth.ExtAuthExtension{
				Spec: &extauth.ExtAuthExtension_ConfigRef{
					ConfigRef: authConfigRef,
				},
			}

			params.Snapshot = &v1snap.ApiSnapshot{
				Upstreams:   v1.UpstreamList{upstream},
				AuthConfigs: extauth.AuthConfigList{authConfig},
			}

		})

		It("should succeed for IntrospectionValidation config", func() {
			translated, err := TranslateExtAuthConfig(context.TODO(), params.Snapshot, authConfigRef)
			Expect(err).NotTo(HaveOccurred())

			Expect(translated.AuthConfigRefName).To(Equal(authConfigRef.Key()))
			Expect(translated.Configs).To(HaveLen(1))

			actual := translated.Configs[0].GetOauth2().GetAccessTokenValidationConfig()
			expected := authConfig.Configs[0].GetOauth2().GetAccessTokenValidation()

			Expect(actual.GetUserinfoUrl()).To(Equal(expected.GetUserinfoUrl()))
			Expect(actual.GetCacheTimeout()).To(Equal(expected.GetCacheTimeout()))
			Expect(actual.GetIntrospection().GetIntrospectionUrl()).To(Equal(expected.GetIntrospection().GetIntrospectionUrl()))
			Expect(actual.GetIntrospection().GetClientId()).To(Equal(expected.GetIntrospection().GetClientId()))
			Expect(actual.GetIntrospection().GetClientSecret()).To(Equal(clientSecret.ClientSecret))
			Expect(actual.GetIntrospection().GetUserIdAttributeName()).To(Equal(expected.GetIntrospection().GetUserIdAttributeName()))
			Expect(actual.GetRequiredScopes().GetScope()).To(Equal(expected.GetRequiredScopes().GetScope()))
		})
	})
})
