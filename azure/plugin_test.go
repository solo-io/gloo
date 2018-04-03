package azure

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/plugin"
)

var _ = Describe("Plugin dependencies", func() {
	Context("with a list of upstreams", func() {
		It("should return a list of secret refs for dependencies", func() {
			cfg := &v1.Config{
				Upstreams: []*v1.Upstream{
					&v1.Upstream{},
					upstream("azure1", "my-appwhos", "azure-secret1"),
					upstream("azure2", "my-azure", "azure-secret2"),
					upstream("azure3", "azurization", "azure-secret1"),
					upstream("azure4", "azoor", ""),
				},
			}
			p := &Plugin{}
			dependencies := p.GetDependencies(cfg)
			Expect(dependencies.SecretRefs).To(HaveLen(4))
			Expect(dependencies.SecretRefs[0]).To(Equal("azure-secret1"))
			Expect(dependencies.SecretRefs[1]).To(Equal("azure-secret2"))
			Expect(dependencies.SecretRefs[2]).To(Equal("azure-secret1"))
			Expect(dependencies.SecretRefs[3]).To(Equal(""))
		})
	})
})

var _ = Describe("Plugin HTTP filters", func() {
	Context("when needed", func() {
		It("should create a filter and reset plugin", func() {
			p := &Plugin{
				isNeeded: true,
				hostname: "previous.host.name",
				apiKeys:  map[string]string{"a": "b"},
			}
			filters := p.HttpFilters(nil)
			Expect(filters).To(HaveLen(1))
			Expect(filters[0].HttpFilter.Name).To(Equal("io.solo.azure_functions"))
			Expect(filters[0].Stage).To(Equal(plugin.OutAuth))
			Expect(p.isNeeded).To(Equal(false))
			Expect(p.hostname).To(Equal(""))
			Expect(p.apiKeys).To(HaveLen(0))
		})
	})
	Context("when not needed", func() {
		It("should reset plugin without creating a filter", func() {
			p := &Plugin{
				isNeeded: false,
				hostname: "previous.host.name",
				apiKeys:  map[string]string{"a": "b"},
			}
			filters := p.HttpFilters(nil)
			Expect(filters).To(HaveLen(0))
			Expect(p.isNeeded).To(Equal(false))
			Expect(p.hostname).To(Equal(""))
			Expect(p.apiKeys).To(HaveLen(0))
		})
	})
})

var _ = Describe("Processing upstream", func() {
	Context("with non-Azure upstream", func() {
		It("should not error and return nothing", func() {
			upstreams := []*v1.Upstream{&v1.Upstream{}, &v1.Upstream{Type: "some-upstream"}}
			p := Plugin{}
			for _, u := range upstreams {
				err := p.ProcessUpstream(nil, u, nil)
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})
	Context("when secret referenced by Azure upstream is missing", func() {
		It("should error", func() {
			upstream := &v1.Upstream{
				Type: UpstreamTypeAzure,
				Spec: upstreamSpec("my-appwhos", "azure-secret1"),
			}
			out := &envoyapi.Cluster{}
			params := &plugin.UpstreamPluginParams{}
			p := Plugin{}
			err := p.ProcessUpstream(params, upstream, out)
			Expect(err).To(HaveOccurred())
		})
	})
	Context("with valid upstream spec", func() {
		var (
			err error
			p   Plugin
			out *envoyapi.Cluster
		)
		BeforeEach(func() {
			p = Plugin{apiKeys: make(map[string]string)}
			upstream := &v1.Upstream{
				Type: UpstreamTypeAzure,
				Spec: upstreamSpec("my-appwhos", "azure-secret1"),
			}
			out = &envoyapi.Cluster{}
			params := &plugin.UpstreamPluginParams{Secrets: map[string]map[string]string{
				"azure-secret1": map[string]string{"_master": "key1", "foo": "key1", "bar": "key2"},
			}}
			err = p.ProcessUpstream(params, upstream, out)
		})
		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		It("should have the correct output", func() {
			Expect(out.Hosts).Should(HaveLen(1))
			Expect(out.Hosts[0].GetSocketAddress().Address).To(Equal("my-appwhos.azurewebsites.net"))
			Expect(out.Hosts[0].GetSocketAddress().PortSpecifier.(*envoycore.SocketAddress_PortValue).PortValue).To(BeEquivalentTo(443))
			Expect(out.TlsContext.Sni).To(Equal("my-appwhos.azurewebsites.net"))
			Expect(out.Type).To(Equal(envoyapi.Cluster_LOGICAL_DNS))
			Expect(out.DnsLookupFamily).To(Equal(envoyapi.Cluster_V4_ONLY))
		})
		It("should have empty Azure metadata in output", func() {
			metadata, ok := out.Metadata.FilterMetadata["io.solo.azure_functions"]
			Expect(ok).To(Equal(true))
			Expect(metadata.Fields).Should(HaveLen(0))
		})
		It("should add the hostname and api key map to the plugin", func() {
			Expect(p.isNeeded).To(Equal(true))
			Expect(p.hostname).To(Equal("my-appwhos.azurewebsites.net"))
			Expect(p.apiKeys).To(Equal(map[string]string{"_master": "key1", "foo": "key1", "bar": "key2"}))
		})
	})
	Context("with valid upstream spec without secrets", func() {
		var (
			err error
			p   Plugin
			out *envoyapi.Cluster
		)
		BeforeEach(func() {
			p = Plugin{apiKeys: make(map[string]string)}
			upstream := &v1.Upstream{
				Type: UpstreamTypeAzure,
				Spec: upstreamSpec("my-appwhos", ""),
			}
			out = &envoyapi.Cluster{}
			params := &plugin.UpstreamPluginParams{Secrets: map[string]map[string]string{
				"some-irrelevant-secret1": map[string]string{"_master": "key1", "foo": "key1", "bar": "key2"},
			}}
			err = p.ProcessUpstream(params, upstream, out)
		})
		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		It("should have the correct output", func() {
			Expect(out.Hosts).Should(HaveLen(1))
			Expect(out.Hosts[0].GetSocketAddress().Address).To(Equal("my-appwhos.azurewebsites.net"))
			Expect(out.Hosts[0].GetSocketAddress().PortSpecifier.(*envoycore.SocketAddress_PortValue).PortValue).To(BeEquivalentTo(443))
			Expect(out.TlsContext.Sni).To(Equal("my-appwhos.azurewebsites.net"))
			Expect(out.Type).To(Equal(envoyapi.Cluster_LOGICAL_DNS))
			Expect(out.DnsLookupFamily).To(Equal(envoyapi.Cluster_V4_ONLY))
		})
		It("should have empty Azure metadata in output", func() {
			metadata, ok := out.Metadata.FilterMetadata["io.solo.azure_functions"]
			Expect(ok).To(Equal(true))
			Expect(metadata.Fields).Should(HaveLen(0))
		})
		It("should add the hostname to the plugin, but the api key map should remain empty", func() {
			Expect(p.isNeeded).To(Equal(true))
			Expect(p.hostname).To(Equal("my-appwhos.azurewebsites.net"))
			Expect(p.apiKeys).To(Equal(make(map[string]string)))
		})
	})
})

var _ = Describe("API key", func() {
	Describe("retrieval", func() {
		Context("of a function key from an empty map", func() {
			It("should fail", func() {
				apiKeys := make(map[string]string)
				_, err := getApiKey(apiKeys, []string{"foo"})
				Expect(err).To(HaveOccurred())
			})
		})
		Context("of a master key from an empty map", func() {
			It("should fail", func() {
				apiKeys := make(map[string]string)
				_, err := getApiKey(apiKeys, []string{MasterKeyName})
				Expect(err).To(HaveOccurred())
			})
		})
		Context("of a function key and a master key from an empty map", func() {
			It("should fail", func() {
				apiKeys := make(map[string]string)
				_, err := getApiKey(apiKeys, []string{"foo", MasterKeyName})
				Expect(err).To(HaveOccurred())
			})
		})
		Context("of an existing empty function key", func() {
			It("should fail", func() {
				apiKeys := map[string]string{"foo": ""}
				_, err := getApiKey(apiKeys, []string{"foo"})
				Expect(err).To(HaveOccurred())
			})
		})
		Context("of an existing empty master key", func() {
			It("should fail", func() {
				apiKeys := map[string]string{MasterKeyName: ""}
				_, err := getApiKey(apiKeys, []string{MasterKeyName})
				Expect(err).To(HaveOccurred())
			})
		})
		Context("of an existing non-empty function key", func() {
			It("should succeed", func() {
				apiKeys := map[string]string{"foo": "functionkey1=="}
				apiKey, err := getApiKey(apiKeys, []string{"foo"})
				Expect(err).NotTo(HaveOccurred())
				Expect(apiKey).To(Equal("functionkey1=="))
			})
		})
		Context("of an existing non-empty master key", func() {
			It("should succeed", func() {
				apiKeys := map[string]string{MasterKeyName: "新신シん양ʃШŞשܫשׁش"}
				apiKey, err := getApiKey(apiKeys, []string{MasterKeyName})
				Expect(err).NotTo(HaveOccurred())
				Expect(apiKey).To(Equal("新신シん양ʃШŞשܫשׁش"))
			})
		})
		Context("of a function key and a master key with both present", func() {
			It("should return the function key", func() {
				apiKeys := map[string]string{MasterKeyName: "key1", "foo": "key2"}
				apiKey, err := getApiKey(apiKeys, []string{"foo", MasterKeyName})
				Expect(err).NotTo(HaveOccurred())
				Expect(apiKey).To(Equal("key2"))
			})
		})
		Context("of a missing function key and an existing master key", func() {
			It("should return the master key", func() {
				apiKeys := map[string]string{MasterKeyName: "key1", "foo": "key2"}
				apiKey, err := getApiKey(apiKeys, []string{"bar", MasterKeyName})
				Expect(err).NotTo(HaveOccurred())
				Expect(apiKey).To(Equal("key1"))
			})
		})
	})
})

var _ = Describe("Path parameters", func() {
	Describe("retrieval", func() {
		Context("with anonymous authorization and an empty map", func() {
			It("should be empty", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "anonymous")
				apiKeys := make(map[string]string)
				pathParameters, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).NotTo(HaveOccurred())
				Expect(pathParameters).To(Equal(""))
			})
		})
		Context("with function-level authorization and an empty map", func() {
			It("should fail", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "function")
				apiKeys := make(map[string]string)
				_, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with admin-level authorization and an empty map", func() {
			It("should fail", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "admin")
				apiKeys := make(map[string]string)
				_, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with an invalid authorization level and an empty map", func() {
			It("should fail", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "invalid")
				apiKeys := make(map[string]string)
				_, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with anonymous authorization and a missing function", func() {
			It("should be empty", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "anonymous")
				apiKeys := map[string]string{"bar": "key1"}
				pathParameters, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).NotTo(HaveOccurred())
				Expect(pathParameters).To(Equal(""))
			})
		})
		Context("with function-level authorization and a missing function", func() {
			It("should fail", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "function")
				apiKeys := map[string]string{"bar": "key1"}
				_, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with admin-level authorization and a missing function", func() {
			It("should fail", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "admin")
				apiKeys := map[string]string{"bar": "key1"}
				_, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with an invalid authorization level and a missing function", func() {
			It("should fail", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "invalid")
				apiKeys := map[string]string{"bar": "key1"}
				_, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with anonymous authorization and an existing function", func() {
			It("should be empty", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "anonymous")
				apiKeys := map[string]string{"foo": "key1"}
				pathParameters, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).NotTo(HaveOccurred())
				Expect(pathParameters).To(Equal(""))
			})
		})
		Context("with function-level authorization and an existing function", func() {
			It("should include key", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "function")
				apiKeys := map[string]string{"foo": "key1"}
				pathParameters, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).NotTo(HaveOccurred())
				Expect(pathParameters).To(Equal("?code=key1"))
			})
		})
		Context("with admin-level authorization and an existing function", func() {
			It("should fail", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "admin")
				apiKeys := map[string]string{"foo": "key1"}
				_, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with an invalid authorization level and an existing function", func() {
			It("should fail", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "invalid")
				apiKeys := map[string]string{"foo": "key1"}
				_, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with anonymous authorization and a master key", func() {
			It("should be empty", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "anonymous")
				apiKeys := map[string]string{"_master": "key1"}
				pathParameters, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).NotTo(HaveOccurred())
				Expect(pathParameters).To(Equal(""))
			})
		})
		Context("with function-level authorization and a master key", func() {
			It("should include the master key", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "function")
				apiKeys := map[string]string{"_master": "key1"}
				pathParameters, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).NotTo(HaveOccurred())
				Expect(pathParameters).To(Equal("?code=key1"))
			})
		})
		Context("with admin-level authorization and a master key", func() {
			It("should include the master key", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "admin")
				apiKeys := map[string]string{"_master": "key1"}
				pathParameters, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).NotTo(HaveOccurred())
				Expect(pathParameters).To(Equal("?code=key1"))
			})
		})
		Context("with an invalid authorization level and a master key", func() {
			It("should fail", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "invalid")
				apiKeys := map[string]string{"_master": "key1"}
				_, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with anonymous authorization and both keys", func() {
			It("should be empty", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "anonymous")
				apiKeys := map[string]string{"_master": "key1", "foo": "key2"}
				pathParameters, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).NotTo(HaveOccurred())
				Expect(pathParameters).To(Equal(""))
			})
		})
		Context("with function-level authorization and both keys", func() {
			It("should include the function key", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "function")
				apiKeys := map[string]string{"_master": "key1", "foo": "key2"}
				pathParameters, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).NotTo(HaveOccurred())
				Expect(pathParameters).To(Equal("?code=key2"))
			})
		})
		Context("with admin-level authorization and both keys", func() {
			It("should include the master key", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "admin")
				apiKeys := map[string]string{"_master": "key1", "foo": "key2"}
				pathParameters, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).NotTo(HaveOccurred())
				Expect(pathParameters).To(Equal("?code=key1"))
			})
		})
		Context("with an invalid authorization level and both keys", func() {
			It("should fail", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "invalid")
				apiKeys := map[string]string{"_master": "key1"}
				_, err := getPathParameters(functionSpec, apiKeys)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

var _ = Describe("Path ", func() {
	Describe("retrieval", func() {
		Context("with anonymous authorization and both keys", func() {
			It("should have no parameters", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "anonymous")
				apiKeys := map[string]string{"_master": "key1", "foo": "key2"}
				pathParameters, err := getPath(functionSpec, apiKeys)
				Expect(err).NotTo(HaveOccurred())
				Expect(pathParameters).To(Equal("/api/foo"))
			})
		})
		Context("with function-level authorization and both keys", func() {
			It("should include the function key", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "function")
				apiKeys := map[string]string{"_master": "key1", "foo": "key2"}
				pathParameters, err := getPath(functionSpec, apiKeys)
				Expect(err).NotTo(HaveOccurred())
				Expect(pathParameters).To(Equal("/api/foo?code=key2"))
			})
		})
		Context("with admin-level authorization and both keys", func() {
			It("should include the master key", func() {
				functionSpec, _ := decodeFunctionSpec("foo", "admin")
				apiKeys := map[string]string{"_master": "key1", "foo": "key2"}
				pathParameters, err := getPath(functionSpec, apiKeys)
				Expect(err).NotTo(HaveOccurred())
				Expect(pathParameters).To(Equal("/api/foo?code=key1"))
			})
		})
	})
})

func upstreamSpec(functionAppName string, secretRef string) v1.UpstreamSpec {
	return &types.Struct{
		Fields: map[string]*types.Value{
			"function_app_name": {Kind: &types.Value_StringValue{StringValue: functionAppName}},
			"secret_ref":        {Kind: &types.Value_StringValue{StringValue: secretRef}},
		},
	}
}

func upstream(name string, functionAppName string, secretRef string) *v1.Upstream {
	return &v1.Upstream{
		Name: name,
		Type: UpstreamTypeAzure,
		Spec: upstreamSpec(functionAppName, secretRef),
	}
}
