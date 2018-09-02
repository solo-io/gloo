package azure_test

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/azure"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/azure"
)

var _ = Describe("Plugin", func() {
	var (
		p      plugins.Plugin
		out    *envoyapi.Cluster
		params plugins.Params
	)

	BeforeEach(func() {
		p = NewAzurePlugin()
		out = &envoyapi.Cluster{}
		params = plugins.Params{}
	})

	Context("with valid upstream spec", func() {
		var (
			err      error
			upstream *v1.Upstream
		)

		BeforeEach(func() {
			upstream = &v1.Upstream{
				Metadata: core.Metadata{
					Name: "test",
					// TODO(yuval-k): namespace
					Namespace: "",
				},
				UpstreamSpec: &v1.UpstreamSpec{
					UpstreamType: &v1.UpstreamSpec_Azure{
						Azure: &azure.UpstreamSpec{
							FunctionAppName: "my-appwhos",
						},
					},
				},
			}
		})
		Context("with secrets", func() {

			BeforeEach(func() {
				upstream.UpstreamSpec.UpstreamType.(*v1.UpstreamSpec_Azure).Azure.SecretRef = "azure-secret1"
				params.Snapshot = &v1.ApiSnapshot{
					SecretList: v1.SecretList{{
						Metadata: core.Metadata{
							Name: "azure-secret1",
							// TODO(yuval-k): namespace
							Namespace: "",
						},
						Data: map[string]string{"_master": "key1", "foo": "key1", "bar": "key2"},
					}},
				}

				err = p.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should have the correct output", func() {
				Expect(out.LoadAssignment.Endpoints).Should(HaveLen(1))
				// Expect(out.Hosts[0].GetSocketAddress().Address).To(Equal("my-appwhos.azurewebsites.net"))
				// Expect(out.Hosts[0].GetSocketAddress().PortSpecifier.(*envoycore.SocketAddress_PortValue).PortValue).To(BeEquivalentTo(443))
				Expect(out.TlsContext.Sni).To(Equal("my-appwhos.azurewebsites.net"))
				Expect(out.Type).To(Equal(envoyapi.Cluster_LOGICAL_DNS))
				Expect(out.DnsLookupFamily).To(Equal(envoyapi.Cluster_V4_ONLY))
			})

		})
		Context("without secrets", func() {
			BeforeEach(func() {
				err = p.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have the correct output", func() {
				Expect(out.LoadAssignment.Endpoints).Should(HaveLen(1))
				// Expect(out.Hosts[0].GetSocketAddress().Address).To(Equal("my-appwhos.azurewebsites.net"))
				// Expect(out.Hosts[0].GetSocketAddress().PortSpecifier.(*envoycore.SocketAddress_PortValue).PortValue).To(BeEquivalentTo(443))
				Expect(out.TlsContext.Sni).To(Equal("my-appwhos.azurewebsites.net"))
				Expect(out.Type).To(Equal(envoyapi.Cluster_LOGICAL_DNS))
				Expect(out.DnsLookupFamily).To(Equal(envoyapi.Cluster_V4_ONLY))
			})
		})
	})
})
