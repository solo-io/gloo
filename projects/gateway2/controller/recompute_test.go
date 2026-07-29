package controller_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Proxy recomputation", func() {
	const (
		timeout  = 10 * time.Second
		interval = 100 * time.Millisecond
	)

	It("does not kick for Pod or EndpointSlice events", func() {
		suffix := fmt.Sprintf("%d", time.Now().UnixNano())

		// Verify the callback is wired through a controller that still requires
		// full proxy recomputation.
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "recompute-control-" + suffix,
				Namespace: "default",
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{Port: 8080}},
			},
		}
		beforeService := kickCount.Load()
		Expect(k8sClient.Create(ctx, service)).To(Succeed())
		Eventually(kickCount.Load, timeout, interval).Should(BeNumerically(">", beforeService))

		beforeEndpoints := kickCount.Load()
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "no-recompute-" + suffix,
				Namespace: "default",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "test",
					Image: "test",
				}},
			},
		}
		endpointSlice := &discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "no-recompute-" + suffix,
				Namespace: "default",
			},
			AddressType: discoveryv1.AddressTypeIPv4,
			Endpoints: []discoveryv1.Endpoint{{
				Addresses: []string{"10.0.0.1"},
			}},
		}

		Expect(k8sClient.Create(ctx, pod)).To(Succeed())
		Expect(k8sClient.Create(ctx, endpointSlice)).To(Succeed())
		Consistently(kickCount.Load, time.Second, interval).Should(Equal(beforeEndpoints))
	})
})
