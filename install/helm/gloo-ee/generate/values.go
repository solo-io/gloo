package generate

import "github.com/solo-io/gloo/install/helm/gloo/generate"

type Config struct {
	Gloo          *generate.Config `json:"gloo,omitempty"`
	Redis         *Redis           `json:"redis,omitempty"`
	RateLimit     *RateLimit       `json:"rateLimit,omitempty"`
	ApiServer     *ApiServer       `json:"apiServer,omitempty"`
	Licensing     *Licensing       `json:"licensing,omitempty"`
	Observability *Observability   `json:"observability,omitempty"`
	Rbac          *Rbac            `json:"rbac"`
	Grafana       interface{}      `json:"grafana,omitempty"`
	Prometheus    interface{}      `json:"prometheus,omitempty"`
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

type RateLimit struct {
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
	Port string `json:"port"`
	Name string `json:"name"`
}

type Redis struct {
	Deployment *RedisDeployment `json:"deployment,omitempty"`
	Service    RedisService     `json:"service,omitempty"`
}

type RedisDeployment struct {
	Image      *generate.Image `json:"image,omitempty"`
	StaticPort string          `json:"staticPort"`
	*generate.DeploymentSpec
}

type RedisService struct {
	Port string `json:"port"`
	Name string `json:"name"`
}

type ApiServer struct {
	NoAuth     bool                 `json:"noAuth"`
	Deployment *ApiServerDeployment `json:"deployment,omitempty"`
	Service    *ApiServerService    `json:"service,omitempty"`
	ConfigMap  *ApiServerConfigMap  `json:"configMap,omitempty"`
}

type ApiServerDeployment struct {
	Server *ApiServerServerDeployment `json:"server,omitempty"`
	Ui     *ApiServerUiDeployment     `json:"ui,omitempty"`
}

type ApiServerServerDeployment struct {
	GraphqlPort string          `json:"graphqlPort"`
	OAuth       *OAuth          `json:"oauth,omitempty"`
	Image       *generate.Image `json:"image"`
	*generate.DeploymentSpec
}

type ApiServerUiDeployment struct {
	StaticPort string          `json:"staticPort"`
	Image      *generate.Image `json:"image,omitempty"`
	*generate.DeploymentSpec
}

type ApiServerService struct {
	Name string `json:"name"`
}

type ApiServerConfigMap struct {
	Name string `json:"name"`
}

type Licensing struct {
	Deployment *LicensingDeployment `json:"deployment,omitempty"`
	Service    *LicensingService    `json:"service,omitempty"`
}

type LicensingDeployment struct {
	Image      *generate.Image `json:"image,omitempty"`
	StaticPort string          `json:"staticPort"`
	*generate.DeploymentSpec
}

type LicensingService struct {
	Port string `json:"port"`
	Name string `json:"name"`
}

type Observability struct {
	Deployment *ObservabilityDeployment `json:"deployment,omitempty"`
}

type ObservabilityDeployment struct {
	Image *generate.Image `json:"image,omitempty"`
	*generate.DeploymentSpec
}
