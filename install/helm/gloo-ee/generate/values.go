package generate

import "github.com/solo-io/gloo/install/helm/gloo/generate"

type HelmConfig struct {
	Config
	Global *generate.Global `json:"global,omitempty"`
}
type Config struct {
	Settings      *generate.Settings `json:"settings,omitempty"`
	LicenseKey    string             `json:"license_key,omitempty"`
	Gloo          *generate.Config   `json:"gloo,omitempty"`
	Redis         *Redis             `json:"redis,omitempty"`
	RateLimit     *RateLimit         `json:"rateLimit,omitempty"`
	ApiServer     *ApiServer         `json:"apiServer,omitempty"`
	Observability *Observability     `json:"observability,omitempty"`
	Rbac          *Rbac              `json:"rbac"`
	Grafana       interface{}        `json:"grafana,omitempty"`
	Prometheus    interface{}        `json:"prometheus,omitempty"`
	Tags          map[string]string  `json:"tags,omitempty"`
}

// Common

type OAuth struct {
	Server string `json:"server"`
	Client string `json:"client"`
}

type Rbac struct {
	Create bool `json:"create"`
}

// Gloo-ee

type GlooEeExtensions struct {
	ExtAuth *ExtAuth `json:"extAuth,omitempty"`
}

type RateLimit struct {
	Enabled    bool                 `json:"enabled"`
	Deployment *RateLimitDeployment `json:"deployment,omitempty"`
	Service    *RateLimitService    `json:"service,omitempty"`
}

type RateLimitDeployment struct {
	RedisUrl    string          `json:"redisUrl"`
	GlooAddress string          `json:"glooAddress"`
	Image       *generate.Image `json:"image,omitempty"`
	*generate.DeploymentSpec
}

type RateLimitService struct {
	Port uint   `json:"port"`
	Name string `json:"name"`
}

type Redis struct {
	Deployment *RedisDeployment `json:"deployment,omitempty"`
	Service    RedisService     `json:"service,omitempty"`
}

type RedisDeployment struct {
	Image      *generate.Image `json:"image,omitempty"`
	StaticPort uint            `json:"staticPort"`
	*generate.DeploymentSpec
}

type RedisService struct {
	Port uint   `json:"port"`
	Name string `json:"name"`
}

type ApiServer struct {
	NoAuth     bool                 `json:"noAuth"`
	Deployment *ApiServerDeployment `json:"deployment,omitempty"`
	Service    *ApiServerService    `json:"service,omitempty"`
	ConfigMap  *ApiServerConfigMap  `json:"configMap,omitempty"`
	EnableBeta bool                 `json:"enableBeta,omitempty"`
}

type ApiServerDeployment struct {
	Server *ApiServerServerDeployment `json:"server,omitempty"`
	Ui     *ApiServerUiDeployment     `json:"ui,omitempty"`
	Envoy  *ApiServerEnvoyDeployment  `json:"envoy,omitempty"`
	*generate.DeploymentSpec
}

type ApiServerServerDeployment struct {
	GrpcPort uint            `json:"grpcPort"`
	OAuth    *OAuth          `json:"oauth,omitempty"`
	Image    *generate.Image `json:"image"`
	*generate.DeploymentSpec
}

type ApiServerEnvoyDeployment struct {
	Image *generate.Image `json:"image"`
	*generate.DeploymentSpec
}

type ApiServerUiDeployment struct {
	StaticPort       uint            `json:"staticPort"`
	StaticPortNoAuth uint            `json:"staticPortNoAuth"`
	Image            *generate.Image `json:"image,omitempty"`
	*generate.DeploymentSpec
}

type ApiServerService struct {
	Name string `json:"name"`
}

type ApiServerConfigMap struct {
	Name string `json:"name"`
}

type Observability struct {
	Deployment *ObservabilityDeployment `json:"deployment,omitempty"`
}

type ObservabilityDeployment struct {
	Image *generate.Image `json:"image,omitempty"`
	*generate.DeploymentSpec
}

type ExtAuth struct {
	UserIdHeader string                    `json:"userIdHeader,omitempty"`
	Deployment   *ExtAuthDeployment        `json:"deployment,omitempty"`
	Service      *ExtAuthService           `json:"service,omitempty"`
	SigningKey   *ExtAuthSigningKey        `json:"signingKey,omitempty"`
	Plugins      map[string]*ExtAuthPlugin `json:"plugins,omitempty"`
	EnvoySidecar bool                      `json:"envoySidecar"`
}

type ExtAuthDeployment struct {
	Name        string          `json:"name"`
	GlooAddress string          `json:"glooAddress,omitempty"`
	Port        uint            `json:"port"`
	Image       *generate.Image `json:"image,omitempty"`
	Stats       bool            `json:"stats" desc:"enable prometheus stats"`
	*generate.DeploymentSpec
}

type ExtAuthService struct {
	Port uint   `json:"port"`
	Name string `json:"name"`
}

type ExtAuthSigningKey struct {
	Name       string `json:"name"`
	SigningKey string `json:"signing-key"`
}

type ExtAuthPlugin struct {
	Image *generate.Image `json:"image,omitempty"`
}
