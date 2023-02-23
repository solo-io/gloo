package test

import (
	"strings"

	"github.com/solo-io/solo-projects/pkg/install"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/k8s-utils/manifesttestutils"
)

var _ = Describe("RBAC Test", func() {
	Context("GlooE", func() {
		Context("implementation-agnostic permissions", func() {
			It("correctly assigns permissions for single-namespace gloo", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"global.glooRbac.namespaced=true", "gloo-fed.enabled=true"},
				})
				Expect(err).NotTo(HaveOccurred())
				permissions := GetGlooEServiceAccountPermissions("gloo-system")
				testManifest.ExpectPermissions(permissions)
			})

			It("correctly assigns permissions for cluster-scoped gloo", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"global.glooRbac.namespaced=false", "gloo-fed.enabled=true"},
				})
				Expect(err).NotTo(HaveOccurred())
				permissions := GetGlooEServiceAccountPermissions("")
				testManifest.ExpectPermissions(permissions)
			})

			It("correctly assigns permissions for gloo-fed.glooFedApiserver.namespaceRestrictedMode=true", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"gloo.settings.singleNamespace=true",
						"gloo-fed.enabled=false",
						"gloo-fed.glooFedApiserver.namespaceRestrictedMode=true"},
				})
				Expect(err).NotTo(HaveOccurred())
				testManifest.Expect("RoleBinding", namespace, "gloo-console")
				testManifest.Expect("RoleBinding", namespace, "gloo-console-envoy")
				testManifest.Expect("Role", namespace, "gloo-console")
				testManifest.Expect("Role", namespace, "gloo-console-envoy")
			})

			It("correctly assigns permissions for gloo-fed.glooFedApiserver.namespaceRestrictedMode=false", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"gloo.settings.singleNamespace=true",
						"gloo-fed.enabled=false",
						"gloo-fed.glooFedApiserver.namespaceRestrictedMode=false"},
				})
				Expect(err).NotTo(HaveOccurred())
				testManifest.Expect("ClusterRoleBinding", namespace, "gloo-console")
				testManifest.Expect("RoleBinding", namespace, "gloo-console")
				testManifest.Expect("ClusterRole", namespace, "gloo-console")
				testManifest.Expect("Role", namespace, "gloo-console")
			})

			It("creates no permissions when rbac is disabled", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"gloo-fed.rbac.create=false",
						"prometheus.rbac.create=false",
						"prometheus.kube-state-metrics.rbac.create=false",
						"grafana.testFramework.enabled=false",
						"global.glooRbac.create=false",
						"gloo-fed.enabled=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				testManifest.ExpectAll(func(resource *unstructured.Unstructured) {
					Expect(resource.GetAPIVersion()).NotTo(ContainSubstring("rbac.authorization.k8s.io"), "Should not contain the RBAC API group")
				})
			})
		})

		Context("all cluster-scoped RBAC resources", func() {

			checkSuffix := func(testManifest TestManifest, suffix string) {
				rbacResources := testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return (resource.GetKind() == "ClusterRole" || resource.GetKind() == "ClusterRoleBinding") &&
						!strings.Contains(resource.GetName(), "glooe-grafana") &&
						!strings.Contains(resource.GetName(), "glooe-prometheus") &&
						!strings.Contains(resource.GetName(), "gloo-fed")
				})

				Expect(rbacResources.NumResources()).NotTo(BeZero())

				rbacResources.ExpectAll(func(resource *unstructured.Unstructured) {
					Expect(resource.GetName()).To(HaveSuffix("-" + suffix))
				})
			}

			It("is all named appropriately when a custom suffix is specified", func() {
				suffix := "test-suffix"
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.glooRbac.nameSuffix=" + suffix,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				checkSuffix(testManifest, suffix)
			})

			It("is all named appropriately in a non-namespaced install", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
				Expect(err).NotTo(HaveOccurred())
				checkSuffix(testManifest, namespace)
			})
		})
	})

})
