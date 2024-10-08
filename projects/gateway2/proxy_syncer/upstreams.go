package proxy_syncer

import (
	"context"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/snapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
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
	upstreamTranslator := &upstreamTranslator{
		newTranslator,
		s,
	}
	return func(kctx krt.HandlerContext, us UpstreamWrapper) *EnvoyCluster {
		c, _ := upstreamTranslator.translateUpstream(ctx, us.Inner)

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

	c, errors := tr.ComputeCluster(params, us, isEds(us))
	if len(errors) > 0 {
		// TODO: ???
	}

	return c, nil
}

func isEds(us *gloov1.Upstream) bool {
	// TODO: either check if we have endpoints for that upstream, like the translator does today;
	// or have a per upstream type check
	panic("unimplemented")
}
