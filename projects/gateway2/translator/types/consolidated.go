package types

import (
	"fmt"

	"github.com/solo-io/gloo/projects/gateway2/utils"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

type ConsolidatedGateway struct {
	Gateway               *gwv1.Gateway
	AllowedListenerSets   []*gwxv1a1.XListenerSet
	DeniedListenerSets    []*gwxv1a1.XListenerSet
	ListenerSetsListeners map[string][]gwv1.Listener
}

func (cgw *ConsolidatedGateway) GetListeners(ls *gwxv1a1.XListenerSet) []gwv1.Listener {

	if cgw.ListenerSetsListeners == nil {
		cgw.ListenerSetsListeners = make(map[string][]gwv1.Listener)
	}

	if listeners, ok := cgw.ListenerSetsListeners[generateListenerSetKey(ls)]; ok {
		return listeners
	}

	cgw.ListenerSetsListeners[generateListenerSetKey(ls)] = utils.ToListenerSlice(ls.Spec.Listeners)
	return cgw.ListenerSetsListeners[generateListenerSetKey(ls)]
}

type ConsolidatedListeners struct {
	GatewayListeners     []gwv1.Listener
	ListenerSetListeners map[string][]gwv1.Listener
}

func (cl *ConsolidatedListeners) SetListenerSetListeners(ls *gwxv1a1.XListenerSet, listeners []gwv1.Listener) {
	if cl.ListenerSetListeners == nil {
		cl.ListenerSetListeners = make(map[string][]gwv1.Listener)
	}
	cl.ListenerSetListeners[generateListenerSetKey(ls)] = listeners
}

func (cl *ConsolidatedListeners) GetListenerSetListeners(ls *gwxv1a1.XListenerSet) []gwv1.Listener {
	if cl.ListenerSetListeners == nil {
		return nil
	}
	return cl.ListenerSetListeners[generateListenerSetKey(ls)]
}

func generateListenerSetKey(ls *gwxv1a1.XListenerSet) string {
	return fmt.Sprintf("%s/%s", ls.GetNamespace(), ls.GetName())
}
