package secretsvc_test

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
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/secretsvc"
	mock_scrub "github.com/solo-io/solo-projects/projects/grpcserver/server/service/secretsvc/scrub/mocks"
)

var (
	apiserver      v1.SecretApiServer
	mockCtrl       *gomock.Controller
	secretClient   *mock_gloo.MockSecretClient
	licenseClient  *mock_license.MockClient
	clientCache    *mocks.MockClientCache
	scrubber       *mock_scrub.MockScrubber
	settingsValues *mock_settings.MockValuesClient
	testErr        = errors.Errorf("test-err")
)

var _ = Describe("ServiceTest", func() {

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		secretClient = mock_gloo.NewMockSecretClient(mockCtrl)
		licenseClient = mock_license.NewMockClient(mockCtrl)
		scrubber = mock_scrub.NewMockScrubber(mockCtrl)
		clientCache = mocks.NewMockClientCache(mockCtrl)
		clientCache.EXPECT().GetSecretClient().Return(secretClient).AnyTimes()
		settingsValues = mock_settings.NewMockValuesClient(mockCtrl)
		apiserver = secretsvc.NewSecretGrpcService(context.TODO(), clientCache, scrubber, licenseClient, settingsValues)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetSecret", func() {
		It("works when the secret client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			secret := gloov1.Secret{
				Kind:     &gloov1.Secret_Aws{Aws: &gloov1.AwsSecret{}},
				Metadata: metadata,
			}

			secretClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(&secret, nil)
			scrubber.EXPECT().Secret(context.Background(), &secret)

			request := &v1.GetSecretRequest{Ref: &ref}
			actual, err := apiserver.GetSecret(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetSecretResponse{Secret: &secret}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the secret client errors", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			secretClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.GetSecretRequest{Ref: &ref}
			_, err := apiserver.GetSecret(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := secretsvc.FailedToReadSecretError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ListSecrets", func() {
		It("works when the secret client works", func() {
			ns1, ns2 := "one", "two"
			secret1 := gloov1.Secret{
				Kind:     &gloov1.Secret_Aws{Aws: &gloov1.AwsSecret{}},
				Metadata: core.Metadata{Namespace: ns1},
			}
			secret2 := gloov1.Secret{
				Kind:     &gloov1.Secret_Azure{Azure: &gloov1.AzureSecret{}},
				Metadata: core.Metadata{Namespace: ns2},
			}

			settingsValues.EXPECT().GetWatchNamespaces().Return([]string{ns1, ns2})
			secretClient.EXPECT().
				List(ns1, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Secret{&secret1}, nil)
			secretClient.EXPECT().
				List(ns2, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Secret{&secret2}, nil)
			scrubber.EXPECT().Secret(context.Background(), &secret1)
			scrubber.EXPECT().Secret(context.Background(), &secret2)

			actual, err := apiserver.ListSecrets(context.TODO(), &v1.ListSecretsRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListSecretsResponse{Secrets: []*gloov1.Secret{&secret1, &secret2}}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the secret client errors", func() {
			ns := "ns"

			settingsValues.EXPECT().GetWatchNamespaces().Return([]string{ns})
			secretClient.EXPECT().
				List(ns, clients.ListOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			_, err := apiserver.ListSecrets(context.TODO(), &v1.ListSecretsRequest{})
			Expect(err).To(HaveOccurred())
			expectedErr := secretsvc.FailedToListSecretsError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("CreateSecret", func() {
		Context("with unified input objects", func() {
			It("works when the secret client works", func() {
				metadata := core.Metadata{
					Namespace: "ns",
					Name:      "name",
				}

				testCases := []gloov1.Secret{
					{
						Kind:     &gloov1.Secret_Aws{Aws: &gloov1.AwsSecret{}},
						Metadata: metadata,
					},
					{
						Kind:     &gloov1.Secret_Azure{Azure: &gloov1.AzureSecret{}},
						Metadata: metadata,
					},
					{
						Kind:     &gloov1.Secret_Tls{Tls: &gloov1.TlsSecret{}},
						Metadata: metadata,
					},
				}

				for _, tc := range testCases {
					secretClient.EXPECT().
						Write(&tc, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: false}).
						Return(&tc, nil)
					scrubber.EXPECT().Secret(context.Background(), &tc)
					licenseClient.EXPECT().IsLicenseValid().Return(nil)

					actual, err := apiserver.CreateSecret(context.TODO(), &v1.CreateSecretRequest{
						Secret: &tc,
					})
					Expect(err).NotTo(HaveOccurred())
					expected := &v1.CreateSecretResponse{Secret: &tc}
					ExpectEqualProtoMessages(actual, expected)
				}
			})

			It("errors when the secret client errors", func() {
				metadata := core.Metadata{
					Namespace: "ns",
					Name:      "name",
				}
				ref := metadata.Ref()
				secret := gloov1.Secret{
					Kind:     &gloov1.Secret_Aws{Aws: &gloov1.AwsSecret{}},
					Metadata: metadata,
				}

				secretClient.EXPECT().
					Write(&secret, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: false}).
					Return(nil, testErr)
				licenseClient.EXPECT().IsLicenseValid().Return(nil)

				request := &v1.CreateSecretRequest{
					Secret: &secret,
				}
				_, err := apiserver.CreateSecret(context.TODO(), request)
				Expect(err).To(HaveOccurred())
				expectedErr := secretsvc.FailedToCreateSecretError(testErr, &ref)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})
	})

	Describe("UpdateSecret", func() {
		Context("with unified input objects", func() {
			It("works when the secret client works", func() {
				metadata := core.Metadata{
					Namespace: "ns",
					Name:      "name",
				}

				testCases := []gloov1.Secret{
					{
						Kind:     &gloov1.Secret_Aws{Aws: &gloov1.AwsSecret{}},
						Metadata: metadata,
					},
					{
						Kind:     &gloov1.Secret_Azure{Azure: &gloov1.AzureSecret{}},
						Metadata: metadata,
					},
					{
						Kind:     &gloov1.Secret_Tls{Tls: &gloov1.TlsSecret{}},
						Metadata: metadata,
					},
				}

				for _, tc := range testCases {
					secretClient.EXPECT().
						Write(&tc, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
						Return(&tc, nil)
					scrubber.EXPECT().Secret(context.Background(), &tc)
					licenseClient.EXPECT().IsLicenseValid().Return(nil)

					actual, err := apiserver.UpdateSecret(context.TODO(), &v1.UpdateSecretRequest{
						Secret: &tc,
					})
					Expect(err).NotTo(HaveOccurred())
					expected := &v1.UpdateSecretResponse{Secret: &tc}
					ExpectEqualProtoMessages(actual, expected)
				}
			})

			It("errors when the secret client errors on write", func() {
				metadata := core.Metadata{}
				ref := metadata.Ref()
				secret := gloov1.Secret{
					Kind: &gloov1.Secret_Aws{
						Aws: &gloov1.AwsSecret{},
					},
					Metadata: metadata,
				}

				secretClient.EXPECT().
					Write(&secret, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
					Return(nil, testErr)
				licenseClient.EXPECT().IsLicenseValid().Return(nil)

				request := &v1.UpdateSecretRequest{
					Secret: &secret,
				}
				_, err := apiserver.UpdateSecret(context.TODO(), request)
				Expect(err).To(HaveOccurred())
				expectedErr := secretsvc.FailedToUpdateSecretError(testErr, &ref)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})
	})

	Describe("DeleteSecret", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})
		It("works when the secret client works", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}

			secretClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(nil)

			request := &v1.DeleteSecretRequest{Ref: &ref}
			actual, err := apiserver.DeleteSecret(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.DeleteSecretResponse{}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the secret client errors", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}

			secretClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(testErr)

			request := &v1.DeleteSecretRequest{Ref: &ref}
			_, err := apiserver.DeleteSecret(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := secretsvc.FailedToDeleteSecretError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
