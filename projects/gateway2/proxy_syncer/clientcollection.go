package proxy_syncer

import (
	"context"
	"fmt"
	"hash/fnv"
	"maps"
	"slices"
	"strings"
	"sync"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_service_discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	xdsserver "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/server"
	"google.golang.org/protobuf/types/known/structpb"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
	"k8s.io/apimachinery/pkg/types"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"

	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
)

type ConnectedClient struct {
	Sid krt.Key[ConnectedClient]

	node   *envoy_config_core_v3.Node
	podRef types.NamespacedName
}

func newConnectedClient(node *envoy_config_core_v3.Node, sid krt.Key[ConnectedClient]) ConnectedClient {
	return ConnectedClient{
		Sid:    sid,
		node:   node,
		podRef: getRef(node),
	}
}

func (w ConnectedClient) ResourceName() string {
	return string(w.Sid)
}

var _ krt.Equaler[ConnectedClient] = new(ConnectedClient)

func (c ConnectedClient) Equals(k ConnectedClient) bool {
	// resource name includes sid; and sid is unique
	return c.Sid == k.Sid
}

type UniqlyConnectedClient struct {
	// name namespace of gateway
	Role   string
	Labels map[string]string

	resourceName string
}

func (c UniqlyConnectedClient) ResourceName() string {
	return c.resourceName
}

var _ krt.Equaler[UniqlyConnectedClient] = new(UniqlyConnectedClient)

func (c UniqlyConnectedClient) Equals(k UniqlyConnectedClient) bool {
	return maps.Equal(c.Labels, k.Labels) && c.Role == k.Role
}

func NewUniqlyConnectedClient(cc ConnectedClient, labels map[string]string) UniqlyConnectedClient {
	role := cc.node.GetMetadata().GetFields()["role"].GetStringValue()
	return UniqlyConnectedClient{
		Role:         role,
		Labels:       labels,
		resourceName: fmt.Sprintf("%s~%d", role, hashLabels(labels)),
	}
}

type callbacksCollection struct {
	ctx        context.Context
	clients    map[krt.Key[ConnectedClient]]ConnectedClient
	fanoutChan chan krt.Event[ConnectedClient]
	stateLock  sync.RWMutex

	eventHandlers handlers[ConnectedClient]
	augmentedPods krt.Collection[augmentedPod]
}

// THIS IS THE SET OF THINGS WE RUN TRANSLATION FOR
// add returned callbacks to the xds server.
func NewConnectedClients(ctx context.Context, augmentedPods krt.Collection[augmentedPod]) (xdsserver.Callbacks, krt.Collection[ConnectedClient]) {

	cb := &callbacksCollection{
		ctx:           ctx,
		clients:       make(map[krt.Key[ConnectedClient]]ConnectedClient),
		fanoutChan:    make(chan krt.Event[ConnectedClient], 100),
		augmentedPods: augmentedPods,
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-cb.fanoutChan:
				cb.eventHandlers.Notify(event)
			}
		}
	}()
	return cb, cb
}

var _ xdsserver.Callbacks = new(callbacksCollection)
var _ krt.Collection[ConnectedClient] = new(callbacksCollection)

// OnStreamOpen is called once an xDS stream is open with a stream ID and the type URL (or "" for ADS).
// Returning an error will end processing and close the stream. OnStreamClosed will still be called.
func (x *callbacksCollection) OnStreamOpen(_ context.Context, _ int64, _ string) error {
	return nil
}

// OnStreamClosed is called immediately prior to closing an xDS stream with a stream ID.
func (x *callbacksCollection) OnStreamClosed(sid int64) {
	c := x.del(sid)
	if c != nil {
		x.Notify(krt.Event[ConnectedClient]{Old: c, Event: controllers.EventDelete})
	}
}

func resouceName(sid int64) krt.Key[ConnectedClient] {
	return krt.Key[ConnectedClient](fmt.Sprintf("%d", sid))
}

func (x *callbacksCollection) del(sid int64) *ConnectedClient {
	x.stateLock.Lock()
	defer x.stateLock.Unlock()

	key := resouceName(sid)
	c, ok := x.clients[key]
	delete(x.clients, key)
	if ok {
		return &c
	}
	return nil
}

func (x *callbacksCollection) add(sid int64, r *envoy_service_discovery_v3.DiscoveryRequest) (*ConnectedClient, error) {
	x.stateLock.Lock()
	defer x.stateLock.Unlock()

	key := resouceName(sid)
	if _, ok := x.clients[key]; !ok && r.Node != nil {

		// TODO: modify request to include the label that are relevant for the client?
		// error if we can get the pod
		podRef := getRef(r.Node)
		k := krt.Key[augmentedPod](krt.Named{Name: podRef.Name, Namespace: podRef.Namespace}.ResourceName())
		pod := x.augmentedPods.GetKey(k)
		if pod == nil {
			return nil, fmt.Errorf("pod not found for node %v", r.Node)
		}

		c := newConnectedClient(r.Node, key)
		x.clients[key] = c
		return &c, nil
	}
	return nil, nil
}

func (x *callbacksCollection) Notify(e krt.Event[ConnectedClient]) {
	for {
		// note: do not use a default block here, we want to block if the channel is full, as otherwise we will have inconsistent state in krt.
		select {
		case x.fanoutChan <- e:
			return
		case <-x.ctx.Done():
			return
		}
	}
}

// OnStreamRequest is called once a request is received on a stream.
// Returning an error will end processing and close the stream. OnStreamClosed will still be called.
func (x *callbacksCollection) OnStreamRequest(sid int64, r *envoy_service_discovery_v3.DiscoveryRequest) error {
	c, err := x.add(sid, r)
	if err != nil {
		return err
	}
	if c != nil {
		x.Notify(krt.Event[ConnectedClient]{New: c, Event: controllers.EventAdd})
	}
	return nil
}

func (x *callbacksCollection) getClients() []ConnectedClient {
	x.stateLock.RLock()
	defer x.stateLock.RUnlock()
	clients := make([]ConnectedClient, 0, len(x.clients))
	for _, c := range x.clients {
		clients = append(clients, c)
	}
	return clients
}

// OnStreamResponse is called immediately prior to sending a response on a stream.
func (x *callbacksCollection) OnStreamResponse(_ int64, _ *envoy_service_discovery_v3.DiscoveryRequest, _ *envoy_service_discovery_v3.DiscoveryResponse) {
}

// OnFetchRequest is called for each Fetch request. Returning an error will end processing of the
// request and respond with an error.
func (x *callbacksCollection) OnFetchRequest(_ context.Context, _ *envoy_service_discovery_v3.DiscoveryRequest) error {
	return nil
}

// OnFetchResponse is called immediately prior to sending a response.
func (x *callbacksCollection) OnFetchResponse(_ *envoy_service_discovery_v3.DiscoveryRequest, _ *envoy_service_discovery_v3.DiscoveryResponse) {
}

func (x *callbacksCollection) Register(f func(o krt.Event[ConnectedClient])) krt.Syncer {
	return x.RegisterBatch(func(events []krt.Event[ConnectedClient], initialSync bool) {
		for _, o := range events {
			f(o)
		}
	}, true)
}

func (x *callbacksCollection) RegisterBatch(f func(o []krt.Event[ConnectedClient], initialSync bool), runExistingState bool) krt.Syncer {
	if runExistingState {
		f(nil, true)
	}
	// notify under lock to make sure that no regular events fire during the initial sync.
	x.eventHandlers.mu.Lock()
	defer x.eventHandlers.mu.Unlock()
	x.eventHandlers.insertLocked(f)
	if runExistingState {
		var events []krt.Event[ConnectedClient]
		for _, v := range x.getClients() {
			events = append(events, krt.Event[ConnectedClient]{
				New:   &v,
				Event: controllers.EventAdd,
			})

		}
		f(events, true)
	}

	return &simpleSyncer{}
}

func (x *callbacksCollection) Synced() krt.Syncer {
	return &simpleSyncer{}
}

// GetKey returns an object by its key, if present. Otherwise, nil is returned.

func (x *callbacksCollection) GetKey(k krt.Key[ConnectedClient]) *ConnectedClient {
	x.stateLock.RLock()
	defer x.stateLock.RUnlock()
	u, ok := x.clients[k]
	if ok {
		return &u
	}
	return nil
}

// List returns all objects in the collection.
// Order of the list is undefined.

func (x *callbacksCollection) List() []ConnectedClient { return x.getClients() }

type simpleSyncer struct{}

func (s *simpleSyncer) WaitUntilSynced(stop <-chan struct{}) bool {
	return true
}

func (s *simpleSyncer) HasSynced() bool {
	return true
}

type handlers[O any] struct {
	mu sync.RWMutex
	h  []func(o []krt.Event[O], initialSync bool)
}

func (o *handlers[O]) Insert(f func(o []krt.Event[O], initialSync bool)) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.h = append(o.h, f)
}
func (o *handlers[O]) insertLocked(f func(o []krt.Event[O], initialSync bool)) {
	o.h = append(o.h, f)
}

func (o *handlers[O]) Get() []func(o []krt.Event[O], initialSync bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return slices.Clone(o.h)
}

func (o *handlers[O]) Notify(e krt.Event[O]) {
	cb := o.Get()
	for _, f := range cb {
		events := [1]krt.Event[O]{e}
		f(events[:], false)
	}
}

func hashLabels(labels map[string]string) uint64 {
	finalHash := uint64(0)
	for k, v := range labels {
		hasher := fnv.New64()
		hasher.Write([]byte(k))
		hasher.Write([]byte{0})
		hasher.Write([]byte(v))
		hasher.Write([]byte{0})
		finalHash ^= hasher.Sum64()
	}
	return finalHash
}

func getLabels(md *structpb.Struct) map[string]string {
	labels := make(map[string]string)
	labelsStruct := md.GetFields()["labels"].GetStructValue()
	for k, v := range labelsStruct.GetFields() {
		labels[k] = v.GetStringValue()
	}
	return labels
}

// /////////////////// sketch of the rest of the syncer ////
type upstreamRef struct {
	name      string
	namespace string
}
type proxyRef struct {
	name       string
	namespace  string
	labelsHash string
}

type Clusters = krt.Index[upstreamRef, envoy_config_cluster_v3.Cluster]

func clusters(krt.Collection[v1.Upstream]) Clusters {
	panic("implement me")
}

type Endpoints = krt.Index[upstreamRef, envoy_config_endpoint_v3.ClusterLoadAssignment]

func endpoints(krt.Collection[v1.Upstream]) Endpoints {
	panic("implement me")
}

type Listeners = krt.Index[proxyRef, envoy_config_listener_v3.Listener]

func listeners(krt.Collection[v1.Proxy]) Listeners {
	panic("implement me")
}

type Routes = krt.Index[proxyRef, envoy_config_route_v3.RouteConfiguration]

func routes(krt.Collection[v1.Proxy]) Routes {
	panic("implement me")
}

func snapForProxy(x krt.Collection[v1.Proxy]) krt.Index[proxyRef, xdsSnapshot] {
	panic("implement me")
}

type xdsSnapshot struct {
	Upstreams krt.Collection[v1.Upstream]
	// Proxies   krt.Collection[v1.Proxy]
	Clusters  Clusters
	Endpoints Endpoints
	Listeners Listeners
	Routes    Routes
}

func defaultSnapshot(proxies krt.Collection[v1.Proxy], snap xdsSnapshot, xdsCache envoycache.SnapshotCache) {
	krt.NewCollection(proxies, func(ctx krt.HandlerContext, o v1.Proxy) *proxyRef {
		// put snapshot in the cache
		upstreams := snap.Upstreams.List()

		clusters := make([]*envoy_config_cluster_v3.Cluster, 0, len(upstreams))
		for i := range upstreams {
			u := upstreams[i]
			c := snap.Clusters.Lookup(upstreamRef{name: u.Name, namespace: u.Namespace})
			for i := range c {
				clusters = append(clusters, &c[i])
			}
		}
		// TODO:
		var endpoints []*envoy_config_endpoint_v3.ClusterLoadAssignment
		var routeConfigs []*envoy_config_route_v3.RouteConfiguration
		listeners := toPtrSlice(snap.Listeners.Lookup(proxyRef{name: o.Name, namespace: o.Namespace}))

		xdsSnapshot := generateXDSSnapshot(clusters, endpoints, routeConfigs, listeners)

		key := xds.SnapshotCacheKey(&o.Spec)
		xdsCache.SetSnapshot(key, xdsSnapshot)
		return &proxyRef{
			name:       o.Name,
			namespace:  o.Namespace,
			labelsHash: "",
		}
	})
}

func toPtrSlice[T any](t []T) []*T {
	out := make([]*T, len(t))
	for i := range t {
		out[i] = &t[i]
	}
	return out
}

func snapshotForProxy(clients krt.Collection[UniqlyConnectedClient], proxies krt.Collection[v1.Proxy], snap xdsSnapshot, xdsCache envoycache.SnapshotCache) {
	krt.NewCollection(clients, func(ctx krt.HandlerContext, o UniqlyConnectedClient) *proxyRef {
		// fetch the proxy for the client

		//grab proxy name from client role
		proxy := krt.FetchOne(ctx, proxies, krt.FilterObjectName(types.NamespacedName{Name: "TODO", Namespace: "TODO"}))

		// get the labels of the unique client, see if there are any dest rules that apply to this proxy
		// (on its namespace, and if selectors match)
		/*
			find the relevant proxy
			see if dest rules are relevant for this proxy
			if so, re-compute upstreams and endpoints.
		*/

		// if so, recompute the upstreams and endpoints

		//		and do the same as above with the client's cache key.
		return &proxyRef{
			name:       proxy.Name,
			namespace:  proxy.Namespace,
			labelsHash: o.ResourceName(),
		}
	})
}

func generateXDSSnapshot(
	clusters []*envoy_config_cluster_v3.Cluster,
	endpoints []*envoy_config_endpoint_v3.ClusterLoadAssignment,
	routeConfigs []*envoy_config_route_v3.RouteConfiguration,
	listeners []*envoy_config_listener_v3.Listener,
) envoycache.Snapshot {

	panic("implement me")
}

//////////////////////////////////////////////////////////////////////

func New2(c krt.Collection[ConnectedClient], augmentedPods krt.Collection[augmentedPod]) krt.Collection[UniqlyConnectedClient] {
	return krt.NewManyFromNothing(
		func(kctx krt.HandlerContext) []UniqlyConnectedClient {
			unqiueClients := make(map[string]struct{})
			var ret []UniqlyConnectedClient
			for _, cc := range krt.Fetch(kctx, c) {
				if cc.podRef.Name == "" {
					continue
				}

				var labels map[string]string
				maybePod := krt.FetchOne(kctx, augmentedPods, krt.FilterObjectName(cc.podRef))
				if maybePod != nil {
					labels = maybePod.podLabels
				}

				ucc := NewUniqlyConnectedClient(cc, labels)
				if _, ok := unqiueClients[ucc.resourceName]; ok {
					continue
				}
				ret = append(ret, ucc)
			}
			return ret
		})
}

func getRef(node *envoy_config_core_v3.Node) types.NamespacedName {
	nns := node.Id
	split := strings.SplitN(nns, ".", 2)
	if len(split) != 2 {
		return types.NamespacedName{}
	}
	return types.NamespacedName{
		Name:      split[0],
		Namespace: split[1],
	}
}

func defaultSnapshot2(proxies krt.Collection[v1.Proxy], snap xdsSnapshot, xdsCache envoycache.SnapshotCache) {
	krt.NewCollection(proxies, func(ctx krt.HandlerContext, o v1.Proxy) *proxyRef {
		// put snapshot in the cache
		upstreams := snap.Upstreams.List()

		clusters := make([]*envoy_config_cluster_v3.Cluster, 0, len(upstreams))
		for i := range upstreams {
			u := &upstreams[i]
			c := snap.Clusters.Lookup(upstreamRef{name: u.Name, namespace: u.Namespace})
			for i := range c {
				clusters = append(clusters, &c[i])
			}
		}
		// TODO:
		var endpoints []*envoy_config_endpoint_v3.ClusterLoadAssignment
		var routeConfigs []*envoy_config_route_v3.RouteConfiguration
		listeners := toPtrSlice(snap.Listeners.Lookup(proxyRef{name: o.Name, namespace: o.Namespace}))

		xdsSnapshot := generateXDSSnapshot(clusters, endpoints, routeConfigs, listeners)

		key := xds.SnapshotCacheKey(&o.Spec)
		xdsCache.SetSnapshot(key, xdsSnapshot)
		return &proxyRef{
			name:       o.Name,
			namespace:  o.Namespace,
			labelsHash: "",
		}
	})
}

func snapshotForProxy2(clients krt.Collection[UniqlyConnectedClient], proxies krt.Collection[v1.Proxy], snap xdsSnapshot, xdsCache envoycache.SnapshotCache) {
	krt.NewCollection(clients, func(ctx krt.HandlerContext, o UniqlyConnectedClient) *proxyRef {
		// fetch the proxy for the client
		//grab proxy name from client role
		proxy := krt.FetchOne(ctx, proxies, krt.FilterObjectName(types.NamespacedName{Name: "TODO", Namespace: "TODO"}))
		// get the labels of the unique client, see if there are any dest rules that apply to this proxy
		// (on its namespace, and if selectors match)
		/*
			find the relevant proxy
			see if dest rules are relevant for this proxy
			if so, re-compute upstreams and endpoints.
		*/

		// set the node hash for the client

		// if so, recompute the upstreams and endpoints

		//		and do the same as above with the client's cache key.
		return &proxyRef{
			name:       proxy.Name,
			namespace:  proxy.Namespace,
			labelsHash: o.ResourceName(),
		}
	})
}

///// on the watch, we somehow need to feed to the cache the specific proxy key
// for this client.... a bit hacky but possible - slightly less hacky,
// modify the request and fix node hasher!
