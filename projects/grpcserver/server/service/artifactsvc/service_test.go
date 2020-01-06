package artifactsvc_test

import (
	"context"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client/mocks"
	mock_settings "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	mock_gloo "github.com/solo-io/gloo/projects/gloo/pkg/mocks"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	mock_license "github.com/solo-io/solo-projects/pkg/license/mocks"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/artifactsvc"
)

var (
	apiserver      v1.ArtifactApiServer
	mockCtrl       *gomock.Controller
	artifactClient *mock_gloo.MockArtifactClient
	licenseClient  *mock_license.MockClient
	valuesClient   *mock_settings.MockValuesClient
	clientCache    *mocks.MockClientCache
	testErr        = errors.Errorf("test-err")
)

var _ = Describe("ServiceTest", func() {

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		artifactClient = mock_gloo.NewMockArtifactClient(mockCtrl)
		licenseClient = mock_license.NewMockClient(mockCtrl)
		valuesClient = mock_settings.NewMockValuesClient(mockCtrl)
		clientCache = mocks.NewMockClientCache(mockCtrl)
		clientCache.EXPECT().GetArtifactClient().Return(artifactClient).AnyTimes()
		apiserver = artifactsvc.NewArtifactGrpcService(context.TODO(), clientCache, licenseClient, valuesClient)
	})

	AfterEach(func() {
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
				Data:     map[string]string{"test": "asdf"},
				Metadata: metadata,
			}

			artifactClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(&artifact, nil)

			request := &v1.GetArtifactRequest{Ref: &ref}
			actual, err := apiserver.GetArtifact(context.TODO(), request)
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
				Return(nil, testErr)

			request := &v1.GetArtifactRequest{Ref: &ref}
			_, err := apiserver.GetArtifact(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := artifactsvc.FailedToReadArtifactError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ListArtifacts", func() {
		It("works when the artifact client works", func() {
			ns1, ns2 := "one", "two"
			artifact1 := gloov1.Artifact{
				Data:     map[string]string{"test": "asdf"},
				Metadata: core.Metadata{Namespace: ns1},
			}
			artifact2 := gloov1.Artifact{
				Data:     map[string]string{"test": "qwerty"},
				Metadata: core.Metadata{Namespace: ns2},
			}

			valuesClient.EXPECT().GetWatchNamespaces().Return([]string{ns1, ns2})
			artifactClient.EXPECT().
				List(ns1, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Artifact{&artifact1}, nil)
			artifactClient.EXPECT().
				List(ns2, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Artifact{&artifact2}, nil)

			actual, err := apiserver.ListArtifacts(context.TODO(), &v1.ListArtifactsRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListArtifactsResponse{Artifacts: []*gloov1.Artifact{&artifact1, &artifact2}}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the artifact client errors", func() {
			ns := "ns"

			valuesClient.EXPECT().GetWatchNamespaces().Return([]string{ns})
			artifactClient.EXPECT().
				List(ns, clients.ListOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			_, err := apiserver.ListArtifacts(context.TODO(), &v1.ListArtifactsRequest{})
			Expect(err).To(HaveOccurred())
			expectedErr := artifactsvc.FailedToListArtifactsError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("CreateArtifact", func() {
		Context("with unified input objects", func() {
			BeforeEach(func() {
				licenseClient.EXPECT().IsLicenseValid().Return(nil)
			})
			buildArtifact := func() *gloov1.Artifact {
				metadata := core.Metadata{
					Namespace:     "ns",
					Name:          "name",
					XXX_sizecache: -1,
				}
				return &gloov1.Artifact{
					Data:     map[string]string{"test": "asdf"},
					Metadata: metadata,
				}
			}

			It("works when the artifact client works", func() {
				artifact := buildArtifact()
				artifactClient.EXPECT().
					Write(artifact, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: false}).
					Return(artifact, nil)

				request := &v1.CreateArtifactRequest{
					Artifact: artifact,
				}
				actual, err := apiserver.CreateArtifact(context.TODO(), request)
				Expect(err).NotTo(HaveOccurred())
				expected := &v1.CreateArtifactResponse{Artifact: buildArtifact()}
				ExpectEqualProtoMessages(actual, expected)
			})

			It("errors when the artifact client errors", func() {
				artifact := buildArtifact()
				artifactClient.EXPECT().
					Write(artifact, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: false}).
					Return(nil, testErr)

				request := &v1.CreateArtifactRequest{
					Artifact: artifact,
				}
				_, err := apiserver.CreateArtifact(context.TODO(), request)
				Expect(err).To(HaveOccurred())
				expectedErr := artifactsvc.FailedToCreateArtifactError(testErr, &core.ResourceRef{
					Namespace: "ns",
					Name:      "name",
				})
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})
	})

	Describe("UpdateArtifact", func() {
		Context("with unified input objects", func() {
			BeforeEach(func() {
				licenseClient.EXPECT().IsLicenseValid().Return(nil)
			})
			buildArtifact := func(testValue string) *gloov1.Artifact {
				metadata := core.Metadata{
					Namespace: "ns",
					Name:      "name",
				}
				return &gloov1.Artifact{
					Data:     map[string]string{"test": testValue},
					Metadata: metadata,
				}
			}
			It("works when the artifact client works", func() {
				testArtifact := buildArtifact("new-test-value")
				artifactClient.EXPECT().
					Write(testArtifact, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
					Return(testArtifact, nil)

				request := &v1.UpdateArtifactRequest{
					Artifact: buildArtifact("new-test-value"),
				}
				actual, err := apiserver.UpdateArtifact(context.TODO(), request)
				Expect(err).NotTo(HaveOccurred())
				expected := &v1.UpdateArtifactResponse{Artifact: buildArtifact("new-test-value")}
				ExpectEqualProtoMessages(actual, expected)
			})

			It("errors when the artifact client errors on write", func() {
				namespace := "ns"
				name := "name"
				testArtifact := buildArtifact("test-value")
				artifactClient.EXPECT().
					Write(testArtifact, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
					Return(nil, testErr)

				request := &v1.UpdateArtifactRequest{
					Artifact: testArtifact,
				}
				_, err := apiserver.UpdateArtifact(context.TODO(), request)
				Expect(err).To(HaveOccurred())
				expectedErr := artifactsvc.FailedToUpdateArtifactError(testErr, &core.ResourceRef{
					Name:      name,
					Namespace: namespace,
				})
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})
	})

	Describe("DeleteArtifact", func() {
		It("works when the artifact client works", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}

			licenseClient.EXPECT().IsLicenseValid().Return(nil)
			artifactClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(nil)

			request := &v1.DeleteArtifactRequest{Ref: &ref}
			actual, err := apiserver.DeleteArtifact(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.DeleteArtifactResponse{}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the artifact client errors", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}

			licenseClient.EXPECT().IsLicenseValid().Return(nil)
			artifactClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(testErr)

			request := &v1.DeleteArtifactRequest{Ref: &ref}
			_, err := apiserver.DeleteArtifact(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := artifactsvc.FailedToDeleteArtifactError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
