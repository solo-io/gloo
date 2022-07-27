package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"unicode"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	. "github.com/solo-io/solo-kit/test/matchers"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/aws/aws-sdk-go/aws"
	bootstrapv3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/config/health_checker/redis/v2"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/redis_proxy/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/health_checkers/redis/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/gomega/format"

	"github.com/golang/protobuf/proto"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/k8s-utils/installutils/kuberesource"
	"github.com/solo-io/solo-projects/pkg/install"
	jobsv1 "k8s.io/api/batch/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8s "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/k8s-utils/manifesttestutils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

var _ = Describe("Helm Test", func() {
	var (
		version string

		backtick = "`"

		normalPromAnnotations = map[string]string{
			"prometheus.io/port":   "9091",
			"prometheus.io/scrape": "true",
			"prometheus.io/path":   "/metrics",
		}

		statsEnvVar = v1.EnvVar{
			Name:  "START_STATS_SERVER",
			Value: "true",
		}
	)

	Describe("gloo-ee helm tests", func() {
		var (
			labels        map[string]string
			selector      map[string]string
			getPullPolicy func() v1.PullPolicy
			manifestYaml  string
		)

		BeforeEach(func() {
			// Disables truncation during error output so we can see the full error message
			format.MaxLength = 0
			version = os.Getenv("TAGGED_VERSION")
			if version == "" {
				version = os.Getenv("VERSION")
				if version == "" {
					version = "0.0.0-dev"
					getPullPolicy = func() v1.PullPolicy { return v1.PullAlways }
				} else {
					fmt.Printf("Using VERSION environment variable for version: %s\n", version)
					getPullPolicy = func() v1.PullPolicy { return v1.PullIfNotPresent }
				}
			} else {
				fmt.Printf("Using TAGGED_VERSION environment variable for version: %s\n", version)
				version = version[1:]
				getPullPolicy = func() v1.PullPolicy { return v1.PullIfNotPresent }
			}
			manifestYaml = ""
		})

		AfterEach(func() {
			if manifestYaml != "" {
				err := os.Remove(manifestYaml)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		Context("observability", func() {
			var (
				observabilityDeployment *appsv1.Deployment
				grafanaDeployment       *appsv1.Deployment
			)
			BeforeEach(func() {
				labels = map[string]string{
					"app":  "gloo",
					"gloo": "observability",
				}
				selector = map[string]string{
					"app":  "gloo",
					"gloo": "observability",
				}

				rb := ResourceBuilder{
					Namespace: namespace,
					Name:      "observability",
					Labels:    labels,
				}

				nonRootUser := int64(10101)
				nonRoot := true

				nonRootSC := &v1.PodSecurityContext{
					RunAsUser:    &nonRootUser,
					RunAsNonRoot: &nonRoot,
				}

				observabilityDeployment = rb.GetDeploymentAppsv1()

				observabilityDeployment.Spec.Template.Spec.Volumes = []v1.Volume{
					{
						Name: "upstream-dashboard-template",
						VolumeSource: v1.VolumeSource{
							ConfigMap: &v1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{Name: "glooe-observability-config"},
								Items: []v1.KeyToPath{
									{
										Key:  "UPSTREAM_DASHBOARD_JSON_TEMPLATE",
										Path: "dashboard-template.json",
									},
								},
							},
						},
					},
					{
						Name: "custom-dashboards",
						VolumeSource: v1.VolumeSource{
							ConfigMap: &v1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{Name: "glooe-grafana-custom-dashboards-v2"},
								Items: []v1.KeyToPath{
									{
										Key:  "envoy.json",
										Path: "envoy.json",
									},
									{
										Key:  "extauth-monitoring.json",
										Path: "extauth-monitoring.json",
									},
									{
										Key:  "gloo-overview.json",
										Path: "gloo-overview.json",
									},
									{
										Key:  "kubernetes.json",
										Path: "kubernetes.json",
									},
									{
										Key:  "upstreams.json",
										Path: "upstreams.json",
									},
								},
							},
						},
					},
				}

				observabilityDeployment.Spec.Template.Spec.Containers = []v1.Container{
					{
						Name:  "observability",
						Image: "quay.io/solo-io/observability-ee:" + version,
						EnvFrom: []v1.EnvFromSource{
							{ConfigMapRef: &v1.ConfigMapEnvSource{LocalObjectReference: v1.LocalObjectReference{Name: "glooe-observability-config"}}},
							{SecretRef: &v1.SecretEnvSource{LocalObjectReference: v1.LocalObjectReference{Name: "glooe-observability-secrets"}}},
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "upstream-dashboard-template",
								ReadOnly:  true,
								MountPath: "/observability",
							},
							{
								Name:      "custom-dashboards",
								ReadOnly:  true,
								MountPath: "/observability/defaults",
							},
						},
						Env: []v1.EnvVar{
							{
								Name: "GLOO_LICENSE_KEY",
								ValueFrom: &v1.EnvVarSource{
									SecretKeyRef: &v1.SecretKeySelector{
										LocalObjectReference: v1.LocalObjectReference{
											Name: "license",
										},
										Key: "license-key",
									},
								},
							},
							{
								Name: "POD_NAMESPACE",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										FieldPath: "metadata.namespace",
									},
								},
							},
							statsEnvVar,
						},
						Resources:       v1.ResourceRequirements{},
						ImagePullPolicy: getPullPolicy(),
					},
				}
				observabilityDeployment.Spec.Template.Spec.ServiceAccountName = "observability"
				observabilityDeployment.Spec.Strategy = appsv1.DeploymentStrategy{}
				observabilityDeployment.Spec.Selector.MatchLabels = selector
				observabilityDeployment.Spec.Template.ObjectMeta.Labels = selector
				annotations := map[string]string{
					"checksum/observability-config": "9d91255a98e28f9b0bfb5d685673e5810fc475d2fe6f9738aae7dd67c8ac5c8d", // observability config checksum
					"checksum/grafana-dashboards":   "1e927634c33379005380b99e746c141d0fa241bf42246305bf9eb6ea29ca6383", // grafana dashboards checksum
				}
				for key, val := range normalPromAnnotations { // deep copy map
					annotations[key] = val
				}
				observabilityDeployment.Spec.Template.ObjectMeta.Annotations = annotations

				observabilityDeployment.Spec.Template.Spec.SecurityContext = nonRootSC
				observabilityDeployment.Spec.Replicas = nil // GetDeploymentAppsv1 explicitly sets it to 1, which we don't want

				grafanaBuilder := ResourceBuilder{
					Namespace: "", // grafana installs to empty namespace during tests
					Name:      "release-name-grafana",
					Labels:    labels,
				}
				grafanaDeployment = grafanaBuilder.GetDeploymentAppsv1()
			})

			It("has valid default dashboards", func() {
				dashboardsDir := "../helm/gloo-ee/dashboards/"
				files, err := ioutil.ReadDir(dashboardsDir)
				Expect(err).NotTo(HaveOccurred(), "Should be able to list files")
				Expect(files).NotTo(HaveLen(0), "Should have dashboard files")
				for _, f := range files {
					if !strings.HasSuffix(f.Name(), ".json") {
						continue // not a JSON file
					}
					bytes, err := ioutil.ReadFile(path.Join(dashboardsDir, f.Name()))
					Expect(err).NotTo(HaveOccurred(), "Should be able to read the Envoy dashboard json file")
					err = json.Unmarshal(bytes, &map[string]interface{}{})
					Expect(err).NotTo(HaveOccurred(), "Should be able to successfully unmarshal the envoy dashboard json")
				}
			})

			It("has valid v2 default dashboards", func() {
				dashboardsDir := "../helm/gloo-ee/dashboards/v2"
				files, err := ioutil.ReadDir(dashboardsDir)
				Expect(err).NotTo(HaveOccurred(), "Should be able to list files")
				Expect(files).NotTo(HaveLen(0), "Should have updated dashboard files")
				for _, f := range files {
					if !strings.HasSuffix(f.Name(), ".json") {
						continue // not a JSON file
					}
					bytes, err := ioutil.ReadFile(path.Join(dashboardsDir, f.Name()))
					Expect(err).NotTo(HaveOccurred(), "Should be able to read the Envoy dashboard json file")
					err = json.Unmarshal(bytes, &map[string]interface{}{})
					Expect(err).NotTo(HaveOccurred(), "Should be able to successfully unmarshal the envoy dashboard json")
				}
			})

			Context("observability deployment", func() {
				It("is installed by default", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
					Expect(err).NotTo(HaveOccurred())

					testManifest.ExpectDeploymentAppsV1(observabilityDeployment)
				})

				It("is not installed when grafana is disabled", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{"grafana.defaultInstallationEnabled=false"},
					})
					Expect(err).NotTo(HaveOccurred())

					testManifest.Expect(observabilityDeployment.Kind, observabilityDeployment.Namespace, observabilityDeployment.Name).To(BeNil())
				})

				It("is installed when a custom grafana instance is present", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{"observability.customGrafana.enabled=true"},
					})
					Expect(err).NotTo(HaveOccurred())

					testManifest.ExpectDeploymentAppsV1(observabilityDeployment)
				})

				It("should support running as arbitrary user", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{"observability.deployment.runAsUser=10102"},
					})
					Expect(err).NotTo(HaveOccurred())

					customUser := int64(10102)
					observabilityDeployment.Spec.Template.Spec.SecurityContext.RunAsUser = &customUser

					testManifest.ExpectDeploymentAppsV1(observabilityDeployment)
				})

				It("should support changing the number of replicas", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{"observability.deployment.replicas=2"},
					})
					Expect(err).NotTo(HaveOccurred())

					customNumReplicas := int32(2)
					observabilityDeployment.Spec.Replicas = &customNumReplicas

					testManifest.ExpectDeploymentAppsV1(observabilityDeployment)
				})

				It("should support setting the log level env", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{"observability.deployment.logLevel=DEBUG"},
					})
					Expect(err).NotTo(HaveOccurred())

					logLevel := "DEBUG"
					envs := observabilityDeployment.Spec.Template.Spec.Containers[0].Env
					for i, env := range envs {
						if env.Name == "LOG_LEVEL" {
							envs[i].Value = logLevel
						}
					}

					testManifest.ExpectDeploymentAppsV1(observabilityDeployment)
				})

				It("correctly sets the GLOO_LICENSE_KEY env", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{"gloo.license_secret_name=custom-license-secret"},
					})
					Expect(err).NotTo(HaveOccurred())

					licenseKeyEnvVarSource := v1.EnvVarSource{
						SecretKeyRef: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "custom-license-secret",
							},
							Key: "license-key",
						},
					}
					envs := observabilityDeployment.Spec.Template.Spec.Containers[0].Env
					for i, env := range envs {
						if env.Name == "GLOO_LICENSE_KEY" {
							envs[i].ValueFrom = &licenseKeyEnvVarSource
						}
					}
					testManifest.ExpectDeploymentAppsV1(observabilityDeployment)
				})

				It("correctly sets resource limits for the observability deployment", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{
							"observability.deployment.resources.limits.cpu=300m",
							"observability.deployment.resources.limits.memory=300Mi",
							"observability.deployment.resources.requests.cpu=30m",
							"observability.deployment.resources.requests.memory=30Mi",
						},
					})
					Expect(err).NotTo(HaveOccurred())

					// UI
					observabilityDeployment.Spec.Template.Spec.Containers[0].Resources = v1.ResourceRequirements{
						Limits: v1.ResourceList{
							v1.ResourceCPU:    resource.MustParse("300m"),
							v1.ResourceMemory: resource.MustParse("300Mi"),
						},
						Requests: v1.ResourceList{
							v1.ResourceCPU:    resource.MustParse("30m"),
							v1.ResourceMemory: resource.MustParse("30Mi"),
						},
					}

					testManifest.ExpectDeploymentAppsV1(observabilityDeployment)
				})

				It("Should have no duplicate resources", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{"observability.deployment.replicas=2"},
					})
					Expect(err).NotTo(HaveOccurred())

					var resources []*unstructured.Unstructured
					// This piece of work is the simplest way to directly access the unstructured resources list backing a testManifest struct
					// without updating go-utils and adding a direct access function to the TestManifest interface.
					// We aren't doing that because updating gloo's go-utils dependency is its own task to be addressed some other time.
					testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						resources = append(resources, unstructured)
						return true
					})
					for idx1, resource1 := range resources {
						for idx2, resource2 := range resources {
							if idx1 == idx2 {
								continue
							}
							Expect(constructResourceID(resource1)).NotTo(Equal(constructResourceID(resource2)))
						}
					}
				})

				Context("pass image pull secrets", func() {
					pullSecretName := "test-pull-secret"
					pullSecret := []v1.LocalObjectReference{
						{Name: pullSecretName},
					}

					It("via global values", func() {
						testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
							valuesArgs: []string{fmt.Sprintf("global.image.pullSecret=%s", pullSecretName)},
						})
						observabilityDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
						testManifest.ExpectDeploymentAppsV1(observabilityDeployment)
						Expect(err).NotTo(HaveOccurred())
					})

					It("via podTemplate values", func() {
						testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
							valuesArgs: []string{
								fmt.Sprintf("observability.deployment.image.pullSecret=%s", pullSecretName),
							},
						})
						Expect(err).NotTo(HaveOccurred())

						observabilityDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
						testManifest.ExpectDeploymentAppsV1(observabilityDeployment)
					})

					It("podTemplate values win over global", func() {
						testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
							valuesArgs: []string{
								"global.image.pullSecret=wrong",
								fmt.Sprintf("observability.deployment.image.pullSecret=%s", pullSecretName),
							},
						})
						Expect(err).NotTo(HaveOccurred())
						observabilityDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
						testManifest.ExpectDeploymentAppsV1(observabilityDeployment)
					})

				})
			})

			Context("observability secret", func() {
				It("it sets the correct entries when given secrets", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{
							"observability.customGrafana.username=a",
							"observability.customGrafana.password=b",
							"observability.customGrafana.apiKey=c",
							"observability.customGrafana.caBundle=d"},
					})
					Expect(err).NotTo(HaveOccurred())

					expectedSecret := &v1.Secret{
						TypeMeta: k8s.TypeMeta{
							Kind:       "Secret",
							APIVersion: "v1",
						},
						ObjectMeta: k8s.ObjectMeta{
							Name:      "glooe-observability-secrets",
							Namespace: namespace,
							Labels: map[string]string{
								"gloo": "glooe-observability-secrets",
								"app":  "gloo",
							},
						},
						Data: map[string][]byte{
							"GRAFANA_USERNAME":  []byte("a"),
							"GRAFANA_PASSWORD":  []byte("b"),
							"GRAFANA_API_KEY":   []byte("c"),
							"GRAFANA_CA_BUNDLE": []byte("d"),
						},
						Type: v1.SecretTypeOpaque,
					}
					testManifest.Expect("Secret", namespace, "glooe-observability-secrets").
						To(BeEquivalentTo(expectedSecret))
				})
			})

			Context("observability RBAC rule", func() {

				It("allows correct operations on upstreams", func() {
					labels = map[string]string{
						"app":  "gloo",
						"gloo": "observability",
					}
					rb := ResourceBuilder{
						Name:   "observability-upstream-role-gloo-system",
						Labels: labels,
					}

					observabilityClusterRole := rb.GetClusterRole()
					observabilityClusterRole.Rules = []rbacv1.PolicyRule{
						{
							Verbs:     []string{"get", "list", "watch"},
							APIGroups: []string{"gloo.solo.io"},
							Resources: []string{"upstreams"},
						},
					}

					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
					Expect(err).NotTo(HaveOccurred())

					clusterRoles := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						return unstructured.GetKind() == "ClusterRole" && unstructured.GetLabels()["gloo"] == "observability"
					})

					clusterRoles.ExpectClusterRole(observabilityClusterRole)
				})
			})

			Context("grafana deployment", func() {
				It("is not installed when grafana is disabled", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{"grafana.defaultInstallationEnabled=false"},
					})
					Expect(err).NotTo(HaveOccurred())

					testManifest.Expect(grafanaDeployment.Kind, grafanaDeployment.Namespace, grafanaDeployment.Name).To(BeNil())
				})

				It("is not installed when using a custom grafana instance", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{
							"grafana.defaultInstallationEnabled=false",
							"observability.customGrafana.enabled=true",
						},
					})
					Expect(err).NotTo(HaveOccurred())

					testManifest.Expect(grafanaDeployment.Kind, grafanaDeployment.Namespace, grafanaDeployment.Name).To(BeNil())
				})
			})
		})

		Context("external auth server", func() {

			var expectedDeployment *appsv1.Deployment

			BeforeEach(func() {
				labels = map[string]string{
					"app":  "gloo",
					"gloo": "extauth",
				}
				selector = map[string]string{
					"gloo": "extauth",
				}

				rb := ResourceBuilder{
					Namespace: namespace,
					Name:      "extauth",
					Labels:    labels,
				}

				nonRootUser := int64(10101)
				nonRoot := true

				nonRootSC := &v1.PodSecurityContext{
					RunAsUser:    &nonRootUser,
					RunAsNonRoot: &nonRoot,
					FSGroup:      &nonRootUser,
				}

				expectedDeployment = rb.GetDeploymentAppsv1()

				expectedDeployment.Spec.Replicas = aws.Int32(1)
				expectedDeployment.Spec.Template.Spec.Containers = []v1.Container{
					{
						Name:            "extauth",
						Image:           "quay.io/solo-io/extauth-ee:" + version,
						ImagePullPolicy: getPullPolicy(),
						Env: []v1.EnvVar{
							{
								Name: "POD_NAMESPACE",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										FieldPath: "metadata.namespace",
									},
								},
							},
							{
								Name:  "SERVICE_NAME",
								Value: "ext-auth",
							},
							{
								Name:  "GLOO_ADDRESS",
								Value: "gloo:9977",
							},
							{
								Name: "SIGNING_KEY",
								ValueFrom: &v1.EnvVarSource{
									SecretKeyRef: &v1.SecretKeySelector{
										LocalObjectReference: v1.LocalObjectReference{
											Name: "extauth-signing-key",
										},
										Key: "signing-key",
									},
								},
							},
							{
								Name: "REDIS_PASSWORD",
								ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "redis",
									},
									Key: "redis-password",
								}},
							},
							{
								Name:  "SERVER_PORT",
								Value: "8083",
							},
							{
								Name:  "USER_ID_HEADER",
								Value: "x-user-id",
							},
							statsEnvVar,
							{
								Name:  "HEALTH_HTTP_PORT",
								Value: "8082",
							},
							{
								Name:  "HEALTH_HTTP_PATH",
								Value: "/healthcheck",
							},
						},
						ReadinessProbe: &v1.Probe{
							Handler: v1.Handler{
								HTTPGet: &v1.HTTPGetAction{
									Path: "/healthcheck",
									Port: intstr.IntOrString{
										Type:   0,
										IntVal: 8082,
									},
								},
							},
							InitialDelaySeconds: 2,
							PeriodSeconds:       5,
							FailureThreshold:    2,
							SuccessThreshold:    1,
						},
						Resources: v1.ResourceRequirements{},
					},
				}
				expectedDeployment.Spec.Strategy = appsv1.DeploymentStrategy{}
				expectedDeployment.Spec.Selector.MatchLabels = selector
				expectedDeployment.Spec.Template.ObjectMeta.Labels = selector
				expectedDeployment.Spec.Template.ObjectMeta.Annotations = normalPromAnnotations

				expectedDeployment.Spec.Template.Spec.SecurityContext = nonRootSC

				expectedDeployment.Spec.Template.Spec.ServiceAccountName = "extauth"

				expectedDeployment.Spec.Template.Spec.Affinity = &v1.Affinity{
					PodAffinity: &v1.PodAffinity{
						PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
							{
								Weight: 100,
								PodAffinityTerm: v1.PodAffinityTerm{
									LabelSelector: &k8s.LabelSelector{
										MatchLabels: map[string]string{
											"gloo": "gateway-proxy",
										},
									},
									TopologyKey: "kubernetes.io/hostname",
								},
							},
						},
					},
				}
				expectedDeployment.Spec.Replicas = nil // GetDeploymentAppsv1 explicitly sets it to 1, which we don't want
			})

			It("should be able to set custom labels for pods", func() {
				// This test checks for labeling in the 5 deployments that are added by Gloo-E.
				// Subchart deployments (like those for graphana) are ignored.
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.rateLimit.enabled=true",
						"global.extensions.extAuth.enabled=true",
						"observability.enabled=true",
						"global.extensions.rateLimit.deployment.extraRateLimitLabels.foo=bar",
						"redis.deployment.extraRedisLabels.foo=bar",
						"observability.deployment.extraObservabilityLabels.foo=bar",
						"global.extensions.extAuth.deployment.extraExtAuthLabels.foo=bar",
					},
				})

				deploymentBlacklist := []string{
					"gloo",
					"discovery",
					"gateway",
					"gateway-proxy",
					"glooe-grafana",
					"glooe-prometheus-kube-state-metrics",
					"glooe-prometheus-server",
					"gloo-fed",
					"gloo-fed-console",
				}

				Expect(err).NotTo(HaveOccurred())
				Expect(testManifest).NotTo(BeNil())
				var resourcesTested = 0
				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "Deployment"
				}).ExpectAll(func(deployment *unstructured.Unstructured) {
					deploymentObject, err := kuberesource.ConvertUnstructured(deployment)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Deployment %+v should be able to convert from unstructured", deployment))
					structuredDeployment, ok := deploymentObject.(*appsv1.Deployment)
					Expect(ok).To(BeTrue(), fmt.Sprintf("Deployment %+v should be able to cast to a structured deployment", deployment))

					// don't test against blacklisted deployments
					deploymentName := structuredDeployment.GetName()
					isCheckable := true
					for _, blacklistedDeployment := range deploymentBlacklist {
						if deploymentName == blacklistedDeployment {
							isCheckable = false
							break
						}
					}
					if !isCheckable {
						return
					}

					deploymentLabels := structuredDeployment.Spec.Template.Labels
					value, ok := deploymentLabels["foo"]
					Expect(ok).To(BeTrue(), fmt.Sprintf("Coundn't find test label 'foo' in deployment %s", deploymentName))
					Expect(value).To(Equal("bar"), fmt.Sprintf("Test label 'foo' in deployment %s, had unexpected value '%s'", deploymentName, value))
					resourcesTested += 1
				})
				// Is there an elegant way to parameterized the expected number of deployments based on the valueArgs?
				Expect(resourcesTested).To(Equal(4), "Tested %d resources when we were expecting 4."+
					" Was a new pod added, or is an existing pod no longer being generated?", resourcesTested)
			})

			It("produces expected default deployment", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
				Expect(err).NotTo(HaveOccurred())

				actualDeployment := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
					return unstructured.GetKind() == "Deployment" && unstructured.GetLabels()["gloo"] == "extauth"
				})

				actualDeployment.ExpectDeploymentAppsV1(expectedDeployment)
			})
			It("Uses fips images", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.image.fips=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				actualDeployment := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
					return unstructured.GetKind() == "Deployment" && unstructured.GetLabels()["gloo"] == "extauth"
				})
				expectedDeployment.Spec.Template.Spec.Containers[0].Image = "quay.io/solo-io/extauth-ee-fips:" + version
				actualDeployment.ExpectDeploymentAppsV1(expectedDeployment)
			})
			It("configures headers to redact", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.headersToRedact=authorize foo",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				actualDeployment := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
					return unstructured.GetKind() == "Deployment" && unstructured.GetLabels()["gloo"] == "extauth"
				})
				expectedDeployment.Spec.Template.Spec.Containers[0].Env = append(expectedDeployment.Spec.Template.Spec.Containers[0].Env, v1.EnvVar{
					Name:  "HEADERS_TO_REDACT",
					Value: "authorize foo",
				})
				actualDeployment.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("should add an anti-injection annotation to all pods when disableAutoinjection is enabled", func() {
				istioAnnotation := "sidecar.istio.io/inject"
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.istioIntegration.disableAutoinjection=true",
						"global.extensions.rateLimit.enabled=true", // check as many as possible
						"global.extensions.extAuth.enabled=true",
						"observability.enabled=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				// don't check stuff from gloo-OS or outside our purview.
				deploymentBlacklist := []string{
					"gloo",
					"discovery",
					"gateway",
					"gateway-proxy",
					"glooe-grafana",
					"glooe-prometheus-kube-state-metrics",
					"glooe-prometheus-server",
					"gloo-fed",
				}

				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "Deployment"
				}).ExpectAll(func(deployment *unstructured.Unstructured) {
					deploymentObject, err := kuberesource.ConvertUnstructured(deployment)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Deployment %+v should be able to convert from unstructured", deployment))
					structuredDeployment, ok := deploymentObject.(*appsv1.Deployment)
					Expect(ok).To(BeTrue(), fmt.Sprintf("Deployment %+v should be able to cast to a structured deployment", deployment))

					deploymentName := deployment.GetName()
					for _, blacklistedDeployment := range deploymentBlacklist {
						if deploymentName == blacklistedDeployment {
							return
						}
					}

					// ensure every deployment has a istio annotation set to false
					val, ok := structuredDeployment.Spec.Template.ObjectMeta.Annotations[istioAnnotation]
					Expect(ok).To(BeTrue(), fmt.Sprintf("Deployment %s should contain an istio injection annotation", deploymentName))
					Expect(val).To(Equal("false"), fmt.Sprintf("Deployment %s should have an istio annotation with value of 'false'", deploymentName))
				})
			})

			Context("dataplane per proxy", func() {

				redisTlsSecretName := "redis-tls-secret"
				redisCACertSecretName := "redis-ca-cert-secret"
				redisRegex := ".*redis.*secret"

				helmOverrideFileContents := func(dataplanePerProxy bool) string {
					return fmt.Sprintf(`
global:
  extensions:
    dataplanePerProxy: %t
  glooStats:
    enabled: true
gloo:
  discovery:
    fdsMode: BLACKLIST
  gateway:
    validation:
      alwaysAcceptResources: false
      allowWarnings: false
      webhook:
        enabled: true
  gatewayProxies:
    gatewayProxy:
      stats:
        enabled: true
      loopBackAddress: 127.0.0.1
      readConfig: true
      gatewaySettings:
        customHttpGateway:
          virtualServices:
          - name: standard-vs
            namespace: gloo-system
          options:
            httpConnectionManagerSettings:
              xffNumTrustedHops: 1
            proxyLatency:
              chargeClusterStat: true
              chargeListenerStat: true
              response: FIRST_INCOMING_FIRST_OUTGOING
        customHttpsGateway:
          virtualServices:
          - name: standard-vs
            namespace: gloo-system
          options:
            httpConnectionManagerSettings:
              xffNumTrustedHops: 1
            proxyLatency:
              chargeClusterStat: true
              chargeListenerStat: true
              response: FIRST_INCOMING_FIRST_OUTGOING
      podTemplate:
        disableNetBind: false
        floatingUserId: false
        httpPort: 8080
        httpsPort: 8443
    customProxy:
      loopBackAddress: 127.0.0.1
      stats:
        enabled: true
      gatewaySettings:
        customHttpGateway: 
          virtualServices:
          - name: custom-vs
            namespace: gloo-system
          options:
            httpConnectionManagerSettings:
              xffNumTrustedHops: 1
            proxyLatency:
              chargeClusterStat: true
              chargeListenerStat: true
              response: FIRST_INCOMING_FIRST_OUTGOING
        customHttpsGateway:
          virtualServices:
          - name: custom-vs
            namespace: gloo-system
          options:
            httpConnectionManagerSettings:
              xffNumTrustedHops: 1
            proxyLatency:
              chargeClusterStat: true
              chargeListenerStat: true
              response: FIRST_INCOMING_FIRST_OUTGOING
      readConfig: true
      configMap:
        data: 
      antiAffinity: false
      kind:
        deployment:
          replicas: 1
      podTemplate:
        disableNetBind: false
        floatingUserId: false
        httpPort: 8180
        httpsPort: 8543
        image:
          repository: gloo-ee-envoy-wrapper
          tag: 1.4.0
        probes: false
        runUnprivileged: false
      service:
        httpPort: 80
        httpsPort: 1443
        type: LoadBalancer
        extraAnnotations:
          prometheus.io/path: "/metrics"
          prometheus.io/port: "8081"
          prometheus.io/scrape: "true"
`, dataplanePerProxy)
				}

				It("allows dataplane per proxy", func() {

					helmOverrideFile := "helm-override-*.yaml"
					tmpFile, err := ioutil.TempFile("", helmOverrideFile)
					Expect(err).ToNot(HaveOccurred())
					_, err = tmpFile.Write([]byte(helmOverrideFileContents(true)))
					Expect(err).NotTo(HaveOccurred())
					defer tmpFile.Close()
					defer os.Remove(tmpFile.Name())
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesFile: tmpFile.Name(),
						valuesArgs: []string{"gloo.gatewayProxies.gatewayProxy.disabled=false"},
					})
					Expect(err).NotTo(HaveOccurred())
					// upstreams are rendered by the rollout job
					job := getJob(testManifest, namespace, "gloo-ee-resource-rollout")

					assertExpectedResourcesForProxy := func(proxyName string) {
						// RateLimit
						gatewayProxyRateLimitResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
							return unstructured.GetLabels()["gloo"] == fmt.Sprintf("rate-limit-%s", proxyName)
						})
						// Deployment and Service
						Expect(gatewayProxyRateLimitResources.NumResources()).To(Equal(2), fmt.Sprintf("%s: Expecting RateLimit Deployment and Service", proxyName))
						// Upstream
						Expect(strings.Count(job.Spec.Template.Spec.Containers[0].Command[2], "gloo: "+fmt.Sprintf("rate-limit-%s", proxyName))).To(Equal(1), fmt.Sprintf("%s: Expecting RateLimit Upstream", proxyName))

						// Redis
						gatewayProxyRedisResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
							return unstructured.GetLabels()["gloo"] == fmt.Sprintf("redis-%s", proxyName)
						})
						// Deployment and Service
						Expect(gatewayProxyRedisResources.NumResources()).To(Equal(2), fmt.Sprintf("%s: Expecting Redis Deployment and Service", proxyName))

						// ExtAuth
						gatewayProxyExtAuthResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
							return unstructured.GetLabels()["gloo"] == fmt.Sprintf("extauth-%s", proxyName)
						})
						// Deployment and Service
						Expect(gatewayProxyExtAuthResources.NumResources()).To(Equal(2), fmt.Sprintf("%s: Expecting Extauth Deployment and Service", proxyName))
						// Upstream
						Expect(strings.Count(job.Spec.Template.Spec.Containers[0].Command[2], "gloo: "+fmt.Sprintf("extauth-%s", proxyName))).To(Equal(1), fmt.Sprintf("%s: Expecting Extauth Upstream", proxyName))
					}

					assertExpectedResourcesForProxy("gateway-proxy")
					assertExpectedResourcesForProxy("custom-proxy")

				})

				It("gateway proxy objects are not created when gatewayProxy is disabled", func() {
					helmOverrideFile := "helm-override-*.yaml"
					tmpFile, err := ioutil.TempFile("", helmOverrideFile)
					Expect(err).ToNot(HaveOccurred())
					_, err = tmpFile.Write([]byte(helmOverrideFileContents(false)))
					Expect(err).NotTo(HaveOccurred())
					defer tmpFile.Close()
					defer os.Remove(tmpFile.Name())
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesFile: tmpFile.Name(),
						valuesArgs: []string{"gloo.gatewayProxies.gatewayProxy.disabled=true"},
					})
					Expect(err).NotTo(HaveOccurred())

					assertExpectedResourcesForDisabledProxy := func(proxyName string) {
						gatewayProxyRateLimitResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
							return unstructured.GetLabels()["gloo"] == fmt.Sprintf("rate-limit-%s", proxyName)
						})

						// Deployment, Service, Upstream
						Expect(gatewayProxyRateLimitResources.NumResources()).To(Equal(0), fmt.Sprintf("%s: Expecting RateLimit Deployment, Service, and Upstream to not be created", proxyName))

						gatewayProxyRedisResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
							return unstructured.GetLabels()["gloo"] == fmt.Sprintf("redis-%s", proxyName)
						})

						// Deployment, Service
						Expect(gatewayProxyRedisResources.NumResources()).To(Equal(0), fmt.Sprintf("%s: Expecting Redis Deployment and Service to not be created", proxyName))

						gatewayProxyExtAuthResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
							return unstructured.GetLabels()["gloo"] == fmt.Sprintf("extauth-%s", proxyName)
						})

						// Deployment, Service, Upstream
						Expect(gatewayProxyExtAuthResources.NumResources()).To(Equal(0), fmt.Sprintf("%s: Expecting Extauth Deployment, Service, and Upstream to not be created", proxyName))

					}

					assertExpectedResourcesForDisabledProxy("gateway-proxy")
				})

				It("Redis objects are not created when .Values.redis.disabled is set", func() {
					// file creation operations to support test
					helmOverrideFile := "helm-override-*.yaml"
					tmpFile, err := ioutil.TempFile("", helmOverrideFile)
					Expect(err).ToNot(HaveOccurred())
					_, err = tmpFile.Write([]byte(helmOverrideFileContents(false)))
					Expect(err).NotTo(HaveOccurred())
					defer tmpFile.Close()
					defer os.Remove(tmpFile.Name())

					proxyName := "gateway-proxy"

					// assert no redis resources exist wtih "redis.disabled=true"
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesFile: tmpFile.Name(),
						valuesArgs: []string{"redis.disabled=true"},
					})

					redisResources := testManifest.SelectResources(func(un *unstructured.Unstructured) bool {
						match, _ := regexp.MatchString(redisRegex, un.GetName())
						return match
					})

					gatewayProxyRedisResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						return unstructured.GetLabels()["gloo"] == fmt.Sprintf("redis-%s", proxyName)
					})
					Expect(gatewayProxyRedisResources.NumResources()).To(Equal(0), fmt.Sprintf("%s: Expecting Redis Deployment and Service to not be created", proxyName))
					Expect(redisResources.NumResources()).To(Equal(0), fmt.Sprintf("%s: Expecting Redis secret to not be created", proxyName))
				})

				It("Redis objects are not built when .Values.redis.disabled is set but rate-limit sets up TLS when .Values.redis.cert.enabled is set", func() {
					// file creation operations to support test
					helmOverrideFile := "helm-override-*.yaml"
					tmpFile, err := ioutil.TempFile("", helmOverrideFile)
					Expect(err).ToNot(HaveOccurred())
					_, err = tmpFile.Write([]byte(helmOverrideFileContents(false)))
					Expect(err).NotTo(HaveOccurred())
					defer tmpFile.Close()
					defer os.Remove(tmpFile.Name())

					proxyName := "gateway-proxy"

					// assert no redis resources exist wtih "redis.disabled=true"
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesFile: tmpFile.Name(),
						valuesArgs: []string{
							"redis.disabled=true",
							"redis.cert.enabled=true",
							"redis.cert.crt=certValue",
							"redis.cert.key=keyValue",
						},
					})
					redisResources := testManifest.SelectResources(func(un *unstructured.Unstructured) bool {
						match, _ := regexp.MatchString(redisRegex, un.GetName())
						return match
					})
					Expect(redisResources.NumResources()).To(Equal(1), fmt.Sprintf("%s: Expecting Redis secret to be created", proxyName))

					redisDeploymentCreated := false
					rateLimitDeploymentCreated := false
					testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
						return resource.GetKind() == "Deployment"
					}).ExpectAll(func(deployment *unstructured.Unstructured) {
						deploymentObject, err := kuberesource.ConvertUnstructured(deployment)
						Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Deployment %+v should be able to convert from unstructured", deployment))
						structuredDeployment, ok := deploymentObject.(*appsv1.Deployment)
						Expect(ok).To(BeTrue(), fmt.Sprintf("Deployment %+v should be able to cast to a structured deployment", deployment))
						if structuredDeployment.GetName() == "redis" {
							redisDeploymentCreated = true
						}
						// should add Redis TLS cert and secret to rate-limit
						if structuredDeployment.GetName() == "rate-limit" {
							rateLimitDeploymentCreated = true
							ex := ExpectContainer{
								Containers: structuredDeployment.Spec.Template.Spec.Containers,
								Name:       "rate-limit",
							}
							Expect(structuredDeployment.Spec.Template.Spec.Containers).To(HaveLen(1), "should have exactly 1 container")
							ex.ExpectToHaveEnv("REDIS_URL", "redis:6379", "should have the redis url for rate-limit")
							ex.ExpectToHaveEnv("REDIS_SOCKET_TYPE", "tls", "should use tls socket for redis url")
							ex.ExpectToHaveEnv("REDIS_CA_CERT", "/etc/tls/ca.crt", "should have tls cert set to secret")
							Expect(structuredDeployment.Spec.Template.Spec.Volumes).To(ContainElement(v1.Volume{
								Name: "ca-cert-volume",
								VolumeSource: v1.VolumeSource{
									Secret: &v1.SecretVolumeSource{
										SecretName:  "glooe-" + redisCACertSecretName,
										Items:       nil,
										DefaultMode: proto.Int(420),
									},
								},
							}))
						}
					})
					Expect(redisDeploymentCreated).To(BeFalse(), "Should not create the redis deployment")
					Expect(rateLimitDeploymentCreated).To(BeTrue(), "Should create the rate-limit deployment")
				})

				It("Be able to attach secret mounts to the ext-auth deployment", func() {
					helmOverrideFile := "helm-override-*.yaml"
					tmpFile, err := ioutil.TempFile("", helmOverrideFile)
					Expect(err).ToNot(HaveOccurred())
					_, err = tmpFile.Write([]byte(helmOverrideFileContents(false)))
					Expect(err).NotTo(HaveOccurred())
					defer tmpFile.Close()
					defer os.Remove(tmpFile.Name())
					secretNamePrefix := "user-session-cert-"
					name1 := "extauthsecret1"
					name2 := "extauthsecret2"
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesFile: tmpFile.Name(),
						valuesArgs: []string{
							fmt.Sprintf("global.extensions.extAuth.deployment.redis.certs[0].secretName=%s", name1),
							fmt.Sprintf("global.extensions.extAuth.deployment.redis.certs[0].mountPath=%s", name1),
							fmt.Sprintf("global.extensions.extAuth.deployment.redis.certs[1].secretName=%s", name2),
							fmt.Sprintf("global.extensions.extAuth.deployment.redis.certs[1].mountPath=%s", name2),
						},
					})

					testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
						return resource.GetKind() == "Deployment"
					}).ExpectAll(func(deployment *unstructured.Unstructured) {
						deploymentObject, err := kuberesource.ConvertUnstructured(deployment)
						Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Deployment %+v should be able to convert from unstructured", deployment))
						structuredDeployment, ok := deploymentObject.(*appsv1.Deployment)
						Expect(ok).To(BeTrue(), fmt.Sprintf("Deployment %+v should be able to cast to a structured deployment", deployment))
						if structuredDeployment.GetName() == "extauth" {
							spec := structuredDeployment.Spec.Template.Spec
							ec := ExpectContainer{
								Containers: spec.Containers,
								Name:       "extauth",
							}
							ev := ExpectVolume{
								Volumes: spec.Volumes,
							}
							ev.ExpectHasName(secretNamePrefix + name1)
							ev.ExpectHasName(secretNamePrefix + name2)
							mp := ec.ExpectToHaveVolumeMount(secretNamePrefix + name1)
							Expect(mp.MountPath).To(Equal(name1))
							mp = ec.ExpectToHaveVolumeMount(secretNamePrefix + name2)
							Expect(mp.MountPath).To(Equal(name2))
						}
					})
				})

				It("Redis objects are built when .Values.redis.disabled is not set and rate-limit sets up TLS when .Values.redis.cert.enabled is set", func() {
					// file creation operations to support test
					helmOverrideFile := "helm-override-*.yaml"
					tmpFile, err := ioutil.TempFile("", helmOverrideFile)
					Expect(err).ToNot(HaveOccurred())
					_, err = tmpFile.Write([]byte(helmOverrideFileContents(false)))
					Expect(err).NotTo(HaveOccurred())
					defer tmpFile.Close()
					defer os.Remove(tmpFile.Name())

					proxyName := "gateway-proxy"

					// assert no redis resources exist wtih "redis.disabled=true"
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesFile: tmpFile.Name(),
						valuesArgs: []string{
							"redis.cert.enabled=true",
							"redis.cert.crt=certValue",
							"redis.cert.key=keyValue",
						},
					})
					redisResources := testManifest.SelectResources(func(un *unstructured.Unstructured) bool {
						match, _ := regexp.MatchString(redisRegex, un.GetName())
						return match
					})

					Expect(redisResources.NumResources()).To(Equal(2), fmt.Sprintf("%s: Expecting Redis secret to be created", proxyName))

					redisDeploymentCreated := false
					rateLimitDeploymentCreated := false
					testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
						return resource.GetKind() == "Deployment"
					}).ExpectAll(func(deployment *unstructured.Unstructured) {
						deploymentObject, err := kuberesource.ConvertUnstructured(deployment)
						Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Deployment %+v should be able to convert from unstructured", deployment))
						structuredDeployment, ok := deploymentObject.(*appsv1.Deployment)
						Expect(ok).To(BeTrue(), fmt.Sprintf("Deployment %+v should be able to cast to a structured deployment", deployment))
						// should have redis deployed with tls enabled
						if structuredDeployment.GetName() == "redis" {
							redisDeploymentCreated = true
							ex := ExpectContainer{
								Containers: structuredDeployment.Spec.Template.Spec.Containers,
								Name:       "redis",
							}
							Expect(structuredDeployment.Spec.Template.Spec.Containers).To(HaveLen(1), "should have exactly 1 container")
							ex.ExpectToHaveArg("--tls-port 6379", "should have tls port enabled at default port")
							ex.ExpectToHaveArg("--port 0", "should not listen to from any port")
							ex.ExpectToHaveArg("--tls-cert-file /etc/tls/tls.crt", "should set the tls cert to the location in the secret volume")
							ex.ExpectToHaveArg("--tls-ca-cert-file /etc/ca-cert/ca.crt", "should set the CA cert to the location in the secret volume")
							ex.ExpectToHaveArg("--tls-key-file /etc/tls/tls.key", "should set the tls key to the location in the secret volume")
							ex.ExpectToHaveArg("--tls-auth-clients no", "should set auth clients to no")
							Expect(structuredDeployment.Spec.Template.Spec.Volumes).To(ContainElement(v1.Volume{
								Name: "tls-volume",
								VolumeSource: v1.VolumeSource{
									Secret: &v1.SecretVolumeSource{
										SecretName:  "glooe-" + redisTlsSecretName,
										Items:       nil,
										DefaultMode: proto.Int(420),
									},
								},
							}))
							Expect(structuredDeployment.Spec.Template.Spec.Volumes).To(ContainElement(v1.Volume{
								Name: "ca-cert-volume",
								VolumeSource: v1.VolumeSource{
									Secret: &v1.SecretVolumeSource{
										SecretName:  "glooe-" + redisCACertSecretName,
										Items:       nil,
										DefaultMode: proto.Int(420),
									},
								},
							}))
						}
						// should add Redis TLS cert and secret to rate-limit
						if structuredDeployment.GetName() == "rate-limit" {
							rateLimitDeploymentCreated = true
							ex := ExpectContainer{
								Containers: structuredDeployment.Spec.Template.Spec.Containers,
								Name:       "rate-limit",
							}
							Expect(structuredDeployment.Spec.Template.Spec.Containers).To(HaveLen(1), "should have exactly 1 container")
							ex.ExpectToHaveEnv("REDIS_URL", "redis:6379", "should have the redis url for rate-limit")
							ex.ExpectToHaveEnv("REDIS_SOCKET_TYPE", "tls", "should use tls socket for redis url")
							ex.ExpectToHaveEnv("REDIS_CA_CERT", "/etc/tls/ca.crt", "should have tls cert set to secret")
							Expect(structuredDeployment.Spec.Template.Spec.Volumes).To(ContainElement(v1.Volume{
								Name: "ca-cert-volume",
								VolumeSource: v1.VolumeSource{
									Secret: &v1.SecretVolumeSource{
										SecretName:  "glooe-" + redisCACertSecretName,
										Items:       nil,
										DefaultMode: proto.Int(420),
									},
								},
							}))
						}
					})
					Expect(redisDeploymentCreated).To(BeTrue(), "Should create the redis deployment")
					Expect(rateLimitDeploymentCreated).To(BeTrue(), "Should create the rate-limit deployment")
				})

				It("doesn't duplicate resources across proxies when dataplane per proxy is false", func() {

					helmOverrideFile := "helm-override-*.yaml"
					tmpFile, err := ioutil.TempFile("", helmOverrideFile)
					Expect(err).ToNot(HaveOccurred())
					_, err = tmpFile.Write([]byte(helmOverrideFileContents(false)))
					Expect(err).NotTo(HaveOccurred())
					defer tmpFile.Close()
					defer os.Remove(tmpFile.Name())
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesFile: tmpFile.Name(),
					})
					Expect(err).NotTo(HaveOccurred())
					// upstreams are rendered by the rollout job
					job := getJob(testManifest, namespace, "gloo-ee-resource-rollout")

					rateLimitResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						return unstructured.GetLabels()["gloo"] == "rate-limit"
					})
					Expect(rateLimitResources.NumResources()).To(Equal(3), "Expecting RateLimit Deployment, Service, and ServiceAccount")
					Expect(strings.Count(job.Spec.Template.Spec.Containers[0].Command[2], "gloo: rate-limit")).To(Equal(1), "Expecting RateLimit Upstream")

					redisResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						return unstructured.GetLabels()["gloo"] == "redis"
					})
					Expect(redisResources.NumResources()).To(Equal(3), "Expecting Redis Deployment, Service, and Secret")

					extAuthResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						return unstructured.GetLabels()["gloo"] == "extauth"
					})
					Expect(extAuthResources.NumResources()).To(Equal(4), "Expecting ExtAuth Deployment, Service, ServiceAccount, and Secret")
					Expect(strings.Count(job.Spec.Template.Spec.Containers[0].Command[2], "gloo: extauth")).To(Equal(1), "Expecting ExtAuth Upstream")
				})
			})

			It("allows setting the number of replicas for the deployment", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"global.extensions.extAuth.deployment.replicas=3"},
				})
				Expect(err).NotTo(HaveOccurred())

				actualDeployment := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
					return unstructured.GetKind() == "Deployment" && unstructured.GetLabels()["gloo"] == "extauth"
				})

				expectedDeployment.Spec.Replicas = aws.Int32(3)
				actualDeployment.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("correctly sets resource limits for the extauth deployment", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.deployment.resources.limits.cpu=300m",
						"global.extensions.extAuth.deployment.resources.limits.memory=300Mi",
						"global.extensions.extAuth.deployment.resources.requests.cpu=30m",
						"global.extensions.extAuth.deployment.resources.requests.memory=30Mi",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				expectedDeployment.Spec.Template.Spec.Containers[0].Resources = v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("300m"),
						v1.ResourceMemory: resource.MustParse("300Mi"),
					},
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("30m"),
						v1.ResourceMemory: resource.MustParse("30Mi"),
					},
				}

				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("allows setting custom runAsUser", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"global.extensions.extAuth.deployment.runAsUser=10102"},
				})
				Expect(err).NotTo(HaveOccurred())

				actualDeployment := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
					return unstructured.GetKind() == "Deployment" && unstructured.GetLabels()["gloo"] == "extauth"
				})

				uid := int64(10102)
				expectedDeployment.Spec.Template.Spec.SecurityContext.RunAsUser = &uid
				actualDeployment.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("allows multiple extauth plugins", func() {
				helmOverrideFileContents := `
global:
  extensions:
    extAuth:
      plugins:
        first-plugin:
          image:
            repository: ext-auth-plugins
            registry: quay.io/solo-io
            pullPolicy: IfNotPresent
            tag: 1.2.3
        second-plugin:
          image:
            repository: foo
            registry: bar
            pullPolicy: IfNotPresent
            tag: 1.2.3`
				helmOverrideFile := "helm-override-*.yaml"
				tmpFile, err := ioutil.TempFile("", helmOverrideFile)
				Expect(err).ToNot(HaveOccurred())
				_, err = tmpFile.Write([]byte(helmOverrideFileContents))
				Expect(err).NotTo(HaveOccurred())
				defer tmpFile.Close()
				defer os.Remove(tmpFile.Name())
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesFile: tmpFile.Name(),
				})
				Expect(err).NotTo(HaveOccurred())

				actualDeployment := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
					return unstructured.GetKind() == "Deployment" && unstructured.GetLabels()["gloo"] == "extauth"
				})

				authPluginVolumeMount := []v1.VolumeMount{
					{
						Name:      "auth-plugins",
						MountPath: "/auth-plugins",
					},
				}
				expectedDeployment.Spec.Template.Spec.InitContainers = []v1.Container{
					{
						Name:            "plugin-first-plugin",
						Image:           "quay.io/solo-io/ext-auth-plugins:1.2.3",
						ImagePullPolicy: getPullPolicy(),
						VolumeMounts:    authPluginVolumeMount,
					},
					{
						Name:            "plugin-second-plugin",
						Image:           "bar/foo:1.2.3",
						ImagePullPolicy: getPullPolicy(),
						VolumeMounts:    authPluginVolumeMount,
					},
				}
				expectedDeployment.Spec.Template.Spec.Volumes = []v1.Volume{
					{
						Name: "auth-plugins",
						VolumeSource: v1.VolumeSource{
							EmptyDir: &v1.EmptyDirVolumeSource{},
						},
					},
				}
				for i, _ := range expectedDeployment.Spec.Template.Spec.Containers {
					expectedDeployment.Spec.Template.Spec.Containers[i].VolumeMounts =
						append(expectedDeployment.Spec.Template.Spec.Containers[i].VolumeMounts, authPluginVolumeMount...)
				}
				actualDeployment.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("can accept extra env vars for extauth deployment", func() {
				expectedDeployment.Spec.Template.Spec.Containers[0].Env = append(
					[]v1.EnvVar{
						{
							Name:  "TEST_EXTRA_ENV_VAR",
							Value: "test",
						},
					},
					expectedDeployment.Spec.Template.Spec.Containers[0].Env...,
				)

				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.deployment.customEnv[0].Name=TEST_EXTRA_ENV_VAR",
						"global.extensions.extAuth.deployment.customEnv[0].Value=test",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("can add extra volume mounts to the extauth deployment", func() {

				expectedDeployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
					expectedDeployment.Spec.Template.Spec.Containers[0].VolumeMounts,
					v1.VolumeMount{
						Name:      "test-path",
						MountPath: "/var/run/sds",
					},
				)

				expectedDeployment.Spec.Template.Spec.Volumes = append(
					expectedDeployment.Spec.Template.Spec.Volumes,
					v1.Volume{
						Name: "test-path",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{
								Path: "/var/run/sds",
							},
						},
					},
				)

				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.deployment.extraVolume[0].Name=test-path",
						"global.extensions.extAuth.deployment.extraVolume[0].HostPath.Path=/var/run/sds",
						"global.extensions.extAuth.deployment.extraVolumeMount[0].Name=test-path",
						"global.extensions.extAuth.deployment.extraVolumeMount[0].MountPath=/var/run/sds",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			Describe("affinity and antiAffinity", func() {
				It("set default affinity rules appropriately", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
					Expect(err).NotTo(HaveOccurred())

					actualDeployment := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						return unstructured.GetKind() == "Deployment" && unstructured.GetLabels()["gloo"] == "extauth"
					})

					expectedDeployment.Spec.Template.Spec.Affinity = &v1.Affinity{
						PodAffinity: &v1.PodAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
								{
									Weight: 100,
									PodAffinityTerm: v1.PodAffinityTerm{
										LabelSelector: &k8s.LabelSelector{
											MatchLabels: map[string]string{
												"gloo": "gateway-proxy",
											},
										},
										TopologyKey: "kubernetes.io/hostname",
									},
								},
							},
						},
					}
					actualDeployment.ExpectDeploymentAppsV1(expectedDeployment)
				})
				It("affinity rules can be set", func() {

					helmOverrideFileContents := `
global:
  extensions:
    extAuth:
      affinity:
        podAffinity: 
          preferredDuringSchedulingIgnoredDuringExecution: 
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchLabels:
                  gloo: test-label
              topologyKey: kubernetes.io/hostname
`
					helmOverrideFile := "helm-override-*.yaml"
					tmpFile, err := ioutil.TempFile("", helmOverrideFile)
					Expect(err).ToNot(HaveOccurred())
					_, err = tmpFile.Write([]byte(helmOverrideFileContents))
					Expect(err).NotTo(HaveOccurred())
					defer tmpFile.Close()
					defer os.Remove(tmpFile.Name())
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesFile: tmpFile.Name(),
					})
					Expect(err).NotTo(HaveOccurred())

					actualDeployment := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						return unstructured.GetKind() == "Deployment" && unstructured.GetLabels()["gloo"] == "extauth"
					})

					expectedDeployment.Spec.Template.Spec.Affinity = &v1.Affinity{
						PodAffinity: &v1.PodAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
								{
									Weight: 100,
									PodAffinityTerm: v1.PodAffinityTerm{
										LabelSelector: &k8s.LabelSelector{
											MatchLabels: map[string]string{
												"gloo": "test-label",
											},
										},
										TopologyKey: "kubernetes.io/hostname",
									},
								},
							},
						},
					}
					actualDeployment.ExpectDeploymentAppsV1(expectedDeployment)
				})
				It("antiAffinity rules can be set", func() {

					helmOverrideFileContents := `
global:
  extensions:
    extAuth:
      antiAffinity:
        podAntiAffinity: 
          preferredDuringSchedulingIgnoredDuringExecution: 
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchLabels:
                  gloo: gateway-proxy
              topologyKey: kubernetes.io/hostname
`
					helmOverrideFile := "helm-override-*.yaml"
					tmpFile, err := ioutil.TempFile("", helmOverrideFile)
					Expect(err).ToNot(HaveOccurred())
					_, err = tmpFile.Write([]byte(helmOverrideFileContents))
					Expect(err).NotTo(HaveOccurred())
					defer tmpFile.Close()
					defer os.Remove(tmpFile.Name())
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesFile: tmpFile.Name(),
					})
					Expect(err).NotTo(HaveOccurred())

					actualDeployment := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						return unstructured.GetKind() == "Deployment" && unstructured.GetLabels()["gloo"] == "extauth"
					})

					expectedDeployment.Spec.Template.Spec.Affinity = &v1.Affinity{
						// default affinity settings
						PodAffinity: &v1.PodAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
								{
									Weight: 100,
									PodAffinityTerm: v1.PodAffinityTerm{
										LabelSelector: &k8s.LabelSelector{
											MatchLabels: map[string]string{
												"gloo": "gateway-proxy",
											},
										},
										TopologyKey: "kubernetes.io/hostname",
									},
								},
							},
						},
						PodAntiAffinity: &v1.PodAntiAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
								{
									Weight: 100,
									PodAffinityTerm: v1.PodAffinityTerm{
										LabelSelector: &k8s.LabelSelector{
											MatchLabels: map[string]string{
												"gloo": "gateway-proxy",
											},
										},
										TopologyKey: "kubernetes.io/hostname",
									},
								},
							},
						},
					}
					actualDeployment.ExpectDeploymentAppsV1(expectedDeployment)
				})
			})

			Context("pass image pull secrets", func() {

				pullSecretName := "test-pull-secret"
				pullSecret := []v1.LocalObjectReference{
					{Name: pullSecretName},
				}

				It("via global values", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{fmt.Sprintf("global.image.pullSecret=%s", pullSecretName)},
					})
					expectedDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
					testManifest.ExpectDeploymentAppsV1(expectedDeployment)
					Expect(err).NotTo(HaveOccurred())
				})

				It("via podTemplate values", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{
							fmt.Sprintf("global.extensions.extAuth.deployment.image.pullSecret=%s", pullSecretName),
						},
					})
					Expect(err).NotTo(HaveOccurred())

					expectedDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
					testManifest.ExpectDeploymentAppsV1(expectedDeployment)
				})

				It("podTemplate values win over global", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{
							"global.image.pullSecret=wrong",
							fmt.Sprintf("global.extensions.extAuth.deployment.image.pullSecret=%s", pullSecretName),
						},
					})
					Expect(err).NotTo(HaveOccurred())
					expectedDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
					testManifest.ExpectDeploymentAppsV1(expectedDeployment)
				})

			})

			It("gwp pdb disabled by default", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
				Expect(err).To(BeNil())
				testManifest.ExpectUnstructured("PodDisruptionBudget", namespace, "ext-auth-pdb").To(BeNil())
			})

			It("can create gwp pdb with minAvailable", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.deployment.podDisruptionBudget.minAvailable=2",
					},
				})
				Expect(err).To(BeNil())

				pdb := makeUnstructured(`
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: ext-auth-pdb
  namespace: gloo-system
spec:
  minAvailable: 2
  selector:
    matchLabels:
      gloo: ext-auth
`)

				testManifest.ExpectUnstructured("PodDisruptionBudget", namespace, "ext-auth-pdb").To(BeEquivalentTo(pdb))
			})

			It("can create gwp pdb with maxUnavailable", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.deployment.podDisruptionBudget.maxUnavailable=2",
					},
				})
				Expect(err).To(BeNil())

				pdb := makeUnstructured(`
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: ext-auth-pdb
  namespace: gloo-system
spec:
  maxUnavailable: 2
  selector:
    matchLabels:
      gloo: ext-auth
`)

				testManifest.ExpectUnstructured("PodDisruptionBudget", namespace, "ext-auth-pdb").To(BeEquivalentTo(pdb))
			})
		})

		Context("gateway-proxy deployment", func() {
			var (
				gatewayProxyDeployment *appsv1.Deployment
			)

			includeStatConfig := func() {
				gatewayProxyDeployment.Spec.Template.ObjectMeta.Annotations["readconfig-stats"] = "/stats"
				gatewayProxyDeployment.Spec.Template.ObjectMeta.Annotations["readconfig-ready"] = "/ready"
				gatewayProxyDeployment.Spec.Template.ObjectMeta.Annotations["readconfig-config_dump"] = "/config_dump"
				gatewayProxyDeployment.Spec.Template.ObjectMeta.Annotations["readconfig-port"] = "8082"
			}

			BeforeEach(func() {
				labels = map[string]string{
					"gloo":             "gateway-proxy",
					"gateway-proxy-id": defaults.GatewayProxyName,
					"app":              "gloo",
				}
				selector = map[string]string{
					"gloo":             "gateway-proxy",
					"gateway-proxy-id": defaults.GatewayProxyName,
				}
				podLabels := map[string]string{
					"gloo":             "gateway-proxy",
					"gateway-proxy-id": defaults.GatewayProxyName,
					"gateway-proxy":    "live",
				}
				podAnnotations := map[string]string{
					"prometheus.io/path":   "/metrics",
					"prometheus.io/port":   "8081",
					"prometheus.io/scrape": "true",
				}
				podname := v1.EnvVar{
					Name: "POD_NAME",
					ValueFrom: &v1.EnvVarSource{
						FieldRef: &v1.ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				}

				container := GetQuayContainerSpec("gloo-ee-envoy-wrapper", version, GetPodNamespaceEnvVar(), podname)
				container.Name = defaults.GatewayProxyName
				container.Args = []string{"--disable-hot-restart"}

				rb := ResourceBuilder{
					Namespace:  namespace,
					Name:       defaults.GatewayProxyName,
					Labels:     labels,
					Containers: []ContainerSpec{container},
				}
				deploy := rb.GetDeploymentAppsv1()
				deploy.Spec.Selector = &k8s.LabelSelector{
					MatchLabels: selector,
				}
				deploy.Spec.Template.ObjectMeta.Labels = podLabels
				deploy.Spec.Template.ObjectMeta.Annotations = podAnnotations
				deploy.Spec.Template.Spec.Volumes = []v1.Volume{{
					Name: "envoy-config",
					VolumeSource: v1.VolumeSource{
						ConfigMap: &v1.ConfigMapVolumeSource{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "gateway-proxy-envoy-config",
							},
						},
					},
				}}
				deploy.Spec.Template.Spec.Containers[0].ImagePullPolicy = getPullPolicy()
				deploy.Spec.Template.Spec.Containers[0].Ports = []v1.ContainerPort{
					{Name: "http", ContainerPort: 8080, Protocol: "TCP"},
					{Name: "https", ContainerPort: 8443, Protocol: "TCP"},
				}
				deploy.Spec.Template.Spec.Containers[0].VolumeMounts = []v1.VolumeMount{{
					Name:      "envoy-config",
					ReadOnly:  false,
					MountPath: "/etc/envoy",
					SubPath:   "",
				}}
				truez := true
				falsez := false
				defaultUser := int64(10101)
				deploy.Spec.Template.Spec.Containers[0].SecurityContext = &v1.SecurityContext{
					Capabilities: &v1.Capabilities{
						Drop: []v1.Capability{"ALL"},
					},
					ReadOnlyRootFilesystem:   &truez,
					AllowPrivilegeEscalation: &falsez,
					RunAsNonRoot:             &truez,
					RunAsUser:                &defaultUser,
				}

				deploy.Spec.Template.Spec.SecurityContext = &v1.PodSecurityContext{
					RunAsUser: &defaultUser,
					FSGroup:   &defaultUser,
				}

				deploy.Spec.Template.Spec.ServiceAccountName = "gateway-proxy"

				gatewayProxyDeployment = deploy
			})

			It("creates a deployment without envoy config annotations", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
				Expect(err).NotTo(HaveOccurred())
				testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
			})

			It("creates a deployment with envoy config annotations", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"gloo.gatewayProxies.gatewayProxy.readConfig=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				includeStatConfig()
				testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
			})

			It("creates a deployment with fips envoy", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.image.fips=true",
					},
				})
				gatewayProxyDeployment.Spec.Template.Spec.Containers[0].Image = "quay.io/solo-io/gloo-ee-envoy-wrapper-fips:" + version
				Expect(err).NotTo(HaveOccurred())
				testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
			})

			It("creates settings with extauth request timeout", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.requestTimeout=1s",
						"global.extensions.extAuth.enabled=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				settings := makeUnstructured(`
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: ` + namespace + `
spec:
  discovery:
    fdsMode: WHITELIST
  extauth:
    transportApiVersion: V3
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
    requestTimeout: "1s"
    userIdHeader: "x-user-id"
  gateway:
    enableGatewayController: true
    readGatewaysFromAllNamespaces: false
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
      disableTransformationValidation: false
      allowWarnings: true
      warnRouteShortCircuiting: false
      validationServerGrpcMaxSizeBytes: 104857600
  gloo:
    enableRestEds: false
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    proxyDebugBindAddr: 0.0.0.0:9966
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    invalidConfigPolicy:
      replaceInvalidRoutes: false
      invalidRouteResponseBody: "Gloo Gateway has invalid configuration. Administrators should run ` + backtick + "glooctl check" + backtick + ` to find and fix config errors."
      invalidRouteResponseCode: 404
      replaceInvalidRoutes: false
  ratelimitServer:
    rateLimitBeforeAuth: false
    ratelimitServerRef:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
  consoleOptions:
    readOnly: false
    apiExplorerEnabled: true
  graphqlOptions:
    schemaChangeValidationOptions:
      rejectBreakingChanges: false
      processingRules: []
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})

			It("creates settings with extauth request body", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.requestBody.maxRequestBytes=64000",
						"global.extensions.extAuth.requestBody.packAsBytes=true",
						"global.extensions.extAuth.enabled=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				settings := makeUnstructured(`
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: ` + namespace + `
spec:
  discovery:
    fdsMode: WHITELIST
  extauth:
    transportApiVersion: V3
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
    requestBody:
      maxRequestBytes: 64000
      packAsBytes: true
    userIdHeader: "x-user-id"
  gateway:
    enableGatewayController: true
    readGatewaysFromAllNamespaces: false
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
      disableTransformationValidation: false
      allowWarnings: true
      warnRouteShortCircuiting: false
      validationServerGrpcMaxSizeBytes: 104857600
  gloo:
    enableRestEds: false
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    proxyDebugBindAddr: 0.0.0.0:9966
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    invalidConfigPolicy:
      replaceInvalidRoutes: false
      invalidRouteResponseBody: "Gloo Gateway has invalid configuration. Administrators should run ` + backtick + "glooctl check" + backtick + ` to find and fix config errors."
      invalidRouteResponseCode: 404
      replaceInvalidRoutes: false
  ratelimitServer:
    rateLimitBeforeAuth: false
    ratelimitServerRef:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
  consoleOptions:
    readOnly: false
    apiExplorerEnabled: true
  graphqlOptions:
    schemaChangeValidationOptions:
      rejectBreakingChanges: false
      processingRules: []
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})

			It("correctly sets the ext auth transport API version", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.transportApiVersion=V2",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				settings := makeUnstructured(`
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: ` + namespace + `
spec:
  discovery:
    fdsMode: WHITELIST
  extauth:
    transportApiVersion: V2
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
    userIdHeader: "x-user-id"
  gateway:
    enableGatewayController: true
    readGatewaysFromAllNamespaces: false
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
      disableTransformationValidation: false
      allowWarnings: true
      warnRouteShortCircuiting: false
      validationServerGrpcMaxSizeBytes: 104857600
  gloo:
    enableRestEds: false
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    proxyDebugBindAddr: 0.0.0.0:9966
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    invalidConfigPolicy:
      replaceInvalidRoutes: false
      invalidRouteResponseBody: "Gloo Gateway has invalid configuration. Administrators should run ` + backtick + "glooctl check" + backtick + ` to find and fix config errors."
      invalidRouteResponseCode: 404
      replaceInvalidRoutes: false
  ratelimitServer:
    rateLimitBeforeAuth: false
    ratelimitServerRef:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
  consoleOptions:
    readOnly: false
    apiExplorerEnabled: true
  graphqlOptions:
    schemaChangeValidationOptions:
      rejectBreakingChanges: false
      processingRules: []
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})

			It("enable default credentials", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"gloo.settings.aws.enableCredentialsDiscovery=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				settings := makeUnstructured(`
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: ` + namespace + `
spec:
  discovery:
    fdsMode: WHITELIST
  extauth:
    transportApiVersion: V3
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
    userIdHeader: "x-user-id"
  gateway:
    enableGatewayController: true
    readGatewaysFromAllNamespaces: false
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
      disableTransformationValidation: false
      allowWarnings: true
      warnRouteShortCircuiting: false
      validationServerGrpcMaxSizeBytes: 104857600
  gloo:
    enableRestEds: false
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    proxyDebugBindAddr: 0.0.0.0:9966
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    invalidConfigPolicy:
      replaceInvalidRoutes: false
      invalidRouteResponseBody: "Gloo Gateway has invalid configuration. Administrators should run ` + backtick + "glooctl check" + backtick + ` to find and fix config errors."
      invalidRouteResponseCode: 404
      replaceInvalidRoutes: false
    awsOptions:
      enableCredentialsDiscovey: true
  ratelimitServer:
    rateLimitBeforeAuth: false
    ratelimitServerRef:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
  consoleOptions:
    readOnly: false
    apiExplorerEnabled: true
  graphqlOptions:
    schemaChangeValidationOptions:
      rejectBreakingChanges: false
      processingRules: []
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})

			It("Allows ratelimit descriptors to be set", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"gloo.settings.rateLimit.descriptors[0].key=generic_key",
						"gloo.settings.rateLimit.descriptors[0].value=per-second",
						"gloo.settings.rateLimit.descriptors[0].rateLimit.requestsPerUnit=2",
						"gloo.settings.rateLimit.descriptors[0].rateLimit.unit=SECOND",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				settings := makeUnstructured(`
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: ` + namespace + `
spec:
  discovery:
    fdsMode: WHITELIST
  extauth:
    transportApiVersion: V3
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
    userIdHeader: "x-user-id"
  gateway:
    enableGatewayController: true
    readGatewaysFromAllNamespaces: false
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
      disableTransformationValidation: false
      allowWarnings: true
      warnRouteShortCircuiting: false
      validationServerGrpcMaxSizeBytes: 104857600
  gloo:
    enableRestEds: false
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    proxyDebugBindAddr: 0.0.0.0:9966
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    invalidConfigPolicy:
      replaceInvalidRoutes: false
      invalidRouteResponseBody: "Gloo Gateway has invalid configuration. Administrators should run ` + backtick + "glooctl check" + backtick + ` to find and fix config errors."
      invalidRouteResponseCode: 404
      replaceInvalidRoutes: false
  ratelimitServer:
    rateLimitBeforeAuth: false
    ratelimitServerRef:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
  ratelimit:
    descriptors:
      - key: generic_key
        value: "per-second"
        rateLimit:
          requestsPerUnit: 2
          unit: SECOND
  consoleOptions:
    readOnly: false
    apiExplorerEnabled: true
  graphqlOptions:
    schemaChangeValidationOptions:
      rejectBreakingChanges: false
      processingRules: []
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})

			It("enable sts discovery", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"gloo.settings.aws.enableServiceAccountCredentials=true",
						"gloo.settings.aws.stsCredentialsRegion=us-east-2",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				settings := makeUnstructured(`
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: ` + namespace + `
spec:
  discovery:
    fdsMode: WHITELIST
  extauth: 
    transportApiVersion: V3
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
    userIdHeader: "x-user-id"
  gateway:
    enableGatewayController: true
    readGatewaysFromAllNamespaces: false
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
      disableTransformationValidation: false
      allowWarnings: true
      warnRouteShortCircuiting: false
      validationServerGrpcMaxSizeBytes: 104857600
  gloo:
    enableRestEds: false
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    proxyDebugBindAddr: 0.0.0.0:9966
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    awsOptions:
      serviceAccountCredentials:
        cluster: aws_sts_cluster
        uri: sts.us-east-2.amazonaws.com
    invalidConfigPolicy:
      replaceInvalidRoutes: false
      invalidRouteResponseBody: "Gloo Gateway has invalid configuration. Administrators should run ` + backtick + "glooctl check" + backtick + ` to find and fix config errors."
      invalidRouteResponseCode: 404
      replaceInvalidRoutes: false
  ratelimitServer:
    rateLimitBeforeAuth: false
    ratelimitServerRef:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
  consoleOptions:
    readOnly: false
    apiExplorerEnabled: true
  graphqlOptions:
    schemaChangeValidationOptions:
      rejectBreakingChanges: false
      processingRules: []
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})

			It("finds resources on all sds and sidecar containers", func() {
				envoySidecarVals := []string{"100Mi", "200m", "300Mi", "400m"}
				sdsVals := []string{"101Mi", "201m", "301Mi", "401m"}

				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.glooMtls.enabled=true", // adds gloo/gateway proxy side containers
						fmt.Sprintf("global.glooMtls.envoySidecarResources.requests.memory=%s", envoySidecarVals[0]),
						fmt.Sprintf("global.glooMtls.envoySidecarResources.requests.cpu=%s", envoySidecarVals[1]),
						fmt.Sprintf("global.glooMtls.envoySidecarResources.limits.memory=%s", envoySidecarVals[2]),
						fmt.Sprintf("global.glooMtls.envoySidecarResources.limits.cpu=%s", envoySidecarVals[3]),
						fmt.Sprintf("global.glooMtls.sdsResources.requests.memory=%s", sdsVals[0]),
						fmt.Sprintf("global.glooMtls.sdsResources.requests.cpu=%s", sdsVals[1]),
						fmt.Sprintf("global.glooMtls.sdsResources.limits.memory=%s", sdsVals[2]),
						fmt.Sprintf("global.glooMtls.sdsResources.limits.cpu=%s", sdsVals[3]),
					},
				})
				Expect(err).NotTo(HaveOccurred())

				// get all deployments for arbitrary examination/testing
				var deployments []*unstructured.Unstructured
				testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
					if unstructured.GetKind() == "Deployment" {
						deployments = append(deployments, unstructured)
					}
					return true
				})
				count := 0

				for _, deployment := range deployments {
					if deployment.GetName() == "gloo" || deployment.GetName() == "gateway-proxy" {
						continue
					}
					// marshall unstructured object into deployment
					rawDeploy, err := deployment.MarshalJSON()
					Expect(err).NotTo(HaveOccurred())
					deploy := appsv1.Deployment{}
					err = json.Unmarshal(rawDeploy, &deploy)
					Expect(err).NotTo(HaveOccurred())

					// look for sidecar and sds containers, then test their resource values.
					for _, container := range deploy.Spec.Template.Spec.Containers {
						// still make sure non-sds/sidecar containers have non-nil resources, since all
						// other containers should have default resources values set in their templates.
						if container.Name == "envoy-sidecar" || container.Name == "sds" {
							var expectedVals = sdsVals
							if container.Name == "envoy-sidecar" {
								expectedVals = envoySidecarVals
							}
							fmt.Printf("\n%s/%s\n", deployment.GetName(), container.Name)

							Expect(container.Resources.Requests.Memory().String()).To(Equal(expectedVals[0]),
								"deployment/container %s/%s had incorrect request memory: expected %s, got %s",
								deployment.GetName(), container.Name, expectedVals[0], container.Resources.Requests.Memory().String())

							Expect(container.Resources.Requests.Cpu().String()).To(Equal(expectedVals[1]),
								"deployment/container %s/%s had incorrect request cpu: expected %s, got %s",
								deployment.GetName(), container.Name, expectedVals[1], container.Resources.Requests.Cpu().String())

							Expect(container.Resources.Limits.Memory().String()).To(Equal(expectedVals[2]),
								"deployment/container %s/%s had incorrect limit memory: expected %s, got %s",
								deployment.GetName(), container.Name, expectedVals[2], container.Resources.Limits.Memory().String())

							Expect(container.Resources.Limits.Cpu().String()).To(Equal(expectedVals[3]),
								"deployment/container %s/%s had incorrect limit cpu: expected %s, got %s",
								deployment.GetName(), container.Name, expectedVals[3], container.Resources.Limits.Cpu().String())
							count += 1
						}
					}
				}
			})

			It("creates a deployment without extauth sidecar", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
				Expect(err).NotTo(HaveOccurred())
				testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
			})

			It("creates a deployment with extauth sidecar", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.envoySidecar=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				gatewayProxyDeployment.Spec.Template.Spec.Volumes = append(
					gatewayProxyDeployment.Spec.Template.Spec.Volumes,
					v1.Volume{
						Name: "shared-data",
						VolumeSource: v1.VolumeSource{
							EmptyDir: &v1.EmptyDirVolumeSource{},
						},
					})

				gatewayProxyDeployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
					gatewayProxyDeployment.Spec.Template.Spec.Containers[0].VolumeMounts,
					v1.VolumeMount{
						Name:      "shared-data",
						MountPath: "/usr/share/shared-data",
					})

				gatewayProxyDeployment.Spec.Template.Spec.Containers = append(
					gatewayProxyDeployment.Spec.Template.Spec.Containers,
					v1.Container{
						Name:            "extauth",
						Image:           "quay.io/solo-io/extauth-ee:" + version,
						Ports:           nil,
						ImagePullPolicy: getPullPolicy(),
						Env: []v1.EnvVar{
							{
								Name: "POD_NAMESPACE",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										FieldPath: "metadata.namespace",
									},
								},
							},
							{
								Name:  "SERVICE_NAME",
								Value: "ext-auth",
							},
							{
								Name:  "GLOO_ADDRESS",
								Value: "gloo:9977",
							},
							{
								Name: "SIGNING_KEY",
								ValueFrom: &v1.EnvVarSource{
									SecretKeyRef: &v1.SecretKeySelector{
										LocalObjectReference: v1.LocalObjectReference{
											Name: "extauth-signing-key",
										},
										Key: "signing-key",
									},
								},
							},
							{
								Name: "REDIS_PASSWORD",
								ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "redis",
									},
									Key: "redis-password",
								}},
							},
							{
								Name:  "SERVER_PORT",
								Value: "8083",
							},
							{
								Name:  "UDS_ADDR",
								Value: "/usr/share/shared-data/.sock",
							},
							{
								Name:  "USER_ID_HEADER",
								Value: "x-user-id",
							},
							{
								Name:  "START_STATS_SERVER",
								Value: "true",
							},
							{
								Name:  "HEALTH_HTTP_PORT",
								Value: "8082",
							},
							{
								Name:  "HEALTH_HTTP_PATH",
								Value: "/healthcheck",
							},
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "shared-data",
								MountPath: "/usr/share/shared-data",
							},
						},
						ReadinessProbe: &v1.Probe{
							Handler: v1.Handler{
								HTTPGet: &v1.HTTPGetAction{
									Path: "/healthcheck",
									Port: intstr.IntOrString{
										Type:   0,
										IntVal: 8082,
									},
								},
							},
							InitialDelaySeconds: 2,
							PeriodSeconds:       5,
							FailureThreshold:    2,
							SuccessThreshold:    1,
						},
					})

				testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
			})

			It("creates a deployment with extauth sidecar and extraVolumeMount", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.envoySidecar=true",
						"global.extensions.extAuth.deployment.extraVolume[0].Name=test-path",
						"global.extensions.extAuth.deployment.extraVolume[0].HostPath.Path=/var/run/sds",
						"global.extensions.extAuth.deployment.extraVolumeMount[0].Name=test-path",
						"global.extensions.extAuth.deployment.extraVolumeMount[0].MountPath=/var/run/sds",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				gatewayProxyDeployment.Spec.Template.Spec.Volumes = append(
					gatewayProxyDeployment.Spec.Template.Spec.Volumes,
					v1.Volume{
						Name: "shared-data",
						VolumeSource: v1.VolumeSource{
							EmptyDir: &v1.EmptyDirVolumeSource{},
						},
					})

				gatewayProxyDeployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
					gatewayProxyDeployment.Spec.Template.Spec.Containers[0].VolumeMounts,
					v1.VolumeMount{
						Name:      "shared-data",
						MountPath: "/usr/share/shared-data",
					})

				gatewayProxyDeployment.Spec.Template.Spec.Containers = append(
					gatewayProxyDeployment.Spec.Template.Spec.Containers,
					v1.Container{
						Name:            "extauth",
						Image:           "quay.io/solo-io/extauth-ee:" + version,
						Ports:           nil,
						ImagePullPolicy: getPullPolicy(),
						Env: []v1.EnvVar{
							{
								Name: "POD_NAMESPACE",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										FieldPath: "metadata.namespace",
									},
								},
							},
							{
								Name:  "SERVICE_NAME",
								Value: "ext-auth",
							},
							{
								Name:  "GLOO_ADDRESS",
								Value: "gloo:9977",
							},
							{
								Name: "SIGNING_KEY",
								ValueFrom: &v1.EnvVarSource{
									SecretKeyRef: &v1.SecretKeySelector{
										LocalObjectReference: v1.LocalObjectReference{
											Name: "extauth-signing-key",
										},
										Key: "signing-key",
									},
								},
							},
							{
								Name: "REDIS_PASSWORD",
								ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "redis",
									},
									Key: "redis-password",
								}},
							},
							{
								Name:  "SERVER_PORT",
								Value: "8083",
							},
							{
								Name:  "UDS_ADDR",
								Value: "/usr/share/shared-data/.sock",
							},
							{
								Name:  "USER_ID_HEADER",
								Value: "x-user-id",
							},
							{
								Name:  "START_STATS_SERVER",
								Value: "true",
							},
							{
								Name:  "HEALTH_HTTP_PORT",
								Value: "8082",
							},
							{
								Name:  "HEALTH_HTTP_PATH",
								Value: "/healthcheck",
							},
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "test-path",
								MountPath: "/var/run/sds",
							},
							{
								Name:      "shared-data",
								MountPath: "/usr/share/shared-data",
							},
						},
						ReadinessProbe: &v1.Probe{
							Handler: v1.Handler{
								HTTPGet: &v1.HTTPGetAction{
									Path: "/healthcheck",
									Port: intstr.IntOrString{
										Type:   0,
										IntVal: 8082,
									},
								},
							},
							InitialDelaySeconds: 2,
							PeriodSeconds:       5,
							FailureThreshold:    2,
							SuccessThreshold:    1,
						},
					})

				testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
			})

			Context("pass image pull secrets", func() {
				pullSecretName := "test-pull-secret"
				pullSecret := []v1.LocalObjectReference{
					{Name: pullSecretName},
				}

				It("via global values", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{fmt.Sprintf("global.image.pullSecret=%s", pullSecretName)},
					})
					gatewayProxyDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
					testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
					Expect(err).NotTo(HaveOccurred())
				})

				It("via podTemplate values", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{
							fmt.Sprintf("gloo.gatewayProxies.gatewayProxy.podTemplate.image.pullSecret=%s", pullSecretName),
						},
					})
					Expect(err).NotTo(HaveOccurred())

					gatewayProxyDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
					testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
				})

				It("podTemplate values win over global", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{
							"global.image.pullSecret=wrong",
							fmt.Sprintf("gloo.gatewayProxies.gatewayProxy.podTemplate.image.pullSecret=%s", pullSecretName),
						},
					})
					Expect(err).NotTo(HaveOccurred())
					gatewayProxyDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
					testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
				})

			})
		})

		Context("stats settings", func() {

			It("exposes http-monitoring port on all relevant services", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"gloo-fed.enabled=true",
						"gloo-fed.glooFedApiserver.enable=true",

						"gloo-fed.glooFedApiserver.stats.serviceMonitorEnabled=true",
						"gloo-fed.glooFed.stats.serviceMonitorEnabled=true",

						"global.glooStats.enabled=true",
						"global.glooStats.serviceMonitorEnabled=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				expectedServicesWithHttpMonitoring := []string{
					"gloo-fed-console",
					"gloo",
					"discovery",
					"gateway-proxy-monitoring-service",
					"extauth",
					"rate-limit",
					"observability",
				}
				var actualServicesWithHttpMonitoring []string

				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "Service"
				}).ExpectAll(func(service *unstructured.Unstructured) {
					serviceObject, err := kuberesource.ConvertUnstructured(service)
					ExpectWithOffset(1, err).NotTo(HaveOccurred(), fmt.Sprintf("Service %+v should be able to convert from unstructured", service))
					structuredService, ok := serviceObject.(*v1.Service)
					ExpectWithOffset(1, ok).To(BeTrue(), fmt.Sprintf("Service %+v should be able to cast to a structured service", service))

					for _, servicePort := range structuredService.Spec.Ports {
						if servicePort.Name == "http-monitoring" {
							actualServicesWithHttpMonitoring = append(actualServicesWithHttpMonitoring, structuredService.GetName())
						}
					}
				})

				Expect(actualServicesWithHttpMonitoring).To(Equal(expectedServicesWithHttpMonitoring))
			})

			It("exposes http-monitoring port on all relevant deployments", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"gloo-fed.enabled=true",
						"gloo-fed.glooFedApiserver.enable=true",

						"gloo-fed.glooFedApiserver.stats.podMonitorEnabled=true",
						"gloo-fed.glooFed.stats.podMonitorEnabled=true",

						"global.glooStats.enabled=true",
						"global.glooStats.podMonitorEnabled=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				expectedDeploymentsWithHttpMonitoring := []string{
					"gloo-fed-console",
					"gloo-fed",
					"gloo",
					"discovery",
					"gateway-proxy",
					"rate-limit",
					"extauth",
					"observability",
				}
				var actualDeploymentsWithHttpMonitoring []string

				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "Deployment"
				}).ExpectAll(func(deployment *unstructured.Unstructured) {
					deploymentObject, err := kuberesource.ConvertUnstructured(deployment)
					ExpectWithOffset(1, err).NotTo(HaveOccurred(), fmt.Sprintf("Deployment %+v should be able to convert from unstructured", deployment))
					structuredDeployment, ok := deploymentObject.(*appsv1.Deployment)
					ExpectWithOffset(1, ok).To(BeTrue(), fmt.Sprintf("Deployment %+v should be able to cast to a structured deployment", deployment))

					for _, container := range structuredDeployment.Spec.Template.Spec.Containers {
						for _, containerPort := range container.Ports {
							if containerPort.Name == "http-monitoring" {
								actualDeploymentsWithHttpMonitoring = append(actualDeploymentsWithHttpMonitoring, structuredDeployment.GetName())
							}
						}
					}
				})

				Expect(actualDeploymentsWithHttpMonitoring).To(Equal(expectedDeploymentsWithHttpMonitoring))
			})

		})

		Context("gloo mtls settings", func() {
			var (
				glooMtlsSecretVolume = v1.Volume{
					Name: "gloo-mtls-certs",
					VolumeSource: v1.VolumeSource{
						Secret: &v1.SecretVolumeSource{
							SecretName:  "gloo-mtls-certs",
							Items:       nil,
							DefaultMode: proto.Int(420),
						},
					},
				}

				haveEnvoySidecar = func(containers []v1.Container) bool {
					for _, c := range containers {
						if c.Name == "envoy-sidecar" {
							return true
						}
					}
					return false
				}

				haveSdsSidecar = func(containers []v1.Container) bool {
					for _, c := range containers {
						if c.Name == "sds" {
							return true
						}
					}
					return false
				}

				haveEnvVariable = func(containers []v1.Container, containerName, env, value string) bool {
					for _, c := range containers {
						if c.Name == containerName {
							Expect(c.Env).To(ContainElement(v1.EnvVar{Name: env, Value: value}))
							return true
						}
					}
					return false
				}
			)

			It("should add or change the correct components in the resulting helm chart", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"global.glooMtls.enabled=true"},
				})
				Expect(err).NotTo(HaveOccurred())

				foundGlooMtlsCertgenJob := false
				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "Job"
				}).ExpectAll(func(job *unstructured.Unstructured) {
					jobObject, err := kuberesource.ConvertUnstructured(job)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Job %+v should be able to convert from unstructured", job))
					structuredJob, ok := jobObject.(*jobsv1.Job)
					Expect(ok).To(BeTrue(), fmt.Sprintf("Job %+v should be able to cast to a structured job", job))

					if structuredJob.GetName() == "gloo-mtls-certgen" {
						foundGlooMtlsCertgenJob = true
					}
				})
				Expect(foundGlooMtlsCertgenJob).To(BeTrue(), "Did not find the gloo-mtls-certgen job")

				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "Deployment"
				}).ExpectAll(func(deployment *unstructured.Unstructured) {
					deploymentObject, err := kuberesource.ConvertUnstructured(deployment)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Deployment %+v should be able to convert from unstructured", deployment))
					structuredDeployment, ok := deploymentObject.(*appsv1.Deployment)
					Expect(ok).To(BeTrue(), fmt.Sprintf("Deployment %+v should be able to cast to a structured deployment", deployment))

					if structuredDeployment.GetName() == "gloo" {
						Expect(haveEnvoySidecar(structuredDeployment.Spec.Template.Spec.Containers)).To(BeTrue())
						Expect(haveSdsSidecar(structuredDeployment.Spec.Template.Spec.Containers)).To(BeTrue())
						Expect(structuredDeployment.Spec.Template.Spec.Volumes).To(ContainElement(glooMtlsSecretVolume))
					}

					if structuredDeployment.GetName() == "gateway-proxy" {
						Expect(haveSdsSidecar(structuredDeployment.Spec.Template.Spec.Containers)).To(BeTrue())
						Expect(structuredDeployment.Spec.Template.Spec.Volumes).To(ContainElement(glooMtlsSecretVolume))
					}

					// should add envoy, sds sidecars to the Extauth and Rate-Limit Deployment
					if structuredDeployment.GetName() == "rate-limit" {
						Expect(haveEnvoySidecar(structuredDeployment.Spec.Template.Spec.Containers)).To(BeTrue())
						Expect(haveSdsSidecar(structuredDeployment.Spec.Template.Spec.Containers)).To(BeTrue())
						Expect(haveEnvVariable(structuredDeployment.Spec.Template.Spec.Containers,
							"rate-limit", "GLOO_ADDRESS", "127.0.0.1:9955")).To(BeTrue())
						Expect(structuredDeployment.Spec.Template.Spec.Volumes).To(ContainElement(glooMtlsSecretVolume))
					}

					if structuredDeployment.GetName() == "extauth" {
						Expect(haveEnvoySidecar(structuredDeployment.Spec.Template.Spec.Containers)).To(BeTrue())
						Expect(haveSdsSidecar(structuredDeployment.Spec.Template.Spec.Containers)).To(BeTrue())
						Expect(haveEnvVariable(structuredDeployment.Spec.Template.Spec.Containers,
							"extauth", "GLOO_ADDRESS", "127.0.0.1:9955")).To(BeTrue())
						Expect(haveEnvVariable(structuredDeployment.Spec.Template.Spec.Containers,
							"extauth", "SERVER_PORT", "8084")).To(BeTrue())
						Expect(structuredDeployment.Spec.Template.Spec.Volumes).To(ContainElement(glooMtlsSecretVolume))
					}

					if structuredDeployment.GetName() == "api-server" {
						Expect(haveSdsSidecar(structuredDeployment.Spec.Template.Spec.Containers)).To(BeTrue())
						Expect(structuredDeployment.Spec.Template.Spec.Volumes).To(ContainElement(glooMtlsSecretVolume))
						Expect(structuredDeployment.Spec.Template.Spec.Containers[2].ReadinessProbe.HTTPGet.Scheme).To(Equal(v1.URISchemeHTTPS))
						Expect(structuredDeployment.Spec.Template.Spec.Containers[2].ReadinessProbe.HTTPGet.Port).To(Equal(intstr.IntOrString{IntVal: 8443}))
					}
				})

				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "Service"
				}).ExpectAll(func(svc *unstructured.Unstructured) {
					serviceObj, err := kuberesource.ConvertUnstructured(svc)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Service %+v should be able to convert from unstructured", svc))
					structuredService, ok := serviceObj.(*v1.Service)
					Expect(ok).To(BeTrue(), fmt.Sprintf("Service %+v should be able to cast to a structured service", svc))

					if structuredService.GetName() == "apiserver-ui" {
						Expect(structuredService.Spec.Ports[0].Port).To(BeEquivalentTo(8443))
					}
				})

				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "ConfigMap"
				}).ExpectAll(func(cfgmap *unstructured.Unstructured) {
					cmObj, err := kuberesource.ConvertUnstructured(cfgmap)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("ConfigMap %+v should be able to convert from unstructured", cfgmap))
					structuredConfigMap, ok := cmObj.(*v1.ConfigMap)
					Expect(ok).To(BeTrue(), fmt.Sprintf("ConfigMap %+v should be able to cast to a structured config map", cfgmap))

					if structuredConfigMap.GetName() == "default-apiserver-envoy-config" {
						bootstrap := bootstrapv3.Bootstrap{}
						Expect(structuredConfigMap.Data["config.yaml"]).NotTo(BeEmpty())
						jsn, err := yaml.YAMLToJSON([]byte(structuredConfigMap.Data["config.yaml"]))
						if err != nil {
							Expect(err).NotTo(HaveOccurred())
						}
						err = jsonpb.Unmarshal(bytes.NewReader(jsn), &bootstrap)
						Expect(err).NotTo(HaveOccurred())
						Expect(bootstrap.Node).To(MatchProto(&corev3.Node{Id: "sds_client", Cluster: "sds_client"}))
						Expect(bootstrap.StaticResources.Listeners[0].FilterChains[0].TransportSocket).NotTo(BeNil())
						tlsContext := tlsv3.DownstreamTlsContext{}
						Expect(ptypes.UnmarshalAny(bootstrap.StaticResources.Listeners[0].FilterChains[0].TransportSocket.GetTypedConfig(), &tlsContext)).NotTo(HaveOccurred())
						Expect(&tlsContext).To(MatchProto(&tlsv3.DownstreamTlsContext{
							CommonTlsContext: &tlsv3.CommonTlsContext{
								TlsCertificateSdsSecretConfigs: []*tlsv3.SdsSecretConfig{
									{
										Name: "server_cert",
										SdsConfig: &corev3.ConfigSource{
											ResourceApiVersion: corev3.ApiVersion_V3,
											ConfigSourceSpecifier: &corev3.ConfigSource_ApiConfigSource{
												ApiConfigSource: &corev3.ApiConfigSource{
													ApiType:             corev3.ApiConfigSource_GRPC,
													TransportApiVersion: corev3.ApiVersion_V3,
													GrpcServices: []*corev3.GrpcService{
														{
															TargetSpecifier: &corev3.GrpcService_EnvoyGrpc_{
																EnvoyGrpc: &corev3.GrpcService_EnvoyGrpc{
																	ClusterName: "gloo_client_sds",
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						}))
					}
				})
			})

			It("should add an additional listener to the gateway-proxy-envoy-config for extauth sidecar", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"global.glooMtls.enabled=true,global.extensions.extAuth.envoySidecar=true"},
				})
				Expect(err).NotTo(HaveOccurred())

				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "ConfigMap"
				}).ExpectAll(func(configMap *unstructured.Unstructured) {
					configMapObject, err := kuberesource.ConvertUnstructured(configMap)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("ConfigMap %+v should be able to convert from unstructured", configMap))
					structuredConfigMap, ok := configMapObject.(*v1.ConfigMap)
					Expect(ok).To(BeTrue(), fmt.Sprintf("ConfigMap %+v should be able to cast to a structured config map", configMap))

					if structuredConfigMap.GetName() == "gateway-proxy-envoy-config" {
						expectedGlooMtlsListener := "    - name: gloo_xds_mtls_listener"
						Expect(structuredConfigMap.Data["envoy.yaml"]).To(ContainSubstring(expectedGlooMtlsListener))
					}
				})
			})

			It("should allow extauth service to handle TLS itself", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"global.extensions.extAuth.tlsEnabled=true,global.extensions.extAuth.certPath=/path/to/custom/cert.crt,global.extensions.extAuth.keyPath=/path/to/custom/key.key"},
				})
				Expect(err).NotTo(HaveOccurred())

				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "Deployment"
				}).ExpectAll(func(deployment *unstructured.Unstructured) {
					deploymentObject, err := kuberesource.ConvertUnstructured(deployment)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Deployment %+v should be able to convert from unstructured", deployment))
					structuredDeployment, ok := deploymentObject.(*appsv1.Deployment)
					Expect(ok).To(BeTrue(), fmt.Sprintf("Deployment %+v should be able to cast to a structured deployment", deployment))

					if structuredDeployment.GetName() == "extauth" {
						Expect(structuredDeployment.Spec.Template.Spec.Containers[0].Env).To(ContainElement(v1.EnvVar{Name: "TLS_ENABLED", Value: "true"}))
						Expect(structuredDeployment.Spec.Template.Spec.Containers[0].Env).To(ContainElement(v1.EnvVar{Name: "CERT_PATH", Value: "/path/to/custom/cert.crt"}))
						Expect(structuredDeployment.Spec.Template.Spec.Containers[0].Env).To(ContainElement(v1.EnvVar{Name: "KEY_PATH", Value: "/path/to/custom/key.key"}))
					}
				})
			})

			It("should allow extauth service to handle TLS itself using a kubernetes secret", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.tlsEnabled=true",
						"global.extensions.extAuth.secretName=my-secret",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "Deployment"
				}).ExpectAll(func(deployment *unstructured.Unstructured) {
					deploymentObject, err := kuberesource.ConvertUnstructured(deployment)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Deployment %+v should be able to convert from unstructured", deployment))
					structuredDeployment, ok := deploymentObject.(*appsv1.Deployment)
					Expect(ok).To(BeTrue(), fmt.Sprintf("Deployment %+v should be able to cast to a structured deployment", deployment))

					if structuredDeployment.GetName() == "extauth" {
						Expect(structuredDeployment.Spec.Template.Spec.Containers[0].Env).To(ContainElement(v1.EnvVar{Name: "TLS_ENABLED", Value: "true"}))
						Expect(structuredDeployment.Spec.Template.Spec.Containers[0].Env).To(ContainElement(v1.EnvVar{Name: "CERT",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: "my-secret"}, Key: "tls.crt"},
							}}))
						Expect(structuredDeployment.Spec.Template.Spec.Containers[0].Env).To(ContainElement(v1.EnvVar{Name: "KEY",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: "my-secret"}, Key: "tls.key"},
							}}))
					}
				})
			})

		})

		Context("redis deployment", func() {

			var expectedDeployment *appsv1.Deployment

			BeforeEach(func() {
				labels = map[string]string{
					"app":  "gloo",
					"gloo": "redis",
				}
				selector = map[string]string{
					"gloo": "redis",
				}

				rb := ResourceBuilder{
					Namespace: namespace,
					Name:      "redis",
					Labels:    labels,
				}

				nonRootUser := int64(999)
				nonRoot := true
				nonRootSC := &v1.PodSecurityContext{
					RunAsUser:    &nonRootUser,
					RunAsGroup:   &nonRootUser,
					RunAsNonRoot: &nonRoot,
					FSGroup:      &nonRootUser,
				}

				mode := int32(420)
				volumes := []v1.Volume{
					{
						Name: "data",
						VolumeSource: v1.VolumeSource{
							EmptyDir: &v1.EmptyDirVolumeSource{},
						},
					},
					{
						Name: "conf",
						VolumeSource: v1.VolumeSource{
							EmptyDir: &v1.EmptyDirVolumeSource{},
						},
					},
					{
						Name: "acl",
						VolumeSource: v1.VolumeSource{
							Secret: &v1.SecretVolumeSource{
								SecretName:  "redis",
								DefaultMode: &mode,
							},
						},
					},
				}

				expectedDeployment = rb.GetDeploymentAppsv1()
				expectedDeployment.Spec.Replicas = nil // GetDeploymentAppsv1 explicitly sets it to 1, which we don't want
				expectedDeployment.Spec.Template.Spec.InitContainers = []v1.Container{
					{
						Name:    "createconf",
						Image:   "docker.io/busybox:1.28",
						Command: []string{"/bin/sh", "-c", "echo 'aclfile /redis-acl/users.acl' > /conf/redis.conf"},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "conf",
								MountPath: "/conf",
							},
						},
					},
				}
				expectedDeployment.Spec.Template.Spec.Containers = []v1.Container{
					{
						Name:            "redis",
						Image:           "docker.io/redis:6.2.4",
						ImagePullPolicy: getPullPolicy(),
						Args:            []string{"redis-server", "/redis-acl/users.acl"},
						Ports: []v1.ContainerPort{
							{
								ContainerPort: 6379,
							},
						},
						Env: []v1.EnvVar{
							{
								Name:  "MASTER",
								Value: "true",
							},
						},
						VolumeMounts: []v1.VolumeMount{
							{
								MountPath: "/redis-master-data",
								Name:      "data",
							},
							{
								MountPath: "/redis-acl",
								Name:      "acl",
							},
							{
								MountPath: "/conf",
								Name:      "conf",
							},
						},
					},
				}
				expectedDeployment.Spec.Strategy = appsv1.DeploymentStrategy{}
				expectedDeployment.Spec.Selector.MatchLabels = selector
				expectedDeployment.Spec.Template.ObjectMeta.Labels = selector

				expectedDeployment.Spec.Template.Spec.SecurityContext = nonRootSC
				expectedDeployment.Spec.Template.Spec.Volumes = volumes
			})

			It("produces expected default deployment", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
				Expect(err).NotTo(HaveOccurred())

				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("can disable pod security context", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"redis.deployment.enablePodSecurityContext=false",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				expectedDeployment.Spec.Template.Spec.SecurityContext = nil
				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("can override redis images", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"redis.deployment.image.registry=X",
						"redis.deployment.image.repository=Y",
						"redis.deployment.image.tag=Z",
						"redis.deployment.initContainer.image.registry=A",
						"redis.deployment.initContainer.image.repository=B",
						"redis.deployment.initContainer.image.tag=C",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				expectedDeployment.Spec.Template.Spec.Containers[0].Image = "X/Y:Z"
				expectedDeployment.Spec.Template.Spec.InitContainers[0].Image = "A/B:C"
				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})
		})

		Context("redis scaled with client-side sharding", func() {
			var (
				glooMtlsSecretVolume = v1.Volume{
					Name: "gloo-mtls-certs",
					VolumeSource: v1.VolumeSource{
						Secret: &v1.SecretVolumeSource{
							SecretName:  "gloo-mtls-certs",
							Items:       nil,
							DefaultMode: proto.Int(420),
						},
					},
				}

				haveEnvoySidecar = func(containers []v1.Container) bool {
					for _, c := range containers {
						if c.Name == "envoy-sidecar" {
							return true
						}
					}
					return false
				}

				haveSdsSidecar = func(containers []v1.Container) bool {
					for _, c := range containers {
						if c.Name == "sds" {
							return true
						}
					}
					return false
				}

				haveEnvVariable = func(containers []v1.Container, containerName, env, value string) bool {
					for _, c := range containers {
						if c.Name == containerName {
							Expect(c.Env).To(ContainElement(v1.EnvVar{Name: env, Value: value}))
							return true
						}
					}
					return false
				}
			)

			It("should add or change the correct components in the resulting helm chart", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.rateLimit.enabled=true",
						"redis.clientSideShardingEnabled=true",
						"redis.deployment.replicas=2",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "Deployment"
				}).ExpectAll(func(deployment *unstructured.Unstructured) {
					deploymentObject, err := kuberesource.ConvertUnstructured(deployment)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Deployment %+v should be able to convert from unstructured", deployment))
					structuredDeployment, ok := deploymentObject.(*appsv1.Deployment)
					Expect(ok).To(BeTrue(), fmt.Sprintf("Deployment %+v should be able to cast to a structured deployment", deployment))

					// should add envoy, but not sds sidecar to the Rate-Limit Deployment
					if structuredDeployment.GetName() == "rate-limit" {
						Expect(haveEnvoySidecar(structuredDeployment.Spec.Template.Spec.Containers)).To(BeTrue(), "should have envoy sidecar")
						Expect(haveSdsSidecar(structuredDeployment.Spec.Template.Spec.Containers)).To(BeFalse(), "shouldn't have SDS sidecar")
						Expect(structuredDeployment.Spec.Template.Spec.Containers).To(HaveLen(2), "should have exactly 2 containers")
						Expect(haveEnvVariable(structuredDeployment.Spec.Template.Spec.Containers,
							"rate-limit", "REDIS_URL", "/var/run/envoy/ratelimit.sock")).To(BeTrue(), "should use unix socket for redis url")
						Expect(haveEnvVariable(structuredDeployment.Spec.Template.Spec.Containers,
							"rate-limit", "REDIS_SOCKET_TYPE", "unix")).To(BeTrue(), "should use unix socket for redis url")
						Expect(structuredDeployment.Spec.Template.Spec.Volumes).NotTo(ContainElement(glooMtlsSecretVolume))
					}

					// Extauth deployment should not have SDS or envoy sidecars
					if structuredDeployment.GetName() == "extauth" {
						Expect(haveEnvoySidecar(structuredDeployment.Spec.Template.Spec.Containers)).To(BeFalse())
						Expect(haveSdsSidecar(structuredDeployment.Spec.Template.Spec.Containers)).To(BeFalse())
						Expect(structuredDeployment.Spec.Template.Spec.Volumes).NotTo(ContainElement(glooMtlsSecretVolume))
					}
				})

				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "Service"
				}).ExpectAll(func(svc *unstructured.Unstructured) {
					serviceObj, err := kuberesource.ConvertUnstructured(svc)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Service %+v should be able to convert from unstructured", svc))
					structuredService, ok := serviceObj.(*v1.Service)
					Expect(ok).To(BeTrue(), fmt.Sprintf("Service %+v should be able to cast to a structured service", svc))

					if structuredService.GetName() == "redis" {
						Expect(structuredService.Spec.ClusterIP).To(BeEquivalentTo("None"), "ClusterIP should be 'None' to indicate headless service")
					}
				})

				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "ConfigMap"
				}).ExpectAll(func(cfgmap *unstructured.Unstructured) {
					cmObj, err := kuberesource.ConvertUnstructured(cfgmap)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("ConfigMap %+v should be able to convert from unstructured", cfgmap))
					structuredConfigMap, ok := cmObj.(*v1.ConfigMap)
					Expect(ok).To(BeTrue(), fmt.Sprintf("ConfigMap %+v should be able to cast to a structured config map", cfgmap))

					if structuredConfigMap.GetName() == "rate-limit-sidecar-config" {
						bootstrap := bootstrapv3.Bootstrap{}
						Expect(structuredConfigMap.Data["envoy-sidecar.yaml"]).NotTo(BeEmpty())
						jsn, err := yaml.YAMLToJSON([]byte(structuredConfigMap.Data["envoy-sidecar.yaml"]))
						if err != nil {
							Expect(err).NotTo(HaveOccurred(), "could not parse envoy sidecar config yaml")
						}
						err = jsonpb.Unmarshal(bytes.NewReader(jsn), &bootstrap)
						Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("unable to unmarshal from json to pb - \n%v\n", string(jsn)))
						Expect(bootstrap.StaticResources.Listeners[0].Name).To(Equal("redis_listener"), "the sidecar envoy should have a redis listener")
						Expect(bootstrap.StaticResources.Clusters[0].Name).To(Equal("redis_cluster"), "the sidecar envoy should have a redis cluster")
						Expect(bootstrap.StaticResources.Clusters[0].LbPolicy).To(Equal(clusterv3.Cluster_MAGLEV), "it should use the maglev algorithm for load balancing")
					}
				})
			})
		})

		Context("rate-limit deployment", func() {

			var expectedDeployment *appsv1.Deployment

			BeforeEach(func() {
				labels = map[string]string{
					"app":  "gloo",
					"gloo": "rate-limit",
				}
				selector = map[string]string{
					"gloo": "rate-limit",
				}

				rb := ResourceBuilder{
					Namespace: namespace,
					Name:      "rate-limit",
					Labels:    labels,
				}

				nonRootUser := int64(10101)
				nonRoot := true

				nonRootSC := &v1.PodSecurityContext{
					RunAsUser:    &nonRootUser,
					RunAsNonRoot: &nonRoot,
				}

				expectedDeployment = rb.GetDeploymentAppsv1()

				expectedDeployment.Spec.Replicas = aws.Int32(1)
				expectedDeployment.Spec.Template.Spec.Containers = []v1.Container{
					{
						Name:            "rate-limit",
						Image:           "quay.io/solo-io/rate-limit-ee:" + version,
						ImagePullPolicy: getPullPolicy(),
						Env: []v1.EnvVar{
							{
								Name: "POD_NAMESPACE",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										FieldPath: "metadata.namespace",
									},
								},
							},
							{
								Name:  "GLOO_ADDRESS",
								Value: "gloo:9977",
							},
							statsEnvVar,
							{
								Name:  "REDIS_URL",
								Value: "redis:6379",
							},
							{
								Name:  "REDIS_SOCKET_TYPE",
								Value: "tcp",
							},
							{
								Name: "REDIS_PASSWORD",
								ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "redis",
									},
									Key: "redis-password",
								}},
							},
							{
								Name:  "READY_PORT_HTTP",
								Value: "18080",
							},
							{
								Name:  "READY_PATH_HTTP",
								Value: "/ready",
							},
						},
						ReadinessProbe: &v1.Probe{
							Handler: v1.Handler{
								HTTPGet: &v1.HTTPGetAction{
									Path: "/ready",
									Port: intstr.IntOrString{
										Type:   0,
										IntVal: 18080,
									},
								},
							},
							InitialDelaySeconds: 2,
							PeriodSeconds:       5,
							FailureThreshold:    2,
							SuccessThreshold:    1,
						},
						Resources: v1.ResourceRequirements{},
					},
				}
				expectedDeployment.Spec.Strategy = appsv1.DeploymentStrategy{}
				expectedDeployment.Spec.Selector.MatchLabels = selector
				expectedDeployment.Spec.Template.ObjectMeta.Labels = selector
				expectedDeployment.Spec.Template.ObjectMeta.Annotations = normalPromAnnotations

				expectedDeployment.Spec.Template.Spec.SecurityContext = nonRootSC

				expectedDeployment.Spec.Template.Spec.ServiceAccountName = "rate-limit"

				expectedDeployment.Spec.Replicas = nil // GetDeploymentAppsv1 explicitly sets it to 1, which we don't want
			})

			It("produces expected default deployment", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
				Expect(err).NotTo(HaveOccurred())

				actualDeployment := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
					return unstructured.GetKind() == "Deployment" && unstructured.GetLabels()["gloo"] == "rate-limit"
				})

				actualDeployment.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("should support setting the log level env", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"global.extensions.rateLimit.deployment.logLevel=DEBUG"},
				})
				Expect(err).NotTo(HaveOccurred())

				logLevel := "DEBUG"
				envs := expectedDeployment.Spec.Template.Spec.Containers[0].Env
				for i, env := range envs {
					if env.Name == "LOG_LEVEL" {
						envs[i].Value = logLevel
					}
				}

				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("should support getting fips images", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"global.image.fips=true"},
				})
				Expect(err).NotTo(HaveOccurred())
				expectedDeployment.Spec.Template.Spec.Containers[0].Image = "quay.io/solo-io/rate-limit-ee-fips:" + version
				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("should support setting beforeAuth", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"global.extensions.rateLimit.beforeAuth=true"},
				})
				Expect(err).NotTo(HaveOccurred())
				settings := makeUnstructured(`
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: ` + namespace + `
spec:
  discovery:
    fdsMode: WHITELIST
  extauth:
    transportApiVersion: V3
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
    userIdHeader: "x-user-id"
  gateway:
    enableGatewayController: true
    readGatewaysFromAllNamespaces: false
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
      disableTransformationValidation: false
      allowWarnings: true
      warnRouteShortCircuiting: false
      validationServerGrpcMaxSizeBytes: 104857600
  gloo:
    enableRestEds: false
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    proxyDebugBindAddr: 0.0.0.0:9966
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    invalidConfigPolicy:
      replaceInvalidRoutes: false
      invalidRouteResponseBody: "Gloo Gateway has invalid configuration. Administrators should run ` + backtick + "glooctl check" + backtick + ` to find and fix config errors."
      invalidRouteResponseCode: 404
      replaceInvalidRoutes: false
  ratelimitServer:
    rateLimitBeforeAuth: true
    ratelimitServerRef:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
  consoleOptions:
    readOnly: false
    apiExplorerEnabled: true
  graphqlOptions:
    schemaChangeValidationOptions:
      rejectBreakingChanges: false
      processingRules: []
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})

			Describe("affinity and antiAffinity", func() {

				It("affinity rules can be set", func() {
					helmOverrideFileContents := `
global:
  extensions:
    rateLimit:
      antiAffinity:
        podAffinity: 
          preferredDuringSchedulingIgnoredDuringExecution: 
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchLabels:
                  gloo: gateway-proxy
              topologyKey: kubernetes.io/hostname
`
					helmOverrideFile := "helm-override-*.yaml"
					tmpFile, err := ioutil.TempFile("", helmOverrideFile)
					Expect(err).ToNot(HaveOccurred())
					_, err = tmpFile.Write([]byte(helmOverrideFileContents))
					Expect(err).NotTo(HaveOccurred())
					defer tmpFile.Close()
					defer os.Remove(tmpFile.Name())
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesFile: tmpFile.Name(),
					})
					Expect(err).NotTo(HaveOccurred())

					actualDeployment := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						return unstructured.GetKind() == "Deployment" && unstructured.GetLabels()["gloo"] == "rate-limit"
					})

					expectedDeployment.Spec.Template.Spec.Affinity = nil
					expectedDeployment.Spec.Template.Spec.Affinity = &v1.Affinity{
						PodAffinity: &v1.PodAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
								{
									Weight: 100,
									PodAffinityTerm: v1.PodAffinityTerm{
										LabelSelector: &k8s.LabelSelector{
											MatchLabels: map[string]string{
												"gloo": "gateway-proxy",
											},
										},
										TopologyKey: "kubernetes.io/hostname",
									},
								},
							},
						},
					}
					actualDeployment.ExpectDeploymentAppsV1(expectedDeployment)
				})
				It("antiAffinity rules can be set", func() {
					helmOverrideFileContents := `
global:
  extensions:
    rateLimit:
      antiAffinity:
        podAntiAffinity: 
          preferredDuringSchedulingIgnoredDuringExecution: 
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchLabels:
                  gloo: gateway-proxy
              topologyKey: kubernetes.io/hostname
`
					helmOverrideFile := "helm-override-*.yaml"
					tmpFile, err := ioutil.TempFile("", helmOverrideFile)
					Expect(err).ToNot(HaveOccurred())
					_, err = tmpFile.Write([]byte(helmOverrideFileContents))
					Expect(err).NotTo(HaveOccurred())
					defer tmpFile.Close()
					defer os.Remove(tmpFile.Name())
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesFile: tmpFile.Name(),
					})
					Expect(err).NotTo(HaveOccurred())

					actualDeployment := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						return unstructured.GetKind() == "Deployment" && unstructured.GetLabels()["gloo"] == "rate-limit"
					})

					expectedDeployment.Spec.Template.Spec.Affinity = nil
					expectedDeployment.Spec.Template.Spec.Affinity = &v1.Affinity{
						PodAntiAffinity: &v1.PodAntiAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
								{
									Weight: 100,
									PodAffinityTerm: v1.PodAffinityTerm{
										LabelSelector: &k8s.LabelSelector{
											MatchLabels: map[string]string{
												"gloo": "gateway-proxy",
											},
										},
										TopologyKey: "kubernetes.io/hostname",
									},
								},
							},
						},
					}
					actualDeployment.ExpectDeploymentAppsV1(expectedDeployment)
				})
			})

			It("gwp pdb disabled by default", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
				Expect(err).To(BeNil())
				testManifest.ExpectUnstructured("PodDisruptionBudget", namespace, "rate-limit-pdb").To(BeNil())
			})

			It("can create gwp pdb with minAvailable", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.rateLimit.deployment.podDisruptionBudget.minAvailable=2",
					},
				})
				Expect(err).To(BeNil())

				pdb := makeUnstructured(`
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: rate-limit-pdb
  namespace: gloo-system
spec:
  minAvailable: 2
  selector:
    matchLabels:
      gloo: rate-limit
`)

				testManifest.ExpectUnstructured("PodDisruptionBudget", namespace, "rate-limit-pdb").To(BeEquivalentTo(pdb))
			})

			It("can create gwp pdb with maxUnavailable", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.rateLimit.deployment.podDisruptionBudget.maxUnavailable=2",
					},
				})
				Expect(err).To(BeNil())

				pdb := makeUnstructured(`
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: rate-limit-pdb
  namespace: gloo-system
spec:
  maxUnavailable: 2
  selector:
    matchLabels:
      gloo: rate-limit
`)

				testManifest.ExpectUnstructured("PodDisruptionBudget", namespace, "rate-limit-pdb").To(BeEquivalentTo(pdb))
			})

			It("can accept extra env vars", func() {
				expectedDeployment.Spec.Template.Spec.Containers[0].Env = append(
					expectedDeployment.Spec.Template.Spec.Containers[0].Env,
					v1.EnvVar{
						Name:  "TEST_EXTRA_ENV_VAR",
						Value: "test",
					})

				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.rateLimit.deployment.customEnv[0].Name=TEST_EXTRA_ENV_VAR",
						"global.extensions.rateLimit.deployment.customEnv[0].Value=test",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})
		})

		Context("caching deployment", func() {

			var expectedDeployment *appsv1.Deployment

			BeforeEach(func() {
				labels = map[string]string{
					"app":  "gloo",
					"gloo": "caching-service",
				}
				selector = map[string]string{
					"gloo": "caching-service",
				}

				rb := ResourceBuilder{
					Namespace: namespace,
					Name:      "caching-service",
					Labels:    labels,
				}

				nonRootUser := int64(10101)
				nonRoot := true

				nonRootSC := &v1.PodSecurityContext{
					RunAsUser:    &nonRootUser,
					RunAsNonRoot: &nonRoot,
				}

				expectedDeployment = rb.GetDeploymentAppsv1()

				expectedDeployment.Spec.Replicas = aws.Int32(1)
				expectedDeployment.Spec.Template.Spec.Containers = []v1.Container{
					{
						Name:            "caching-service",
						Image:           "quay.io/solo-io/caching-ee:" + version,
						ImagePullPolicy: getPullPolicy(),
						Env: []v1.EnvVar{
							{
								Name: "POD_NAMESPACE",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										FieldPath: "metadata.namespace",
									},
								},
							},
							{
								Name:  "REDIS_URL",
								Value: "redis:6379",
							},
							{
								Name:  "REDIS_SOCKET_TYPE",
								Value: "tcp",
							},
							{
								Name: "REDIS_PASSWORD",
								ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "redis",
									},
									Key: "redis-password",
								}},
							},
							{
								Name:  "START_STATS_SERVER",
								Value: "true",
							},
							{
								Name:  "READY_PORT",
								Value: "18080",
							},
							{
								Name:  "READY_PATH",
								Value: "/ready",
							},
						},
						ReadinessProbe: &v1.Probe{
							Handler: v1.Handler{
								HTTPGet: &v1.HTTPGetAction{
									Path: "/ready",
									Port: intstr.IntOrString{
										Type:   0,
										IntVal: 18080,
									},
								},
							},
							InitialDelaySeconds: 2,
							PeriodSeconds:       5,
							FailureThreshold:    2,
							SuccessThreshold:    1,
						},
						Resources: v1.ResourceRequirements{},
					},
				}
				expectedDeployment.Spec.Strategy = appsv1.DeploymentStrategy{}
				expectedDeployment.Spec.Selector.MatchLabels = selector
				expectedDeployment.Spec.Template.ObjectMeta.Labels = selector
				expectedDeployment.Spec.Template.ObjectMeta.Annotations = normalPromAnnotations

				expectedDeployment.Spec.Template.Spec.SecurityContext = nonRootSC

				expectedDeployment.Spec.Template.Spec.ServiceAccountName = "caching-service"

				expectedDeployment.Spec.Replicas = nil // GetDeploymentAppsv1 explicitly sets it to 1, which we don't want
			})

			It("produces expected default deployment", func() {
				//TODO: once rc3 oss is out pull that in so this succeeds in ci
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace,
					helmValues{valuesArgs: []string{"global.extensions.caching.enabled=true"}})
				Expect(err).NotTo(HaveOccurred())

				actualDeployment := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
					return unstructured.GetKind() == "Deployment" && unstructured.GetLabels()["gloo"] == "caching-service"
				})

				actualDeployment.ExpectDeploymentAppsV1(expectedDeployment)
			})

		})

		Context("gloo-fed apiserver deployment", func() {
			const defaultBootstrapConfigMapName = "gloo-fed-default-apiserver-envoy-config"

			var expectedDeployment *appsv1.Deployment

			BeforeEach(func() {
				labels = map[string]string{
					"app":      "gloo-fed",
					"gloo-fed": "console",
				}
				selector = map[string]string{
					"app":      "gloo-fed",
					"gloo-fed": "console",
				}

				podname := v1.EnvVar{
					Name: "POD_NAME",
					ValueFrom: &v1.EnvVarSource{
						FieldRef: &v1.ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				}

				truez := true
				falsez := false

				uiContainer := v1.Container{
					Name:            "console",
					Image:           "quay.io/solo-io/gloo-federation-console:" + version,
					ImagePullPolicy: getPullPolicy(),
					VolumeMounts: []v1.VolumeMount{
						{Name: "empty-cache", MountPath: "/var/cache/nginx"},
						{Name: "empty-run", MountPath: "/var/run"},
					},
					SecurityContext: &v1.SecurityContext{
						Capabilities: &v1.Capabilities{
							Drop: []v1.Capability{"ALL"},
						},
						RunAsNonRoot:             &truez,
						RunAsUser:                aws.Int64(101),
						ReadOnlyRootFilesystem:   &truez,
						AllowPrivilegeEscalation: &falsez,
					},
					Ports: []v1.ContainerPort{{Name: "static", ContainerPort: 8090, Protocol: v1.ProtocolTCP}},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							v1.ResourceCPU:    resource.MustParse("125m"),
							v1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
				}

				licenseEnvVar := v1.EnvVar{
					Name: "GLOO_LICENSE_KEY",
					ValueFrom: &v1.EnvVarSource{
						SecretKeyRef: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "license",
							},
							Key: "license-key",
						},
					},
				}

				grpcServerContainer := v1.Container{
					Name:            "apiserver",
					Image:           "quay.io/solo-io/gloo-fed-apiserver:" + version,
					ImagePullPolicy: getPullPolicy(),
					Ports: []v1.ContainerPort{
						{Name: "grpc", ContainerPort: 10101, Protocol: v1.ProtocolTCP},
						{Name: "healthcheck", HostPort: 0, ContainerPort: 8081, Protocol: v1.ProtocolTCP}},
					Env: []v1.EnvVar{
						GetPodNamespaceEnvVar(),
						podname,
						{
							Name:  "WRITE_NAMESPACE",
							Value: "gloo-system",
						},
						licenseEnvVar,
						statsEnvVar,
					},
					SecurityContext: &v1.SecurityContext{
						Capabilities: &v1.Capabilities{
							Drop: []v1.Capability{"ALL"},
						},
						RunAsNonRoot:             &truez,
						RunAsUser:                aws.Int64(101),
						ReadOnlyRootFilesystem:   &truez,
						AllowPrivilegeEscalation: &falsez,
					},
					VolumeMounts: []v1.VolumeMount{
						{Name: "empty-cache", MountPath: "/var/cache/nginx"},
						{Name: "empty-run", MountPath: "/var/run"},
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							v1.ResourceCPU:    resource.MustParse("125m"),
							v1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
				}

				envoyContainer := v1.Container{
					Name:            "envoy",
					Image:           "quay.io/solo-io/gloo-fed-apiserver-envoy:" + version,
					ImagePullPolicy: getPullPolicy(),
					VolumeMounts: []v1.VolumeMount{
						{Name: "envoy-config", MountPath: "/etc/envoy", ReadOnly: true},
					},
					Env: []v1.EnvVar{
						{
							Name:  "ENVOY_UID",
							Value: "0",
						},
					},
					SecurityContext: &v1.SecurityContext{
						Capabilities: &v1.Capabilities{
							Drop: []v1.Capability{"ALL"},
						},
						RunAsNonRoot:             &truez,
						RunAsUser:                aws.Int64(101),
						ReadOnlyRootFilesystem:   &truez,
						AllowPrivilegeEscalation: &falsez,
					},
					ReadinessProbe: &v1.Probe{
						Handler: v1.Handler{HTTPGet: &v1.HTTPGetAction{
							Path: "/",
							Port: intstr.IntOrString{IntVal: 8090},
						}},
						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							v1.ResourceCPU:    resource.MustParse("125m"),
							v1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
				}

				rb := ResourceBuilder{
					Namespace: namespace,
					Name:      "gloo-fed-console",
					Labels:    labels,
				}
				expectedDeployment = rb.GetDeploymentAppsv1()
				expectedDeployment.Spec.Selector.MatchLabels = selector
				expectedDeployment.Spec.Template.ObjectMeta.Labels = selector
				expectedDeployment.Spec.Template.Spec.Volumes = []v1.Volume{
					{Name: "empty-cache", VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}}},
					{Name: "empty-run", VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}}},
					{Name: "envoy-config", VolumeSource: v1.VolumeSource{ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: defaultBootstrapConfigMapName,
						},
					}}},
				}
				expectedDeployment.Spec.Template.Spec.Containers = []v1.Container{grpcServerContainer, uiContainer, envoyContainer}
				expectedDeployment.Spec.Template.Spec.ServiceAccountName = "gloo-fed-console"
				expectedDeployment.Spec.Template.ObjectMeta.Annotations = normalPromAnnotations
				expectedDeployment.Spec.Template.Spec.SecurityContext = nil
				expectedDeployment.Spec.Replicas = nil // GetDeploymentAppsv1 explicitly sets it to 1, which we don't want
			})

			It("is there by default", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
				Expect(err).NotTo(HaveOccurred())
				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("enables apiserver even if gloo-fed is disabled", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"gloo-fed.enabled=false",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("disables apiserver if explicitly disabled", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"gloo-fed.glooFedApiserver.enable=false",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				testManifest.Expect(expectedDeployment.Kind, expectedDeployment.Namespace, expectedDeployment.Name).To(BeNil())
			})

			It("can override image registry globally", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.image.registry=myregistry",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				expectedDeployment.Spec.Template.Spec.Containers[0].Image = "myregistry/gloo-fed-apiserver:" + version
				expectedDeployment.Spec.Template.Spec.Containers[1].Image = "myregistry/gloo-federation-console:" + version
				expectedDeployment.Spec.Template.Spec.Containers[2].Image = "myregistry/gloo-fed-apiserver-envoy:" + version
				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("does render the default bootstrap config map for the envoy sidecar", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
				Expect(err).NotTo(HaveOccurred())
				testManifest.Expect("ConfigMap", namespace, defaultBootstrapConfigMapName).NotTo(BeNil())
			})

			It("correctly sets resource limits", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"gloo-fed.glooFedApiserver.console.resources.limits.cpu=300m",
						"gloo-fed.glooFedApiserver.console.resources.limits.memory=300Mi",
						"gloo-fed.glooFedApiserver.console.resources.requests.cpu=30m",
						"gloo-fed.glooFedApiserver.console.resources.requests.memory=30Mi",
						"gloo-fed.glooFedApiserver.envoy.resources.limits.cpu=100m",
						"gloo-fed.glooFedApiserver.envoy.resources.limits.memory=100Mi",
						"gloo-fed.glooFedApiserver.envoy.resources.requests.cpu=10m",
						"gloo-fed.glooFedApiserver.envoy.resources.requests.memory=10Mi",
						"gloo-fed.glooFedApiserver.resources.limits.cpu=200m",
						"gloo-fed.glooFedApiserver.resources.limits.memory=200Mi",
						"gloo-fed.glooFedApiserver.resources.requests.cpu=20m",
						"gloo-fed.glooFedApiserver.resources.requests.memory=20Mi",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				// Apiserver
				expectedDeployment.Spec.Template.Spec.Containers[0].Resources = v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("200m"),
						v1.ResourceMemory: resource.MustParse("200Mi"),
					},
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("20m"),
						v1.ResourceMemory: resource.MustParse("20Mi"),
					},
				}

				// Console
				expectedDeployment.Spec.Template.Spec.Containers[1].Resources = v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("300m"),
						v1.ResourceMemory: resource.MustParse("300Mi"),
					},
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("30m"),
						v1.ResourceMemory: resource.MustParse("30Mi"),
					},
				}

				// Envoy
				expectedDeployment.Spec.Template.Spec.Containers[2].Resources = v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("100m"),
						v1.ResourceMemory: resource.MustParse("100Mi"),
					},
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("10m"),
						v1.ResourceMemory: resource.MustParse("10Mi"),
					},
				}

				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("allows setting custom runAsUser", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"gloo-fed.glooFedApiserver.runAsUser=10102"},
				})
				Expect(err).NotTo(HaveOccurred())

				uid := int64(10102)
				// Apiserver container
				expectedDeployment.Spec.Template.Spec.Containers[0].SecurityContext.RunAsUser = &uid
				expectedDeployment.Spec.Template.Spec.Containers[1].SecurityContext.RunAsUser = &uid
				expectedDeployment.Spec.Template.Spec.Containers[2].SecurityContext.RunAsUser = &uid
				testManifest.ExpectDeploymentAppsV1(expectedDeployment)

			})

			It("allows setting a custom number of replicas", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"gloo-fed.glooFedApiserver.replicas=2"},
				})
				Expect(err).NotTo(HaveOccurred())

				customNumReplicas := int32(2)
				expectedDeployment.Spec.Replicas = &customNumReplicas
				testManifest.ExpectDeploymentAppsV1(expectedDeployment)

			})

			It("correctly sets the GLOO_LICENSE_KEY env", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"gloo-fed.license_secret_name=custom-license-secret",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				licenseKeyEnvVarSource := v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "custom-license-secret",
						},
						Key: "license-key",
					},
				}

				for _, container := range expectedDeployment.Spec.Template.Spec.Containers {
					if container.Name == "apiserver" {
						for i, env := range container.Env {
							if env.Name == "GLOO_LICENSE_KEY" {
								container.Env[i].ValueFrom = &licenseKeyEnvVarSource
							}
						}
					}
				}

				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			When("a custom bootstrap config for the API server envoy sidecar is provided", func() {
				const customConfigMapName = "custom-bootstrap-config"
				var actualManifest TestManifest

				BeforeEach(func() {
					var err error
					actualManifest, err = BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{
							"gloo-fed.glooFedApiserver.envoy.bootstrapConfig.configMapName=" + customConfigMapName,
						},
					})
					Expect(err).NotTo(HaveOccurred())
				})

				It("adds the custom config map to the API server deployment volume mounts instead of the default one", func() {
					expectedDeployment.Spec.Template.Spec.Volumes = []v1.Volume{
						{Name: "empty-cache", VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}}},
						{Name: "empty-run", VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}}},
						{Name: "envoy-config", VolumeSource: v1.VolumeSource{ConfigMap: &v1.ConfigMapVolumeSource{
							LocalObjectReference: v1.LocalObjectReference{
								Name: customConfigMapName,
							},
						}}},
					}
					actualManifest.ExpectDeploymentAppsV1(expectedDeployment)
				})

				It("does not render the default config map", func() {
					actualManifest.Expect("ConfigMap", namespace, defaultBootstrapConfigMapName).To(BeNil())
				})

			})

			It("sits behind a service that is not exposed outside of the cluster", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
				Expect(err).NotTo(HaveOccurred())
				apiServerService := testManifest.SelectResources(func(u *unstructured.Unstructured) bool {
					if u.GetKind() != "Service" {
						return false
					}
					runtimeObj, err := kuberesource.ConvertUnstructured(u)
					Expect(err).NotTo(HaveOccurred())

					service, isService := runtimeObj.(*v1.Service)
					if isService && service.GetName() == "gloo-fed-console" {
						Expect(service.Spec.Type).To(Equal(v1.ServiceTypeClusterIP), "The gloo-fed-console service should be of type ClusterIP so it is not exposed outside the cluster")
						return true
					} else if !isService {
						Fail("Unexpected casting error")
						return false
					} else {
						return false
					}
				})

				Expect(apiServerService.NumResources()).To(Equal(1), "Should have found the gloo-fed-console service")
			})

			Context("pass image pull secrets", func() {
				pullSecretName := "test-pull-secret"
				pullSecret := []v1.LocalObjectReference{
					{Name: pullSecretName},
				}

				It("via global values", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{fmt.Sprintf("gloo-fed.glooFedApiserver.image.pullSecret=%s", pullSecretName)},
					})
					Expect(err).NotTo(HaveOccurred())

					expectedDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
					testManifest.ExpectDeploymentAppsV1(expectedDeployment)

				})

				It("via podTemplate values", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{
							fmt.Sprintf("gloo-fed.glooFedApiserver.image.pullSecret=%s", pullSecretName),
						},
					})
					Expect(err).NotTo(HaveOccurred())

					expectedDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
					testManifest.ExpectDeploymentAppsV1(expectedDeployment)
				})

				It("podTemplate values win over global", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{
							"global.image.pullSecret=wrong",
							fmt.Sprintf("gloo-fed.glooFedApiserver.image.pullSecret=%s", pullSecretName),
						},
					})
					Expect(err).NotTo(HaveOccurred())
					expectedDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
					testManifest.ExpectDeploymentAppsV1(expectedDeployment)
				})

			})
		})

		Context("ui console settings", func() {
			It("writes default console options to settings manifest", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{},
				})
				Expect(err).NotTo(HaveOccurred())
				settings := makeUnstructured(`
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: ` + namespace + `
spec:
  discovery:
    fdsMode: WHITELIST
  extauth:
    transportApiVersion: V3
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
    userIdHeader: "x-user-id"
  gateway:
    enableGatewayController: true
    readGatewaysFromAllNamespaces: false
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
      disableTransformationValidation: false
      allowWarnings: true
      warnRouteShortCircuiting: false
      validationServerGrpcMaxSizeBytes: 104857600
  gloo:
    enableRestEds: false
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    proxyDebugBindAddr: 0.0.0.0:9966
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    invalidConfigPolicy:
      replaceInvalidRoutes: false
      invalidRouteResponseBody: "Gloo Gateway has invalid configuration. Administrators should run ` + backtick + "glooctl check" + backtick + ` to find and fix config errors."
      invalidRouteResponseCode: 404
      replaceInvalidRoutes: false
  ratelimitServer:
    rateLimitBeforeAuth: false
    ratelimitServerRef:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
  consoleOptions:
    readOnly: false
    apiExplorerEnabled: true
  graphqlOptions:
    schemaChangeValidationOptions:
      rejectBreakingChanges: false
      processingRules: []
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})
			It("can override console options", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.console.readOnly=true",
						"global.console.apiExplorerEnabled=false",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				settings := makeUnstructured(`
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: ` + namespace + `
spec:
  discovery:
    fdsMode: WHITELIST
  extauth:
    transportApiVersion: V3
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
    userIdHeader: "x-user-id"
  gateway:
    enableGatewayController: true
    readGatewaysFromAllNamespaces: false
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
      disableTransformationValidation: false
      allowWarnings: true
      warnRouteShortCircuiting: false
      validationServerGrpcMaxSizeBytes: 104857600
  gloo:
    enableRestEds: false
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    proxyDebugBindAddr: 0.0.0.0:9966
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    invalidConfigPolicy:
      replaceInvalidRoutes: false
      invalidRouteResponseBody: "Gloo Gateway has invalid configuration. Administrators should run ` + backtick + "glooctl check" + backtick + ` to find and fix config errors."
      invalidRouteResponseCode: 404
      replaceInvalidRoutes: false
  ratelimitServer:
    rateLimitBeforeAuth: false
    ratelimitServerRef:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
  consoleOptions:
    readOnly: true
    apiExplorerEnabled: false
  graphqlOptions:
    schemaChangeValidationOptions:
      rejectBreakingChanges: false
      processingRules: []
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})
		})

		Context("graphql settings", func() {
			It("writes default graphql options to settings manifest", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{},
				})
				Expect(err).NotTo(HaveOccurred())
				settings := makeUnstructured(`
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: ` + namespace + `
spec:
  discovery:
    fdsMode: WHITELIST
  extauth:
    transportApiVersion: V3
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
    userIdHeader: "x-user-id"
  gateway:
    enableGatewayController: true
    readGatewaysFromAllNamespaces: false
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
      disableTransformationValidation: false
      allowWarnings: true
      warnRouteShortCircuiting: false
      validationServerGrpcMaxSizeBytes: 104857600
  gloo:
    enableRestEds: false
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    proxyDebugBindAddr: 0.0.0.0:9966
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    invalidConfigPolicy:
      replaceInvalidRoutes: false
      invalidRouteResponseBody: "Gloo Gateway has invalid configuration. Administrators should run ` + backtick + "glooctl check" + backtick + ` to find and fix config errors."
      invalidRouteResponseCode: 404
      replaceInvalidRoutes: false
  ratelimitServer:
    rateLimitBeforeAuth: false
    ratelimitServerRef:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
  consoleOptions:
    readOnly: false
    apiExplorerEnabled: true
  graphqlOptions:
    schemaChangeValidationOptions:
      rejectBreakingChanges: false
      processingRules: []
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})
			It("can override graphql options", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.graphql.changeValidation.rejectBreaking=true",
						"global.graphql.changeValidation.rules.ignoreUnreachable=true",
						"global.graphql.changeValidation.rules.deprecatedFieldRemovalDangerous=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				settings := makeUnstructured(`
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: ` + namespace + `
spec:
  discovery:
    fdsMode: WHITELIST
  extauth:
    transportApiVersion: V3
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
    userIdHeader: "x-user-id"
  gateway:
    enableGatewayController: true
    readGatewaysFromAllNamespaces: false
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
      disableTransformationValidation: false
      allowWarnings: true
      warnRouteShortCircuiting: false
      validationServerGrpcMaxSizeBytes: 104857600
  gloo:
    enableRestEds: false
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    proxyDebugBindAddr: 0.0.0.0:9966
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    invalidConfigPolicy:
      replaceInvalidRoutes: false
      invalidRouteResponseBody: "Gloo Gateway has invalid configuration. Administrators should run ` + backtick + "glooctl check" + backtick + ` to find and fix config errors."
      invalidRouteResponseCode: 404
      replaceInvalidRoutes: false
  ratelimitServer:
    rateLimitBeforeAuth: false
    ratelimitServerRef:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
  consoleOptions:
    readOnly: false
    apiExplorerEnabled: true
  graphqlOptions:
    schemaChangeValidationOptions:
      rejectBreakingChanges: true
      processingRules:
      - RULE_DEPRECATED_FIELD_REMOVAL_DANGEROUS
      - RULE_IGNORE_UNREACHABLE
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})
			It("does not add rules that are set to false", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.graphql.changeValidation.rules.ignoreUnreachable=false",
						"global.graphql.changeValidation.rules.dangerousToBreaking=true",
						"global.graphql.changeValidation.rules.ignoreDescriptionChanges=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				settings := makeUnstructured(`
apiVersion: gloo.solo.io/v1
kind: Settings
metadata:
  labels:
    app: gloo
    gloo: settings
  name: default
  namespace: ` + namespace + `
spec:
  discovery:
    fdsMode: WHITELIST
  extauth:
    transportApiVersion: V3
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
    userIdHeader: "x-user-id"
  gateway:
    enableGatewayController: true
    readGatewaysFromAllNamespaces: false
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
      disableTransformationValidation: false
      allowWarnings: true
      warnRouteShortCircuiting: false
      validationServerGrpcMaxSizeBytes: 104857600
  gloo:
    enableRestEds: false
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    proxyDebugBindAddr: 0.0.0.0:9966
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    invalidConfigPolicy:
      replaceInvalidRoutes: false
      invalidRouteResponseBody: "Gloo Gateway has invalid configuration. Administrators should run ` + backtick + "glooctl check" + backtick + ` to find and fix config errors."
      invalidRouteResponseCode: 404
      replaceInvalidRoutes: false
  ratelimitServer:
    rateLimitBeforeAuth: false
    ratelimitServerRef:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
  consoleOptions:
    readOnly: false
    apiExplorerEnabled: true
  graphqlOptions:
    schemaChangeValidationOptions:
      rejectBreakingChanges: false
      processingRules:
      - RULE_DANGEROUS_TO_BREAKING
      - RULE_IGNORE_DESCRIPTION_CHANGES
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})
		})

		Describe("Standard k8s values", func() {
			DescribeTable("PodSpec affinity, tolerations, nodeName, hostAliases, nodeSelector, restartPolicy on Deployments and Jobs",
				func(kind string, resourceName string, value string, extraArgs ...string) {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: append([]string{
							value + ".nodeSelector.label=someLabel",
							value + ".nodeName=someNodeName",
							value + ".tolerations[0].operator=someToleration",
							value + ".hostAliases[0]=someHostAlias",
							value + ".affinity.nodeAffinity=someNodeAffinity",
							value + ".restartPolicy=someRestartPolicy",
						}, extraArgs...),
					})
					Expect(err).NotTo(HaveOccurred())
					resources := testManifest.SelectResources(func(u *unstructured.Unstructured) bool {
						if u.GetKind() == kind && u.GetName() == resourceName {
							a := getFieldFromUnstructured(u, "spec", "template", "spec", "nodeSelector")
							Expect(a).To(Equal(map[string]interface{}{"label": "someLabel"}))
							a = getFieldFromUnstructured(u, "spec", "template", "spec", "nodeName")
							Expect(a).To(Equal("someNodeName"))
							a = getFieldFromUnstructured(u, "spec", "template", "spec", "tolerations")
							Expect(a).To(Equal([]interface{}{map[string]interface{}{"operator": "someToleration"}}))
							a = getFieldFromUnstructured(u, "spec", "template", "spec", "hostAliases")
							Expect(a).To(Equal([]interface{}{"someHostAlias"}))
							a = getFieldFromUnstructured(u, "spec", "template", "spec", "affinity")
							Expect(a).To(Equal(map[string]interface{}{"nodeAffinity": "someNodeAffinity"}))
							a = getFieldFromUnstructured(u, "spec", "template", "spec", "restartPolicy")
							Expect(a).To(Equal("someRestartPolicy"))
							return true
						}
						return false
					})
					Expect(resources.NumResources()).To(Equal(1))
				},
				Entry("redis deployment", "Deployment", "redis", "redis.deployment"),
				Entry("rate limit deployment", "Deployment", "rate-limit", "global.extensions.rateLimit.deployment"),
				Entry("observability deployment", "Deployment", "observability", "observability.deployment"),
				Entry("extauth deployment", "Deployment", "extauth", "global.extensions.extAuth.deployment"),
			)
		})

		Context("Kube resource overrides", func() {
			DescribeTable("overrides YAML in generated sources", func(overrideProperty string, extraArgs ...string) {
				// Override property should be the path to `kubeResourceOverride`, like gloo.deployment.kubeResourceOverride
				valueArg := fmt.Sprintf("%s.metadata.labels.overriddenLabel=label", overrideProperty)
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: append(extraArgs, valueArg),
				})
				Expect(err).NotTo(HaveOccurred())
				// We are overriding the generated yaml by adding our own label to the metadata
				resources := testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetLabels()["overriddenLabel"] == "label" && resource.GetKind() != ""
				})
				// some resources are contained directly in the manifest, and custom resources (upstreams) are
				// applied by the rollout job, so we need to get the total from both places
				countFromResources := resources.NumResources()
				job := getJob(testManifest, namespace, "gloo-ee-resource-rollout")
				countFromJob := strings.Count(job.Spec.Template.Spec.Containers[0].Command[2], "overriddenLabel: label")
				Expect(countFromResources + countFromJob).To(Equal(1))
			},
				Entry("0-redis-service", "redis.service.kubeResourceOverride"),
				Entry("1-redis-deployment", "redis.deployment.kubeResourceOverride"),
				Entry("2-rate-limit-deployment", "global.extensions.rateLimit.deployment.kubeResourceOverride"),
				Entry("3-rate-limit-service", "global.extensions.rateLimit.service.kubeResourceOverride"),
				Entry("4-ratelimit-upstream", "global.extensions.rateLimit.upstream.kubeResourceOverride"),
				Entry("8-observability-service-account", "observability.serviceAccount.kubeResourceOverride", "observability.enabled=true"),
				Entry("9-observability-configmap", "observability.configMap.kubeResourceOverride"),
				Entry("9-observability-deployment", "observability.deployment.kubeResourceOverride"),
				Entry("9-observability-secret", "observability.secret.kubeResourceOverride"),
				Entry("20-extauth-secret", "global.extensions.extAuth.secret.kubeResourceOverride"),
				Entry("21-extauth-deployment", "global.extensions.extAuth.deployment.kubeResourceOverride"),
				Entry("22-extauth-service", "global.extensions.extAuth.service.kubeResourceOverride"),
			)
		})

		Context("custom resource lifecycles", func() {

			It("creates migration, rollout, and cleanup jobs", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
				Expect(err).NotTo(HaveOccurred())

				// getJob will fail if the job doesn't exist
				_ = getJob(testManifest, namespace, "gloo-resource-migration")
				_ = getJob(testManifest, namespace, "gloo-resource-rollout")
				_ = getJob(testManifest, namespace, "gloo-ee-resource-rollout")
				_ = getJob(testManifest, namespace, "gloo-resource-cleanup")

			})

			It("applies extauth and ratelimit upstreams", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.envoySidecar=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				job := getJob(testManifest, namespace, "gloo-ee-resource-rollout")
				Expect(job.Spec.Template.Spec.Containers[0].Command[2]).To(ContainSubstring(`apiVersion: gloo.solo.io/v1
kind: Upstream
metadata:
  name: extauth
  namespace: ` + namespace))
				Expect(job.Spec.Template.Spec.Containers[0].Command[2]).To(ContainSubstring(`apiVersion: gloo.solo.io/v1
kind: Upstream
metadata:
  name: extauth-sidecar
  namespace: ` + namespace))
				Expect(job.Spec.Template.Spec.Containers[0].Command[2]).To(ContainSubstring(`apiVersion: gloo.solo.io/v1
kind: Upstream
metadata:
  name: rate-limit
  namespace: ` + namespace))
			})

			It("does not call kubectl apply when extauth and ratelimit upstreams are disabled", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.enabled=false",
						"global.extensions.rateLimit.enabled=false",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				job := getJob(testManifest, namespace, "gloo-ee-resource-rollout")
				Expect(job.Spec.Template.Spec.Containers[0].Command[2]).NotTo(ContainSubstring("kubectl apply"))
				Expect(job.Spec.Template.Spec.Containers[0].Command[2]).To(ContainSubstring("no custom resources to apply"))
			})
		})

		// Lines ending with whitespace causes malformatted config map (https://github.com/solo-io/gloo/issues/4645)
		It("Should not contain trailing whitespace", func() {
			out, err := exec.Command("helm", "template", "../helm/gloo-ee").CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), string(out))

			lines := strings.Split(string(out), "\n")
			// more descriptive fail message that prints out the manifest that includes the trailing whitespace
			manifestStartingLine := 0
			skip := false
			for idx, line := range lines {
				if strings.Contains(line, "---") {
					manifestStartingLine = idx
					continue
				}
				// skip all the content within kubectl apply commands (used in the rollout job)
				// since there is extra whitespace that can't be removed
				if strings.Contains(line, "kubectl apply -f - <<EOF") {
					skip = true
					continue
				}
				if strings.TrimSpace(line) == "EOF" {
					skip = false
					continue
				}
				if !skip && strings.TrimRightFunc(line, unicode.IsSpace) != line {
					// Ensure that we are only checking this for Gloo charts, and not our subcharts
					manifest := strings.Join(lines[manifestStartingLine:idx+1], "\n")
					if strings.Contains(manifest, "# Source: gloo-ee/templates") {
						Fail(manifest + "\n ^^^ the above line has whitespace")
					}

				}
			}
		})
	})

})

func constructResourceID(resource *unstructured.Unstructured) string {
	// technically vulnerable to resources that have commas in their names, but that's not a big concern
	return fmt.Sprintf("%s,%s,%s", resource.GetNamespace(), resource.GetName(), resource.GroupVersionKind().String())
}

// gets value of field nested within an Unstructured struct.
// fieldPath is the path to the value, so the value foo.bar.baz would be passed in as "foo", "bar, "baz"
func getFieldFromUnstructured(uns *unstructured.Unstructured, fieldPath ...string) interface{} {
	if len(fieldPath) < 1 {
		return nil
	}
	obj := uns.Object[fieldPath[0]]
	for _, field := range fieldPath[1:] {
		obj = obj.(map[string]interface{})[field]
	}
	return obj
}

func getJob(testManifest TestManifest, jobNamespace string, jobName string) *jobsv1.Job {
	jobUns := testManifest.ExpectCustomResource("Job", jobNamespace, jobName)
	jobObj, err := kuberesource.ConvertUnstructured(jobUns)
	Expect(err).NotTo(HaveOccurred())
	Expect(jobObj).To(BeAssignableToTypeOf(&jobsv1.Job{}))
	return jobObj.(*jobsv1.Job)
}
