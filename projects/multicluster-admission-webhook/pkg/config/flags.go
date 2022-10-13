package config

import (
	"strings"

	"github.com/spf13/viper"
)

/*
	Viper offers a way to unify global config from multiple sources.
	Part of that contract is having string keys to access the data, all known keys should be defined below.
*/
const (
	PodNamespace                       = "pod_namespace"
	ServiceName                        = "service_name"
	TlsCertSecretName                  = "secret_name"
	ValidatingWebhookConfigurationName = "validating_webhook_configuration_name"
	CertDir                            = "cert_dir"
	AdmissionWebhookPath               = "webhook_path"
)

// These env vars are defined in the Helm chart for the deployment,
// in multicluster-admission-webhook/install/helm/chart.go
var (
	PodNamespaceEnvVar                       = strings.ToUpper(PodNamespace)
	ServiceNameEnvVar                        = strings.ToUpper(ServiceName)
	TlsCertSecretNameEnvVar                  = strings.ToUpper(TlsCertSecretName)
	ValidatingWebhookConfigurationNameEnvVar = strings.ToUpper(ValidatingWebhookConfigurationName)
	CertDirEnvVar                            = strings.ToUpper(CertDir)
	AdmissionWebhookPathEnvVar               = strings.ToUpper(AdmissionWebhookPath)
)

type Config struct {
	*viper.Viper
}

func (c *Config) setDefaults() {
	c.SetDefault(PodNamespace, "multicluster-admission")
	c.SetDefault(ServiceName, "multicluster-admission-webhook")
	c.SetDefault(TlsCertSecretName, "multicluster-admission-webhook")
	c.SetDefault(ValidatingWebhookConfigurationName, "multicluster-admission-webhook")
	c.SetDefault(CertDir, "/etc/certs/admission")
	c.SetDefault(AdmissionWebhookPath, "/admission")
}

func NewConfig() *Config {
	newConfig := &Config{
		Viper: viper.New(),
	}
	newConfig.setDefaults()
	// Bind all values to env variables
	// This will convert the keys used to upper case before searching, ex: pod_namespace -> POD_NAMESPACE
	newConfig.AutomaticEnv()
	return newConfig
}
