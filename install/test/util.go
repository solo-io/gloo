package test

import (
	. "github.com/onsi/gomega"
	glooTest "github.com/solo-io/gloo/install/test"
	"github.com/solo-io/k8s-utils/manifesttestutils"
	v1 "k8s.io/api/core/v1"
)

// ExpectContainer is used for helm containers. This will help with finding and expecting args and environments
type ExpectContainer struct {
	Containers []v1.Container
	Name       string
}

func (c *ExpectContainer) ExpectToHaveArg(arg, errorMsg string) {
	Expect(c.hasArgument(arg)).To(BeTrue(), errorMsg)
}

func (c *ExpectContainer) ExpectToHaveEnv(envName, envValue, errorMsg string) {
	Expect(c.hasEnvVar(envName, envValue)).To(BeTrue(), errorMsg)
}

func (c *ExpectContainer) ExpectToNotHaveEnv(envName, errorMsg string) {
	Expect(c.doesNotHaveEnvVar(envName)).To(BeTrue(), errorMsg)
}

func (ec *ExpectContainer) ExpectToHaveVolumeMount(name string) *v1.VolumeMount {
	vm := ec.getVolumeMount(name)
	Expect(vm).NotTo(BeNil())
	return vm
}

func (ec *ExpectContainer) ExpectToNotHaveVolumeMount(name string) *v1.VolumeMount {
	vm := ec.getVolumeMount(name)
	Expect(vm).To(BeNil())
	return vm
}

func (ec *ExpectContainer) getContainer() *v1.Container {
	for _, c := range ec.Containers {
		if c.Name == ec.Name {
			return &c
		}
	}
	return nil
}

func (ec *ExpectContainer) getVolumeMount(name string) *v1.VolumeMount {
	c := ec.getContainer()
	for _, vm := range c.VolumeMounts {
		if vm.Name == name {
			return &vm
		}
	}
	return nil
}

func (ec *ExpectContainer) hasArgument(arg string) bool {
	c := ec.getContainer()
	if c == nil {
		return false
	}
	Expect(c.Args).To(ContainElement(arg))
	return true
}

func (ec *ExpectContainer) hasEnvVar(env, value string) bool {
	c := ec.getContainer()
	if c == nil {
		return false
	}
	Expect(c.Env).To(ContainElement(v1.EnvVar{Name: env, Value: value}))
	return true
}

func (ec *ExpectContainer) doesNotHaveEnvVar(env string) bool {
	c := ec.getContainer()
	if c == nil {
		return true
	}
	for _, e := range c.Env {
		if e.Name == env {
			Expect(e.Name).ToNot(Equal(env))
			return false
		}
	}
	return true
}

type ExpectVolume struct {
	Volumes []v1.Volume
}

func (ev *ExpectVolume) ExpectHasName(name string) {
	Expect(ev.getByName(name)).ToNot(BeNil())
}

func (ev *ExpectVolume) getByName(name string) *v1.Volume {
	for _, v := range ev.Volumes {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func GetGlooEServiceAccountPermissions(namespace string) *manifesttestutils.ServiceAccountPermissions {
	// build off of the permissions imported from Gloo
	permissions := glooTest.GetServiceAccountPermissions(namespace)
	ApplyPermissionsForGlooEServiceAccounts(namespace, permissions)
	ApplyPermissionsForPrometheusServiceAccounts(permissions)
	ApplyPermissionsForGlooFedServiceAccounts(permissions)
	ApplyPermissionsForGlooFedConsoleServiceAccounts(permissions)
	return permissions
}

func ApplyPermissionsForGlooEServiceAccounts(namespace string, permissions *manifesttestutils.ServiceAccountPermissions) {
	// Observability
	permissions.AddExpectedPermission(
		"gloo-system.observability",
		namespace,
		[]string{"gloo.solo.io"},
		[]string{"settings"},
		[]string{"get", "list", "watch"})
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

func ApplyPermissionsForGlooFedServiceAccounts(permissions *manifesttestutils.ServiceAccountPermissions) {
	permissions.AddExpectedPermission("gloo-system.gloo-fed",
		"",
		[]string{""},
		[]string{"secrets"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed",
		"",
		[]string{"apps"},
		[]string{"deployments"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed",
		"",
		[]string{"gloo.solo.io"},
		[]string{"settings"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed",
		"",
		[]string{"fed.enterprise.gloo.solo.io"},
		[]string{"federatedauthconfigs", "federatedauthconfigs/status"},
		[]string{"*"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed",
		"",
		[]string{"fed.gateway.solo.io"},
		[]string{
			"federatedgateways",
			"federatedgateways/status",
			"federatedmatchablehttpgateways",
			"federatedmatchablehttpgateways/status",
			"federatedroutetables",
			"federatedroutetables/status",
			"federatedvirtualservices",
			"federatedvirtualservices/status",
		},
		[]string{"*"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed",
		"",
		[]string{"fed.gloo.solo.io"},
		[]string{
			"federatedauthconfigs",
			"federatedauthconfigs/status",
			"federatedsettings",
			"federatedsettings/status",
			"federatedupstreamgroups",
			"federatedupstreamgroups/status",
			"federatedupstreams",
			"federatedupstreams/status",
		},
		[]string{"*"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed",
		"",
		[]string{"fed.ratelimit.solo.io"},
		[]string{"federatedratelimitconfigs", "federatedratelimitconfigs/status"},
		[]string{"*"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed",
		"",
		[]string{"fed.solo.io"},
		[]string{"failoverschemes"},
		[]string{"*"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed",
		"",
		[]string{"fed.solo.io"},
		[]string{"failoverschemes/status"},
		[]string{"get", "update"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed",
		"",
		[]string{"fed.solo.io"},
		[]string{"glooinstances"},
		[]string{"*"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed",
		"",
		[]string{"fed.solo.io"},
		[]string{"glooinstances/status"},
		[]string{"get", "update"})
}

func ApplyPermissionsForGlooFedConsoleServiceAccounts(permissions *manifesttestutils.ServiceAccountPermissions) {
	permissions.AddExpectedPermission("gloo-system.gloo-fed-console",
		"",
		[]string{""},
		[]string{"secrets"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed-console",
		"",
		[]string{"apps"},
		[]string{"deployments"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed-console",
		"",
		[]string{"gloo.solo.io"},
		[]string{"settings"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed-console",
		"",
		[]string{"fed.enterprise.gloo.solo.io"},
		[]string{"federatedauthconfigs", "federatedauthconfigs/status"},
		[]string{"*"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed-console",
		"",
		[]string{"fed.gateway.solo.io"},
		[]string{
			"federatedgateways",
			"federatedgateways/status",
			"federatedmatchablehttpgateways",
			"federatedmatchablehttpgateways/status",
			"federatedroutetables",
			"federatedroutetables/status",
			"federatedvirtualservices",
			"federatedvirtualservices/status",
		},
		[]string{"*"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed-console",
		"",
		[]string{"fed.gloo.solo.io"},
		[]string{
			"federatedauthconfigs",
			"federatedauthconfigs/status",
			"federatedsettings",
			"federatedsettings/status",
			"federatedupstreamgroups",
			"federatedupstreamgroups/status",
			"federatedupstreams",
			"federatedupstreams/status",
		},
		[]string{"*"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed-console",
		"",
		[]string{"fed.ratelimit.solo.io"},
		[]string{"federatedratelimitconfigs", "federatedratelimitconfigs/status"},
		[]string{"*"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed-console",
		"",
		[]string{"fed.solo.io"},
		[]string{"failoverschemes"},
		[]string{"*"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed-console",
		"",
		[]string{"fed.solo.io"},
		[]string{"failoverschemes/status"},
		[]string{"get", "update"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed-console",
		"",
		[]string{"fed.solo.io"},
		[]string{"glooinstances"},
		[]string{"*"})
	permissions.AddExpectedPermission("gloo-system.gloo-fed-console",
		"",
		[]string{"fed.solo.io"},
		[]string{"glooinstances/status"},
		[]string{"get", "update"})
}
