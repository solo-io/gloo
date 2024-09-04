package serviceentry

import (
	"context"
	"fmt"
	"slices"
	"time"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams/NoOpUpstreamClient"
	"github.com/solo-io/go-utils/contextutils"
	skclients "github.com/solo-io/solo-kit/pkg/api/v1/clients"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"google.golang.org/protobuf/types/known/wrapperspb"

	networking "istio.io/api/networking/v1beta1"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1beta1"
	istioclient "istio.io/client-go/pkg/clientset/versioned"
)

const UpstreamNamePrefix = "istio-se:"

var _ upstreams.ClientPlugin = &sePlugin{}

func (s *sePlugin) Client() v1.UpstreamClient {
	return NewServiceEntryUpstreamClient(s.client.Istio())
}

func (s *sePlugin) SourceName() string {
	return PluginName
}

func NewServiceEntryUpstreamClient(client istioclient.Interface) v1.UpstreamClient {
	return &serviceEntryClient{istio: client}
}

type serviceEntryClient struct {
	istio istioclient.Interface
}

func (s *serviceEntryClient) List(namespace string, opts skclients.ListOpts) (v1.UpstreamList, error) {
	listOpts := buildListOptions(opts.ExpressionSelector, opts.Selector)
	items, err := s.istio.NetworkingV1beta1().ServiceEntries(namespace).List(opts.Ctx, listOpts)
	if err != nil {
		return nil, err
	}
	return ConvertUpstreamList(items.Items), nil
}

// Watch is implemented based onon the KubeResourceWatch impl in solo-kit
// to keep behaviors consistent.
func (s *serviceEntryClient) Watch(namespace string, opts skclients.WatchOpts) (<-chan v1.UpstreamList, <-chan error, error) {
	// TODO use informer + lister here?
	listOpts := buildListOptions(opts.ExpressionSelector, opts.Selector)

	us, errs := make(chan v1.UpstreamList), make(chan error)

	// periodically list and push updates based on ticker
	var previous v1.UpstreamList
	initialized := false
	updateList := func() {
		list, err := s.istio.NetworkingV1beta1().ServiceEntries(namespace).List(opts.Ctx, listOpts)
		if err != nil {
			errs <- err
			return
		}
		upstreams := ConvertUpstreamList(list.Items)
		// TODO only push if upstreams != previous
		unchanged := !initialized && slices.EqualFunc(upstreams, previous, func(a, b *v1.Upstream) bool {
			return a.Equal(b)
		})
		if unchanged {
			return
		}
		us <- upstreams
		previous = upstreams
		initialized = true
	}

	// setup a kube watch to help us only list when something changes
	watch, err := s.istio.NetworkingV1beta1().ServiceEntries(namespace).Watch(opts.Ctx, listOpts)
	if err != nil {
		return nil, nil, err
	}

	// on a timer, if the watch notified us before the last tick, trigger listing
	go func() {
		defer close(us)
		defer close(errs)
		defer watch.Stop()

		timer := time.NewTicker(time.Second)
		defer timer.Stop()

		update := false
		for {
			select {
			case _, ok := <-watch.ResultChan():
				// only bother doing the fetch if we saw the watch change
				if !ok {
					return
				}
				update = true
			case <-timer.C:
				if update {
					updateList()
					update = false
				}
			case <-opts.Ctx.Done():
				return
			}
		}
	}()

	return us, errs, nil
}

func ConvertUpstreamList(servicEntries []*networkingclient.ServiceEntry) v1.UpstreamList {
	var out v1.UpstreamList
	for _, se := range servicEntries {
		out = append(out, ConvertUpstreams(se)...)
	}
	return out
}

func ConvertUpstreams(se *networkingclient.ServiceEntry) []*v1.Upstream {
	// TODO How to handle these fields:
	// - Reslution: DNS; will gloo automatically use an envoy DNS cluster when it sees non IP endpoints?
	// - Reslution: None; I don't think we can do original_dst/passthrough.
	// - Location: MESH_INTERNAL/MESH_EXTERNAL. Decide whether to use mTLS
	// - exportTo: we're not respecting this but we should
	var out []*v1.Upstream
	for _, port := range se.Spec.Ports {
		labels := make(map[string]string, len(se.Labels))
		for k, v := range se.Labels {
			labels[k] = v
		}
		labels[plugins.UpstreamSourceLabel] = "ServiceEntry"

		us := &v1.Upstream{
			Metadata: &core.Metadata{
				Name:        UpstreamNamePrefix + se.Name,
				Namespace:   se.Namespace,
				Cluster:     "", // can we populate this for the local cluster?
				Labels:      labels,
				Annotations: se.Annotations,
			},
		}

		if se.Spec.WorkloadSelector != nil {
			us.UpstreamType = buildKube(se, port)
		} else {
			us.UpstreamType = buildStatic(se, port)
		}

		out = append(out, us)
	}
	return out
}

func buildKube(se *networkingclient.ServiceEntry, port *networking.ServicePort) *v1.Upstream_Kube {
	return &v1.Upstream_Kube{
		Kube: &kubernetes.UpstreamSpec{
			ServiceName:      se.Name,
			ServiceNamespace: se.Namespace,
			ServicePort:      port.Number,
			Selector:         se.Spec.WorkloadSelector.Labels,
		},
	}
}

func buildStatic(se *networkingclient.ServiceEntry, port *networking.ServicePort) *v1.Upstream_Static {
	var hosts []*static.Host
	for _, ep := range se.Spec.Endpoints {
		hosts = append(hosts, &static.Host{
			Addr: ep.Address,
			Port: targetPort(port.Number, port.TargetPort, wePort(ep, port)),
			// TODO do we want to set this for non-IP endpoints?
			// SniAddr:             "",
			// TODO how structured is this? The inline-WE have Labels we could put here.
			// Metadata:            map[string]*structpb.Struct{},
		})
	}

	return &v1.Upstream_Static{
		Static: &static.UpstreamSpec{
			Hosts:  hosts,
			UseTls: &wrapperspb.BoolValue{Value: isProtocolTLS(port.Protocol)},
			// TODO do we want this on?
			// AutoSniRewrite: &wrapperspb.BoolValue{},
		},
	}
}

func wePort(we *networking.WorkloadEntry, svcPort *networking.ServicePort) uint32 {
	if we.Ports == nil {
		return 0
	}
	return we.Ports[svcPort.Name]
}

// We don't actually use the following in thi shim, but we must satisfy the UpstreamClient interface.

const notImplementedErrMsg = "this operation is not supported by this client"

func (c *serviceEntryClient) BaseClient() skclients.ResourceClient {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return &NoOpUpstreamClient.NoOpUpstreamClient{}
}

func (c *serviceEntryClient) Register() error {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return fmt.Errorf(notImplementedErrMsg)
}

func (c *serviceEntryClient) Read(namespace, name string, opts skclients.ReadOpts) (*v1.Upstream, error) {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return nil, fmt.Errorf(notImplementedErrMsg)
}

func (c *serviceEntryClient) Write(resource *v1.Upstream, opts skclients.WriteOpts) (*v1.Upstream, error) {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return nil, fmt.Errorf(notImplementedErrMsg)
}

func (c *serviceEntryClient) Delete(namespace, name string, opts skclients.DeleteOpts) error {
	contextutils.LoggerFrom(context.Background()).DPanic(notImplementedErrMsg)
	return fmt.Errorf(notImplementedErrMsg)
}
