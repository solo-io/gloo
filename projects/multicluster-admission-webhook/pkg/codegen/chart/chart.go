package chart

import (
	"os"
	"strings"

	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/config"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/gertd/go-pluralize"
	"github.com/gobuffalo/packr"
	"github.com/solo-io/go-utils/stringutils"
	"github.com/solo-io/skv2/codegen/model"
	"google.golang.org/protobuf/proto"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

/*
	Skv2 model.Chart definition, imported into downstream projects to generate
	Helm charts for the multicluster-admission-webhook.
*/

const (
	WebhookPort          = 8443                   // This port must be > 1024 so that the webhook pod can be run as a non-root user.
	AdmissionWebhookPath = "/admission"           // URL path for contacting the admission webhook
	certDir              = "/etc/certs/admission" // The directory in the volume mount containing the webhook TLS cert
	defaultRegistry      = "soloio"               // Default Docker registry for image
)

func GenerateChart(
	data model.Data,
	groups []model.Group,
	resourcesToSelect map[schema.GroupVersion][]string,
) *model.Chart {

	return &model.Chart{
		Operators: []model.Operator{
			makeOperator(data),
		},
		Data: data,
		// Exclude the default namespace and configmap template
		FilterTemplate: func(outPath string) bool {
			return stringutils.ContainsString(outPath, []string{
				"templates/namespace.yaml",
				"templates/configmap.yaml",
				"templates/deployment.yaml",
			},
			)
		},
		CustomTemplates: makeCustomTemplates(groups, resourcesToSelect),
	}
}

func makeCustomTemplates(
	groups []model.Group,
	resourcesToSelect map[schema.GroupVersion][]string,
) model.CustomTemplates {
	templatesBox := packr.NewBox("./")
	webhookConfig := "validating-webhook-configuration.yaml"
	webhookConfigFname := webhookConfig + "tmpl"
	webhookConfigContents, err := templatesBox.FindString(webhookConfigFname)
	if err != nil {
		panic(err)
	}
	deployServiceConfig := "deployment-service.yaml"
	deployServiceConfigFname := deployServiceConfig + "tmpl"
	deployServiceConfigContents, err := templatesBox.FindString(deployServiceConfigFname)
	if err != nil {
		panic(err)
	}
	values := "values.yaml"
	valuesFname := values + "tmpl"
	valuesContents, err := templatesBox.FindString(valuesFname)
	if err != nil {
		panic(err)
	}
	return model.CustomTemplates{
		Templates: map[string]string{
			"templates/" + webhookConfig:       webhookConfigContents,
			"templates/" + deployServiceConfig: deployServiceConfigContents,
			values:                             valuesContents,
		},
		Funcs: map[string]interface{}{
			"getAdmissionRules":       makeAdmissionRules(groups, resourcesToSelect),
			"getAdmissionWebhookPath": func() string { return AdmissionWebhookPath },
			"getWebhookPort":          func() int { return WebhookPort },
		},
	}
}

func makeAdmissionRules(
	groups []model.Group,
	resourcesToSelect map[schema.GroupVersion][]string,
) func() []admissionv1.RuleWithOperations {
	// Just in case arg was nil
	if resourcesToSelect == nil {
		resourcesToSelect = make(map[schema.GroupVersion][]string)
	}
	var admissionRules []admissionv1.RuleWithOperations
	for _, group := range groups {
		rule := admissionv1.RuleWithOperations{
			Operations: []admissionv1.OperationType{admissionv1.Create, admissionv1.Update, admissionv1.Delete},
			Rule: admissionv1.Rule{
				APIGroups:   []string{group.Group},
				APIVersions: []string{group.Version},
			},
		}
		if gv, ok := resourcesToSelect[group.GroupVersion]; ok {
			// If the groupversion has fine grained resources
			for _, groupResource := range group.Resources {
				if stringutils.ContainsString(groupResource.Kind, gv) {
					rule.Rule.Resources = append(rule.Rule.Resources,
						strings.ToLower(pluralize.NewClient().Plural(groupResource.Kind)),
					)
				}
			}
		} else {
			rule.Rule.Resources = []string{"*"}
		}
		admissionRules = append(admissionRules, rule)
	}
	return func() []admissionv1.RuleWithOperations {
		return admissionRules
	}
}

func makeOperator(data model.Data) model.Operator {
	return model.Operator{
		Name: data.Name,
		Deployment: model.Deployment{
			Container: model.Container{
				Image: makeImage(data),
				Resources: &v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("125m"),
						v1.ResourceMemory: resource.MustParse("256Mi"),
					},
				},
				VolumeMounts: []v1.VolumeMount{
					{
						Name:      "admission-certs",
						ReadOnly:  true,
						MountPath: certDir,
					},
				},
				Env: []v1.EnvVar{
					{
						Name: config.PodNamespaceEnvVar,
						ValueFrom: &v1.EnvVarSource{
							FieldRef: &v1.ObjectFieldSelector{
								FieldPath: "metadata.namespace",
							},
						},
					},
					{
						Name:  config.ServiceNameEnvVar,
						Value: data.Name,
					},
					{
						Name:  config.TlsCertSecretNameEnvVar,
						Value: data.Name,
					},
					{
						Name:  config.ValidatingWebhookConfigurationNameEnvVar,
						Value: data.Name,
					},
					{
						Name:  config.CertDirEnvVar,
						Value: certDir,
					},
					{
						Name:  config.AdmissionWebhookPathEnvVar,
						Value: AdmissionWebhookPath,
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: "admission-certs",
					VolumeSource: v1.VolumeSource{Secret: &v1.SecretVolumeSource{
						SecretName: data.Name,
						Optional:   proto.Bool(true),
					}},
				},
			},
		},
		Service: model.Service{
			Type: v1.ServiceTypeClusterIP,
			Ports: []model.ServicePort{
				{
					Name:        "webhook",
					DefaultPort: WebhookPort,
				},
			},
		},
		Rbac: []rbacv1.PolicyRule{
			// RBAC multicluster role and rolebindings
			{
				APIGroups: []string{"multicluster.solo.io"},
				Verbs:     []string{"get", "list", "watch"},
				Resources: []string{
					"multiclusterroles",
					"multiclusterrolebindings",
				},
			},
			// Webhook TLS cert secret
			{
				APIGroups: []string{""},
				Verbs:     []string{"*"},
				Resources: []string{"secrets"},
			},
			// Webhook configuration
			{
				APIGroups: []string{"admissionregistration.k8s.io"},
				Verbs:     []string{"get", "list", "watch", "update"},
				Resources: []string{"validatingwebhookconfigurations"},
			},
		},
	}
}

// cache and operator share same image
func makeImage(data model.Data) model.Image {
	registry := os.Getenv("IMAGE_REGISTRY")
	if registry == "" {
		registry = defaultRegistry
	}
	return model.Image{
		Registry:   registry,
		Repository: data.Name,
		Tag:        data.Version,
		PullPolicy: v1.PullIfNotPresent,
	}
}
