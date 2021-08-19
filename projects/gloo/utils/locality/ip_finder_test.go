package locality_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	mock_k8s_core_clients "github.com/solo-io/external-apis/pkg/api/k8s/core/v1/mocks"
	"github.com/solo-io/solo-projects/projects/gloo/utils/locality"
	core_v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("IpFinder", func() {
	var (
		ctrl *gomock.Controller
		ctx  context.Context

		podClient  *mock_k8s_core_clients.MockPodClient
		nodeClient *mock_k8s_core_clients.MockNodeClient

		cluster = "new.cluster.who.dis"
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())

		podClient = mock_k8s_core_clients.NewMockPodClient(ctrl)
		nodeClient = mock_k8s_core_clients.NewMockNodeClient(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("can return endpoints for a LoadBalancer Service", func() {
		externalIpGetter := locality.NewExternalIpFinder(cluster, podClient, nodeClient)

		svc := &core_v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "service-name",
			},
			Spec: core_v1.ServiceSpec{
				Ports: []core_v1.ServicePort{
					{
						Name: "new.port.who.dis",
						Port: 9000,
					},
				},
				Type: core_v1.ServiceTypeLoadBalancer,
			},
			Status: core_v1.ServiceStatus{
				LoadBalancer: core_v1.LoadBalancerStatus{
					Ingress: []core_v1.LoadBalancerIngress{
						{
							IP: "new.ip.who.dis",
						},
					},
				},
			},
		}

		expected := &locality.IngressEndpoint{
			Address: svc.Status.LoadBalancer.Ingress[0].IP,
			Ports: []*locality.Port{
				{
					Port: uint32(svc.Spec.Ports[0].Port),
					Name: svc.Spec.Ports[0].Name,
				},
			},
			ServiceName: svc.Name,
		}

		eps, err := externalIpGetter.GetExternalIps(ctx, []*core_v1.Service{svc})
		Expect(err).NotTo(HaveOccurred())
		Expect(eps).To(HaveLen(1))
		Expect(eps[0]).To(Equal(expected))
	})

	It("can return endpoints for a NodePort Service", func() {
		externalIpGetter := locality.NewExternalIpFinder(cluster, podClient, nodeClient)

		svc := &core_v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-name",
				Namespace: "service-namespace",
			},
			Spec: core_v1.ServiceSpec{
				Selector: map[string]string{
					"select": "this",
				},
				Ports: []core_v1.ServicePort{
					{
						Name:     "new.port.who.dis",
						Port:     9000,
						NodePort: 32000,
					},
				},
				Type: core_v1.ServiceTypeNodePort,
			},
		}

		node := &core_v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-name",
			},
			Status: core_v1.NodeStatus{
				Addresses: []core_v1.NodeAddress{
					{
						Address: "node.address.who.dis",
					},
				},
			},
		}

		pod := &core_v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pod-name",
			},
			Spec: core_v1.PodSpec{
				NodeName: node.Name,
			},
		}

		expected := &locality.IngressEndpoint{
			Address: node.Status.Addresses[0].Address,
			Ports: []*locality.Port{
				{
					Port: uint32(svc.Spec.Ports[0].NodePort),
					Name: svc.Spec.Ports[0].Name,
				},
			},
			ServiceName: svc.Name,
		}

		podClient.EXPECT().
			ListPod(ctx, &client.ListOptions{
				LabelSelector: labels.SelectorFromSet(svc.Spec.Selector),
				Namespace:     svc.Namespace,
			}).Return(&core_v1.PodList{Items: []core_v1.Pod{*pod}}, nil)

		nodeClient.EXPECT().
			GetNode(ctx, pod.Spec.NodeName).
			Return(node, nil)

		eps, err := externalIpGetter.GetExternalIps(ctx, []*core_v1.Service{svc})
		Expect(err).NotTo(HaveOccurred())
		Expect(eps).To(HaveLen(1))
		Expect(eps[0]).To(Equal(expected))
	})

})
