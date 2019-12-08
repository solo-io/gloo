package test

import (
	"os"

	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/solo-projects/install/helm/gloo-ee/generate"
	"github.com/solo-io/solo-projects/pkg/install"
	k8s "k8s.io/apimachinery/pkg/apis/meta/v1"
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
						},
						Resources:       v1.ResourceRequirements{},
						ImagePullPolicy: "Always",
					},
				}
				observabilityDeployment.Spec.Template.Spec.ServiceAccountName = "observability"
				observabilityDeployment.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{
					{
						Name: "solo-io-readerbot-pull-secret",
					},
				}
				observabilityDeployment.Spec.Strategy = appsv1.DeploymentStrategy{}
				observabilityDeployment.Spec.Selector.MatchLabels = selector
				observabilityDeployment.Spec.Template.ObjectMeta.Labels = selector

				grafanaBuilder := ResourceBuilder{
					Namespace: "", // grafana installs to empty namespace during tests
					Name:      "release-name-grafana",
					Labels:    labels,
				}
				grafanaDeployment = grafanaBuilder.GetDeploymentAppsv1()
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
					deploy.Spec.Template.Spec.Containers[0].SecurityContext = &v1.SecurityContext{
						Capabilities: &v1.Capabilities{
							Add:  []v1.Capability{"NET_BIND_SERVICE"},
							Drop: []v1.Capability{"ALL"},
						},
						ReadOnlyRootFilesystem:   &truez,
						AllowPrivilegeEscalation: &falsez,
					}

					deploy.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{
						{
							Name: "solo-io-readerbot-pull-secret",
						},
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
				Context("apiserver deployment", func() {
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
								grpcPortEnvVar,
								noAuthEnvVar,
								licenseEnvVar,
							},
						}
						envoyContainer := v1.Container{
							Name:            "gloo-grpcserver-envoy",
							Image:           "quay.io/solo-io/grpcserver-envoy:" + version,
							ImagePullPolicy: v1.PullAlways,
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{HTTPGet: &v1.HTTPGetAction{
									Path: "/",
									Port: intstr.IntOrString{IntVal: 8080},
								}},
								InitialDelaySeconds: 5,
								PeriodSeconds:       10,
							},
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
						}
						deploy.Spec.Template.Spec.Containers = []v1.Container{uiContainer, grpcServerContainer, envoyContainer}
						deploy.Spec.Template.Spec.ServiceAccountName = "apiserver-ui"
						deploy.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{
							{
								Name: "solo-io-readerbot-pull-secret",
							},
						}
					})

					It("is there by default", func() {
						testManifest, err := BuildTestManifest(install.GlooEnterpriseChartName, namespace, helmValues{})
						Expect(err).NotTo(HaveOccurred())
						testManifest.ExpectDeploymentAppsV1(deploy)
					})
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
			glooOsVersion, err = generate.GetGlooOsVersion("../..", glooEGenerationFiles, glooOsWithReadOnlyUiGenerationFiles)
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
					deploy.Spec.Template.Spec.Containers[0].SecurityContext = &v1.SecurityContext{
						Capabilities: &v1.Capabilities{
							Add:  []v1.Capability{"NET_BIND_SERVICE"},
							Drop: []v1.Capability{"ALL"},
						},
						ReadOnlyRootFilesystem:   &truez,
						AllowPrivilegeEscalation: &falsez,
					}

					deploy.Spec.Template.Spec.ServiceAccountName = "gateway-proxy"

					gatewayProxyDeployment = deploy
				})

				It("creates a deployment without envoy config annotations", func() {
					testManifest, err := BuildTestManifest(install.GlooOsWithUiChartName, namespace, helmValues{})
					Expect(err).NotTo(HaveOccurred())
					testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
				})

				It("creates a deployment with envoy config annotations", func() {
					testManifest, err := BuildTestManifest(install.GlooOsWithUiChartName, namespace, helmValues{
						valuesArgs: []string{"gloo.gatewayProxies.gatewayProxy.readConfig=true"},
					})
					Expect(err).NotTo(HaveOccurred())
					includeStatConfig()
					testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
				})

				Context("apiserver deployment", func() {
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
								grpcPortEnvVar,
								noAuthEnvVar,
							},
						}
						envoyContainer := v1.Container{
							Name:            "gloo-grpcserver-envoy",
							Image:           "quay.io/solo-io/grpcserver-envoy:" + version,
							ImagePullPolicy: v1.PullAlways,
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{HTTPGet: &v1.HTTPGetAction{
									Path: "/",
									Port: intstr.IntOrString{IntVal: 8080},
								}},
								InitialDelaySeconds: 5,
								PeriodSeconds:       10,
							},
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
						}
						deploy.Spec.Template.Spec.Containers = []v1.Container{uiContainer, grpcServerContainer, envoyContainer}
						deploy.Spec.Template.Spec.ServiceAccountName = "apiserver-ui"
					})

					It("is there by default", func() {
						testManifest, err := BuildTestManifest(install.GlooOsWithUiChartName, namespace, helmValues{})
						Expect(err).NotTo(HaveOccurred())
						testManifest.ExpectDeploymentAppsV1(deploy)
					})
				})
			})

		})
	})
})
