package generate

import (
	glooGen "github.com/solo-io/gloo/install/helm/gloo/generate"
	v1 "k8s.io/api/core/v1"
)

type HelmConfig struct {
	Config
	Global *glooGen.Global `json:"global,omitempty"`
}
type Config struct {
	Settings            *glooGen.Settings `json:"settings,omitempty"`
	LicenseKey          string            `json:"license_key,omitempty"`
	CreateLicenseSecret bool              `json:"create_license_secret"`
	LicenseSecretName   string            `json:"license_secret_name"`
	Gloo                *glooGen.Config   `json:"gloo,omitempty"`
	Redis               *Redis            `json:"redis,omitempty"`
	RateLimit           *RateLimit        `json:"rateLimit,omitempty"`
	Observability       *Observability    `json:"observability,omitempty"`
	Rbac                *Rbac             `json:"rbac"`
	Grafana             interface{}       `json:"grafana,omitempty"`
	Prometheus          interface{}       `json:"prometheus,omitempty"`
	Tags                map[string]string `json:"tags,omitempty"`
	GlooFed             *GlooFed          `json:"gloo-fed,omitempty"`
}

// Common

type Rbac struct {
	Create bool `json:"create"`
}

// Gloo-ee

type GlooEeExtensions struct {
	ExtAuth   *ExtAuth   `json:"extAuth,omitempty"`
	RateLimit *RateLimit `json:"rateLimit,omitempty"`
	GlooRedis *GlooRedis `json:"glooRedis,omitempty"`
}

type GlooRedis struct {
	EnableAcl *bool `json:"enableAcl,omitempty" desc:"Whether to include the ACL policy on redis install. Defaults to true."`
}

type RateLimit struct {
	Enabled         bool                          `json:"enabled,omitempty" desc:"if true, deploy rate limit service (default true)"`
	Deployment      *RateLimitDeployment          `json:"deployment,omitempty"`
	Service         *RateLimitService             `json:"service,omitempty"`
	Upstream        *glooGen.KubeResourceOverride `json:"upstream,omitempty"`
	CustomRateLimit interface{}                   `json:"customRateLimit,omitempty"`
	BeforeAuth      bool                          `json:"beforeAuth,omitempty" desc:"If true, rate limiting checks occur before auth (default false)."`
	Affinity        []map[string]interface{}      `json:"affinity,omitempty" desc:"Affinity rules to be applied"`
	AntiAffinity    []map[string]interface{}      `json:"antiAffinity,omitempty" desc:"Anti-affinity rules to be applied"`
}

type DynamoDb struct {
	Region             string `json:"region" desc:"aws region to run DynamoDB requests in"`
	SecretName         string `json:"secretName,omitempty" desc:"name of the aws secret in gloo's installation namespace that has aws creds (if provided, uses DynamoDB to back rate-limiting service instead of Redis)"`
	RateLimitTableName string `json:"tableName" desc:"DynamoDB table name used to back rate limit service (default rate-limits)"`
	ConsistentReads    bool   `json:"consistentReads" desc:"if true, reads from DynamoDB will be strongly consistent (default false)"`
	BatchSize          uint8  `json:"batchSize" desc:"batch size for get requests to DynamoDB (max 100, default 100)"`
}

type RateLimitDeployment struct {
	Name                 string               `json:"name"`
	GlooAddress          string               `json:"glooAddress"`
	GlooPort             uint                 `json:"glooPort" desc:"Sets the port of the gloo xDS server in the ratelimit sidecar envoy bootstrap config"`
	DynamoDb             DynamoDb             `json:"dynamodb"`
	Image                *glooGen.Image       `json:"image,omitempty"`
	Stats                *glooGen.Stats       `json:"stats"`
	RunAsUser            float64              `json:"runAsUser" desc:"Explicitly set the user ID for the container to run as. Default is 10101"`
	FloatingUserId       bool                 `json:"floatingUserId" desc:"set to true to allow the cluster to dynamically assign a user ID"`
	ExtraRateLimitLabels map[string]string    `json:"extraRateLimitLabels,omitempty" desc:"Optional extra key-value pairs to add to the spec.template.metadata.labels data of the rateLimit deployment."`
	LogLevel             *string              `json:"logLevel,omitempty" desc:"Level at which the pod should log. Options include \"info\", \"debug\", \"warn\", \"error\", \"panic\" and \"fatal\". Default level is info"`
	PodDisruptionBudget  *PodDisruptionBudget `json:"podDisruptionBudget,omitempty" desc:"PodDisruptionBudget is an object to define the max disruption that can be caused to the rate-limit pods"`
	*glooGen.KubeResourceOverride
	*glooGen.DeploymentSpec
}

type RateLimitService struct {
	Port uint   `json:"port"`
	Name string `json:"name"`
	*glooGen.KubeResourceOverride
}

type GlooFed struct {
	Enabled             bool `json:"enabled,omitempty" desc:"If true, deploy federation service (default true)."`
	CreateLicenseSecret bool `json:"create_license_secret"`
}

type Redis struct {
	Deployment                *RedisDeployment `json:"deployment,omitempty"`
	Service                   RedisService     `json:"service,omitempty"`
	ClientSideShardingEnabled bool             `json:"clientSideShardingEnabled" desc:"If set to true, Envoy will be used as a Redis proxy and load balance requests between redis instances scaled via replicas. Default is false."`
	AclPrefix                 *string          `json:"aclPrefix,omitempty" desc:"The ACL policy for the default redis user. This is the prefix only, and if overridden, should end with < to signal the password."`
}

type RedisDeployment struct {
	Image                    *glooGen.Image    `json:"image,omitempty"`
	Name                     string            `json:"name"`
	StaticPort               uint              `json:"staticPort"`
	RunAsUser                float64           `json:"runAsUser" desc:"Explicitly set the user ID for the container to run as. Default is 999"`
	RunAsGroup               float64           `json:"runAsGroup" desc:"Explicitly set the group ID for the container to run as. Default is 999"`
	FsGroup                  float64           `json:"fsGroup" desc:"Explicitly set the fsGroup ID for the container to run as. Default is 999"`
	FloatingUserId           bool              `json:"floatingUserId" desc:"set to true to allow the cluster to dynamically assign a user ID"`
	ExtraRedisLabels         map[string]string `json:"extraRedisLabels,omitempty" desc:"Optional extra key-value pairs to add to the spec.template.metadata.labels data of the redis deployment."`
	EnablePodSecurityContext *bool             `json:"enablePodSecurityContext,omitempty" desc:"Whether or not to render the pod security context. Default is true"`
	*glooGen.DeploymentSpec
	*glooGen.KubeResourceOverride
}

type RedisService struct {
	Port uint   `json:"port"`
	Name string `json:"name"`
	*glooGen.KubeResourceOverride
}

type Observability struct {
	Enabled                   bool                          `json:"enabled,omitempty" desc:"if true, deploy observability service (default true)"`
	Deployment                *ObservabilityDeployment      `json:"deployment,omitempty"`
	CustomGrafana             *CustomGrafana                `json:"customGrafana" desc:"Configure a custom grafana deployment to work with Gloo observability, rather than the default Gloo grafana"`
	UpstreamDashboardTemplate string                        `json:"upstreamDashboardTemplate" desc:"Provide a custom dashboard template to use when generating per-upstream dashboards. The only variables available for use in this template are: {{.Uid}} and {{.EnvoyClusterName}}. Recommended to use Helm's --set-file to provide this value."`
	Rbac                      *glooGen.KubeResourceOverride `json:"rbac,omitempty"`
	ServiceAccount            *glooGen.KubeResourceOverride `json:"serviceAccount,omitempty"`
	ConfigMap                 *glooGen.KubeResourceOverride `json:"configMap,omitempty"`
	Secret                    *glooGen.KubeResourceOverride `json:"secret,omitempty"`
}

type ObservabilityDeployment struct {
	Image                    *glooGen.Image    `json:"image,omitempty"`
	Stats                    *glooGen.Stats    `json:"stats"`
	RunAsUser                float64           `json:"runAsUser" desc:"Explicitly set the user ID for the container to run as. Default is 10101"`
	FloatingUserId           bool              `json:"floatingUserId" desc:"set to true to allow the cluster to dynamically assign a user ID"`
	ExtraObservabilityLabels map[string]string `json:"extraObservabilityLabels,omitempty" desc:"Optional extra key-value pairs to add to the spec.template.metadata.labels data of the Observability deployment."`
	LogLevel                 *string           `json:"logLevel,omitempty" desc:"Level at which the pod should log. Options include \"info\", \"debug\", \"warn\", \"error\", \"panic\" and \"fatal\". Default level is info"`

	*glooGen.DeploymentSpec
	*glooGen.KubeResourceOverride
}

type CustomGrafana struct {
	Enabled  bool   `json:"enabled,omitempty" desc:"Set to true to indicate that the observability pod should talk to a custom grafana instance"`
	Username string `json:"username,omitempty" desc:"Set this and the 'password' field to authenticate to the custom grafana instance using basic auth"`
	Password string `json:"password,omitempty" desc:"Set this and the 'username' field to authenticate to the custom grafana instance using basic auth"`
	ApiKey   string `json:"apiKey,omitempty" desc:"Authenticate to the custom grafana instance using this api key"`
	Url      string `json:"url,omitempty" desc:"The URL for the custom grafana instance"`
	CaBundle string `json:"caBundle,omitempty" desc:"The Certificate Authority used to verify the server certificates."`
	*glooGen.KubeResourceOverride
}

type ExtAuth struct {
	Enabled              bool                          `json:"enabled,omitempty" desc:"if true, deploy ExtAuth service (default true)"`
	UserIdHeader         string                        `json:"userIdHeader,omitempty"`
	Deployment           *ExtAuthDeployment            `json:"deployment,omitempty"`
	Service              *ExtAuthService               `json:"service,omitempty"`
	SigningKey           *ExtAuthSigningKey            `json:"signingKey,omitempty"`
	TlsEnabled           bool                          `json:"tlsEnabled" desc:"if true, have extauth terminate TLS itself (whereas Gloo mTLS mode runs an Envoy and SDS sidecars to do TLS termination and cert rotation)"`
	CertPath             string                        `json:"certPath,omitempty" desc:"location of tls termination cert, if omitted defaults to /etc/envoy/ssl/tls.crt"`
	KeyPath              string                        `json:"keyPath,omitempty" desc:"location of tls termination key, if omitted defaults to /etc/envoy/ssl/tls.key"`
	Plugins              map[string]*ExtAuthPlugin     `json:"plugins,omitempty"`
	EnvoySidecar         bool                          `json:"envoySidecar" desc:"if true, deploy ExtAuth as a sidecar with envoy (defaults to false)"`
	StandaloneDeployment bool                          `json:"standaloneDeployment" desc:"if true, create a standalone ExtAuth deployment (defaults to true)"`
	TransportApiVersion  string                        `json:"transportApiVersion" desc:"Determines the API version for the ext_authz transport protocol that will be used by Envoy to communicate with the auth server. Defaults to 'V3''"`
	ServiceName          string                        `json:"serviceName,omitempty"`
	RequestTimeout       string                        `json:"requestTimeout,omitempty" desc:"Timeout for the ext auth service to respond (defaults to 200ms)"`
	HeadersToRedact      string                        `json:"headersToRedact,omitempty" desc:"Space separated list of headers to redact from the logs. To avoid the default redactions, specify '-' as the value"`
	Secret               *glooGen.KubeResourceOverride `json:"secret,omitempty"`
	Upstream             *glooGen.KubeResourceOverride `json:"upstream,omitempty"`
	RequestBody          *BufferSettings               `json:"requestBody,omitempty" desc:"Set in order to send the body of the request, and not just the headers"`
	Affinity             []map[string]interface{}      `json:"affinity,omitempty" desc:"Affinity rules to be applied. If unset, require extAuth pods to be scheduled on nodes with already-running gateway-proxy pods"`
	AntiAffinity         []map[string]interface{}      `json:"antiAffinity,omitempty" desc:"Anti-affinity rules to be applied"`
}

type ExtAuthDeployment struct {
	Name                string               `json:"name"`
	GlooAddress         string               `json:"glooAddress,omitempty"`
	GlooPort            uint                 `json:"glooPort" desc:"Sets the port of the gloo xDS server in the ratelimit sidecar envoy bootstrap config"`
	Port                uint                 `json:"port"`
	Image               *glooGen.Image       `json:"image,omitempty"`
	Stats               *glooGen.Stats       `json:"stats"`
	RunAsUser           float64              `json:"runAsUser" desc:"Explicitly set the user ID for the container to run as. Default is 10101"`
	FsGroup             float64              `json:"fsGroup" desc:"Explicitly set the group ID for volume ownership. Default is 10101"`
	FloatingUserId      bool                 `json:"floatingUserId" desc:"set to true to allow the cluster to dynamically assign a user ID"`
	ExtraExtAuthLabels  map[string]string    `json:"extraExtAuthLabels,omitempty" desc:"Optional extra key-value pairs to add to the spec.template.metadata.labels data of the ExtAuth deployment."`
	ExtraVolume         []v1.Volume          `json:"extraVolume,omitempty" desc:"custom defined yaml for allowing extra volume on the extauth container"`
	ExtraVolumeMount    []v1.VolumeMount     `json:"extraVolumeMount,omitempty" desc:"custom defined yaml for allowing extra volume mounts on the extauth container"`
	PodDisruptionBudget *PodDisruptionBudget `json:"podDisruptionBudget,omitempty" desc:"PodDisruptionBudget is an object to define the max disruption that can be caused to the ExtAuth pods"`
	*glooGen.DeploymentSpec
	*glooGen.KubeResourceOverride
}

type ExtAuthService struct {
	Port uint   `json:"port"`
	Name string `json:"name"`
	*glooGen.KubeResourceOverride
}

type ExtAuthSigningKey struct {
	Name       string `json:"name"`
	SigningKey string `json:"signing-key"`
}

type ExtAuthPlugin struct {
	Image *glooGen.Image `json:"image,omitempty"`
}

type BufferSettings struct {
	MaxRequestBytes     uint32 `json:"max_request_bytes,omitempty" desc:"Sets the maximum size of a message body that the filter will hold in memory, returning 413 and *not* initiating the authorization process when reaching the maximum (defaults to 4KB)"`
	AllowPartialMessage bool   `json:"allow_partial_message,omitempty" desc:"if true, Envoy will buffer the message until max_request_bytes is reached, dispatch the authorization request, and not return an error"`
	PackAsBytes         bool   `json:"pack_as_bytes,omitempty" desc:"if true, Envoy will send the body sent to the external authorization service with raw bytes"`
}

type OAuth struct {
	Server string `json:"server"`
	Client string `json:"client"`
}

type PodDisruptionBudget struct {
	minAvailable   int32 `json:"minAvailable,omitempty" desc:"An eviction is allowed if at least \"minAvailable\" pods selected by \"selector\" will still be available after the eviction, i.e. even in the absence of the evicted pod. So for example you can prevent all voluntary evictions by specifying \"100%\"."`
	maxUnavailable int32 `json:"maxUnavailable,omitempty" desc:"An eviction is allowed if at most \"maxUnavailable\" pods selected by \"selector\" are unavailable after the eviction, i.e. even in absence of the evicted pod. For example, one can prevent all voluntary evictions by specifying 0. This is a mutually exclusive setting with \"minAvailable\"."`
}
