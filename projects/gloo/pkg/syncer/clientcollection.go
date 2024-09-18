package syncer

import (
	"context"
	"slices"
	"sync"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	xdsserver "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/server"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/maps"
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

type callbacksCollection struct {
	clients     map[int64]ConnectedClient
	clientsLock sync.RWMutex

	eventHandlers handlers[ConnectedClient]
}

// add returned callbacks to the xds server.
func NewConnectedClients() (xdsserver.Callbacks, krt.Collection[ConnectedClient]) {

	cb := &callbacksCollection{
		clients: make(map[int64]ConnectedClient),
	}
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
	x.clientsLock.Lock()
	c, ok := x.clients[sid]
	delete(x.clients, sid)
	x.clientsLock.Unlock()

	if ok {
		// TODO: should this be in a goroutine?
		x.eventHandlers.Notify(krt.Event[ConnectedClient]{Old: &c, Event: controllers.EventDelete})
	}
}

func (x *callbacksCollection) add(sid int64, r *envoy_service_discovery_v3.DiscoveryRequest) *ConnectedClient {
	x.clientsLock.Lock()
	defer x.clientsLock.Unlock()

	if _, ok := x.clients[sid]; !ok && r.Node != nil {
		c := ConnectedClient{
			Node: r.Node,
		}
		x.clients[sid] = c

		return &c
	}
	return nil
}

// OnStreamRequest is called once a request is received on a stream.
// Returning an error will end processing and close the stream. OnStreamClosed will still be called.
func (x *callbacksCollection) OnStreamRequest(sid int64, r *envoy_service_discovery_v3.DiscoveryRequest) error {
	c := x.add(sid, r)
	if c != nil {
		// TODO: should this be in a goroutine?
		x.eventHandlers.Notify(krt.Event[ConnectedClient]{New: c, Event: controllers.EventAdd})
	}
	return nil
}

func (x *callbacksCollection) getClients() []ConnectedClient {
	x.clientsLock.RLock()
	defer x.clientsLock.RUnlock()
	return maps.Values(x.clients)

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
	x.eventHandlers.Insert(f)
	if runExistingState {
		for _, v := range x.getClients() {

			f([]krt.Event[ConnectedClient]{{
				New:   &v,
				Event: controllers.EventAdd,
			}}, true)
		}
	}

	return &simpleSyncer{}
}

func (x *callbacksCollection) Synced() krt.Syncer {
	return &simpleSyncer{}
}

// GetKey returns an object by its key, if present. Otherwise, nil is returned.

func (x *callbacksCollection) GetKey(k krt.Key[ConnectedClient]) *ConnectedClient {
	clients := x.getClients()
	for _, c := range clients {
		if string(k) == c.ResourceName() {
			return &c
		}
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
