package webhook

import (
	"context"

	k8sArV1 "github.com/solo-io/external-apis/pkg/api/k8s/admissionregistration.k8s.io/v1"
	k8sCoreV1Clients "github.com/solo-io/external-apis/pkg/api/k8s/core/v1"
	multicluster_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/codegen/chart"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/config"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/internal/certgen"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/internal/handler"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/internal/placement"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/internal/validation"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/rbac"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Register the webhook with the k8s webhook server associated with the provided manager,
// using the provided webhook config and placement parser.
func InitializeWebhook(
	ctx context.Context,
	mgr manager.Manager,
	cfg *config.Config,
	parser rbac.Parser,
) error {
	if err := multicluster_v1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		return err
	}
	return mgr.Add(manager.RunnableFunc(func(ctx context.Context) error {
		if err := initWebhookCa(ctx, mgr.GetClient(), cfg); err != nil {
			return err
		}

		admissionClientset := multicluster_v1alpha1.NewClientset(mgr.GetClient())
		placementMatcher := placement.NewMatcher()
		admissionValidator := validation.NewMultiClusterAdmissionValidator(admissionClientset, placementMatcher, parser)

		admissionWebhook := &admission.Webhook{
			Handler: handler.NewAdmissionWebhookHandler(admissionValidator),
		}
		webhookServer := mgr.GetWebhookServer()

		// Set the cert directory to the one specified
		webhookServer.CertDir = cfg.GetString(config.CertDir)

		// Override default webhook port of 443 to run as non-root user.
		webhookServer.Port = chart.WebhookPort
		webhookServer.Register(cfg.GetString(config.AdmissionWebhookPath), admissionWebhook)
		return nil
	}))
}

// Ensure a CA bundle for TLS when communicating to the webhook.
func initWebhookCa(
	ctx context.Context,
	client client.Client,
	cfg *config.Config,
) error {
	caManager := certgen.NewSelfSignedWebhookCAManager(
		k8sCoreV1Clients.NewClientset(client),
		k8sArV1.NewClientset(client),
		cfg,
	)
	return caManager.EnsureCaCerts(
		ctx,
		cfg.GetString(config.TlsCertSecretName),
		cfg.GetString(config.PodNamespace),
		cfg.GetString(config.ServiceName),
		cfg.GetString(config.PodNamespace),
	)
}
