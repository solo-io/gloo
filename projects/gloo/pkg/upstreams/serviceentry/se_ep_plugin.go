package serviceentry

import (
	"net"
	"strconv"

	"github.com/solo-io/gloo/projects/gloo/pkg/discovery"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
	"istio.io/istio/pkg/config/schema/gvr"
	kubeclient "istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/kube/kubetypes"
	corev1 "k8s.io/api/core/v1"
)

var _ discovery.DiscoveryPlugin = &sePlugin{}

const (
	podEndpointPrefix = "istio-se-pod-ep"
	weEndpointPrefix  = "istio-se-we-ep"
)

func epName(seName, resKind, resName, ns string) string {
	return "se-" + seName + "-ep-" + resKind + "-" + resName + "-" + ns
}

type epBuildFunc[T any] func(
	se serviceEntryUpstreams,
	o T,
	writeNamespace string,
	svcPort uint32,
	portName string,
) *v1.Endpoint

func buildPodEndpoint(
	se serviceEntryUpstreams,
	pod *corev1.Pod,
	writeNamespace string,
	svcPort uint32,
	_ string,
) *v1.Endpoint {
	port := svcPort
	if sePort := se.targetPorts[svcPort]; sePort > 0 {
		port = sePort
	}
	return &v1.Endpoint{
		Metadata: &core.Metadata{
			Name:        epName(se.Name, "Pod", pod.Name, pod.Namespace),
			Namespace:   writeNamespace,
			Labels:      pod.Labels,
			Annotations: pod.Annotations,
		},
		Address: pod.Status.PodIP,
		Port:    port,
		// Hostname:    "",
	}
}

func buildWorkloadEntryEndpoint(
	se serviceEntryUpstreams,
	we *networkingclient.WorkloadEntry,
	writeNamespace string,
	svcPort uint32,
	portName string,
) *v1.Endpoint {
	port := svcPort
	if wePort := we.Spec.Ports[portName]; wePort > 0 {
		port = wePort
	} else if sePort := se.targetPorts[svcPort]; sePort > 0 {
		port = sePort
	}

	return &v1.Endpoint{
		Metadata: &core.Metadata{
			Name:        epName(we.Name, "WorkloadEntry", we.Name, we.Namespace),
			Namespace:   writeNamespace,
			Labels:      we.Labels,
			Annotations: we.Annotations,
		},
		Address: we.Spec.Address,
		Port:    port,
		// Hostname:    "",
	}
}

func (s *sePlugin) WatchEndpoints(
	writeNamespace string,
	upstreamsToTrack v1.UpstreamList,
	opts clients.WatchOpts,
) (<-chan v1.EndpointList, <-chan error, error) {
	defaultFilter := kclient.Filter{ObjectFilter: s.client.ObjectFilter()}

	// watch service entries everywhere, they're referenced by global hostnames
	// we don't know if the VirtualDestination will cause SE to be created in
	// system namespace or user namespace.
	seInformer := kclient.NewDelayedInformer[*networkingclient.ServiceEntry](s.client,
		gvr.ServiceEntry, kubetypes.StandardInformer, defaultFilter)
	ServiceEntries := krt.WrapClient(seInformer, krt.WithName("ServiceEntries"))

	// since we watch ServiceEntry everywhere, we must do the same for WorkloadEntries
	weInformer := kclient.NewDelayedInformer[*networkingclient.WorkloadEntry](s.client,
		gvr.ServiceEntry, kubetypes.StandardInformer, defaultFilter)
	WorkloadEntries := krt.WrapClient(weInformer, krt.WithName("WorkloadEntries"))

	// since we watch ServiceEntry everywhere, we must do the same for Pods
	podInformer := kclient.NewFiltered[*corev1.Pod](s.client, kclient.Filter{
		ObjectTransform: kubeclient.StripPodUnusedFields,
	})
	Pods := krt.WrapClient(podInformer, krt.WithName("Pods"))

	// Upstreams by host:port
	Upstreams, MeshUpstreamsIndex := meshUpstreams(upstreamsToTrack)

	// ServiceEntries with upstream references
	upstreamServiceEntries, seNsIndex := serviceEntriesWithUpstreams(ServiceEntries, Upstreams, MeshUpstreamsIndex)

	// generate endpoints
	podEndpoints := podEndpoints(
		Pods,
		writeNamespace,
		upstreamServiceEntries, seNsIndex,
	)
	weEndpoints := weEndpoints(
		WorkloadEntries,
		writeNamespace,
		upstreamServiceEntries, seNsIndex,
	)
	seInlineEndpoints := seInlineEndpoints(writeNamespace, upstreamServiceEntries)
	allEndpoints := krt.JoinCollection([]krt.Collection[glooEndpoint]{
		podEndpoints,
		weEndpoints,
		seInlineEndpoints,
	})

	// TODO inline endpoints from serviceentry
	// seEndpoints := krt.NewManyCollection[*networkingclient.ServiceEntry, *v1.Endpoint](ServiceEntries, weEndpoints(
	// 	...
	// ))

	// TODO also handle Upstream_Kube for Waypoint usecase:
	// * ServiceEntry will gen Upstream_Kube so users can use it just like in Istio
	// * We also need to allow Service selects WorkloadEntry for Istio parity

	// push updates to watcher
	endpointsChan := make(chan v1.EndpointList)
	errs := make(chan error)
	allEndpoints.RegisterBatch(func(_ []krt.Event[glooEndpoint], _ bool) {
		ep := endpoints(podEndpoints.List()).Unwrap()
		select {
		case endpointsChan <- ep:
		default:
		}
	}, true)
	go func() {
		<-opts.Ctx.Done()
		close(endpointsChan)
		close(errs)
	}()

	go seInformer.Start(opts.Ctx.Done())
	go podInformer.Start(opts.Ctx.Done())
	// TODO the client doesn't need to be run on every WatchEndpoints call
	// TODO maybe it's better to only set up new Upstream-dependent collections for each Watch
	go s.client.RunAndWait(opts.Ctx.Done())

	return endpointsChan, errs, nil
}

// meshUpstreams are Static Upstreams with isMesh set to true.
// The given index is keyed by host:port pairs from the static upstream `hosts` field.
func meshUpstreams(upstreamsToTrack v1.UpstreamList) (krt.Collection[upstream], krt.Index[string, upstream]) {
	var upstreams []upstream
	for _, us := range upstreamsToTrack {
		upstreams = append(upstreams, upstream{us})
	}

	// index mesh upstreams by the host/port combos in their hosts
	Upstreams := krt.NewStaticCollection(upstreams)
	MeshUpstreamsIndex := krt.NewIndex(Upstreams, func(us upstream) []string {
		if !us.GetStatic().GetIsMesh() {
			return nil
		}
		var hostPorts []string
		for _, h := range us.GetStatic().GetHosts() {
			hostPorts = append(hostPorts, net.JoinHostPort(h.GetAddr(), strconv.Itoa(int(h.GetPort()))))
		}
		return hostPorts
	})
	return Upstreams, MeshUpstreamsIndex
}

// serviceEntriesWithUpstreams is a collection of ServiceEntry along with all of the
// Upstreams that reference this ServiceEntry by hostname.
// We also return a namespace index to speed up other lookups.
func serviceEntriesWithUpstreams(
	ServiceEntries krt.Collection[*networkingclient.ServiceEntry],
	Upstreams krt.Collection[upstream], MeshUpstreamsIndex krt.Index[string, upstream],
) (krt.Collection[serviceEntryUpstreams], krt.Index[string, serviceEntryUpstreams]) {
	upstreamServiceEntries := krt.NewCollection(ServiceEntries, func(ctx krt.HandlerContext, se *networkingclient.ServiceEntry) *serviceEntryUpstreams {
		referencingUs := make(map[uint32][]upstream, len(se.Spec.Ports))
		portMapping := make(map[uint32]uint32, len(se.Spec.Ports))
		for _, host := range se.Spec.Hosts {
			for _, port := range se.Spec.Ports {
				key := net.JoinHostPort(host, strconv.Itoa(int(port.Number)))
				referencedBy := krt.Fetch(ctx, Upstreams, krt.FilterIndex(MeshUpstreamsIndex, key))
				referencingUs[port.Number] = append(referencingUs[port.Number], referencedBy...)

				targetPort := port.Number
				if port.TargetPort > 0 {
					targetPort = port.TargetPort
				}
				portMapping[port.Number] = targetPort
			}
		}
		if len(referencingUs) == 0 {
			return nil
		}
		return &serviceEntryUpstreams{
			ServiceEntry: se,
			upstreams:    referencingUs,
			targetPorts:  portMapping,
		}
	})
	seNsIndex := krt.NewNamespaceIndex(upstreamServiceEntries)
	return upstreamServiceEntries, seNsIndex
}

// seInlineEndpoints builds Gloo Endpoints for the inline endpoints field.
func seInlineEndpoints(
	writeNamespace string,
	upstreamServiceEntries krt.Collection[serviceEntryUpstreams],
) krt.Collection[glooEndpoint] {
	return krt.NewManyCollection(upstreamServiceEntries, func(ctx krt.HandlerContext, se serviceEntryUpstreams) []glooEndpoint {
		if se.Spec.GetWorkloadSelector().GetLabels() != nil && len(se.Spec.GetWorkloadSelector().GetLabels()) > 0 {
			// workloadSelector and inline endpoints can't be used together
			return nil
		}

		eps := map[uint32]*v1.Endpoint{}
		for i, we := range se.Spec.GetEndpoints() {
			weNamed := &networkingclient.WorkloadEntry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      se.Name + "inline-" + strconv.Itoa(i),
					Namespace: se.Namespace,
					Labels:    we.Labels,
				},
				// TODO lock copy...
				Spec: *we,
			}
			generateEndpointsForServiceEntry(writeNamespace, se, eps, weNamed, buildWorkloadEntryEndpoint)
		}
		return nil
	})
}

// podEndpoints generates Endpoints for a single Pod.
// See buildSelectedEndpoints for details.
func podEndpoints(
	Pods krt.Collection[*corev1.Pod],
	writeNamespace string,
	upstreamServiceEntries krt.Collection[serviceEntryUpstreams],
	seNsIndex krt.Index[string, serviceEntryUpstreams],
) krt.Collection[glooEndpoint] {
	return krt.NewManyCollection(Pods, func(ctx krt.HandlerContext, pod *corev1.Pod) []glooEndpoint {
		return buildSelectedEndpoints(
			ctx,
			upstreamServiceEntries, seNsIndex,
			writeNamespace,
			pod, pod.Namespace, pod.GetLabels(),
			buildPodEndpoint,
		)
	})
}

// weEndpoints generates Endpoints for a single WorkloadEntry.
// See buildSelectedEndpoints for details.
func weEndpoints(
	WorkloadEntries krt.Collection[*networkingclient.WorkloadEntry],
	writeNamespace string,
	upstreamServiceEntries krt.Collection[serviceEntryUpstreams],
	seNsIndex krt.Index[string, serviceEntryUpstreams],
) krt.Collection[glooEndpoint] {
	return krt.NewManyCollection(WorkloadEntries, func(ctx krt.HandlerContext, we *networkingclient.WorkloadEntry) []glooEndpoint {
		if we.Spec.Address == "" {
			// TODO handle empty address WorkloadEntry for cross-network?
			// I think we end up using the VIP outbound and zTunnel handles it.
			return nil
		}
		return buildSelectedEndpoints(
			ctx,
			upstreamServiceEntries, seNsIndex,
			writeNamespace,
			we, we.Namespace, we.GetLabels(),
			buildWorkloadEntryEndpoint,
		)
	})
}

// buildSelectedEndpoints applies the ServiceEntry's workloadSelector to o.
// The output will contain and endpoint for each port reachable by a ServiceEntry,
// only for ServiceEntry that static in-mesh Upstreams reference by Hostname.
func buildSelectedEndpoints[T any](
	ctx krt.HandlerContext,
	upstreamServiceEntries krt.Collection[serviceEntryUpstreams],
	seNsIndex krt.Index[string, serviceEntryUpstreams],
	writeNamespace string,
	o T,
	ns string, labels map[string]string,
	mappingFunc epBuildFunc[T],
) []glooEndpoint {
	// find all the ServiceEntries that select us
	selectedBy := krt.Fetch(ctx, upstreamServiceEntries, krt.FilterIndex(seNsIndex, ns), krt.FilterSelectsNonEmpty(labels))
	if len(selectedBy) == 0 {
		return nil
	}

	endpointsByPort := map[uint32]*v1.Endpoint{}

	for _, se := range selectedBy {
		if len(se.upstreams) == 0 {
			continue
		}
		generateEndpointsForServiceEntry(writeNamespace, se, endpointsByPort, o, mappingFunc)
	}

	var out []glooEndpoint
	for _, ep := range endpointsByPort {
		out = append(out, glooEndpoint{ep})
	}
	return out
}

// generateEndpointsForServiceEntry will build an endpoint for each reachable
// port. because different upstreams may point to a different port on the same
// endpoint. Each Endpoint may be shared by multiple upstreams.
func generateEndpointsForServiceEntry[T any](
	writeNamespace string,
	se serviceEntryUpstreams,
	endpointsByPort map[uint32]*v1.Endpoint,
	o T,
	mappingFunc epBuildFunc[T],
) {
	for svcPort, upstreams := range se.upstreams {
		// init an endpoint instance for this port if necessary
		endpoint := endpointsByPort[svcPort]
		if endpoint == nil {
			portName := se.portNames[svcPort]
			endpoint = mappingFunc(se, o, writeNamespace, svcPort, portName)
			endpointsByPort[svcPort] = endpoint
		}
		for _, us := range upstreams {
			endpoint.Upstreams = append(endpoint.Upstreams, &core.ResourceRef{
				Name:      us.GetMetadata().GetName(),
				Namespace: us.GetMetadata().GetNamespace(),
			})
		}
	}
}

func getNamespace(obj runtime.Object) string {
	metaObj, err := meta.Accessor(obj)
	if err != nil {
		return ""
	}
	return metaObj.GetNamespace()
}

// satisfy interface but don't implement; only implement hybrid client for ServiceEntry Upstreams.
func (s *sePlugin) DiscoverUpstreams(watchNamespaces []string, writeNamespace string, opts clients.WatchOpts, discOpts discovery.Opts) (chan v1.UpstreamList, chan error, error) {
	us, errs := make(chan v1.UpstreamList), make(chan error)
	close(us)
	close(errs)
	return us, errs, nil
}

// satisfy interface but don't implement; only implement hybrid client for ServiceEntry Upstreams.
func (s *sePlugin) UpdateUpstream(original *v1.Upstream, desired *v1.Upstream) (bool, error) {
	return false, nil
}
