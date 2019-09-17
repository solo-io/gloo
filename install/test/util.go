package test

import (
	glooTest "github.com/solo-io/gloo/install/test"
	"github.com/solo-io/go-utils/manifesttestutils"
)

func GetGlooWithReadOnlyUiServiceAccountPermissions(namespace string) *manifesttestutils.ServiceAccountPermissions {

	// build off of the permissions imported from Gloo
	permissions := glooTest.GetServiceAccountPermissions(namespace)

	// Apiserver
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{""},
		[]string{"pods", "services", "configmaps", "namespaces", "secrets"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"apiextensions.k8s.io"},
		[]string{"customresourcedefinitions"},
		[]string{"get"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gloo.solo.io"},
		[]string{"artifacts", "upstreams", "upstreamgroups", "proxies", "secrets"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gloo.solo.io"},
		[]string{"settings"},
		[]string{"get", "list", "watch", "create"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gateway.solo.io.v2"},
		[]string{"gateways"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gateway.solo.io"},
		[]string{"virtualservices"},
		[]string{"get", "list", "watch"})

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
		[]string{"list"})

	return permissions
}
func GetGlooEServiceAccountPermissions(namespace string) *manifesttestutils.ServiceAccountPermissions {

	// build off of the permissions imported from Gloo
	permissions := glooTest.GetServiceAccountPermissions(namespace)

	// Apiserver
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{""},
		[]string{"pods", "services", "configmaps", "namespaces", "secrets"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{""},
		[]string{"secrets"},
		[]string{"get", "list", "watch", "create", "update", "delete"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"apiextensions.k8s.io"},
		[]string{"customresourcedefinitions"},
		[]string{"get", "create"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gloo.solo.io"},
		[]string{"artifacts", "upstreams", "upstreamgroups", "proxies", "secrets"},
		[]string{"create", "delete", "update", "get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gloo.solo.io"},
		[]string{"settings"},
		[]string{"get", "list", "watch", "create", "update", "delete"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gateway.solo.io"},
		[]string{"routetables"},
		[]string{"list"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gateway.solo.io.v2"},
		[]string{"gateways"},
		[]string{"get", "list", "watch", "create", "update", "delete"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gateway.solo.io"},
		[]string{"virtualservices"},
		[]string{"get", "list", "watch", "create", "update", "delete"})

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
		[]string{"list"})

	return permissions
}
