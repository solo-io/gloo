package translator

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/control-plane/translator/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
)

func (t *translator) computeClusters(snap *v1.Snapshot, resourceErrs reporter.ResourceErrors) []*envoyapi.Cluster {
	var (
		clusters []*envoyapi.Cluster
	)
	for _, upstream := range snap.UpstreamList {
		cluster := t.computeCluster(snap, upstream, resourceErrs)
		clusters = append(clusters, cluster)
	}
	return clusters
}

func (t *translator) computeCluster(snap *v1.Snapshot, upstream *v1.Upstream, resourceErrs reporter.ResourceErrors) *envoyapi.Cluster {
	out := initializeCluster(upstream, snap.EndpointList)

	params := plugins.Params{
		Secrets:   snap.SecretList,
		Artifacts: snap.ArtifactList,
	}

	for _, plug := range t.plugins {
		upstreamPlugin, ok := plug.(plugins.UpstreamPlugin)
		if !ok {
			continue
		}

		if err := upstreamPlugin.ProcessUpstream(params, upstream, out); err != nil {
			resourceErrs.AddError(upstream, err)
		}
	}
	if err := validateCluster(out); err != nil {
		resourceErrs.AddError(upstream, err)
	}
	return out
}

func initializeCluster(upstream *v1.Upstream, endpoints []*v1.Endpoint) *envoyapi.Cluster {
	out := &envoyapi.Cluster{
		Name:     upstream.Metadata.Name,
		Metadata: new(envoycore.Metadata),
	}
	// set Type = EDS if we have endpoints for the upstream
	if len(endpointsForUpstream(upstream, endpoints)) > 0 {
		out.Type = envoyapi.Cluster_EDS
	}
	// this field can be overridden by plugins
	out.ConnectTimeout = defaults.ClusterConnectionTimeout
	return out
}

// TODO: add more validation here
func validateCluster(c *envoyapi.Cluster) error {
	if c.Type == envoyapi.Cluster_STATIC || c.Type == envoyapi.Cluster_STRICT_DNS || c.Type == envoyapi.Cluster_LOGICAL_DNS {
		if (len(c.Hosts) == 0) && (c.LoadAssignment == nil || len(c.LoadAssignment.Endpoints) == 0) {
			return errors.Errorf("cluster type %v specified but hosts were empty", c.Type.String())
		}
	}
	return nil
}
