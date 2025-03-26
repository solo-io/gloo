package types

import (
	"fmt"

	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

type ConsolidatedGateway struct {
	Gateway      *gwv1.Gateway
	ListenerSets []*gwxv1a1.XListenerSet
}

func (cgw *ConsolidatedGateway) GetListeners(ls *gwxv1a1.XListenerSet) []gwv1.Listener {
	var listeners []gwv1.Listener

	toListener := func(le *gwxv1a1.ListenerEntry) *gwv1.Listener {
		copy := le.DeepCopy()
		return &gwv1.Listener{
			Name:          copy.Name,
			Hostname:      copy.Hostname,
			Port:          copy.Port,
			Protocol:      copy.Protocol,
			TLS:           copy.TLS,
			AllowedRoutes: copy.AllowedRoutes,
		}
	}

	for _, l := range ls.Spec.Listeners {
		listeners = append(listeners, *toListener(&l))
	}

	return listeners
}

type ConsolidatedListeners struct {
	GatewayListeners     []gwv1.Listener
	ListenerSetListeners map[string][]gwv1.Listener
}

func (cl *ConsolidatedListeners) SetListenerSetListeners(ls *gwxv1a1.XListenerSet, listeners []gwv1.Listener) {
	if cl.ListenerSetListeners == nil {
		cl.ListenerSetListeners = make(map[string][]gwv1.Listener)
	}
	cl.ListenerSetListeners[cl.ListenerSetListenersKey(ls)] = listeners
}

func (cl *ConsolidatedListeners) GetListenerSetListeners(ls *gwxv1a1.XListenerSet) []gwv1.Listener {
	if cl.ListenerSetListeners == nil {
		return nil
	}
	return cl.ListenerSetListeners[cl.ListenerSetListenersKey(ls)]
}

func (cl *ConsolidatedListeners) ListenerSetListenersKey(ls *gwxv1a1.XListenerSet) string {
	return fmt.Sprintf("%s/%s", ls.GetNamespace(), ls.GetName())
}
