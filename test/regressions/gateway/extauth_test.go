package gateway_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"time"

	kubeerrors "k8s.io/apimachinery/pkg/api/errors"

	envoyutil "github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth"
	"github.com/solo-io/go-utils/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"

	"github.com/solo-io/go-utils/testutils"
	"k8s.io/client-go/kubernetes"

	"github.com/solo-io/go-utils/testutils/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v2"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	skerrors "github.com/solo-io/solo-kit/pkg/errors"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"k8s.io/client-go/rest"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

var _ = Describe("External auth", func() {

	var (
		ctx        context.Context
		cancel     context.CancelFunc
		cfg        *rest.Config
		kubeClient kubernetes.Interface

		gatewayClient        v2.GatewayClient
		virtualServiceClient v1.VirtualServiceClient

		err error
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())

		kubeClient, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		cache := kube.NewKubeCache(ctx)
		gatewayClientFactory := &factory.KubeResourceClientFactory{
			Crd:         v2.GatewayCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}
		virtualServiceClientFactory := &factory.KubeResourceClientFactory{
			Crd:         v1.VirtualServiceCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}

		gatewayClient, err = v2.NewGatewayClient(gatewayClientFactory)
		Expect(err).NotTo(HaveOccurred())

		virtualServiceClient, err = v1.NewVirtualServiceClient(virtualServiceClientFactory)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("authenticate requests via LDAP", func() {

		const (
			ldapAssetDir               = "./../assets/ldap"
			ldapServerConfigDirName    = "ldif"
			ldapServerManifestFilename = "ldap-server-manifest.yaml"
			gatewayPort                = int(80)
			response401                = "HTTP/1.1 401 Unauthorized"
			response403                = "HTTP/1.1 403 Forbidden"
			response200                = "HTTP/1.1 200 OK"
		)

		BeforeEach(func() {

			By("create a config map containing the bootstrap configuration for the LDAP server", func() {
				err = testutils.Kubectl(
					"create", "configmap", "ldap", "-n", testHelper.InstallNamespace, "--from-file", filepath.Join(ldapAssetDir, ldapServerConfigDirName))
				Expect(err).NotTo(HaveOccurred())

				Eventually(func() error {
					_, err := kubeClient.CoreV1().ConfigMaps(testHelper.InstallNamespace).Get("ldap", metav1.GetOptions{})
					return err
				}, "15s", "0.5s").Should(BeNil())
			})

			By("deploy an the LDAP server with the correspondent service", func() {
				err = testutils.Kubectl("apply", "-n", testHelper.InstallNamespace, "-f", filepath.Join(ldapAssetDir, ldapServerManifestFilename))
				Expect(err).NotTo(HaveOccurred())

				Eventually(func() error {
					_, err := kubeClient.CoreV1().Services(testHelper.InstallNamespace).Get("ldap", metav1.GetOptions{})
					return err
				}, "15s", "0.5s").Should(BeNil())

				Eventually(func() error {
					deployment, err := kubeClient.AppsV1().Deployments(testHelper.InstallNamespace).Get("ldap", metav1.GetOptions{})
					if err != nil {
						return err
					}
					if deployment.Status.AvailableReplicas == 0 {
						return errors.New("no available replicas for LDAP server deployment")
					}
					return nil
				}, "30s", "0.5s").Should(BeNil())
			})

			By("create an LDAP-secured route to the test upstream", func() {
				extAuthConfig, err := envoyutil.MessageToStruct(&extauthapi.VhostExtension{
					Configs: []*extauthapi.AuthConfig{
						{
							AuthConfig: &extauthapi.AuthConfig_Ldap{
								Ldap: &extauthapi.Ldap{
									Address:        fmt.Sprintf("ldap.%s.svc.cluster.local:389", testHelper.InstallNamespace),
									UserDnTemplate: "uid=%s,ou=people,dc=solo,dc=io",
									AllowedGroups: []string{
										"cn=managers,ou=groups,dc=solo,dc=io",
									},
								},
							},
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())

				virtualHostExtensions := &gloov1.Extensions{
					Configs: map[string]*types.Struct{
						extauth.ExtensionName: extAuthConfig,
					},
				}

				writeVhost(virtualServiceClient, virtualHostExtensions, nil, nil)

				defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
				// wait for default gateway to be created
				Eventually(func() (*v2.Gateway, error) {
					return gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
				}, "15s", "0.5s").Should(Not(BeNil()))
			})
		})

		AfterEach(func() {
			cancel()

			isNotFound := func(err error) bool {
				return err != nil && kubeerrors.IsNotFound(err)
			}

			// Delete config map
			err := kubeClient.CoreV1().ConfigMaps(testHelper.InstallNamespace).Delete("ldap", &metav1.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())
			Eventually(func() bool {
				_, err := kubeClient.CoreV1().ConfigMaps(testHelper.InstallNamespace).Get("ldap", metav1.GetOptions{})
				return isNotFound(err)
			}, "15s", "0.5s").Should(BeTrue())

			// Delete LDAP server deployment and service
			err = testutils.Kubectl("delete", "-n", testHelper.InstallNamespace, "-f", filepath.Join(ldapAssetDir, ldapServerManifestFilename))
			Expect(err).NotTo(HaveOccurred())
			Eventually(func() bool {
				_, err := kubeClient.CoreV1().Services(testHelper.InstallNamespace).Get("ldap", metav1.GetOptions{})
				return isNotFound(err)
			}, "15s", "0.5s").Should(BeTrue())
			Eventually(func() bool {
				_, err := kubeClient.AppsV1().Deployments(testHelper.InstallNamespace).Get("ldap", metav1.GetOptions{})
				return isNotFound(err)
			}, "15s", "0.5s").Should(BeTrue())

			// Delete virtual service
			err = virtualServiceClient.Delete(testHelper.InstallNamespace, "vs", clients.DeleteOpts{Ctx: ctx})
			Expect(err).NotTo(HaveOccurred())
			Eventually(func() bool {
				_, err := virtualServiceClient.Read(testHelper.InstallNamespace, "vs", clients.ReadOpts{Ctx: ctx})
				if err != nil && skerrors.IsNotExist(err) {
					return true
				}
				return false
			}, "15s", "0.5s").Should(BeTrue())
		})

		// Credentials must be in the <username>:<password> format
		buildAuthHeader := func(credentials string) map[string]string {
			encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
			return map[string]string{
				"Authorization": fmt.Sprintf("Basic %s", encodedCredentials),
			}
		}

		curlAndAssertResponse := func(headers map[string]string, expectedResponseSubstring string) {
			testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
				Protocol:          "http",
				Path:              "/",
				Method:            "GET",
				Headers:           headers,
				Host:              defaults.GatewayProxyName,
				Service:           defaults.GatewayProxyName,
				Port:              gatewayPort,
				ConnectionTimeout: 10,   // this is important, as the first curl call sometimes hangs indefinitely
				Verbose:           true, // this is important, as curl will only output status codes with verbose output
			}, expectedResponseSubstring, 1, time.Minute)
		}

		It("works as expected", func() {

			By("returns 401 if no authentication header is provided", func() {
				curlAndAssertResponse(nil, response401)
			})

			By("returns 401 if the user is unknown", func() {
				curlAndAssertResponse(buildAuthHeader("john:doe"), response401)
			})

			By("returns 403 if the user does not belong to the allowed groups", func() {
				curlAndAssertResponse(buildAuthHeader("marco:marcopwd"), response403)
			})

			By("returns 200 if the user belongs to one of the allowed groups", func() {
				curlAndAssertResponse(buildAuthHeader("rick:rickpwd"), response200)
			})
		})
	})
})
