package syncer

import (
	"context"
	"fmt"
	"hash/fnv"
	"slices"
	"sync"

	"maps"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	xdsserver "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/server"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
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
	ctx        context.Context
	clients    map[int64]ConnectedClient
	fanoutChan chan krt.Event[ConnectedClient]
	stateLock  sync.RWMutex

	eventHandlers handlers[ConnectedClient]
}

// add returned callbacks to the xds server.
func NewConnectedClients(ctx context.Context) (xdsserver.Callbacks, krt.Collection[ConnectedClient]) {

	cb := &callbacksCollection{
		ctx:        ctx,
		clients:    make(map[int64]ConnectedClient),
		fanoutChan: make(chan krt.Event[ConnectedClient], 100),
	}
	go func() {
	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			case event := <-cb.fanoutChan:
				cb.eventHandlers.Notify(event)
			}
		}
		// nil out the channel, incase notify is called after the context is done
		cb.stateLock.Lock()
		cb.fanoutChan = nil
		cb.stateLock.Unlock()
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
	x.stateLock.Lock()
	c, ok := x.clients[sid]
	delete(x.clients, sid)
	x.stateLock.Unlock()

	if ok {
		x.Notify(krt.Event[ConnectedClient]{Old: &c, Event: controllers.EventDelete})
	}
}

func (x *callbacksCollection) add(sid int64, r *envoy_service_discovery_v3.DiscoveryRequest) *ConnectedClient {
	x.stateLock.Lock()
	defer x.stateLock.Unlock()

	if _, ok := x.clients[sid]; !ok && r.Node != nil {
		c := ConnectedClient{
			Node: r.Node,
		}
		x.clients[sid] = c

		return &c
	}
	return nil
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
	c := x.add(sid, r)
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

/////////////////////////////////////////////////// collection of unique clients (i.e. same gateway and same labels)

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

// THIS IS THE SET OF THINGS WE RUN TRANSLATION FOR
func NewUniqlyConnectedClientCollection(c krt.Collection[ConnectedClient]) krt.Collection[UniqlyConnectedClient] {
	if false {
		// ideally i want to some something like this, but it is not currently supported.
		// as the mapping must by 1:1.
		// i.e. N->N. This mapping is N->M (where N >= M). as we may de-duplicated
		return krt.NewCollection[ConnectedClient, UniqlyConnectedClient](
			c,
			func(ctx krt.HandlerContext, cc ConnectedClient) *UniqlyConnectedClient {
				ucc := NewUniqlyConnectedClient(cc)
				return &ucc
			},
		)
	} else {
		// this works, but is less efficient, as any connected client change will trigger a recompute
		return krt.NewManyFromNothing[UniqlyConnectedClient](
			c,
			func(ctx krt.HandlerContext) []UniqlyConnectedClient {
				unqiueClients := make(map[string]struct{})
				var ret []UniqlyConnectedClient
				for _, cc := range krt.Fetch(ctx, c) {
					ucc := NewUniqlyConnectedClient(cc)
					if _, ok := unqiueClients[ucc.resourceName]; ok {
						continue
					}
					ret = append(ret, ucc)
				}
				return ret
			})
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
		hasher.Reset()
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
