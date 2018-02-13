package kubernetes

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/pkg/errors"

	"github.com/solo-io/glue/internal/bootstrap"
	"github.com/solo-io/glue/internal/plugins"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/endpointdiscovery"
	"github.com/solo-io/glue/pkg/plugin"
	"github.com/solo-io/glue/pkg/secretwatcher"
)

func init() {
	plugins.Register(&Plugin{}, startEndpointDiscovery)
}

func startEndpointDiscovery(opts bootstrap.Options, stopCh <-chan struct{}) (endpointdiscovery.Interface, error) {
	kubeConfig := opts.KubeOptions.KubeConfig
	masterUrl := opts.KubeOptions.MasterURL
	syncFrequency := opts.ConfigWatcherOptions.SyncFrequency
	return NewEndpointDiscovery(kubeConfig, masterUrl, syncFrequency, stopCh)
}

type Plugin struct{}

const (
	// define Upstream type name
	UpstreamTypeKube v1.UpstreamType = "kubernetes"
)

func (p *Plugin) GetDependencies(_ *v1.Config) *plugin.Dependencies {
	return nil
}

func (p *Plugin) ProcessUpstream(in *v1.Upstream, secrets secretwatcher.SecretMap, out *envoyapi.Cluster) error {
	if in.Type != UpstreamTypeKube {
		return nil
	}
	// decode does validation for us
	if _, err := DecodeUpstreamSpec(in.Spec); err != nil {
		return errors.Wrap(err, "invalid kubernetes upstream spec")
	}

	// just configure the cluster to use EDS:ADS and call it a day
	out.Type = envoyapi.Cluster_EDS
	out.EdsClusterConfig = &envoyapi.Cluster_EdsClusterConfig{
		EdsConfig: &envoycore.ConfigSource{
			ConfigSourceSpecifier: &envoycore.ConfigSource_Ads{
				Ads: &envoycore.AggregatedConfigSource{},
			},
		},
	}
	return nil
}
