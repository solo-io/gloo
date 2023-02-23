package certgen_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	mock_k8s_ar_clients "github.com/solo-io/external-apis/pkg/api/k8s/admissionregistration.k8s.io/v1/mocks"
	mock_k8s_core_clients "github.com/solo-io/external-apis/pkg/api/k8s/core/v1/mocks"
	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/config"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/internal/certgen"
	admission_v1 "k8s.io/api/admissionregistration/v1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("SelfSigned", func() {

	var (
		ctx           context.Context
		ctrl          *gomock.Controller
		cfg           *config.Config
		vwcClient     *mock_k8s_ar_clients.MockValidatingWebhookConfigurationClient
		arClientset   *mock_k8s_ar_clients.MockClientset
		secretClient  *mock_k8s_core_clients.MockSecretClient
		coreClientset *mock_k8s_core_clients.MockClientset

		testErr = eris.New("hello")
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())
		cfg = config.NewConfig()
		arClientset = mock_k8s_ar_clients.NewMockClientset(ctrl)
		vwcClient = mock_k8s_ar_clients.NewMockValidatingWebhookConfigurationClient(ctrl)
		coreClientset = mock_k8s_core_clients.NewMockClientset(ctrl)
		secretClient = mock_k8s_core_clients.NewMockSecretClient(ctrl)

		arClientset.EXPECT().
			ValidatingWebhookConfigurations().
			Return(vwcClient)
		coreClientset.EXPECT().
			Secrets().
			Return(secretClient)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("EnsureCaCerts", func() {

		It("will fail if secret get returns an error that isn't NotFound", func() {
			caManager := certgen.NewSelfSignedWebhookCAManager(coreClientset, arClientset, cfg)
			secretName, secretNamespace := "secret-name", "secret-namespace"

			secretClient.EXPECT().
				GetSecret(ctx, client.ObjectKey{
					Name:      secretName,
					Namespace: secretNamespace,
				}).
				Return(nil, testErr)
			err := caManager.EnsureCaCerts(ctx, secretName, secretNamespace, "", "")
			Expect(err).To(HaveOccurred())
			Expect(err).To(testutils.HaveInErrorChain(certgen.CaSecretError(testErr, secretName, secretNamespace)))
		})

		It("will fail if ValidatingWebhookConfiguration cannot be found", func() {
			caManager := certgen.NewSelfSignedWebhookCAManager(coreClientset, arClientset, cfg)
			secretName, secretNamespace := "secret-name", "secret-namespace"

			secretClient.EXPECT().
				GetSecret(ctx, client.ObjectKey{
					Name:      secretName,
					Namespace: secretNamespace,
				}).
				Return(&core_v1.Secret{Data: map[string][]byte{}}, nil)

			vwcClient.EXPECT().
				GetValidatingWebhookConfiguration(ctx, client.ObjectKey{
					Name: cfg.GetString(config.ValidatingWebhookConfigurationName),
				}).
				Return(nil, testErr)
			err := caManager.EnsureCaCerts(ctx, secretName, secretNamespace, "", "")
			Expect(err).To(HaveOccurred())
			Expect(err).To(testutils.HaveInErrorChain(testErr))
		})

		It("will update ValidatingWebhookConfiguration with existing CA from secret if it exists", func() {
			caManager := certgen.NewSelfSignedWebhookCAManager(coreClientset, arClientset, cfg)
			secretName, secretNamespace := "secret-name", "secret-namespace"
			svcName, svcNamespace := "svc-name", "svc-namespace"

			caBundle := []byte("ca-bundle")
			secretClient.EXPECT().
				GetSecret(ctx, client.ObjectKey{
					Name:      secretName,
					Namespace: secretNamespace,
				}).
				Return(&core_v1.Secret{
					Data: map[string][]byte{
						core_v1.ServiceAccountRootCAKey: caBundle,
					},
				}, nil)

			vwc := &admission_v1.ValidatingWebhookConfiguration{
				Webhooks: []admission_v1.ValidatingWebhook{
					{
						ClientConfig: admission_v1.WebhookClientConfig{
							Service: &admission_v1.ServiceReference{
								Namespace: svcNamespace,
								Name:      svcName,
							},
						},
					},
				},
			}

			vwcClient.EXPECT().
				GetValidatingWebhookConfiguration(ctx, client.ObjectKey{
					Name: cfg.GetString(config.ValidatingWebhookConfigurationName),
				}).
				Return(vwc, nil)

			update := vwc.DeepCopyObject().(*admission_v1.ValidatingWebhookConfiguration)
			update.Webhooks[0].ClientConfig.CABundle = caBundle
			vwcClient.EXPECT().
				UpdateValidatingWebhookConfiguration(ctx, update).
				Return(nil)

			err := caManager.EnsureCaCerts(ctx, secretName, secretNamespace, svcName, svcNamespace)
			Expect(err).NotTo(HaveOccurred())
		})

		It("will update ValidatingWebhookConfiguration new generated CA if secret not found", func() {
			caManager := certgen.NewSelfSignedWebhookCAManager(coreClientset, arClientset, cfg)
			secretName, secretNamespace := "secret-name", "secret-namespace"
			svcName, svcNamespace := "svc-name", "svc-namespace"

			secretClient.EXPECT().
				GetSecret(ctx, client.ObjectKey{
					Name:      secretName,
					Namespace: secretNamespace,
				}).
				Return(nil, errors.NewNotFound(schema.GroupResource{}, ""))

			var caBundle []byte
			secretClient.EXPECT().
				// Testing result in DoAndReturn
				CreateSecret(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, secret *core_v1.Secret, opts ...client.CreateOption) error {
					Expect(secret.GetName()).To(Equal(secretName))
					Expect(secret.GetNamespace()).To(Equal(secretNamespace))
					Expect(secret.Type).To(Equal(core_v1.SecretTypeTLS))
					Expect(secret.Data[core_v1.TLSPrivateKeyKey]).NotTo(HaveLen(0))
					Expect(secret.Data[core_v1.TLSCertKey]).NotTo(HaveLen(0))
					Expect(secret.Data[core_v1.ServiceAccountRootCAKey]).NotTo(HaveLen(0))
					caBundle = secret.Data[core_v1.ServiceAccountRootCAKey]
					return nil
				})

			vwc := &admission_v1.ValidatingWebhookConfiguration{
				Webhooks: []admission_v1.ValidatingWebhook{
					{
						ClientConfig: admission_v1.WebhookClientConfig{
							Service: &admission_v1.ServiceReference{
								Namespace: svcNamespace,
								Name:      svcName,
							},
						},
					},
				},
			}

			vwcClient.EXPECT().
				GetValidatingWebhookConfiguration(ctx, client.ObjectKey{
					Name: cfg.GetString(config.ValidatingWebhookConfigurationName),
				}).
				Return(vwc, nil)

			update := vwc.DeepCopyObject().(*admission_v1.ValidatingWebhookConfiguration)
			vwcClient.EXPECT().
				UpdateValidatingWebhookConfiguration(ctx, gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, vwc *admission_v1.ValidatingWebhookConfiguration, opts ...client.UpdateOption) error {
					update.Webhooks[0].ClientConfig.CABundle = caBundle
					Expect(update).To(Equal(vwc))
					return nil
				})

			err := caManager.EnsureCaCerts(ctx, secretName, secretNamespace, svcName, svcNamespace)
			Expect(err).NotTo(HaveOccurred())
		})

	})

})
