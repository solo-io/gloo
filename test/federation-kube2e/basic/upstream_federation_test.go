package basic_test

import (
	"context"
	"time"

	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation/placement"

	"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/options/static"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloo_types "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1"
	gloo_fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1/types"
	fed_core_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/core/v1"
	mc_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/core/v1"
	multicluster_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1/types"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Upstream federation", func() {

	var (
		ctx              context.Context
		placementManager placement.Manager
		upstreamSpec     *gloo_types.UpstreamSpec
		meta             *fed_core_v1.TemplateMetadata
		fedUpstream      *v1.FederatedUpstream
	)

	BeforeEach(func() {
		ctx = context.TODO()

		placementManager = placement.NewManager(namespace, "gloo-fed-pod")
	})

	AfterEach(func() {
		// Just in case test fails earlier than expected
		if fedUpstream != nil {
			clientset, err := v1.NewClientsetFromConfig(managementClusterConfig.RestConfig)
			Expect(err).NotTo(HaveOccurred())
			_ = clientset.FederatedUpstreams().DeleteFederatedUpstream(ctx, client.ObjectKey{
				Namespace: fedUpstream.GetNamespace(),
				Name:      fedUpstream.GetName(),
			})
		}
	})

	It("throws validation error when missing Placement", func() {
		fedUpstream = &v1.FederatedUpstream{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      "fed-upstream",
			},
			Spec: gloo_fed_types.FederatedUpstreamSpec{
				Template: &gloo_fed_types.FederatedUpstreamSpec_Template{
					Spec:     upstreamSpec,
					Metadata: meta,
				},
				Placement: nil, // placement set to `nil` to cause error
			},
		}

		// register fedUpstream with kind-mgmt
		clientset, err := v1.NewClientsetFromConfig(managementClusterConfig.RestConfig)
		Expect(err).NotTo(HaveOccurred())
		err = clientset.FederatedUpstreams().CreateFederatedUpstream(ctx, fedUpstream)
		Expect(err).NotTo(HaveOccurred())

		// wait for INVALID placement status, per `federation_reconcilers.go`
		Eventually(func(g Gomega) {
			resultingFedUpstream, err := clientset.FederatedUpstreams().GetFederatedUpstream(
				ctx,
				types.NamespacedName{
					Name:      fedUpstream.GetObjectMeta().GetName(),
					Namespace: fedUpstream.GetObjectMeta().GetNamespace(),
				},
			)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(placementManager.GetPlacementStatus(&resultingFedUpstream.Status).GetState()).To(Equal(mc_types.PlacementStatus_INVALID))
		}, 10*time.Second).Should(Succeed())
	})

	It("works", func() {
		clientset, err := v1.NewClientsetFromConfig(managementClusterConfig.RestConfig)
		Expect(err).NotTo(HaveOccurred())

		upstreamSpec = &gloo_types.UpstreamSpec{
			UpstreamType: &gloo_types.UpstreamSpec_Static{
				Static: &static.UpstreamSpec{
					Hosts: []*static.Host{{
						Addr: "solo.io",
						Port: 80,
					}},
				},
			},
		}

		meta = &fed_core_v1.TemplateMetadata{
			Annotations: map[string]string{"anno": "tation"},
			Labels:      map[string]string{"label": "printer"},
			Name:        "charles",
		}

		fedUpstream = &v1.FederatedUpstream{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      "fed-upstream",
			},
			Spec: gloo_fed_types.FederatedUpstreamSpec{
				Template: &gloo_fed_types.FederatedUpstreamSpec_Template{
					Spec:     upstreamSpec,
					Metadata: meta,
				},
				Placement: &multicluster_types.Placement{
					Namespaces: []string{namespace},
					Clusters:   []string{remoteClusterConfig.KubeContext},
				},
			},
		}

		err = clientset.FederatedUpstreams().CreateFederatedUpstream(ctx, fedUpstream)
		Expect(err).NotTo(HaveOccurred())

		remoteClientSet, err := gloo_v1.NewClientsetFromConfig(remoteClusterConfig.RestConfig)
		Expect(err).NotTo(HaveOccurred())
		var resultingUpstream *gloo_v1.Upstream
		Eventually(func() *gloo_v1.Upstream {
			resultingUpstream, _ = remoteClientSet.Upstreams().
				GetUpstream(ctx, types.NamespacedName{Name: meta.Name, Namespace: namespace})
			return resultingUpstream
		}, 10*time.Second).ShouldNot(BeNil())

		Expect(resultingUpstream.Spec).To(Equal(*upstreamSpec))
		Expect(resultingUpstream.Annotations).To(Equal(meta.Annotations))
		Expect(resultingUpstream.Labels).To(Equal(map[string]string{
			"label":             "printer",
			federation.HubOwner: namespace + ".fed-upstream",
		}))

		Eventually(func() bool {
			resultingUpstream, _ = remoteClientSet.Upstreams().
				GetUpstream(ctx, types.NamespacedName{Name: meta.Name, Namespace: namespace})
			return resultingUpstream.Status.State == gloo_types.UpstreamStatus_Accepted
		}, 10*time.Second).Should(BeTrue(), "remote upstream should be marked accepted by gloo")

		err = clientset.FederatedUpstreams().
			DeleteFederatedUpstream(ctx, client.ObjectKey{
				Name:      fedUpstream.GetName(),
				Namespace: fedUpstream.GetNamespace(),
			})
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() bool {
			_, err = remoteClientSet.Upstreams().
				GetUpstream(ctx, types.NamespacedName{Name: meta.Name, Namespace: namespace})
			if err != nil {
				return errors.IsNotFound(err)
			}
			return false
		}, 10*time.Second).Should(BeTrue(), "remote upstream should be cleaned up when federated upstream is deleted")
	})
})
