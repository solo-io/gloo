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
	Observability *Observability     `json:"observability,omitempty"`
	Rbac          *Rbac              `json:"rbac"`
	Grafana       interface{}        `json:"grafana,omitempty"`
	Prometheus    interface{}        `json:"prometheus,omitempty"`
	Tags          map[string]string  `json:"tags,omitempty"`
}

// Common

type Rbac struct {
	Create bool `json:"create"`
}

// Gloo-ee

type GlooEeExtensions struct {
	ExtAuth *ExtAuth `json:"extAuth,omitempty"`
}

type RateLimit struct {
	Enabled    bool                 `json:"enabled" desc:"if true, deploy rate limit service (default true)"`
	Deployment *RateLimitDeployment `json:"deployment,omitempty"`
	Service    *RateLimitService    `json:"service,omitempty"`
}

type DynamoDb struct {
	Region             string `json:"region" desc:"aws region to run DynamoDB requests in"`
	SecretName         string `json:"secretName,omitempty" desc:"name of the aws secret in gloo's installation namespace that has aws creds (if provided, uses DynamoDB to back rate-limiting service instead of Redis)"`
	RateLimitTableName string `json:"tableName" desc:"DynamoDB table name used to back rate limit service (default rate-limits)"`
	ConsistentReads    bool   `json:"consistentReads" desc:"if true, reads from DynamoDB will be strongly consistent (default false)"`
	BatchSize          uint8  `json:"batchSize" desc:"batch size for get requests to DynamoDB (max 100, default 100)"`
}

type RateLimitDeployment struct {
	RedisUrl    string          `json:"redisUrl"`
	GlooAddress string          `json:"glooAddress"`
	DynamoDb    DynamoDb        `json:"dynamodb"`
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
	ServiceName  string                    `json:"serviceName,omitempty"`
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
