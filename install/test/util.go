package test

import (
	glooTest "github.com/solo-io/gloo/install/test"
	"github.com/solo-io/k8s-utils/manifesttestutils"
)

func GetGlooEServiceAccountPermissions(namespace string) *manifesttestutils.ServiceAccountPermissions {
	// build off of the permissions imported from Gloo
	permissions := glooTest.GetServiceAccountPermissions(namespace)
	ApplyPermissionsForGlooEServiceAccounts(namespace, permissions)
	ApplyPermissionsForPrometheusServiceAccounts(permissions)
	return permissions
}

func ApplyPermissionsForGlooEServiceAccounts(namespace string, permissions *manifesttestutils.ServiceAccountPermissions) {
	// Observability
	permissions.AddExpectedPermission(
		"gloo-system.observability",
		namespace,
		[]string{"gloo.solo.io"},
		[]string{"settings"},
		[]string{"create", "get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.observability",
		namespace,
		[]string{"gloo.solo.io"},
		[]string{"upstreams"},
		[]string{"get", "list", "watch"})
}

func ApplyPermissionsForPrometheusServiceAccounts(permissions *manifesttestutils.ServiceAccountPermissions) {
	// Prometheus
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-server",
		"",
		[]string{""},
		[]string{"configmaps"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-server",
		"",
		[]string{""},
		[]string{"endpoints"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-server",
		"",
		[]string{""},
		[]string{"ingresses"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-server",
		"",
		[]string{""},
		[]string{"nodes"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-server",
		"",
		[]string{""},
		[]string{"nodes/metrics"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-server",
		"",
		[]string{""},
		[]string{"nodes/proxy"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-server",
		"",
		[]string{""},
		[]string{"pods"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-server",
		"",
		[]string{""},
		[]string{"services"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-server",
		"",
		[]string{"extensions"},
		[]string{"ingresses"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-server",
		"",
		[]string{"extensions"},
		[]string{"ingresses/status"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-server",
		"",
		[]string{"networking.k8s.io"},
		[]string{"ingresses"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-server",
		"",
		[]string{"networking.k8s.io"},
		[]string{"ingresses/status"},
		[]string{"get", "list", "watch"})

	// Kube state metrics
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{""},
		[]string{"configmaps"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{""},
		[]string{"endpoints"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{""},
		[]string{"limitranges"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{""},
		[]string{"namespaces"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{""},
		[]string{"nodes"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{""},
		[]string{"persistentvolumeclaims"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{""},
		[]string{"persistentvolumes"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{""},
		[]string{"pods"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{""},
		[]string{"replicationcontrollers"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{""},
		[]string{"resourcequotas"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{""},
		[]string{"secrets"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{""},
		[]string{"services"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"admissionregistration.k8s.io"},
		[]string{"mutatingwebhookconfigurations"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"admissionregistration.k8s.io"},
		[]string{"validatingwebhookconfigurations"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"apps"},
		[]string{"daemonsets"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"apps"},
		[]string{"deployments"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"apps"},
		[]string{"replicasets"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"apps"},
		[]string{"statefulsets"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"apps"},
		[]string{"replicasets"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"autoscaling"},
		[]string{"horizontalpodautoscalers"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"batch"},
		[]string{"cronjobs"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"batch"},
		[]string{"jobs"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"certificates.k8s.io"},
		[]string{"certificatesigningrequests"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"extensions"},
		[]string{"daemonsets"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"extensions"},
		[]string{"deployments"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"extensions"},
		[]string{"ingresses"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"extensions"},
		[]string{"replicasets"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"apps"},
		[]string{"replicasets"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"networking.k8s.io"},
		[]string{"ingresses"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"networking.k8s.io"},
		[]string{"networkpolicies"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"policy"},
		[]string{"poddisruptionbudgets"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"storage.k8s.io"},
		[]string{"storageclasses"},
		[]string{"list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.glooe-prometheus-kube-state-metrics",
		"",
		[]string{"storage.k8s.io"},
		[]string{"volumeattachments"},
		[]string{"list", "watch"})
}
