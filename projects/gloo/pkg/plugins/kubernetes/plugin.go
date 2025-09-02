package kubernetes

import (
	"context"
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

	clusterDomain := network.GetClusterDomainName()
	resolvedURL := fmt.Sprintf("tcp://%v.%v.svc.%s:%v",
		kubeSpec.Kube.GetServiceName(),
		kubeSpec.Kube.GetServiceNamespace(),
		clusterDomain,
		kubeSpec.Kube.GetServicePort(),
	)

	// Use background context since we don't have access to params.Ctx here
	contextutils.LoggerFrom(context.Background()).Infow("Resolving Kubernetes upstream URL",
		"issue", "8539",
		"upstream_name", u.GetMetadata().GetName(),
		"service_name", kubeSpec.Kube.GetServiceName(),
		"service_namespace", kubeSpec.Kube.GetServiceNamespace(),
		"service_port", kubeSpec.Kube.GetServicePort(),
		"cluster_domain", clusterDomain,
		"resolved_url", resolvedURL)

	return url.Parse(resolvedURL)
}

var _ plugins.UpstreamPlugin = &plugin{}

func (p *plugin) ProcessUpstream(params plugins.Params, in *v1.Upstream, out *clusterv3.Cluster) error {
	logger := contextutils.LoggerFrom(params.Ctx)

	kube, ok := in.GetUpstreamType().(*v1.Upstream_Kube)
	if !ok {
		// don't process any non-kube upstreams.
		logger.Infow("Skipping non-kubernetes upstream",
			"issue", "8539",
			"upstream_name", in.GetMetadata().GetName(),
			"upstream_type", fmt.Sprintf("%T", in.GetUpstreamType()))
		return nil
	}

	logger.Infow("Processing Kubernetes upstream",
		"issue", "8539",
		"upstream_name", in.GetMetadata().GetName(),
		"upstream_namespace", in.GetMetadata().GetNamespace(),
		"service_name", kube.Kube.GetServiceName(),
		"service_namespace", kube.Kube.GetServiceNamespace(),
		"service_port", kube.Kube.GetServicePort())

	// configure the cluster to use EDS:ADS and call it a day. huh?
	xds.SetEdsOnCluster(out, p.settings)
	upstreamRef := in.GetMetadata().Ref()

	// if we are in ggv2 / krt mode, we won't have the kubeCoreCache set.
	// instead we will use the krt collection to fetch the service.
	// in a future PR plugins will have access to krt context, so they can use fetch.
	if p.svcCollection != nil {
		logger.Infow("Using KRT collection for service lookup",
			"issue", "8539",
			"service_name", kube.Kube.GetServiceName(),
			"service_namespace", kube.Kube.GetServiceNamespace())

		// TODO: change this to fetch once we have krt context in plugins in a follow-up
		if p.svcCollection.GetKey(krt.Key[*corev1.Service](krt.Named{
			Name:      kube.Kube.GetServiceName(),
			Namespace: kube.Kube.GetServiceNamespace(),
		}.ResourceName())) != nil {
			logger.Infow("Service found in KRT collection",
				"issue", "8539",
				"service_name", kube.Kube.GetServiceName(),
				"service_namespace", kube.Kube.GetServiceNamespace())
			return nil
		}

	} else {
		logger.Infow("Using KubeCore cache for service lookup",
			"issue", "8539",
			"service_name", kube.Kube.GetServiceName(),
			"service_namespace", kube.Kube.GetServiceNamespace())

		lister := p.kubeCoreCache.NamespacedServiceLister(kube.Kube.GetServiceNamespace())
		if lister == nil {
			logger.Infow("Invalid service namespace for upstream",
				"issue", "8539",
				"upstream_ref", upstreamRef.String(),
				"service_name", kube.Kube.GetServiceName(),
				"service_namespace", kube.Kube.GetServiceNamespace())
			return errors.Errorf("Upstream %s references the service '%s' which has an invalid ServiceNamespace '%s'.",
				upstreamRef.String(),
				kube.Kube.GetServiceName(),
				kube.Kube.GetServiceNamespace(),
			)
		}

		svcs, err := lister.List(labels.NewSelector())
		if err != nil {
			logger.Infow("Error listing services from cache",
				"issue", "8539",
				"service_namespace", kube.Kube.GetServiceNamespace(),
				"error", err.Error())
			return err
		}

		logger.Infow("Searching for service in cache",
			"issue", "8539",
			"service_name", kube.Kube.GetServiceName(),
			"services_found", len(svcs))

		for _, s := range svcs {
			if s.Name == kube.Kube.GetServiceName() {
				logger.Infow("Service found in cache",
					"issue", "8539",
					"service_name", kube.Kube.GetServiceName(),
					"service_namespace", kube.Kube.GetServiceNamespace())
				return nil
			}
		}
	}

	logger.Debug("service does not exist", zap.String("upstream", upstreamRef.String()), zap.String("service", kube.Kube.GetServiceName()), zap.String("namespace", kube.Kube.GetServiceNamespace()))
	logger.Infow("Service not found, returning error",
		"issue", "8539",
		"upstream_ref", upstreamRef.String(),
		"service_name", kube.Kube.GetServiceName(),
		"service_namespace", kube.Kube.GetServiceNamespace())

	return errors.Errorf("Upstream %s references the service '%s' which does not exist in namespace '%s'",
		upstreamRef.String(),
		kube.Kube.GetServiceName(),
		kube.Kube.GetServiceNamespace(),
	)
}
