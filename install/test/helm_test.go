package test

import (
	"io/ioutil"
	"os"

	"k8s.io/apimachinery/pkg/util/intstr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gateway/pkg/translator"
	. "github.com/solo-io/go-utils/manifesttestutils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Helm Test", func() {

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

		prepareMakefile := func(helmFlags string) {
			makefileSerializer.Lock()
			defer makefileSerializer.Unlock()

			f, err := ioutil.TempFile("", "*.yaml")
			Expect(err).NotTo(HaveOccurred())
			err = f.Close()
			Expect(err).ToNot(HaveOccurred())
			manifestYaml = f.Name()

			MustMake(".", "-C", "../..", "install/glooe-gateway.yaml", "HELMFLAGS="+helmFlags, "OUTPUT_YAML="+manifestYaml)
			testManifest = NewTestManifest(manifestYaml)
		}

		Context("gateway", func() {
			BeforeEach(func() {
				labels = map[string]string{
					"gloo":             "gateway-proxy",
					"gateway-proxy-id": translator.GatewayProxyName,
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
						"gateway-proxy-id": translator.GatewayProxyName,
					}
					podLabels := map[string]string{
						"gloo":             "gateway-proxy",
						"gateway-proxy-id": translator.GatewayProxyName,
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
					container.Name = translator.GatewayProxyName
					container.Args = []string{"--disable-hot-restart"}

					rb := ResourceBuilder{
						Namespace:  namespace,
						Name:       translator.GatewayProxyName,
						Labels:     labels,
						Containers: []ContainerSpec{container},
					}
					deploy := rb.GetDeploymentAppsv1()
					deploy.Spec.Selector = &metav1.LabelSelector{
						MatchLabels: selector,
					}
					deploy.Spec.Template.ObjectMeta.Labels = podLabels
					deploy.Spec.Template.ObjectMeta.Annotations = podAnnotations
					deploy.Spec.Template.Spec.Volumes = []v1.Volume{{
						Name: "envoy-config",
						VolumeSource: v1.VolumeSource{
							ConfigMap: &v1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{
									Name: "gateway-proxy-v2-envoy-config",
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
					helmFlags := "--namespace " + namespace + " --set namespace.create=true"

					prepareMakefile(helmFlags)
					testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
				})

				It("creates a deployment with envoy config annotations", func() {
					helmFlags := "--namespace " + namespace + " --set namespace.create=true " +
						"--set gloo.gatewayProxies.gatewayProxyV2.readConfig=true"

					prepareMakefile(helmFlags)
					includeStatConfig()
					testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
				})

				It("creates a deployment without extauth sidecar", func() {
					helmFlags := "--namespace " + namespace + " --set namespace.create=true"

					prepareMakefile(helmFlags)
					testManifest.ExpectDeploymentAppsV1(gatewayProxyDeployment)
				})

				It("creates a deployment with extauth sidecar", func() {
					helmFlags := "--namespace " + namespace + " --set namespace.create=true " +
						"--set global.extensions.extAuth.envoySidecar=true "
					prepareMakefile(helmFlags)

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
						helmFlags := "--namespace " + namespace + " --set namespace.create=true"
						prepareMakefile(helmFlags)
						testManifest.ExpectDeploymentAppsV1(deploy)
					})
				})
			})

		})
	})
})
