package failover_scheme_handler_test

import (
	"context"

	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"

	rpc_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.rpc/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/failover_scheme_handler"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	mock_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/mocks"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("failover scheme handler", func() {
	var (
		ctrl                   *gomock.Controller
		failoverSchemeClient   *mock_v1.MockFailoverSchemeClient
		testFailoverSchemeList *fedv1.FailoverSchemeList
	)

	BeforeEach(func() {
		ctrl, _ = gomock.WithContext(context.TODO(), GinkgoT())

		failover1 := fedv1.FailoverScheme{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      "failover-test",
				Namespace: "gloo-fed",
			},
			Spec: types.FailoverSchemeSpec{
				Primary: &v1.ClusterObjectRef{
					Name:        "test-name",
					Namespace:   "test-namespace",
					ClusterName: "test-cluster",
				},
			},
			Status: types.FailoverSchemeStatus{},
		}
		failover2 := fedv1.FailoverScheme{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      "failover-test-2",
				Namespace: "gloo-fed",
			},
			Spec: types.FailoverSchemeSpec{
				Primary: &v1.ClusterObjectRef{
					Name:        "test-name-2",
					Namespace:   "test-namespace-2",
					ClusterName: "test-cluster-2",
				},
			},
			Status: types.FailoverSchemeStatus{},
		}
		testFailoverSchemeList = &fedv1.FailoverSchemeList{
			Items: []fedv1.FailoverScheme{
				failover1,
				failover2,
			},
		}
		failoverSchemeClient = mock_v1.NewMockFailoverSchemeClient(ctrl)
		failoverSchemeClient.EXPECT().ListFailoverScheme(gomock.Any()).Return(testFailoverSchemeList, nil).AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("returns nil if there are no failover schemes for an upstream", func() {
		failoverSchemeServer := failover_scheme_handler.NewFailoverSchemeHandler(failoverSchemeClient)
		resp, err := failoverSchemeServer.GetFailoverScheme(context.TODO(), &rpc_v1.GetFailoverSchemeRequest{
			UpstreamRef: &v1.ClusterObjectRef{
				Name: "nonexistent-test-name",
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).To(Equal(&rpc_v1.GetFailoverSchemeResponse{
			FailoverScheme: &rpc_v1.FailoverScheme{},
		}))
	})

	It("returns a failover if there is a failover scheme that matches an upstream", func() {
		failoverSchemeServer := failover_scheme_handler.NewFailoverSchemeHandler(failoverSchemeClient)
		resp, err := failoverSchemeServer.GetFailoverScheme(context.TODO(), &rpc_v1.GetFailoverSchemeRequest{
			UpstreamRef: &v1.ClusterObjectRef{
				Name:        "test-name",
				Namespace:   "test-namespace",
				ClusterName: "test-cluster",
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).To(Equal(&rpc_v1.GetFailoverSchemeResponse{
			FailoverScheme: &rpc_v1.FailoverScheme{
				Metadata: &rpc_v1.ObjectMeta{
					Name:      "failover-test",
					Namespace: "gloo-fed",
				},
				Spec: &types.FailoverSchemeSpec{
					Primary: &v1.ClusterObjectRef{
						Name:        "test-name",
						Namespace:   "test-namespace",
						ClusterName: "test-cluster",
					},
				},
				Status: &types.FailoverSchemeStatus{},
			}}))
	})

})
