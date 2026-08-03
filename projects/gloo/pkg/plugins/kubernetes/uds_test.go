package kubernetes_test

import (
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options"
	gloov1kube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"
)

var _ = Describe("Uds", func() {

	It("should preseve ssl config when updating upstreams", func() {
		desired := &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &gloov1kube.UpstreamSpec{
					ServiceName: "test",
				},
			},
		}
		original := &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &gloov1kube.UpstreamSpec{},
			},
			SslConfig: &ssl.UpstreamSslConfig{Sni: "testsni"},
		}
		updated, err := UpdateUpstream(original, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeTrue())
		Expect(desired.SslConfig).To(BeIdenticalTo(original.SslConfig))
	})

	It("should update ssl config when one is desired", func() {
		desiredSslConfig := &ssl.UpstreamSslConfig{
			SslSecrets: &ssl.UpstreamSslConfig_SecretRef{
				SecretRef: &core.ResourceRef{Name: "hi", Namespace: "there"},
			},
		}
		desired := &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &gloov1kube.UpstreamSpec{
					ServiceName: "test",
				},
			},
			SslConfig: desiredSslConfig,
		}
		original := &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &gloov1kube.UpstreamSpec{},
			},
			SslConfig: &ssl.UpstreamSslConfig{Sni: "testsni"},
		}
		updated, err := UpdateUpstream(original, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeTrue())
		Expect(desired.SslConfig).To(BeIdenticalTo(desiredSslConfig))
	})

	It("should update the selector when the service selector changes", func() {
		desired := &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &gloov1kube.UpstreamSpec{
					ServiceName: "test",
					Selector:    map[string]string{"app": "echo", "version": "green"},
				},
			},
		}
		original := &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &gloov1kube.UpstreamSpec{
					ServiceName: "test",
					Selector:    map[string]string{"app": "echo", "version": "blue"},
				},
			},
		}
		updated, err := UpdateUpstream(original, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeTrue())
		Expect(desired.GetKube().GetSelector()).To(Equal(map[string]string{"app": "echo", "version": "green"}))
	})

	It("should not report a change when the selector is unchanged", func() {
		selector := map[string]string{"app": "echo", "version": "blue"}
		desired := &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &gloov1kube.UpstreamSpec{
					ServiceName: "test",
					Selector:    selector,
				},
			},
		}
		original := &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &gloov1kube.UpstreamSpec{
					ServiceName: "test",
					Selector:    selector,
				},
			},
		}
		updated, err := UpdateUpstream(original, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeFalse())
		Expect(desired.GetKube().GetSelector()).To(Equal(selector))
	})

	It("should clear the selector when the service no longer has one", func() {
		desired := &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &gloov1kube.UpstreamSpec{
					ServiceName: "test",
				},
			},
		}
		original := &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &gloov1kube.UpstreamSpec{
					ServiceName: "test",
					Selector:    map[string]string{"app": "echo", "version": "blue"},
				},
			},
		}
		updated, err := UpdateUpstream(original, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeTrue())
		Expect(desired.GetKube().GetSelector()).To(BeEmpty())
	})

	It("should preserve the service spec when updating upstreams", func() {
		desired := &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &gloov1kube.UpstreamSpec{
					ServiceName: "test",
				},
			},
		}
		original := &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &gloov1kube.UpstreamSpec{
					ServiceName: "test",
					ServiceSpec: &options.ServiceSpec{
						PluginType: &options.ServiceSpec_GrpcJsonTranscoder{},
					},
				},
			},
		}
		updated, err := UpdateUpstream(original, desired)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).To(BeFalse())
		Expect(desired.GetKube().GetServiceSpec()).NotTo(BeNil())
	})

})
