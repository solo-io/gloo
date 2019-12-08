package test

import (
	"strings"

	"github.com/solo-io/solo-projects/pkg/install"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/go-utils/manifesttestutils"
)

var _ = Describe("RBAC Test", func() {
	Context("GlooE", func() {
		Context("implementation-agnostic permissions", func() {
			It("correctly assigns permissions for single-namespace gloo", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"global.glooRbac.namespaced=true"},
				})
				Expect(err).NotTo(HaveOccurred())
				permissions := GetGlooEServiceAccountPermissions("gloo-system")
				testManifest.ExpectPermissions(permissions)
			})

			It("correctly assigns permissions for cluster-scoped gloo", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"global.glooRbac.namespaced=false"},
				})
				Expect(err).NotTo(HaveOccurred())
				permissions := GetGlooEServiceAccountPermissions("")
				testManifest.ExpectPermissions(permissions)
			})

			It("creates no permissions when rbac is disabled", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"prometheus.rbac.create=false",
						"grafana.testFramework.enabled=false",
						"global.glooRbac.create=false",
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
						!strings.Contains(resource.GetName(), "glooe-prometheus")
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

	Context("Gloo OS with read-only UI", func() {
		Context("implementation-agnostic permissions", func() {
			It("correctly assigns permissions for single-namespace gloo", func() {
				testManifest, err := BuildTestManifest(install.GlooOsWithUiChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.glooRbac.namespaced=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				permissions := GetGlooWithReadOnlyUiServiceAccountPermissions("gloo-system")
				testManifest.ExpectPermissions(permissions)
			})

			It("correctly assigns permissions for cluster-scoped gloo", func() {
				testManifest, err := BuildTestManifest(install.GlooOsWithUiChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.glooRbac.namespaced=false",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				permissions := GetGlooWithReadOnlyUiServiceAccountPermissions("")
				testManifest.ExpectPermissions(permissions)
			})
		})
	})
})
