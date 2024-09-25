package syncer

import (
	"context"
	"fmt"
	"hash/fnv"
	"maps"
	"slices"
	"sync"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_service_discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	xdsserver "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/server"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
	"k8s.io/apimachinery/pkg/types"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"

	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
)

type ConnectedClient struct {
	Node *envoy_config_core_v3.Node
}

func (w ConnectedClient) ResourceName() string {
	return w.Node.Id
}

var _ krt.Equaler[ConnectedClient] = new(ConnectedClient)

func (c ConnectedClient) Equals(k ConnectedClient) bool {
	return proto.Equal(c.Node, c.Node)
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

func NewUniqlyConnectedClient(cc ConnectedClient) UniqlyConnectedClient {
	labels := getLabels(cc.Node.GetMetadata())
	role := cc.Node.GetMetadata().GetFields()["role"].GetStringValue()
	return UniqlyConnectedClient{
		Role:         role,
		Labels:       labels,
		resourceName: fmt.Sprintf("%s~%d", role, hashLabels(labels)),
	}
}

type callbacksCollection struct {
	ctx              context.Context
	clients          map[int64]ConnectedClient
	uniqClientsCount map[string]uint64
	uniqClients      map[string]UniqlyConnectedClient
	fanoutChan       chan krt.Event[UniqlyConnectedClient]
	stateLock        sync.RWMutex

	eventHandlers handlers[UniqlyConnectedClient]
}

// THIS IS THE SET OF THINGS WE RUN TRANSLATION FOR
// add returned callbacks to the xds server.
func NewConnectedClients(ctx context.Context) (xdsserver.Callbacks, krt.Collection[UniqlyConnectedClient]) {

	cb := &callbacksCollection{
		ctx:              ctx,
		clients:          make(map[int64]ConnectedClient),
		uniqClientsCount: make(map[string]uint64),
		uniqClients:      make(map[string]UniqlyConnectedClient),
		fanoutChan:       make(chan krt.Event[UniqlyConnectedClient], 100),
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
var _ krt.Collection[UniqlyConnectedClient] = new(callbacksCollection)

// OnStreamOpen is called once an xDS stream is open with a stream ID and the type URL (or "" for ADS).
// Returning an error will end processing and close the stream. OnStreamClosed will still be called.
func (x *callbacksCollection) OnStreamOpen(_ context.Context, _ int64, _ string) error {
	return nil
}

// OnStreamClosed is called immediately prior to closing an xDS stream with a stream ID.
func (x *callbacksCollection) OnStreamClosed(sid int64) {
	ucc := x.del(sid)
	if ucc != nil {
		x.Notify(krt.Event[UniqlyConnectedClient]{Old: ucc, Event: controllers.EventDelete})
	}
}

func (x *callbacksCollection) del(sid int64) *UniqlyConnectedClient {
	x.stateLock.Lock()
	defer x.stateLock.Unlock()

	c, ok := x.clients[sid]
	delete(x.clients, sid)
	if ok {
		ucc := NewUniqlyConnectedClient(c)
		current := x.uniqClientsCount[ucc.resourceName]
		if current == 1 {
			delete(x.uniqClientsCount, ucc.resourceName)
			delete(x.uniqClients, ucc.resourceName)
			return &ucc
		} else {
			x.uniqClientsCount[ucc.resourceName] = current - 1
		}
	}
	return nil
}

func (x *callbacksCollection) add(sid int64, r *envoy_service_discovery_v3.DiscoveryRequest) *UniqlyConnectedClient {
	x.stateLock.Lock()
	defer x.stateLock.Unlock()

	if _, ok := x.clients[sid]; !ok && r.Node != nil {
		c := ConnectedClient{
			Node: r.Node,
		}
		x.clients[sid] = c
		ucc := NewUniqlyConnectedClient(c)
		current := x.uniqClientsCount[ucc.resourceName]
		x.uniqClientsCount[ucc.resourceName] = current + 1
		if current == 0 {
			x.uniqClients[ucc.resourceName] = ucc
			return &ucc
		}

	}
	return nil
}

func (x *callbacksCollection) Notify(e krt.Event[UniqlyConnectedClient]) {
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
	c := x.add(sid, r)
	if c != nil {
		x.Notify(krt.Event[UniqlyConnectedClient]{New: c, Event: controllers.EventAdd})
	}
	return nil
}

func (x *callbacksCollection) getClients() []UniqlyConnectedClient {
	x.stateLock.RLock()
	defer x.stateLock.RUnlock()
	clients := make([]UniqlyConnectedClient, 0, len(x.uniqClients))
	for _, c := range x.uniqClients {
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

func (x *callbacksCollection) Register(f func(o krt.Event[UniqlyConnectedClient])) krt.Syncer {
	return x.RegisterBatch(func(events []krt.Event[UniqlyConnectedClient], initialSync bool) {
		for _, o := range events {
			f(o)
		}
	}, true)
}

func (x *callbacksCollection) RegisterBatch(f func(o []krt.Event[UniqlyConnectedClient], initialSync bool), runExistingState bool) krt.Syncer {
	if runExistingState {
		f(nil, true)
	}
	// notify under lock to make sure that no regular events fire during the initial sync.
	x.eventHandlers.mu.Lock()
	defer x.eventHandlers.mu.Unlock()
	x.eventHandlers.insertLocked(f)
	if runExistingState {
		var events []krt.Event[UniqlyConnectedClient]
		for _, v := range x.getClients() {
			events = append(events, krt.Event[UniqlyConnectedClient]{
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

func (x *callbacksCollection) GetKey(k krt.Key[UniqlyConnectedClient]) *UniqlyConnectedClient {
	x.stateLock.RLock()
	defer x.stateLock.RUnlock()

	u, ok := x.uniqClients[string(k)]
	if ok {
		return &u
	}
	return nil
}

// List returns all objects in the collection.
// Order of the list is undefined.

func (x *callbacksCollection) List() []UniqlyConnectedClient { return x.getClients() }

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

type Clusters = krt.Index[envoy_config_cluster_v3.Cluster, upstreamRef]

func clusters(krt.Collection[v1.Upstream]) Clusters {
	panic("implement me")
}

type Endpoints = krt.Index[envoy_config_endpoint_v3.ClusterLoadAssignment, upstreamRef]

func endpoints(krt.Collection[v1.Upstream]) Endpoints {
	panic("implement me")
}

type Listeners = krt.Index[envoy_config_listener_v3.Listener, proxyRef]

func listeners(krt.Collection[v1.Proxy]) Listeners {
	panic("implement me")
}

type Routes = krt.Index[envoy_config_route_v3.RouteConfiguration, proxyRef]

func routes(krt.Collection[v1.Proxy]) Routes {
	panic("implement me")
}

func snapForProxy(x krt.Collection[v1.Proxy]) krt.Index[xdsSnapshot, proxyRef] {
	return krt.Index[xdsSnapshot, proxyRef]{}
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
