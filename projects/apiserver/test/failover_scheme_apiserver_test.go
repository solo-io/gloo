package apiserver_test

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	rpc_fed_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.rpc/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/failover_scheme_handler"
	fed_solo_io_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	mock_fed_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/mocks"
)

var _ = Describe("FailoverSchemeApiServer", func() {

	var (
		ctrl                    *gomock.Controller
		ctx                     context.Context
		failoverSchemeClient    *mock_fed_v1.MockFailoverSchemeClient
		failoverSchemeApiServer rpc_fed_v1.FailoverSchemeApiServer
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())
		failoverSchemeClient = mock_fed_v1.NewMockFailoverSchemeClient(ctrl)
		failoverSchemeApiServer = failover_scheme_handler.NewFailoverSchemeHandler(failoverSchemeClient)

	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("can ListFailoverSchemes", func() {
		exampleFailoverYaml := `
apiVersion: fed.solo.io/v1
kind: FailoverScheme
metadata:
  name: failover
  namespace: gloo-system
spec:
  primary:
    cluster_name: kind-remote
    name: default-service-blue-test
    namespace: gloo-system
  failover_groups:
    - priority_group:
        - cluster: kind-remote
          upstreams:
            - name: default-service-green-test
              namespace: gloo-system
`
		testFailoverScheme := &fed_solo_io_v1.FailoverScheme{}
		err := yaml.Unmarshal([]byte(exampleFailoverYaml), testFailoverScheme)
		Expect(err).NotTo(HaveOccurred())
		failoverSchemeList := &fed_solo_io_v1.FailoverSchemeList{
			Items: []fed_solo_io_v1.FailoverScheme{
				*testFailoverScheme,
			},
		}
		failoverSchemeClient.EXPECT().ListFailoverScheme(ctx).Return(failoverSchemeList, nil)
		resp, err := failoverSchemeApiServer.GetFailoverScheme(ctx, &rpc_fed_v1.GetFailoverSchemeRequest{
			UpstreamRef: &skv2v1.ClusterObjectRef{
				Name:        "default-service-blue-test",
				Namespace:   "gloo-system",
				ClusterName: "kind-remote",
			},
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(resp.FailoverScheme).NotTo(BeNil())
		Expect(len(resp.FailoverScheme.GetSpec().GetFailoverGroups()) > 0).To(BeTrue())
	})

})
