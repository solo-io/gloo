package artifactsvc_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	mock_gloo "github.com/solo-io/gloo/projects/gloo/pkg/mocks"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/artifactsvc"
	. "github.com/solo-io/solo-projects/projects/grpcserver/server/service/internal/testutils"
	"google.golang.org/grpc"
)

var (
	grpcServer     *grpc.Server
	conn           *grpc.ClientConn
	apiserver      v1.ArtifactApiServer
	client         v1.ArtifactApiClient
	mockCtrl       *gomock.Controller
	artifactClient *mock_gloo.MockArtifactClient
	testErr        = errors.Errorf("test-err")
)

var _ = Describe("ServiceTest", func() {

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		artifactClient = mock_gloo.NewMockArtifactClient(mockCtrl)
		apiserver = artifactsvc.NewArtifactGrpcService(context.TODO(), artifactClient)

		grpcServer, conn = MustRunGrpcServer(func(s *grpc.Server) { v1.RegisterArtifactApiServer(s, apiserver) })
		client = v1.NewArtifactApiClient(conn)
	})

	AfterEach(func() {
		grpcServer.Stop()
		mockCtrl.Finish()
	})

	Describe("GetArtifact", func() {
		It("works when the artifact client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			artifact := gloov1.Artifact{
				Data:     "asdf",
				Metadata: metadata,
			}

			artifactClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(&artifact, nil).
				Times(1)

			request := &v1.GetArtifactRequest{Ref: &ref}
			actual, err := client.GetArtifact(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetArtifactResponse{Artifact: &artifact}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the artifact client errors", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			artifactClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr).
				Times(1)

			request := &v1.GetArtifactRequest{Ref: &ref}
			_, err := client.GetArtifact(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := artifactsvc.FailedToReadArtifactError(testErr, ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ListArtifacts", func() {
		It("works when the artifact client works", func() {
			ns1, ns2 := "one", "two"
			artifact1 := gloov1.Artifact{
				Data:     "asdf",
				Metadata: core.Metadata{Namespace: ns1},
			}
			artifact2 := gloov1.Artifact{
				Data:     "qwerty",
				Metadata: core.Metadata{Namespace: ns2},
			}

			artifactClient.EXPECT().
				List(ns1, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Artifact{&artifact1}, nil).
				Times(1)
			artifactClient.EXPECT().
				List(ns2, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Artifact{&artifact2}, nil).
				Times(1)

			request := &v1.ListArtifactsRequest{NamespaceList: []string{ns1, ns2}}
			actual, err := client.ListArtifacts(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListArtifactsResponse{ArtifactList: []*gloov1.Artifact{&artifact1, &artifact2}}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the artifact client errors", func() {
			ns := "ns"

			artifactClient.EXPECT().
				List(ns, clients.ListOpts{Ctx: context.TODO()}).
				Return(nil, testErr).
				Times(1)

			request := &v1.ListArtifactsRequest{NamespaceList: []string{ns}}
			_, err := client.ListArtifacts(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := artifactsvc.FailedToListArtifactsError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("CreateArtifact", func() {
		It("works when the artifact client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			artifact := gloov1.Artifact{
				Data:     "asdf",
				Metadata: metadata,
			}

			artifactClient.EXPECT().
				Write(&artifact, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: false}).
				Return(&artifact, nil).
				Times(1)

			request := &v1.CreateArtifactRequest{Ref: &ref, Data: artifact.Data}
			actual, err := client.CreateArtifact(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.CreateArtifactResponse{Artifact: &artifact}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the artifact client errors", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			artifact := gloov1.Artifact{
				Data:     "asdf",
				Metadata: metadata,
			}

			artifactClient.EXPECT().
				Write(&artifact, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: false}).
				Return(nil, testErr).
				Times(1)

			request := &v1.CreateArtifactRequest{Ref: &ref, Data: artifact.Data}
			_, err := client.CreateArtifact(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := artifactsvc.FailedToCreateArtifactError(testErr, ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("UpdateArtifact", func() {
		It("works when the artifact client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			oldArtifact := gloov1.Artifact{
				Data:     "asdf",
				Metadata: metadata,
			}
			newArtifact := gloov1.Artifact{
				Data:     "qwerty",
				Metadata: metadata,
			}

			artifactClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(&oldArtifact, nil).
				Times(1)
			artifactClient.EXPECT().
				Write(&newArtifact, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(&newArtifact, nil).
				Times(1)

			request := &v1.UpdateArtifactRequest{Ref: &ref, Data: newArtifact.Data}
			actual, err := client.UpdateArtifact(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.UpdateArtifactResponse{Artifact: &newArtifact}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the artifact client errors on read", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			artifact := gloov1.Artifact{
				Data:     "asdf",
				Metadata: metadata,
			}

			artifactClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr).
				Times(1)

			request := &v1.UpdateArtifactRequest{Ref: &ref, Data: artifact.Data}
			_, err := client.UpdateArtifact(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := artifactsvc.FailedToUpdateArtifactError(testErr, ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})

		It("errors when the artifact client errors on write", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			artifact := gloov1.Artifact{
				Data:     "asdf",
				Metadata: metadata,
			}

			artifactClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(&artifact, nil).
				Times(1)
			artifactClient.EXPECT().
				Write(&artifact, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(nil, testErr).
				Times(1)

			request := &v1.UpdateArtifactRequest{Ref: &ref, Data: artifact.Data}
			_, err := client.UpdateArtifact(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := artifactsvc.FailedToUpdateArtifactError(testErr, ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("DeleteArtifact", func() {
		It("works when the artifact client works", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}

			artifactClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(nil).
				Times(1)

			request := &v1.DeleteArtifactRequest{Ref: &ref}
			actual, err := client.DeleteArtifact(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.DeleteArtifactResponse{}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the artifact client errors", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}

			artifactClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(testErr).
				Times(1)

			request := &v1.DeleteArtifactRequest{Ref: &ref}
			_, err := client.DeleteArtifact(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := artifactsvc.FailedToDeleteArtifactError(testErr, ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
