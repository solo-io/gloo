package secretsvc_test

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
	. "github.com/solo-io/solo-projects/projects/grpcserver/server/service/internal/testutils"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/secretsvc"
	"google.golang.org/grpc"
)

var (
	grpcServer   *grpc.Server
	conn         *grpc.ClientConn
	apiserver    v1.SecretApiServer
	client       v1.SecretApiClient
	mockCtrl     *gomock.Controller
	secretClient *mock_gloo.MockSecretClient
	testErr      = errors.Errorf("test-err")
)

var _ = Describe("ServiceTest", func() {

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		secretClient = mock_gloo.NewMockSecretClient(mockCtrl)
		apiserver = secretsvc.NewSecretGrpcService(context.TODO(), secretClient)

		grpcServer, conn = MustRunGrpcServer(func(s *grpc.Server) { v1.RegisterSecretApiServer(s, apiserver) })
		client = v1.NewSecretApiClient(conn)
	})

	AfterEach(func() {
		grpcServer.Stop()
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

			request := &v1.GetSecretRequest{Ref: &ref}
			actual, err := client.GetSecret(context.TODO(), request)
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
			_, err := client.GetSecret(context.TODO(), request)
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

			secretClient.EXPECT().
				List(ns1, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Secret{&secret1}, nil)

			secretClient.EXPECT().
				List(ns2, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Secret{&secret2}, nil)

			request := &v1.ListSecretsRequest{NamespaceList: []string{ns1, ns2}}
			actual, err := client.ListSecrets(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListSecretsResponse{SecretList: []*gloov1.Secret{&secret1, &secret2}}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the secret client errors", func() {
			ns := "ns"

			secretClient.EXPECT().
				List(ns, clients.ListOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.ListSecretsRequest{NamespaceList: []string{ns}}
			_, err := client.ListSecrets(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := secretsvc.FailedToListSecretsError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("CreateSecret", func() {
		It("works when the secret client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			testCases := []struct {
				request v1.CreateSecretRequest
				secret  gloov1.Secret
			}{
				{
					request: v1.CreateSecretRequest{
						Ref:  &ref,
						Kind: &v1.CreateSecretRequest_Aws{Aws: &gloov1.AwsSecret{}},
					},
					secret: gloov1.Secret{
						Kind:     &gloov1.Secret_Aws{Aws: &gloov1.AwsSecret{}},
						Metadata: metadata,
					},
				},
				{
					request: v1.CreateSecretRequest{
						Ref:  &ref,
						Kind: &v1.CreateSecretRequest_Azure{Azure: &gloov1.AzureSecret{}},
					},
					secret: gloov1.Secret{
						Kind:     &gloov1.Secret_Azure{Azure: &gloov1.AzureSecret{}},
						Metadata: metadata,
					},
				},
				{
					request: v1.CreateSecretRequest{
						Ref:  &ref,
						Kind: &v1.CreateSecretRequest_Extension{Extension: &gloov1.Extension{}},
					},
					secret: gloov1.Secret{
						Kind:     &gloov1.Secret_Extension{Extension: &gloov1.Extension{}},
						Metadata: metadata,
					},
				},
				{
					request: v1.CreateSecretRequest{
						Ref:  &ref,
						Kind: &v1.CreateSecretRequest_Tls{Tls: &gloov1.TlsSecret{}},
					},
					secret: gloov1.Secret{
						Kind:     &gloov1.Secret_Tls{Tls: &gloov1.TlsSecret{}},
						Metadata: metadata,
					},
				},
			}

			for _, tc := range testCases {
				secretClient.EXPECT().
					Write(&tc.secret, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: false}).
					Return(&tc.secret, nil)

				actual, err := client.CreateSecret(context.TODO(), &tc.request)
				Expect(err).NotTo(HaveOccurred())
				expected := &v1.CreateSecretResponse{Secret: &tc.secret}
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

			request := &v1.CreateSecretRequest{Ref: &ref, Kind: &v1.CreateSecretRequest_Aws{Aws: &gloov1.AwsSecret{}}}
			_, err := client.CreateSecret(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := secretsvc.FailedToCreateSecretError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("UpdateSecret", func() {
		It("works when the secret client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			oldSecret := gloov1.Secret{
				Kind:     &gloov1.Secret_Aws{},
				Metadata: metadata,
			}

			testCases := []struct {
				request   v1.UpdateSecretRequest
				newSecret gloov1.Secret
			}{
				{
					request: v1.UpdateSecretRequest{
						Ref:  &ref,
						Kind: &v1.UpdateSecretRequest_Aws{Aws: &gloov1.AwsSecret{}},
					},
					newSecret: gloov1.Secret{
						Kind:     &gloov1.Secret_Aws{Aws: &gloov1.AwsSecret{}},
						Metadata: metadata,
					},
				},
				{
					request: v1.UpdateSecretRequest{
						Ref:  &ref,
						Kind: &v1.UpdateSecretRequest_Azure{Azure: &gloov1.AzureSecret{}},
					},
					newSecret: gloov1.Secret{
						Kind:     &gloov1.Secret_Azure{Azure: &gloov1.AzureSecret{}},
						Metadata: metadata,
					},
				},
				{
					request: v1.UpdateSecretRequest{
						Ref:  &ref,
						Kind: &v1.UpdateSecretRequest_Extension{Extension: &gloov1.Extension{}},
					},
					newSecret: gloov1.Secret{
						Kind:     &gloov1.Secret_Extension{Extension: &gloov1.Extension{}},
						Metadata: metadata,
					},
				},
				{
					request: v1.UpdateSecretRequest{
						Ref:  &ref,
						Kind: &v1.UpdateSecretRequest_Tls{Tls: &gloov1.TlsSecret{}},
					},
					newSecret: gloov1.Secret{
						Kind:     &gloov1.Secret_Tls{Tls: &gloov1.TlsSecret{}},
						Metadata: metadata,
					},
				},
			}

			for _, tc := range testCases {
				secretClient.EXPECT().
					Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
					Return(&oldSecret, nil)
				secretClient.EXPECT().
					Write(&tc.newSecret, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
					Return(&tc.newSecret, nil)

				actual, err := client.UpdateSecret(context.TODO(), &tc.request)
				Expect(err).NotTo(HaveOccurred())
				expected := &v1.UpdateSecretResponse{Secret: &tc.newSecret}
				ExpectEqualProtoMessages(actual, expected)
			}
		})

		It("errors when the secret client errors on read", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			secretClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.UpdateSecretRequest{Ref: &ref, Kind: &v1.UpdateSecretRequest_Azure{Azure: &gloov1.AzureSecret{}}}
			_, err := client.UpdateSecret(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := secretsvc.FailedToUpdateSecretError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})

		It("errors when the secret client errors on write", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			secret := gloov1.Secret{
				Kind:     &gloov1.Secret_Aws{},
				Metadata: metadata,
			}

			secretClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(&secret, nil)

			secretClient.EXPECT().
				Write(&secret, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(nil, testErr)

			request := &v1.UpdateSecretRequest{Ref: &ref, Kind: &v1.UpdateSecretRequest_Azure{Azure: &gloov1.AzureSecret{}}}
			_, err := client.UpdateSecret(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := secretsvc.FailedToUpdateSecretError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("DeleteSecret", func() {
		It("works when the secret client works", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}

			secretClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(nil)

			request := &v1.DeleteSecretRequest{Ref: &ref}
			actual, err := client.DeleteSecret(context.TODO(), request)
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
			_, err := client.DeleteSecret(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := secretsvc.FailedToDeleteSecretError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
