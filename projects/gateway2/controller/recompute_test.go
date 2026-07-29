package controller_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

var _ = Describe("Proxy recomputation", Ordered, func() {
	const timeout = 10 * time.Second

	It("does not kick proxy translation for endpoint-only events", func() {
		suffix := fmt.Sprintf("%d", time.Now().UnixNano())
		serviceName := "kick-control-" + suffix

		beforeService := kickCount.Load()
		Expect(k8sClient.Create(ctx, &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: "default",
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{
					Name: "http",
					Port: 8080,
				}},
			},
		})).To(Succeed())
		Eventually(kickCount.Load, timeout).Should(BeNumerically(">", beforeService),
			"the positive control should prove that Kick is wired to the controller")

		baseline := kickCount.Load()
		Expect(k8sClient.Create(ctx, &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "no-kick-pod-" + suffix,
				Namespace: "default",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  "test",
					Image: "test",
				}},
			},
		})).To(Succeed())
		Expect(k8sClient.Create(ctx, &discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "no-kick-slice-" + suffix,
				Namespace: "default",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: serviceName,
				},
			},
			AddressType: discoveryv1.AddressTypeIPv4,
			Endpoints: []discoveryv1.Endpoint{{
				Addresses: []string{"10.0.0.1"},
			}},
			Ports: []discoveryv1.EndpointPort{{
				Name: ptr.To("http"),
				Port: ptr.To(int32(8080)),
			}},
		})).To(Succeed())

		Consistently(kickCount.Load, time.Second).Should(Equal(baseline))
	})
})
