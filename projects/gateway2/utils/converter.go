package utils

import (
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

func ToListener(listenerEntry gwxv1a1.ListenerEntry) gwv1.Listener {
	copy := listenerEntry.DeepCopy()
	return gwv1.Listener{
		Name:          copy.Name,
		Hostname:      copy.Hostname,
		Port:          copy.Port,
		Protocol:      copy.Protocol,
		TLS:           copy.TLS,
		AllowedRoutes: copy.AllowedRoutes,
	}
}

func ToListenerSlice(listenerEntries []gwxv1a1.ListenerEntry) []gwv1.Listener {
	listeners := make([]gwv1.Listener, len(listenerEntries))
	for _, l := range listenerEntries {
		listeners = append(listeners, ToListener(l))
	}
	return listeners
}
