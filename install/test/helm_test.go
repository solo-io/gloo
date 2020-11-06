package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/setup"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/aws/aws-sdk-go/aws"
	bootstrapv3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/config/health_checker/redis/v2"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/redis_proxy/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/go-utils/installutils/kuberesource"
	"github.com/solo-io/solo-projects/install/helm/gloo-ee/generate"
	"github.com/solo-io/solo-projects/pkg/install"
	jobsv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8s "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/go-utils/manifesttestutils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

var _ = Describe("Helm Test", func() {
	var (
		version string

		normalPromAnnotations = map[string]string{
			"prometheus.io/path":   "/metrics",
			"prometheus.io/port":   "9091",
			"prometheus.io/scrape": "true",
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
			version = os.Getenv("TAGGED_VERSION")
			if version == "" {
				version = "dev"
				getPullPolicy = func() v1.PullPolicy { return v1.PullAlways }
			} else {
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
										Key:  "DASHBOARD_JSON_TEMPLATE",
										Path: "dashboard-template.json",
									},
								},
							},
						},
					},
				}
				observabilityDeployment.Spec.Template.Spec.Containers = []v1.Container{
					{
						Name:  "observability",
						Image: "quay.io/solo-io/observability-ee:dev",
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
						ImagePullPolicy: "Always",
					},
				}
				observabilityDeployment.Spec.Template.Spec.ServiceAccountName = "observability"
				observabilityDeployment.Spec.Strategy = appsv1.DeploymentStrategy{}
				observabilityDeployment.Spec.Selector.MatchLabels = selector
				observabilityDeployment.Spec.Template.ObjectMeta.Labels = selector
				observabilityDeployment.Spec.Template.ObjectMeta.Annotations = normalPromAnnotations

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

				It("correctly sets the GLOO_LICENSE_KEY env", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{"license_secret_name=custom-license-secret"},
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
						Image:           "quay.io/solo-io/extauth-ee:dev",
						ImagePullPolicy: "Always",
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
								Name:  "SERVER_PORT",
								Value: "8083",
							},
							{
								Name:  "USER_ID_HEADER",
								Value: "x-user-id",
							},
							statsEnvVar,
						},
						ReadinessProbe: &v1.Probe{
							Handler: v1.Handler{
								Exec: &v1.ExecAction{
									Command: []string{"/bin/sh", "-c", "nc -z localhost 8083"},
								},
							},
							InitialDelaySeconds: 1,
							FailureThreshold:    3,
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
						"rateLimit.enabled=true",
						"apiServer.enable=true",
						"global.extensions.extAuth.enabled=true",
						"observability.enabled=true",
						"rateLimit.deployment.extraRateLimitLabels.foo=bar",
						"redis.deployment.extraRedisLabels.foo=bar",
						"apiServer.deployment.extraApiServerLabels.foo=bar",
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
				Expect(resourcesTested).To(Equal(5), "Tested %d resources when we were expecting 5."+
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

			Context("dataplane per proxy", func() {

				helmOverrideFileContents := func(dataplanePerProxy bool) string {
					return fmt.Sprintf(`
global:
  extensions:
    dataplanePerProxy: %t
  glooStats:
    enabled: true
gloo:
  global:
    glooStats:
      enabled: true
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
      kind:
        deployment:
          antiAffinity: false
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
        annotations:
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

					assertExpectedResourcesForProxy := func(proxyName string) {
						gatewayProxyRateLimitResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
							return unstructured.GetLabels()["gloo"] == fmt.Sprintf("rate-limit-%s", proxyName)
						})

						// Deployment, Service, Upstream
						Expect(gatewayProxyRateLimitResources.NumResources()).To(Equal(3), fmt.Sprintf("%s: Expecting RateLimit Deployment, Service, and Upstream", proxyName))

						gatewayProxyRedisResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
							return unstructured.GetLabels()["gloo"] == fmt.Sprintf("redis-%s", proxyName)
						})

						// Deployment, Service
						Expect(gatewayProxyRedisResources.NumResources()).To(Equal(2), fmt.Sprintf("%s: Expecting Redis Deployment and Service", proxyName))

						gatewayProxyExtAuthResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
							return unstructured.GetLabels()["gloo"] == fmt.Sprintf("extauth-%s", proxyName)
						})

						// Deployment, Service, Upstream
						Expect(gatewayProxyExtAuthResources.NumResources()).To(Equal(3), fmt.Sprintf("%s: Expecting Extauth Deployment, Service, and Upstream", proxyName))

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

					rateLimitResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						return unstructured.GetLabels()["gloo"] == "rate-limit"
					})
					Expect(rateLimitResources.NumResources()).To(Equal(4), "Expecting RateLimit Deployment, Service, Upstream and ServiceAccount")

					redisResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						return unstructured.GetLabels()["gloo"] == "redis"
					})
					Expect(redisResources.NumResources()).To(Equal(2), "Expecting Redis Deployment and Service")

					extAuthResources := testManifest.SelectResources(func(unstructured *unstructured.Unstructured) bool {
						return unstructured.GetLabels()["gloo"] == "extauth"
					})
					Expect(extAuthResources.NumResources()).To(Equal(5), "Expecting ExtAuth Deployment, Service, Upstream, ServiceAccount, and Secret")
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
						ImagePullPolicy: v1.PullIfNotPresent,
						VolumeMounts:    authPluginVolumeMount,
					},
					{
						Name:            "plugin-second-plugin",
						Image:           "bar/foo:1.2.3",
						ImagePullPolicy: v1.PullIfNotPresent,
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

			It("creates settings with extauth request timeout", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.extensions.extAuth.requestTimeout=1s",
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
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
    requestTimeout: "1s"
  gateway:
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
  gloo:
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
  ratelimitServer:
    ratelimit_server_ref: 
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})

			It("enable default credentials", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"settings.aws.enableCredentialsDiscovery=true",
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
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
  gateway:
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
  gloo:
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    awsOptions:
      enableCredentialsDiscovey: true
  ratelimitServer:
    ratelimit_server_ref:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})

			It("Allows ratelimit descriptors to be set", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"settings.rateLimit.descriptors[0].key=generic_key",
						"settings.rateLimit.descriptors[0].value=per-second",
						"settings.rateLimit.descriptors[0].rateLimit.requestsPerUnit=2",
						"settings.rateLimit.descriptors[0].rateLimit.unit=SECOND",
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
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
  gateway:
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
  gloo:
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
  ratelimitServer:
    ratelimit_server_ref:
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
  rateLimit:
    descriptors:
      - key: generic_key
        value: "per-second"
        rateLimit:
          requestsPerUnit: 2
          unit: SECOND
`)
				testManifest.ExpectUnstructured(settings.GetKind(), settings.GetNamespace(), settings.GetName()).To(BeEquivalentTo(settings))
			})

			It("enable sts discovery", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"settings.aws.enableServiceAccountCredentials=true",
						"settings.aws.stsCredentialsRegion=us-east-2",
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
    extauthzServerRef:
      name: extauth
      namespace: ` + namespace + `
  gateway:
    validation:
      alwaysAccept: true
      proxyValidationServerAddr: gloo:9988
  gloo:
    xdsBindAddr: 0.0.0.0:9977
    restXdsBindAddr: 0.0.0.0:9976
    disableKubernetesDestinations: false
    disableProxyGarbageCollection: false
    awsOptions:
      serviceAccountCredentials:
        cluster: aws_sts_cluster
        uri: sts.us-east-2.amazonaws.com
  ratelimitServer:
    ratelimit_server_ref: 
      namespace: ` + namespace + `
      name: rate-limit
  kubernetesArtifactSource: {}
  kubernetesConfigSource: {}
  kubernetesSecretSource: {}
  refreshRate: 60s
  discoveryNamespace: ` + namespace + `
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
						Image:           "quay.io/solo-io/extauth-ee:dev",
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
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "shared-data",
								MountPath: "/usr/share/shared-data",
							},
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

		Context("apiserver deployment", func() {
			const defaultBootstrapConfigMapName = "default-apiserver-envoy-config"

			var expectedDeployment *appsv1.Deployment

			BeforeEach(func() {
				labels = map[string]string{
					"gloo": "apiserver-ui",
					"app":  "gloo",
				}
				selector = map[string]string{
					"app":  "gloo",
					"gloo": "apiserver-ui",
				}
				grpcPortEnvVar := v1.EnvVar{
					Name:  "GRPC_PORT",
					Value: "10101",
				}
				rbacNamespacedEnvVar := v1.EnvVar{
					Name:  setup.NamespacedRbacEnvName,
					Value: "false",
				}
				noAuthEnvVar := v1.EnvVar{
					Name:  "NO_AUTH",
					Value: "1",
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
				uiContainer := v1.Container{
					Name:            "apiserver-ui",
					Image:           "quay.io/solo-io/grpcserver-ui:" + version,
					ImagePullPolicy: v1.PullAlways,
					VolumeMounts: []v1.VolumeMount{
						{Name: "empty-cache", MountPath: "/var/cache/nginx"},
						{Name: "empty-run", MountPath: "/var/run"},
					},
					Ports: []v1.ContainerPort{{Name: "static", ContainerPort: 8080, Protocol: v1.ProtocolTCP}},
				}
				grpcServerContainer := v1.Container{
					Name:            "apiserver",
					Image:           "quay.io/solo-io/grpcserver-ee:" + version,
					ImagePullPolicy: v1.PullAlways,
					Ports:           []v1.ContainerPort{{Name: "grpcport", ContainerPort: 10101, Protocol: v1.ProtocolTCP}},
					Env: []v1.EnvVar{
						GetPodNamespaceEnvVar(),
						rbacNamespacedEnvVar,
						grpcPortEnvVar,
						statsEnvVar,
						noAuthEnvVar,
						licenseEnvVar,
					},
				}
				envoyContainer := v1.Container{
					Name:            "gloo-grpcserver-envoy",
					Image:           "quay.io/solo-io/grpcserver-envoy:" + version,
					ImagePullPolicy: v1.PullAlways,
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
						RunAsUser: aws.Int64(101),
					},
					ReadinessProbe: &v1.Probe{
						Handler: v1.Handler{HTTPGet: &v1.HTTPGetAction{
							Path: "/",
							Port: intstr.IntOrString{IntVal: 8080},
						}},
						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
					},
				}

				nonRootUser := int64(10101)
				nonRoot := true

				nonRootSC := &v1.PodSecurityContext{
					RunAsUser:    &nonRootUser,
					RunAsNonRoot: &nonRoot,
				}

				rb := ResourceBuilder{
					Namespace: namespace,
					Name:      "api-server",
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
				expectedDeployment.Spec.Template.Spec.Containers = []v1.Container{uiContainer, grpcServerContainer, envoyContainer}
				expectedDeployment.Spec.Template.Spec.ServiceAccountName = "apiserver-ui"
				expectedDeployment.Spec.Template.ObjectMeta.Annotations = normalPromAnnotations
				expectedDeployment.Spec.Template.Spec.SecurityContext = nonRootSC
				expectedDeployment.Spec.Replicas = nil // GetDeploymentAppsv1 explicitly sets it to 1, which we don't want
			})

			It("is there by default", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
				Expect(err).NotTo(HaveOccurred())
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
						"apiServer.deployment.ui.resources.limits.cpu=300m",
						"apiServer.deployment.ui.resources.limits.memory=300Mi",
						"apiServer.deployment.ui.resources.requests.cpu=30m",
						"apiServer.deployment.ui.resources.requests.memory=30Mi",
						"apiServer.deployment.envoy.resources.limits.cpu=100m",
						"apiServer.deployment.envoy.resources.limits.memory=100Mi",
						"apiServer.deployment.envoy.resources.requests.cpu=10m",
						"apiServer.deployment.envoy.resources.requests.memory=10Mi",
						"apiServer.deployment.server.resources.limits.cpu=200m",
						"apiServer.deployment.server.resources.limits.memory=200Mi",
						"apiServer.deployment.server.resources.requests.cpu=20m",
						"apiServer.deployment.server.resources.requests.memory=20Mi",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				// UI
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

				// Server
				expectedDeployment.Spec.Template.Spec.Containers[1].Resources = v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("200m"),
						v1.ResourceMemory: resource.MustParse("200Mi"),
					},
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("20m"),
						v1.ResourceMemory: resource.MustParse("20Mi"),
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
					valuesArgs: []string{"apiServer.deployment.runAsUser=10102"},
				})
				Expect(err).NotTo(HaveOccurred())

				uid := int64(10102)
				expectedDeployment.Spec.Template.Spec.SecurityContext.RunAsUser = &uid

				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("allows setting a custom number of replicas", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"apiServer.deployment.replicas=2"},
				})
				Expect(err).NotTo(HaveOccurred())

				customNumReplicas := int32(2)
				expectedDeployment.Spec.Replicas = &customNumReplicas

				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("correctly sets the GLOO_LICENSE_KEY env", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"apiServer.enterprise=true",
						"license_secret_name=custom-license-secret",
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
				envs := expectedDeployment.Spec.Template.Spec.Containers[1].Env
				for i, env := range envs {
					if env.Name == "GLOO_LICENSE_KEY" {
						envs[i].ValueFrom = &licenseKeyEnvVarSource
					}
				}
				testManifest.ExpectDeploymentAppsV1(expectedDeployment)
			})

			It("correctly sets the RBAC_NAMESPACED env", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"global.glooRbac.namespaced=true",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				envs := expectedDeployment.Spec.Template.Spec.Containers[1].Env
				for i, env := range envs {
					if env.Name == setup.NamespacedRbacEnvName {
						envs[i].Value = "true"
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
							"apiServer.deployment.envoy.bootstrapConfig.configMapName=" + customConfigMapName,
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

			It("can be set as NodePort service", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{
						"apiServer.service.serviceType=NodePort",
					},
				})
				Expect(err).NotTo(HaveOccurred())
				apiServerService := testManifest.SelectResources(func(u *unstructured.Unstructured) bool {
					if u.GetKind() != "Service" {
						return false
					}
					runtimeObj, err := kuberesource.ConvertUnstructured(u)
					Expect(err).NotTo(HaveOccurred())

					service, isService := runtimeObj.(*v1.Service)
					if isService && service.GetName() == "apiserver-ui" {
						Expect(service.Spec.Type).To(Equal(v1.ServiceTypeNodePort), "The apiserver-ui service should be of type NodePort so it is not exposed outside the cluster")
						return true
					} else if !isService {
						Fail("Unexpected casting error")
						return false
					} else {
						return false
					}
				})

				Expect(apiServerService.NumResources()).To(Equal(1), "Should have found the apiserver-ui service")
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
					if isService && service.GetName() == "apiserver-ui" {
						Expect(service.Spec.Type).To(Equal(v1.ServiceTypeClusterIP), "The apiserver-ui service should be of type ClusterIP so it is not exposed outside the cluster")
						return true
					} else if !isService {
						Fail("Unexpected casting error")
						return false
					} else {
						return false
					}
				})

				Expect(apiServerService.NumResources()).To(Equal(1), "Should have found the apiserver-ui service")
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
					Expect(err).NotTo(HaveOccurred())

					expectedDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
					testManifest.ExpectDeploymentAppsV1(expectedDeployment)

				})

				It("via podTemplate values", func() {
					testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
						valuesArgs: []string{
							fmt.Sprintf("apiServer.deployment.server.image.pullSecret=%s", pullSecretName),
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
							fmt.Sprintf("apiServer.deployment.server.image.pullSecret=%s", pullSecretName),
						},
					})
					Expect(err).NotTo(HaveOccurred())
					expectedDeployment.Spec.Template.Spec.ImagePullSecrets = pullSecret
					testManifest.ExpectDeploymentAppsV1(expectedDeployment)
				})

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
						Expect(bootstrap.Node).To(Equal(&corev3.Node{Id: "sds_client", Cluster: "sds_client"}))
						Expect(bootstrap.StaticResources.Listeners[0].FilterChains[0].TransportSocket).NotTo(BeNil())
						tlsContext := tlsv3.DownstreamTlsContext{}
						Expect(ptypes.UnmarshalAny(bootstrap.StaticResources.Listeners[0].FilterChains[0].TransportSocket.GetTypedConfig(), &tlsContext)).NotTo(HaveOccurred())
						Expect(tlsContext).To(Equal(tlsv3.DownstreamTlsContext{
							CommonTlsContext: &tlsv3.CommonTlsContext{
								TlsCertificateSdsSecretConfigs: []*tlsv3.SdsSecretConfig{
									{
										Name: "server_cert",
										SdsConfig: &corev3.ConfigSource{
											ConfigSourceSpecifier: &corev3.ConfigSource_ApiConfigSource{
												ApiConfigSource: &corev3.ApiConfigSource{
													ApiType: corev3.ApiConfigSource_GRPC,
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
					valuesArgs: []string{"global.extensions.extAuth.tlsEnabled=true,global.extensions.extAuth.secretName=my-secret"},
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

			It("should allow apiserver service to handle TLS itself using a kubernetes secret", func() {
				testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{
					valuesArgs: []string{"apiServer.sslSecretName=ssl-secret"},
				})
				Expect(err).NotTo(HaveOccurred())

				testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return resource.GetKind() == "Deployment"
				}).ExpectAll(func(deployment *unstructured.Unstructured) {
					deploymentObject, err := kuberesource.ConvertUnstructured(deployment)
					Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Deployment %+v should be able to convert from unstructured", deployment))
					structuredDeployment, ok := deploymentObject.(*appsv1.Deployment)
					Expect(ok).To(BeTrue(), fmt.Sprintf("Deployment %+v should be able to cast to a structured deployment", deployment))

					if structuredDeployment.GetName() == "api-server" {
						mode := int32(420)
						Expect(structuredDeployment.Spec.Template.Spec.Volumes[3].Name).To(Equal("apiserver-ssl-certs"))
						Expect(structuredDeployment.Spec.Template.Spec.Volumes[3].VolumeSource.Secret).To(Equal(&corev1.SecretVolumeSource{
							SecretName:  "ssl-secret",
							DefaultMode: &mode,
						}))
						Expect(structuredDeployment.Spec.Template.Spec.Containers[2].VolumeMounts[1]).To(Equal(corev1.VolumeMount{
							MountPath: "/etc/apiserver/ssl",
							ReadOnly:  true,
							Name:      "apiserver-ssl-certs",
						}))
						Expect(structuredDeployment.Spec.Template.Spec.Containers[2].ReadinessProbe.HTTPGet.Scheme).To(Equal(corev1.URISchemeHTTPS))
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

						Expect(bootstrap.StaticResources.Listeners[0].FilterChains[0].TransportSocket).NotTo(BeNil())
						tlsContext := tlsv3.DownstreamTlsContext{}
						Expect(ptypes.UnmarshalAny(bootstrap.StaticResources.Listeners[0].FilterChains[0].TransportSocket.GetTypedConfig(), &tlsContext)).NotTo(HaveOccurred())
						Expect(tlsContext).To(Equal(tlsv3.DownstreamTlsContext{
							CommonTlsContext: &tlsv3.CommonTlsContext{
								TlsCertificates: []*tlsv3.TlsCertificate{
									{
										CertificateChain: &corev3.DataSource{
											Specifier: &corev3.DataSource_Filename{
												Filename: "/etc/apiserver/ssl/tls.crt",
											},
										},
										PrivateKey: &corev3.DataSource{
											Specifier: &corev3.DataSource_Filename{
												Filename: "/etc/apiserver/ssl/tls.key",
											},
										},
									},
								},
							},
						}))
					}
				})

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
						"rateLimit.enabled=true",
						"ratelimitServer.ratelimitServerRef.name=rate-limit",
						"ratelimitServer.ratelimitServerRef.namespace=gloo-system",
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

	})

	Describe("gloo with read-only ui helm tests", func() {
		var (
			labels           map[string]string
			selector         map[string]string
			manifestYaml     string
			glooOsVersion    string
			glooOsPullPolicy v1.PullPolicy
		)

		BeforeEach(func() {

			var err error
			var glooEGenerationFiles = &generate.GenerationFiles{
				Artifact:             generate.GlooE,
				RequirementsTemplate: "../../install/helm/gloo-ee/requirements-template.yaml",
			}
			var glooOsWithReadOnlyUiGenerationFiles = &generate.GenerationFiles{
				Artifact:             generate.GlooWithRoUi,
				RequirementsTemplate: "../../install/helm/gloo-os-with-ui/requirements-template.yaml",
			}
			glooOsVersion, err = generate.GetGlooOsVersion(glooEGenerationFiles, glooOsWithReadOnlyUiGenerationFiles)
			Expect(err).NotTo(HaveOccurred())
			glooOsPullPolicy = v1.PullAlways

			version = os.Getenv("TAGGED_VERSION")
			if version == "" {
				version = "dev"
			} else {
				version = version[1:]
			}
			manifestYaml = ""
		})

		AfterEach(func() {
			if manifestYaml != "" {
				err := os.Remove(manifestYaml)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		Context("gateway", func() {
			BeforeEach(func() {
				labels = map[string]string{
					"gloo":             "gateway-proxy",
					"gateway-proxy-id": defaults.GatewayProxyName,
					"app":              "gloo",
				}
				selector = map[string]string{
					"gateway-proxy": "live",
				}
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

					container := GetQuayContainerSpec("gloo-envoy-wrapper", glooOsVersion, GetPodNamespaceEnvVar(), podname)
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
					deploy.Spec.Template.Spec.Containers[0].ImagePullPolicy = glooOsPullPolicy
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

				It("creates a deployment without envoy config annotations that contains the Settings resource", func() {
					testManifest, err := BuildTestManifest(install.GlooOsWithUiChartName, namespace, helmValues{})
					Expect(err).NotTo(HaveOccurred())
					testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
					testManifest.ExpectCustomResource("Settings", namespace, "default")
				})

				It("creates a deployment with envoy config annotations that contains the Settings resource", func() {
					testManifest, err := BuildTestManifest(install.GlooOsWithUiChartName, namespace, helmValues{
						valuesArgs: []string{"gloo.gatewayProxies.gatewayProxy.readConfig=true"},
					})
					Expect(err).NotTo(HaveOccurred())
					includeStatConfig()
					testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
					testManifest.ExpectCustomResource("Settings", namespace, "default")
				})

				Context("apiserver deployment", func() {

					const defaultBootstrapConfigMapName = "default-apiserver-envoy-config"

					var deploy *appsv1.Deployment

					BeforeEach(func() {
						labels = map[string]string{
							"gloo": "apiserver-ui",
							"app":  "gloo",
						}
						selector = map[string]string{
							"app":  "gloo",
							"gloo": "apiserver-ui",
						}
						grpcPortEnvVar := v1.EnvVar{
							Name:  "GRPC_PORT",
							Value: "10101",
						}
						rbacNamespacedEnvVar := v1.EnvVar{
							Name:  setup.NamespacedRbacEnvName,
							Value: "false",
						}
						noAuthEnvVar := v1.EnvVar{
							Name:  "NO_AUTH",
							Value: "1",
						}
						uiContainer := v1.Container{
							Name:            "apiserver-ui",
							Image:           "quay.io/solo-io/grpcserver-ui:" + version,
							ImagePullPolicy: v1.PullAlways,
							VolumeMounts: []v1.VolumeMount{
								{Name: "empty-cache", MountPath: "/var/cache/nginx"},
								{Name: "empty-run", MountPath: "/var/run"},
							},
							Ports: []v1.ContainerPort{{Name: "static", ContainerPort: 8080, Protocol: v1.ProtocolTCP}},
						}
						grpcServerContainer := v1.Container{
							Name:            "apiserver",
							Image:           "quay.io/solo-io/grpcserver-ee:" + version,
							ImagePullPolicy: v1.PullAlways,
							Ports:           []v1.ContainerPort{{Name: "grpcport", ContainerPort: 10101, Protocol: v1.ProtocolTCP}},
							Env: []v1.EnvVar{
								GetPodNamespaceEnvVar(),
								rbacNamespacedEnvVar,
								grpcPortEnvVar,
								statsEnvVar,
								noAuthEnvVar,
							},
						}
						envoyContainer := v1.Container{
							Name:            "gloo-grpcserver-envoy",
							Image:           "quay.io/solo-io/grpcserver-envoy:" + version,
							ImagePullPolicy: v1.PullAlways,
							VolumeMounts: []v1.VolumeMount{
								{Name: "envoy-config", MountPath: "/etc/envoy", ReadOnly: true},
							},
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{HTTPGet: &v1.HTTPGetAction{
									Path: "/",
									Port: intstr.IntOrString{IntVal: 8080},
								}},
								InitialDelaySeconds: 5,
								PeriodSeconds:       10,
							},
							Env: []v1.EnvVar{
								{
									Name:  "ENVOY_UID",
									Value: "0",
								},
							},
							SecurityContext: &v1.SecurityContext{
								RunAsUser: aws.Int64(101),
							},
						}

						nonRootUser := int64(10101)
						nonRoot := true

						nonRootSC := &v1.PodSecurityContext{
							RunAsUser:    &nonRootUser,
							RunAsNonRoot: &nonRoot,
						}

						rb := ResourceBuilder{
							Namespace: namespace,
							Name:      "api-server",
							Labels:    labels,
						}
						deploy = rb.GetDeploymentAppsv1()
						deploy.Spec.Selector.MatchLabels = selector
						deploy.Spec.Template.ObjectMeta.Labels = selector
						deploy.Spec.Template.Spec.Volumes = []v1.Volume{
							{Name: "empty-cache", VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}}},
							{Name: "empty-run", VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}}},
							{Name: "envoy-config", VolumeSource: v1.VolumeSource{ConfigMap: &v1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{
									Name: defaultBootstrapConfigMapName,
								},
							}}},
						}
						deploy.Spec.Template.Spec.Containers = []v1.Container{uiContainer, grpcServerContainer, envoyContainer}
						deploy.Spec.Template.Spec.ServiceAccountName = "apiserver-ui"
						deploy.Spec.Template.ObjectMeta.Annotations = normalPromAnnotations
						deploy.Spec.Template.Spec.SecurityContext = nonRootSC
						deploy.Spec.Replicas = nil
					})

					It("is there by default", func() {
						testManifest, err := BuildTestManifest(install.GlooOsWithUiChartName, namespace, helmValues{})
						Expect(err).NotTo(HaveOccurred())
						testManifest.ExpectDeploymentAppsV1(deploy)
					})

					It("can customize number of replicas", func() {
						testManifest, err := BuildTestManifest(install.GlooOsWithUiChartName, namespace, helmValues{
							valuesArgs: []string{"apiServer.deployment.replicas=2"},
						})

						Expect(err).NotTo(HaveOccurred())

						customNumReplicas := int32(2)
						deploy.Spec.Replicas = &customNumReplicas

						testManifest.ExpectDeploymentAppsV1(deploy)

					})

					It("does render the default bootstrap config map for the envoy sidecar", func() {
						testManifest, err := BuildTestManifest(install.GlooOsWithUiChartName, namespace, helmValues{})
						Expect(err).NotTo(HaveOccurred())
						testManifest.Expect("ConfigMap", namespace, defaultBootstrapConfigMapName).NotTo(BeNil())
					})

					When("a custom bootstrap config for the API server envoy sidecar is provided", func() {
						const customConfigMapName = "custom-bootstrap-config"
						var actualManifest TestManifest

						BeforeEach(func() {
							var err error
							actualManifest, err = BuildTestManifest(install.GlooOsWithUiChartName, namespace, helmValues{
								valuesArgs: []string{
									"apiServer.deployment.envoy.bootstrapConfig.configMapName=" + customConfigMapName,
								},
							})
							Expect(err).NotTo(HaveOccurred())
						})

						It("adds the custom config map to the API server deployment volume mounts instead of the default one", func() {
							deploy.Spec.Template.Spec.Volumes = []v1.Volume{
								{Name: "empty-cache", VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}}},
								{Name: "empty-run", VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}}},
								{Name: "envoy-config", VolumeSource: v1.VolumeSource{ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: customConfigMapName,
									},
								}}},
							}
							actualManifest.ExpectDeploymentAppsV1(deploy)
						})

						It("does not render the default config map", func() {
							actualManifest.Expect("ConfigMap", namespace, defaultBootstrapConfigMapName).To(BeNil())
						})
					})
				})
			})
		})
	})
})

func constructResourceID(resource *unstructured.Unstructured) string {
	//technically vulnerable to resources that have commas in their names, but that's not a big concern
	return fmt.Sprintf("%s,%s,%s", resource.GetNamespace(), resource.GetName(), resource.GroupVersionKind().String())
}
