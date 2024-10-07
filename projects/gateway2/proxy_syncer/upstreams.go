package proxy_syncer

import (
	"context"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/snapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube/krt"
)

type EnvoyCluster struct {
	*envoy_config_cluster_v3.Cluster
}

func (c EnvoyCluster) ResourceName() string {
	return c.Cluster.Name
}
func (c EnvoyCluster) Equals(in EnvoyCluster) bool {
	return proto.Equal(c.Cluster, in.Cluster)
}

func NewGlooK8sUpstreams(ctx context.Context, s *v1.Settings, secrets krt.Collection[ResourceWrapper[*gloov1.Secret]], finalUpstreams krt.Collection[UpstreamWrapper], newTranslator func() translator.Translator) krt.Collection[EnvoyCluster] {

	return krt.NewCollection(finalUpstreams, TransformUpstreamBuilder(ctx, s, secrets, newTranslator), krt.WithName("EnvoyCluster"))
}

func TransformUpstreamBuilder(ctx context.Context, s *v1.Settings, secrets krt.Collection[ResourceWrapper[*gloov1.Secret]], newTranslator func() translator.Translator) func(kctx krt.HandlerContext, us UpstreamWrapper) *EnvoyCluster {
	return func(kctx krt.HandlerContext, us UpstreamWrapper) *EnvoyCluster {
		snap := snapshot.Snapshot{}
		snap.Upstreams = snapshot.SliceCollection[*gloov1.Upstream]([]*gloov1.Upstream{us.Inner})
		snap.Secrets = KrtWrappedCollection[*gloov1.Secret]{
			C:    secrets,
			Kctx: kctx,
		}

		params := plugins.Params{
			Ctx:      ctx,
			Settings: s,
			Snapshot: &snap,
			Messages: map[*core.ResourceRef][]string{},
		}
		t := newTranslator()
		proxy := &v1.Proxy{
			Metadata: &core.Metadata{
				Name:      "temp-proxy",
				Namespace: "temp-ns",
			},
		}

		xdsSnapshot, reports, _ := t.Translate(params, proxy)

		if report, ok := reports[us.Inner]; ok {

			// TODO: report error
			if report.Errors != nil {
				// DO SOMETHING
			}

		}
		c := getCluster(xdsSnapshot.(*xds.EnvoySnapshot).Clusters)

		return &EnvoyCluster{Cluster: c}
	}
}

type upstreamTranslator struct {
	newTranslator func() translator.Translator
	s             *v1.Settings
}

func (t *upstreamTranslator) translateUpstream(ctx context.Context, us *gloov1.Upstream) (*envoy_config_cluster_v3.Cluster, error) {
	snap := snapshot.Snapshot{}
	snap.Upstreams = snapshot.SliceCollection[*gloov1.Upstream]([]*gloov1.Upstream{us})
	params := plugins.Params{
		Ctx:      ctx,
		Settings: t.s,
		Snapshot: &snap,
		Messages: map[*core.ResourceRef][]string{},
	}
	tr := t.newTranslator()
	proxy := &v1.Proxy{
		Metadata: &core.Metadata{
			Name:      "temp-proxy",
			Namespace: "temp-ns",
		},
	}

	xdsSnapshot, reports, _ := tr.Translate(params, proxy)

	if report, ok := reports[us]; ok {

		// TODO: report error
		if report.Errors != nil {
			// DO SOMETHING
		}

	}
	c := getCluster(xdsSnapshot.(*xds.EnvoySnapshot).Clusters)
	return c, nil
}

func getCluster(r envoycache.Resources) *envoy_config_cluster_v3.Cluster {
	if len(r.Items) > 1 {
		// TODO: log error.. this should never happen.
	}
	for _, cluster := range r.Items {
		if c, ok := cluster.ResourceProto().(*envoy_config_cluster_v3.Cluster); ok {
			return c
		}
	}
	return nil
}
