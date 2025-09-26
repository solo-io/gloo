package kubernetes

import (
	"context"
	"strings"
	"time"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes/serviceconverter"
	corecache "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
)

var _ = Describe("Plugin", func() {
	var (
		ctx      context.Context
		kube     *fake.Clientset
		params   plugins.Params
		plugin   plugins.Plugin
		upstream *v1.Upstream
		out      *envoy_config_cluster_v3.Cluster
	)
	BeforeEach(func() {
		kube = fake.NewSimpleClientset()
		kubeCoreCache, err := corecache.NewKubeCoreCache(context.Background(), kube)
		Expect(err).NotTo(HaveOccurred())
		plugin = NewPlugin(kube, kubeCoreCache, nil)
		plugin.Init(plugins.InitParams{})
		upstream = &v1.Upstream{
			Metadata: &core.Metadata{
				Name:      "myUpstream",
				Namespace: "ns",
			},
			UpstreamType: &v1.Upstream_Kube{
				Kube: &kubernetes.UpstreamSpec{
					ServiceName:      "mySvc",
					ServiceNamespace: "ns",
				},
			},
		}
		out = &envoy_config_cluster_v3.Cluster{}
	})

	Context("upstreams", func() {

		It("should error upstream with nonexistent service", func() {

			err := plugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "does not exist in namespace")).To(BeTrue())
		})

		It("should error upstream with nonexistent serviceNamespace", func() {
			// Reset plugin to not watch all namespaces.
			kubeCoreCache, err := corecache.NewKubeCoreCacheWithOptions(context.Background(), kube, time.Duration(1), []string{"ns"})
			Expect(err).NotTo(HaveOccurred())
			plugin = NewPlugin(kube, kubeCoreCache, nil)
			plugin.Init(plugins.InitParams{})
			upstream.UpstreamType = &v1.Upstream_Kube{
				Kube: &kubernetes.UpstreamSpec{
					ServiceName:      "mySvc",
					ServiceNamespace: "ns-fake",
				},
			}
			err = plugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "invalid ServiceNamespace")).To(BeTrue())
		})

	})

	Context("determine ip family", func() {

		var (
			svcName           = "mySvc"
			svcAnnotationName = "mySvcWithAnnotation"
			svcNs             = "ns"

			svc = &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      svcName,
					Namespace: svcNs,
				},
				Spec: corev1.ServiceSpec{},
			}

			svcAnnotations = &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      svcAnnotationName,
					Namespace: svcNs,
					Annotations: map[string]string{
						serviceconverter.GlooDnsIpFamilyAnnotation: "v6",
					},
				},
				Spec: corev1.ServiceSpec{},
			}
		)

		BeforeEach(func() {
			ctx = context.TODO()
			kube = fake.NewClientset()
			_, err := kube.CoreV1().Services(svcNs).Create(ctx, svc, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
			_, err = kube.CoreV1().Services(svcNs).Create(ctx, svcAnnotations, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			kubeCoreCache, err := corecache.NewKubeCoreCacheWithOptions(ctx, kube, time.Duration(1), []string{svcNs})
			Expect(err).NotTo(HaveOccurred())

			out = new(envoy_config_cluster_v3.Cluster)

			plugin = NewPlugin(kube, kubeCoreCache, nil)
			upstream = &v1.Upstream{
				Metadata: &core.Metadata{
					Name:      "myUpstream",
					Namespace: svcNs,
				},
				UpstreamType: &v1.Upstream_Kube{
					Kube: &kubernetes.UpstreamSpec{
						ServiceName:      svcName,
						ServiceNamespace: svcNs,
					},
				},
			}
		})

		JustBeforeEach(func() {
			plugin.Init(plugins.InitParams{})
		})

		It("should process upstream with default ip family", func() {
			err := plugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			Expect(err).To(Not(HaveOccurred()))
			Expect(out.GetDnsLookupFamily()).To(Equal(envoy_config_cluster_v3.Cluster_AUTO))
		})

		It("should process upstream with custom ip family", func() {
			uc = DefaultUpstreamConverter()
			port := corev1.ServicePort{
				Port: 123,
			}
			upstream := uc.CreateUpstream(ctx, svcAnnotations, port)
			err := plugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			Expect(err).To(Not(HaveOccurred()))
			Expect(out.GetDnsLookupFamily()).To(Equal(envoy_config_cluster_v3.Cluster_V6_ONLY))
		})
	})

})
