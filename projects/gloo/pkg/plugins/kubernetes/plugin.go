package kubernetes

import (
	"fmt"
	"net/url"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	errors "github.com/rotisserie/eris"
	"go.uber.org/zap"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"knative.dev/pkg/network"

	"github.com/solo-io/go-utils/contextutils"
	corecache "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
)

const (
	ExtensionName = "kubernetes"
)

type plugin struct {
	UpstreamConverter

	kube          kubernetes.Interface
	kubeCoreCache corecache.KubeCoreCache
	svcCollection krt.Collection[*corev1.Service]
	settings      *v1.Settings
}

func NewPlugin(kube kubernetes.Interface, kubeCoreCache corecache.KubeCoreCache, svcCollection krt.Collection[*corev1.Service]) plugins.Plugin {
	return &plugin{
		kube:              kube,
		kubeCoreCache:     kubeCoreCache,
		svcCollection:     svcCollection,
		UpstreamConverter: DefaultUpstreamConverter(),
	}
}

func (p *plugin) Name() string { return ExtensionName }

func (p *plugin) Init(params plugins.InitParams) {
	p.settings = params.Settings
}

// Resolve returns the URL of the upstream. Where is this used? Remove?
func (p *plugin) Resolve(u *v1.Upstream) (*url.URL, error) {
	kubeSpec, ok := u.GetUpstreamType().(*v1.Upstream_Kube)
	if !ok {
		return nil, nil
	}
	return url.Parse(fmt.Sprintf("tcp://%v.%v.svc.%s:%v",
		kubeSpec.Kube.GetServiceName(),
		kubeSpec.Kube.GetServiceNamespace(),
		network.GetClusterDomainName(),
		kubeSpec.Kube.GetServicePort(),
	))
}

var _ plugins.UpstreamPlugin = &plugin{}

func (p *plugin) ProcessUpstream(params plugins.Params, in *v1.Upstream, out *clusterv3.Cluster) error {
	kube, ok := in.GetUpstreamType().(*v1.Upstream_Kube)
	if !ok {
		// don't process any non-kube upstreams.
		return nil
	}

	// configure the cluster to use EDS:ADS and call it a day. huh?
	xds.SetEdsOnCluster(out, p.settings)
	upstreamRef := in.GetMetadata().Ref()

	// if we are in ggv2 / krt mode, we won't have the kubeCoreCache set.
	// instead we will use the krt collection to fetch the service.
	// in a future PR plugins will have access to krt context, so they can use fetch.
	if p.svcCollection != nil {
		// TODO: change this to fetch once we have krt context in plugins in a follow-up
		if p.svcCollection.GetKey(krt.Named{
			Name:      kube.Kube.GetServiceName(),
			Namespace: kube.Kube.GetServiceNamespace(),
		}.ResourceName()) != nil {
			return nil
		}
	} else {

		lister := p.kubeCoreCache.NamespacedServiceLister(kube.Kube.GetServiceNamespace())
		if lister == nil {
			return errors.Errorf("Upstream %s references the service '%s' which has an invalid ServiceNamespace '%s'.",
				upstreamRef.String(),
				kube.Kube.GetServiceName(),
				kube.Kube.GetServiceNamespace(),
			)
		}

		svcs, err := lister.List(labels.NewSelector())
		if err != nil {
			return err
		}

		for _, s := range svcs {
			if s.Name == kube.Kube.GetServiceName() {
				return nil
			}
		}
	}

	logger := contextutils.LoggerFrom(params.Ctx)
	logger.Debug("service does not exist", zap.String("upstream", upstreamRef.String()), zap.String("service", kube.Kube.GetServiceName()), zap.String("namespace", kube.Kube.GetServiceNamespace()))

	return errors.Errorf("Upstream %s references the service '%s' which does not exist in namespace '%s'",
		upstreamRef.String(),
		kube.Kube.GetServiceName(),
		kube.Kube.GetServiceNamespace(),
	)
}
