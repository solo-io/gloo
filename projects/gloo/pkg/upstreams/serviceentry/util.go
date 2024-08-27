package serviceentry

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func filterSlice[T any](items []T, f func(T) bool) []T {
	var result []T
	for _, item := range items {
		if f(item) {
			result = append(result, item)
		}
	}
	return result
}

func targetPort(base, svc, ep uint32) uint32 {
	if ep > 0 {
		return ep
	}
	if svc > 0 {
		return svc
	}
	return base
}

// parseIstioProtocol always gives the all-caps protocol part of a port name.
// Example: http-foobar would be HTTP.
func parseIstioProtocol(protocol string) string {
	protocol = strings.ToUpper(protocol)
	if idx := strings.Index(protocol, "-"); idx != -1 {
		protocol = protocol[:idx]
	}
	return protocol
}

func isProtocolTLS(protocol string) bool {
	p := parseIstioProtocol(protocol)
	return p == "HTTPS" || p == "TLS"
}

func buildListOptions(expressionSelector string, selector map[string]string) metav1.ListOptions {
	if expressionSelector != "" {
		return metav1.ListOptions{
			LabelSelector: expressionSelector,
		}
	}
	sel := labels.NewSelector()
	for k, v := range selector {
		req, _ := labels.NewRequirement(k, selection.Equals, strings.Split(v, ","))
		sel = sel.Add(*req)
	}
	return metav1.ListOptions{
		LabelSelector: sel.String(),
	}
}
