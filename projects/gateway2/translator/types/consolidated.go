package types

import (
	"fmt"
	"slices"

	"github.com/solo-io/gloo/projects/gateway2/reports"
	"github.com/solo-io/gloo/projects/gateway2/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

type ConsolidatedGateway struct {
	Gateway               *gwv1.Gateway
	AllowedListenerSets   []*gwxv1a1.XListenerSet
	DeniedListenerSets    []*gwxv1a1.XListenerSet
	consolidatedListeners []ConsolidatedListener
	listenerSetsListeners map[string][]gwv1.Listener
}

func (cgw *ConsolidatedGateway) GetListeners(ls *gwxv1a1.XListenerSet) []gwv1.Listener {

	if cgw.listenerSetsListeners == nil {
		cgw.listenerSetsListeners = make(map[string][]gwv1.Listener)
	}

	if listeners, ok := cgw.listenerSetsListeners[generateListenerSetKey(ls)]; ok {
		return listeners
	}

	if !slices.Contains(cgw.AllowedListenerSets, ls) {
		cgw.listenerSetsListeners[generateListenerSetKey(ls)] = nil
		return nil
	}

	cgw.listenerSetsListeners[generateListenerSetKey(ls)] = utils.ToListenerSlice(ls.Spec.Listeners)
	return cgw.listenerSetsListeners[generateListenerSetKey(ls)]
}

func (cgw *ConsolidatedGateway) GetConsolidatedListeners() []ConsolidatedListener {
	if cgw.consolidatedListeners == nil {
		var consolidatedListeners []ConsolidatedListener
		for _, listener := range cgw.Gateway.Spec.Listeners {
			consolidatedListeners = append(consolidatedListeners, ConsolidatedListener{
				Listener:    &listener,
				Gateway:     cgw.Gateway,
				ListenerSet: nil,
			})
		}
		for _, ls := range cgw.AllowedListenerSets {
			for _, listener := range cgw.GetListeners(ls) {
				consolidatedListeners = append(consolidatedListeners, ConsolidatedListener{
					Listener:    &listener,
					Gateway:     cgw.Gateway,
					ListenerSet: ls,
				})
			}
		}
		cgw.consolidatedListeners = consolidatedListeners
	}
	return cgw.consolidatedListeners
}

type ConsolidatedListener struct {
	Listener    *gwv1.Listener
	Gateway     *gwv1.Gateway
	ListenerSet *gwxv1a1.XListenerSet
}

func (cl ConsolidatedListener) GetParentReporter(reporter reports.Reporter) reports.GatewayReporter {
	parentReporter := reporter.Gateway(cl.Gateway)
	if cl.ListenerSet != nil {
		parentReporter = reporter.ListenerSet(cl.ListenerSet)
	}
	return parentReporter
}

func (cl ConsolidatedListener) GetParent() client.Object {
	if cl.ListenerSet != nil {
		return cl.ListenerSet
	}
	return cl.Gateway
}

func generateListenerSetKey(ls *gwxv1a1.XListenerSet) string {
	return fmt.Sprintf("%s/%s", ls.GetNamespace(), ls.GetName())
}
