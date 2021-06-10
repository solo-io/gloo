package gloo_fed_e2e_test

import (
	"context"
	"time"

	"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/options/static"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/skv2-enterprise/multicluster-admission-webhook/pkg/api/multicluster.solo.io/v1alpha1"
	"github.com/solo-io/skv2/test"
	gloo_types "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1"
	gloo_fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1/types"
	fed_core_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/core/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Upstream federation", func() {

	var (
		ctx          context.Context
		upstreamSpec *gloo_types.UpstreamSpec
		meta         *fed_core_v1.TemplateMetadata

		localGlooHubNamespace = "gloo-system"
		remoteGlooNamespace   = "gloo-system"

		fedUpstream *v1.FederatedUpstream
	)

	BeforeEach(func() {
		ctx = context.TODO()
	})

	AfterEach(func() {
		// Just in case test fails earlier than expected
		if fedUpstream != nil {
			clientset, err := v1.NewClientsetFromConfig(test.MustConfig(""))
			Expect(err).NotTo(HaveOccurred())
			clientset.FederatedUpstreams().DeleteFederatedUpstream(ctx, client.ObjectKey{
				Namespace: fedUpstream.GetNamespace(),
				Name:      fedUpstream.GetName(),
			})
		}
	})

	It("works", func() {
		clientset, err := v1.NewClientsetFromConfig(test.MustConfig(""))
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
				Namespace: localGlooHubNamespace,
				Name:      "fed-upstream",
			},
			Spec: gloo_fed_types.FederatedUpstreamSpec{
				Template: &gloo_fed_types.FederatedUpstreamSpec_Template{
					Spec:     upstreamSpec,
					Metadata: meta,
				},
				Placement: &v1alpha1.Placement{
					Namespaces: []string{remoteGlooNamespace},
					Clusters:   []string{remoteClusterContext},
				},
			},
		}

		err = clientset.FederatedUpstreams().CreateFederatedUpstream(ctx, fedUpstream)
		Expect(err).NotTo(HaveOccurred())

		remoteClientSet, err := gloo_v1.NewClientsetFromConfig(test.MustConfig(remoteClusterContext))
		Expect(err).NotTo(HaveOccurred())
		var resultingUpstream *gloo_v1.Upstream
		Eventually(func() *gloo_v1.Upstream {
			resultingUpstream, _ = remoteClientSet.Upstreams().
				GetUpstream(ctx, types.NamespacedName{Name: meta.Name, Namespace: remoteGlooNamespace})
			return resultingUpstream
		}, 10*time.Second).ShouldNot(BeNil())

		Expect(resultingUpstream.Spec).To(Equal(*upstreamSpec))
		Expect(resultingUpstream.Annotations).To(Equal(meta.Annotations))
		Expect(resultingUpstream.Labels).To(Equal(map[string]string{
			"label":             "printer",
			federation.HubOwner: localGlooHubNamespace + ".fed-upstream",
		}))

		Eventually(func() bool {
			resultingUpstream, _ = remoteClientSet.Upstreams().
				GetUpstream(ctx, types.NamespacedName{Name: meta.Name, Namespace: remoteGlooNamespace})
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
				GetUpstream(ctx, types.NamespacedName{Name: meta.Name, Namespace: remoteGlooNamespace})
			if err != nil {
				return errors.IsNotFound(err)
			}
			return false
		}, 10*time.Second).Should(BeTrue(), "remote upstream should be cleaned up when federated upstream is deleted")
	})
})
