package generate

import (
	glooGen "github.com/solo-io/gloo/install/helm/gloo/generate"
)

type HelmConfig struct {
	Global                 *GlobalConfig        `json:"global,omitempty"`
	Enabled                *bool                `json:"enabled,omitempty" desc:"If true, deploy federation service (default true)."`
	CreateLicenseSecret    *bool                `json:"create_license_secret,omitempty" desc:"Create a secret for the license specified in 'license_key'. Set to 'false' if you use 'license_secret_name' instead."`
	LicenseSecretName      *string              `json:"license_secret_name,omitempty" desc:"The name of a secret that contains your Gloo Edge license key. Set 'create_license_key' to 'false' to disable use of the default license secret."`
	LicenseKey             *string              `json:"license_key,omitempty" desc:"Your Gloo Edge license key."`
	EnableMultiClusterRbac *bool                `json:"enableMultiClusterRbac,omitempty"`
	GlooFedApiServer       *ApiServerDeployment `json:"glooFedApiserver,omitempty"`
	GlooFed                *GlooFedDeployment   `json:"glooFed,omitempty"`
	Rbac                   *RbacConfiguration   `json:"rbac,omitempty"`
	RbacWebhook            *RbacWebhook         `json:"rbacWebhook,omitempty"`
}

type GlobalConfig struct {
	*glooGen.Global
	Console *ConsoleOptions `json:"console,omitempty" desc:"Configuration options for the Enterprise Console (UI)."`
}

type GlooFedDeployment struct {
	Replicas  *int                          `json:"replicas,omitempty"`
	Resources *glooGen.ResourceRequirements `json:"resources,omitempty"`
	Image     *glooGen.Image                `json:"image,omitempty"`
	Stats     *glooGen.Stats                `json:"stats,omitempty"`

	Retries *GlooFedRetries `json:"retries,omitempty" desc:"Retry options for failures reconciling cluster events."`
}

type GlooFedRetries struct {
	ClusterWatcherRemote *RetryOptions `json:"clusterWatcherRemote,omitempty" desc:"Retry options used in case of failure when creating/starting a remote cluster manager (e.g. if a remote cluster is unreachable)."`
	ClusterWatcherLocal  *RetryOptions `json:"clusterWatcherLocal,omitempty" desc:"Retry options used in case of failure reconciling federated resources on the management cluster, in response to the addition/removal of a cluster (e.g. if there is an error reading or updating the federated resources on the management cluster)."`
}

type RetryOptions struct {
	Type      *string `json:"type,omitempty" desc:"The type of delay to use when retrying. Must be either 'backoff' (for exponential backoff) or 'fixed' (for fixed intervals)."`
	Delay     *string `json:"delay,omitempty" desc:"The delay between retries. For exponential backoff, this is the delay for the initial retry. This must be a [Duration string](https://pkg.go.dev/time#Duration.String), e.g. '100ms' or '1m5s'."`
	MaxDelay  *string `json:"maxDelay,omitempty" desc:"The maximum delay between retries. This can be used to cap the retry interval when exponential backoff is used. If set to 0, there will be no maximum delay. This must be a [Duration string](https://pkg.go.dev/time#Duration.String), e.g. '100ms' or '1m5s'."`
	MaxJitter *string `json:"maxJitter,omitempty" desc:"The maximum amount of random jitter to add between retries. If this value is greater than 0, retries will be done with a random amount of jitter, up to maxJitter. If this value is 0, then no randomness will be added to retries. This must be a [Duration string](https://pkg.go.dev/time#Duration.String), e.g. '100ms' or '1m5s'."`
	Attempts  *uint   `json:"attempts,omitempty" desc:"The maximum number of attempts to make. Set to 0 to retry forever."`
}

type ApiServerDeployment struct {
	Enable                  *bool                         `json:"enable,omitempty"`
	Replicas                *int                          `json:"replicas,omitempty"`
	Image                   *glooGen.Image                `json:"image,omitempty"`
	Port                    *int                          `json:"port,omitempty"`
	HealthCheckPort         *int                          `json:"healthCheckPort,omitempty"`
	Resources               *glooGen.ResourceRequirements `json:"resources,omitempty"`
	Stats                   *glooGen.Stats                `json:"stats,omitempty"`
	FloatingUserId          *bool                         `json:"floatingUserId,omitempty"`
	RunAsUser               *float64                      `json:"runAsUser,omitempty"`
	Console                 *ConsoleContainer             `json:"console,omitempty"`
	Envoy                   *EnvoyContainer               `json:"envoy,omitempty"`
	NamespaceRestrictedMode *bool                         `json:"namespaceRestrictedMode,omitempty" desc:"If true:  Convert the ClusterRole used in apiserver to Role.  Useful in single-namespace deployments of gloo-ee where permissions can be more restrictive--recommended to not set in a multi-cluster deployment.  Default is false."`
}

type ConsoleContainer struct {
	Image     *glooGen.Image                `json:"image,omitempty"`
	Port      *int                          `json:"port,omitempty"`
	Resources *glooGen.ResourceRequirements `json:"resources,omitempty"`
}

type EnvoyContainer struct {
	Image                 *glooGen.Image                `json:"image,omitempty"`
	Resources             *glooGen.ResourceRequirements `json:"resources,omitempty"`
	BoostrapConfiguration *EnvoyBootstrapConfiguration  `json:"bootstrapConfig,omitempty""`
}

type EnvoyBootstrapConfiguration struct {
	ConfigMapName *string `json:"configMapName,omitempty"`
}

type RbacWebhook struct {
	Image     *glooGen.Image                `json:"image,omitempty"`
	Resources *glooGen.ResourceRequirements `json:"resources,omitempty"`
}

type RbacConfiguration struct {
	Create *bool `json:"create,omitempty"`
}

type ConsoleOptions struct {
	ReadOnly           *bool `json:"readOnly,omitempty" desc:"If true, then custom resources can only be viewed in read-only mode in the UI. If false, then resources can be created, updated, and deleted via the UI (default false)."`
	ApiExplorerEnabled *bool `json:"apiExplorerEnabled,omitempty" desc:"Whether the GraphQL API Explorer is enabled (default true)."`
}
