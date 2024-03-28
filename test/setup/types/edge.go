package types

import (
	"fmt"
	"os"
	"strings"

	"github.com/solo-io/gloo/test/setup/defaults"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type (
	// This follows the values schema for the Gloo Edge Helm chart
	GlooEdge struct {
		Gateway2  *Gateway2  `json:"gateway2,omitempty" desc:"Gloo Gateway settings, subject to change."`
		Settings  *Settings  `json:"settings,omitempty"`
		Gloo      *Gloo      `json:"gloo,omitempty"`
		Discovery *Discovery `json:"discovery,omitempty"`
		Gateway   *Gateway   `json:"gateway,omitempty"`
		Global    *Global    `json:"global,omitempty"`
		// TODO(npolshak): Add accessLogger, ingress, ingressProxy, GatewayProxies
	}

	Gateway2 struct {
		ControlPlane *ControlPlane `json:"controlPlane,omitempty"`
	}

	ControlPlane struct {
		Enabled bool `json:"enabled,omitempty"`
	}

	Settings struct {
		WatchNamespaces               []string `json:"watchNamespaces,omitempty" desc:"whitelist of namespaces for Gloo Edge to watch for services and CRDs. Empty list means all namespaces"`
		WriteNamespace                string   `json:"writeNamespace,omitempty" desc:"namespace where intermediary CRDs will be written to, e.g. Upstreams written by Gloo Edge Discovery."`
		Create                        bool     `json:"create,omitempty" desc:"create a Settings CRD which provides bootstrap configuration to Gloo Edge controllers"`
		SingleNamespace               bool     `json:"singleNamespace,omitempty" desc:"Enable to use install namespace as WatchNamespace and WriteNamespace"`
		Linkerd                       bool     `json:"linkerd,omitempty" desc:"Enable automatic Linkerd integration in Gloo Edge"`
		DisableProxyGarbageCollection bool     `json:"disableProxyGarbageCollection,omitempty" desc:"Set this option to determine the state of an Envoy listener when the corresponding Proxy resource has no routes. If false (default), Gloo Edge will propagate the state of the Proxy to Envoy, resetting the listener to a clean slate with no routes. If true, Gloo Edge will keep serving the routes from the last applied valid configuration."`
		RegexMaxProgramSize           uint32   `json:"regexMaxProgramSize,omitempty" desc:"Set this field to specify the RE2 default max program size which is a rough estimate of how complex the compiled regex is to evaluate. If not specified, this defaults to 1024."`
		DisableKubernetesDestinations bool     `json:"disableKubernetesDestinations,omitempty" desc:"Enable or disable Gloo Edge to scan Kubernetes services in the cluster and create in-memory Upstream resources to represent them. These resources enable Gloo Edge to route requests to a Kubernetes service. Note that if you have a large number of services in your cluster and you do not restrict the namespaces that Gloo Edge watches, the API snapshot increases which can have a negative impact on the Gloo Edge translation time. In addition, load balancing is done in kube-proxy which can have further performance impacts. Using Gloo Upstreams as a routing destination bypasses kube-proxy as the request is routed to the pod directly. Alternatively, you can use [Kubernetes](https://docs.solo.io/gloo-edge/latest/reference/api/github.com/solo-io/gloo/projects/gloo/api/v1/options/kubernetes/kubernetes.proto.sk/) Upstream resources as a routing destination to forward requests to the pod directly. For more information, see the [docs](https://docs.solo.io/gloo-edge/latest/guides/traffic_management/destination_types/kubernetes_services/)."`
		EnableRestEds                 bool     `json:"enableRestEds,omitempty" desc:"Whether or not to use rest xds for all EDS by default. Defaults to false."`
		DevMode                       bool     `json:"devMode,omitempty" desc:"Whether or not to enable dev mode. Defaults to false. Setting to true at install time will expose the gloo dev admin endpoint on port 10010. Not recommended for production."`
	}

	Gloo struct {
		LogLevel string `json:"logLevel,omitempty" desc:"Level at which the pod should log. Options include \"info\", \"debug\", \"warn\", \"error\", \"panic\" and \"fatal\". Default level is info"`
	}

	Discovery struct {
		Enabled bool   `json:"enabled,omitempty" desc:"enable Discovery features"`
		Image   *Image `json:"image,omitempty"`
	}

	Gateway struct {
		Enabled    bool        `json:"enabled,omitempty" desc:"enable Gloo Edge API Gateway features"`
		Image      *Image      `json:"image,omitempty"`
		CertGenJob *CertGenJob `json:"certGenJob,omitempty" desc:"generate self-signed certs with this job to be used with the gateway validation webhook. this job will only run if validation is enabled for the gateway"`
		RolloutJob *RolloutJob `json:"rolloutJob,omitempty" desc:"This job waits for the 'gloo' deployment to successfully roll out (if the validation webhook is enabled), and then applies the Gloo Edge custom resources."`
		CleanupJob *CleanupJob `json:"cleanupJob,omitempty" desc:"This job cleans up resources that are not deleted by Helm when Gloo Edge is uninstalled."`
	}

	CertGenJob struct {
		Enabled bool   `json:"enabled,omitempty" desc:"enable the job that generates the certificates for the validating webhook at install time (default true)"`
		Image   *Image `json:"image,omitempty"`
	}

	RolloutJob struct {
		Enabled bool   `json:"enabled,omitempty" desc:"Enable the job that applies default Gloo Edge custom resources at install and upgrade time (default true)."`
		Image   *Image `json:"image,omitempty"`
	}

	CleanupJob struct {
		Enabled bool   `json:"enabled,omitempty" desc:"Enable the job that removes Gloo Edge custom resources when Gloo Edge is uninstalled (default true)."`
		Image   *Image `json:"image,omitempty"`
	}

	Global struct {
		Image    *Image    `json:"image,omitempty"`
		IstioSDS *IstioSDS `json:"istioSDS,omitempty" desc:"Config used for installing Gloo Edge with Istio SDS cert rotation features to facilitate Istio mTLS"`
		GlooMtls *Mtls     `json:"glooMtls,omitempty" desc:"Config used to enable internal mtls authentication"`
	}

	IstioSDS struct {
		Enabled bool `json:"enabled,omitempty" desc:"Enables SDS cert-rotator sidecar for istio mTLS cert rotation"`
	}

	Mtls struct {
		Enabled      bool                   `json:"enabled,omitempty" desc:"Enables internal mtls authentication"`
		Sds          *SdsContainer          `json:"sds,omitempty"`
		EnvoySidecar *EnvoySidecarContainer `json:"envoy,omitempty"`
		IstioProxy   *IstioProxyContainer   `json:"istioProxy,omitempty" desc:"Istio-proxy container"`
	}

	SdsContainer struct {
		Image *Image `json:"image,omitempty"`
	}

	EnvoySidecarContainer struct {
		Image *Image `json:"image,omitempty"`
	}

	IstioProxyContainer struct {
		Image *Image `json:"image,omitempty" desc:"Istio-proxy image to use for mTLS"`
	}

	Image struct {
		Tag        string `yaml:"tag,omitempty" json:"tag,omitempty"`
		Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
		Registry   string `yaml:"registry,omitempty" json:"registry,omitempty"`
		Load       bool   `yaml:"load,omitempty" json:"load,omitempty"`
	}

	App struct {
		ConfigMaps     []*corev1.ConfigMap    `yaml:"configMaps,omitempty" json:"configMaps,omitempty"`
		Service        *corev1.Service        `yaml:"service,omitempty" json:"service,omitempty"`
		ServiceAccount *corev1.ServiceAccount `yaml:"serviceAccount,omitempty" json:"serviceAccount,omitempty"`
		Deployment     *appsv1.Deployment     `yaml:"deployment,omitempty" json:"deployment,omitempty"`
		Versions       []string               `yaml:"versions,omitempty" json:"versions,omitempty"`
	}

	RepositoryImage struct {
		Tag        string `yaml:"tag,omitempty" json:"tag,omitempty"`
		Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
		Load       bool   `yaml:"load,omitempty" json:"load,omitempty"`
	}
)

func (p GlooEdge) Images() (imageRefs []string) {

	defaultImage := &Image{}
	if p.Global != nil {
		if p.Global.Image != nil {
			// load override
			defaultImage = p.Global.Image
		}

		glooMtlsEnabled := p.Global.GlooMtls != nil && p.Global.GlooMtls.Enabled
		istioSDSEnabled := p.Global.IstioSDS != nil && p.Global.IstioSDS.Enabled

		// Check if either GlooMtls or IstioSDS is enabled.
		// Note: Image overrides are only defined on GlooMtls.
		if glooMtlsEnabled || istioSDSEnabled {
			// Add Sds image IstioProxy override if GlooMtls is enabled
			if p.Global.GlooMtls != nil && p.Global.GlooMtls.IstioProxy != nil {
				imageRefs = append(imageRefs, p.Global.GlooMtls.Sds.Image.Ref(defaults.DefaultSdsImageName))
			} else {
				// Add default global image if GlooMtls image overrides are not set
				imageRefs = append(imageRefs, defaultImage.Ref(defaults.DefaultSdsImageName))
			}

			// Add EnvoySidecar image override if GlooMtls is enabled
			if p.Global.GlooMtls != nil && p.Global.GlooMtls.EnvoySidecar != nil {
				imageRefs = append(imageRefs, p.Global.GlooMtls.EnvoySidecar.Image.Ref(defaults.DefaultGatewayImageName))
			} else {
				// Add default global image if GlooMtls image overrides are not set
				imageRefs = append(imageRefs, defaultImage.Ref(defaults.DefaultGatewayImageName))
			}

			// Add IstioProxy image override if IstioSDS is enabled
			if istioSDSEnabled {
				if p.Global.GlooMtls != nil && p.Global.GlooMtls.IstioProxy != nil {
					imageRefs = append(imageRefs, p.Global.GlooMtls.IstioProxy.Image.Ref(defaults.DefaultIstioProxyImageName))
				} else {
					// Add default image
					imageRefs = append(imageRefs, fmt.Sprintf("%s/%s:%s", defaults.DefaultIstioImageRegistry, defaults.DefaultIstioProxyImageName, defaults.DefaultIstioTag))
				}
			}
		} else if istioSDSEnabled {
			// Add default global image if IstioSDS is enabled, but GlooMtls image overrides are not set
			imageRefs = append(imageRefs, defaultImage.Ref(defaults.DefaultSdsImageName))
			imageRefs = append(imageRefs, fmt.Sprintf("%s/%s:%s", defaults.DefaultIstioImageRegistry, defaults.DefaultIstioProxyImageName, defaults.DefaultIstioTag))
		}

	}

	// gloo image is configured via Global.Image
	imageRefs = append(imageRefs, defaultImage.Ref(defaults.DefaultGlooImageName))

	if p.Gateway != nil && p.Gateway.Enabled {
		// use overrides
		imageRefs = append(imageRefs, p.Gateway.Image.Ref(defaults.DefaultGatewayImageName))
		if p.Gateway.CertGenJob != nil && p.Gateway.CertGenJob.Enabled {
			imageRefs = append(imageRefs, p.Gateway.CertGenJob.Image.Ref(defaults.DefaultCertGenJobImageName))
		} else {
			imageRefs = append(imageRefs, defaultImage.Ref(defaults.DefaultCertGenJobImageName))
		}
		if p.Gateway.RolloutJob != nil && p.Gateway.RolloutJob.Enabled {
			imageRefs = append(imageRefs, p.Gateway.RolloutJob.Image.Ref(defaults.DefaultRolloutJobImageName))
		} else {
			imageRefs = append(imageRefs, defaultImage.Ref(defaults.DefaultRolloutJobImageName))
		}
		if p.Gateway.CleanupJob != nil && p.Gateway.CleanupJob.Enabled {
			imageRefs = append(imageRefs, p.Gateway.CleanupJob.Image.Ref(defaults.DefaultCleanupJobImageName))
		} else {
			imageRefs = append(imageRefs, defaultImage.Ref(defaults.DefaultCleanupJobImageName))
		}
	} else {
		// use default global image
		imageRefs = append(imageRefs, defaultImage.Ref(defaults.DefaultGatewayImageName))
		imageRefs = append(imageRefs, defaultImage.Ref(defaults.DefaultCertGenJobImageName))
		imageRefs = append(imageRefs, defaultImage.Ref(defaults.DefaultRolloutJobImageName))
		imageRefs = append(imageRefs, defaultImage.Ref(defaults.DefaultCleanupJobImageName))
	}

	if p.Gateway2 != nil && p.Gateway2.ControlPlane != nil && p.Gateway2.ControlPlane.Enabled {
		// Note: Gateway2 uses same image as Gateway
		imageRefs = append(imageRefs, p.Gateway.Image.Ref(defaults.DefaultGatewayImageName))
	} else {
		// use default global image
		imageRefs = append(imageRefs, defaultImage.Ref(defaults.DefaultGatewayImageName))
	}

	if p.Discovery != nil && p.Discovery.Enabled {
		imageRefs = append(imageRefs, p.Discovery.Image.Ref(defaults.DefaultDiscoveryImageName))
	} else {
		// use default global image
		imageRefs = append(imageRefs, defaultImage.Ref(defaults.DefaultDiscoveryImageName))
	}

	return imageRefs
}

func (i Image) Ref(component string) string {
	return fmt.Sprintf(
		"%s/%s:%s",
		withDefault(i.Registry, os.Getenv("IMAGE_REGISTRY")),
		withDefault(i.Repository, component),
		withDefault(i.Tag, os.Getenv("VERSION")),
	)
}

func (i RepositoryImage) Ref(component string) string {
	return fmt.Sprintf(
		"%s:%s",
		withDefault(i.Repository, strings.Join([]string{os.Getenv("REGISTRY"), component}, "/")),
		withDefault(i.Tag, os.Getenv("VERSION")),
	)
}

func withDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
