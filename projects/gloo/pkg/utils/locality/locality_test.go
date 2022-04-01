package locality_test

import (
	"context"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/locality"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	mock_k8s_core_clients "github.com/solo-io/external-apis/pkg/api/k8s/core/v1/mocks"
	"github.com/solo-io/go-utils/testutils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("LocalityFinder", func() {

	var (
		ctrl *gomock.Controller
		ctx  context.Context

		nodeClient *mock_k8s_core_clients.MockNodeClient
		podClient  *mock_k8s_core_clients.MockPodClient

		nodeRegion = "region-1"
		nodeZone   = "region-1a"

		stableLabels     = [2]string{corev1.LabelZoneRegionStable, corev1.LabelZoneFailureDomainStable}
		deprecatedLabels = [2]string{corev1.LabelZoneRegion, corev1.LabelZoneFailureDomain}
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())

		nodeClient = mock_k8s_core_clients.NewMockNodeClient(ctrl)
		podClient = mock_k8s_core_clients.NewMockPodClient(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Zones", func() {

		var checkLabels = func(keys [2]string, selector *metav1.LabelSelector, namespace string) {
			node := &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "node",
					Namespace: "namespace",
					Labels: map[string]string{
						keys[0]: nodeRegion,
						keys[1]: nodeZone,
					},
				},
			}

			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod",
					Namespace: "namespace",
				},
				Spec: corev1.PodSpec{
					NodeName: node.GetName(),
				},
			}

			ls, err := metav1.LabelSelectorAsSelector(selector)
			Expect(err).NotTo(HaveOccurred())

			podClient.EXPECT().
				ListPod(
					ctx,
					client.InNamespace(namespace),
					client.MatchingLabelsSelector{Selector: ls},
				).
				Return(&corev1.PodList{Items: []corev1.Pod{*pod}}, nil)

			nodeClient.EXPECT().
				GetNode(ctx, node.GetName()).
				Return(node, nil)
		}

		Context("Deployment", func() {

			var (
				deployment = &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "deployment",
						Namespace: "namespace",
					},
					Spec: appsv1.DeploymentSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels:      map[string]string{"hello": "world"},
							MatchExpressions: nil,
						},
					},
				}
			)

			It("will work with new stable labels", func() {
				checkLabels(
					stableLabels,
					deployment.Spec.Selector,
					deployment.GetNamespace(),
				)

				localityFinder := locality.NewLocalityFinder(nodeClient, podClient)
				zones, err := localityFinder.ZonesForDeployment(ctx, deployment)
				Expect(err).NotTo(HaveOccurred())
				Expect(zones).To(Equal([]string{nodeZone}))
			})

			It("will work with deprecated labels", func() {
				checkLabels(
					deprecatedLabels,
					deployment.Spec.Selector,
					deployment.GetNamespace(),
				)

				localityFinder := locality.NewLocalityFinder(nodeClient, podClient)
				zones, err := localityFinder.ZonesForDeployment(ctx, deployment)
				Expect(err).NotTo(HaveOccurred())
				Expect(zones).To(Equal([]string{nodeZone}))
			})

		})

		Context("DaemonSet", func() {

			var (
				daemonset = &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "daemonset",
						Namespace: "namespace",
					},
					Spec: appsv1.DaemonSetSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels:      map[string]string{"hello": "world"},
							MatchExpressions: nil,
						},
					},
				}
			)

			It("will work with new stable labels", func() {
				checkLabels(
					stableLabels,
					daemonset.Spec.Selector,
					daemonset.GetNamespace(),
				)

				localityFinder := locality.NewLocalityFinder(nodeClient, podClient)
				zones, err := localityFinder.ZonesForDaemonSet(ctx, daemonset)
				Expect(err).NotTo(HaveOccurred())
				Expect(zones).To(Equal([]string{nodeZone}))
			})

			It("will work with deprecated labels", func() {
				checkLabels(
					deprecatedLabels,
					daemonset.Spec.Selector,
					daemonset.GetNamespace(),
				)

				localityFinder := locality.NewLocalityFinder(nodeClient, podClient)
				zones, err := localityFinder.ZonesForDaemonSet(ctx, daemonset)
				Expect(err).NotTo(HaveOccurred())
				Expect(zones).To(Equal([]string{nodeZone}))
			})

		})

	})

	Context("Region", func() {

		It("will fail if nodes cannot be listed", func() {
			testErr := eris.New("hello")

			nodeClient.EXPECT().
				ListNode(ctx).
				Return(nil, testErr)

			localityFinder := locality.NewLocalityFinder(nodeClient, podClient)
			_, err := localityFinder.GetRegion(ctx)
			Expect(err).To(HaveOccurred())
			Expect(err).To(testutils.HaveInErrorChain(testErr))
		})

		It("will return the stable label", func() {
			region := "us-east-1"
			node := &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						stableLabels[0]: region,
					},
				},
			}

			nodeClient.EXPECT().
				ListNode(ctx).
				Return(&corev1.NodeList{Items: []corev1.Node{*node}}, nil)

			localityFinder := locality.NewLocalityFinder(nodeClient, podClient)
			expectedRegion, err := localityFinder.GetRegion(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(expectedRegion).To(Equal(region))
		})

		It("will return the deprecated label", func() {
			region := "us-east-1"
			node := &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						deprecatedLabels[0]: region,
					},
				},
			}

			nodeClient.EXPECT().
				ListNode(ctx).
				Return(&corev1.NodeList{Items: []corev1.Node{*node}}, nil)

			localityFinder := locality.NewLocalityFinder(nodeClient, podClient)
			expectedRegion, err := localityFinder.GetRegion(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(expectedRegion).To(Equal(region))
		})

	})

})
