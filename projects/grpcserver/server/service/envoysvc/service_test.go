package envoysvc_test

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/go-utils/testutils"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc"
	mock_envoy_details "github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc/envoydetails/mocks"
)

var (
	apiServer     v1.EnvoyApiServer
	mockCtrl      *gomock.Controller
	detailsClient *mock_envoy_details.MockClient
	podNamespace  = "test-pod-ns"
	testErr       = errors.New("test-err")
)

var _ = Describe("Envoy Service Test", func() {

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		detailsClient = mock_envoy_details.NewMockClient(mockCtrl)
		apiServer = envoysvc.NewEnvoyGrpcService(context.Background(), detailsClient, podNamespace)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("ListEnvoyDetails", func() {
		It("works when the envoy details client works", func() {
			detailsList := []*v1.EnvoyDetails{{Name: "test"}}
			detailsClient.EXPECT().List(context.Background(), podNamespace).Return(detailsList, nil)

			actual, err := apiServer.ListEnvoyDetails(context.Background(), &v1.ListEnvoyDetailsRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListEnvoyDetailsResponse{EnvoyDetails: detailsList}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the envoy details client errors", func() {
			detailsClient.EXPECT().List(context.Background(), podNamespace).Return(nil, testErr)

			_, err := apiServer.ListEnvoyDetails(context.Background(), &v1.ListEnvoyDetailsRequest{})
			Expect(err).To(HaveOccurred())
			expectedErr := envoysvc.FailedToListEnvoyDetailsError(testErr)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
