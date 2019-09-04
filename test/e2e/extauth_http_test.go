package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	envoyutil "github.com/envoyproxy/go-control-plane/pkg/util"
	extauth2 "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"

	"github.com/gogo/protobuf/types"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/test/services"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth"

	"github.com/fgrosse/zaptest"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

var _ = Describe("External http", func() {

	var (
		ctx             context.Context
		cancel          context.CancelFunc
		testClients     services.TestClients
		extauthSettings *extauth.Settings
		extauthport     uint32
		extauthserver   *gloov1.Upstream
	)

	BeforeEach(func() {
		extauthport = atomic.AddUint32(&baseExtauthPort, 1) + uint32(config.GinkgoConfig.ParallelNode*1000)

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())

		extauthserver = &gloov1.Upstream{
			Metadata: core.Metadata{
				Name:      "extauth-server",
				Namespace: "default",
			},
			UpstreamSpec: &gloov1.UpstreamSpec{
				UpstreamType: &gloov1.UpstreamSpec_Static{
					Static: &gloov1static.UpstreamSpec{
						Hosts: []*gloov1static.Host{{
							Addr: "localhost",
							Port: extauthport,
						}},
					},
				},
			},
		}

		ref := extauthserver.Metadata.Ref()
		extauthSettings = &extauth.Settings{
			ExtauthzServerRef: &ref,
			HttpService:       &extauth.HttpService{},
		}

	})
	JustBeforeEach(func() {
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		_, err := testClients.UpstreamClient.Write(extauthserver, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		settingsStruct, err := envoyutil.MessageToStruct(extauthSettings)
		Expect(err).NotTo(HaveOccurred())

		extensions := &gloov1.Extensions{
			Configs: map[string]*types.Struct{
				extauth2.ExtensionName: settingsStruct,
			},
		}

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}
		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, extensions)
		go func(testctx context.Context) {
			defer GinkgoRecover()
			srvMux := http.NewServeMux()
			fail := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			}
			ok := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}
			srvMux.HandleFunc("/deny", fail)
			srvMux.HandleFunc("/allow", ok)
			srvMux.HandleFunc("/prefix/prefixdeny", fail)
			srvMux.HandleFunc("/prefix/prefixallow", ok)

			var srv http.Server
			srv.Handler = srvMux
			srv.Addr = fmt.Sprintf(":%d", extauthport)
			go func() {
				<-testctx.Done()
				srv.Close()
			}()
			err := srv.ListenAndServe()
			if err != http.ErrServerClosed {
				Expect(err).NotTo(HaveOccurred())
			}
		}(ctx)
	})

	AfterEach(func() {
		cancel()
	})

	Context("With envoy", func() {

		var (
			envoyInstance *services.EnvoyInstance
			testUpstream  *v1helpers.TestUpstream
			envoyPort     = uint32(8080)
		)

		JustBeforeEach(func() {
			var err error
			envoyInstance, err = envoyFactory.NewEnvoyInstance()
			Expect(err).NotTo(HaveOccurred())

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())

			testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

			var opts clients.WriteOpts
			up := testUpstream.Upstream
			_, err = testClients.UpstreamClient.Write(up, opts)
			Expect(err).NotTo(HaveOccurred())

		})

		AfterEach(func() {
			if envoyInstance != nil {
				envoyInstance.Clean()
			}
		})

		Context("custom sanity tests", func() {
			JustBeforeEach(func() {

				// drain channel as we dont care about it
				go func() {
					for range testUpstream.C {
					}
				}()

				proxy := getProxyCustomAuth(envoyPort, testUpstream.Upstream.Metadata.Ref())

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				Eventually(func() (core.Status, error) {
					proxy, err := testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
					if err != nil {
						return core.Status{}, err
					}

					return proxy.Status, nil
				}, "5s", "0.1s").Should(MatchFields(IgnoreExtras, Fields{
					"Reason": BeEmpty(),
					"State":  Equal(core.Status_Accepted),
				}))
			})

			It("should deny ext auth envoy", func() {
				Eventually(func() (int, error) {
					resp, err := http.Get(fmt.Sprintf("http://%s:%d/deny", "localhost", envoyPort))
					if err != nil {
						return 0, err
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusUnauthorized))
			})

			It("should allow ext auth envoy", func() {
				Eventually(func() (int, error) {
					resp, err := http.Get(fmt.Sprintf("http://%s:%d/allow", "localhost", envoyPort))
					if err != nil {
						return 0, err
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))
			})
			Context("with prefix", func() {
				BeforeEach(func() {
					ref := extauthserver.Metadata.Ref()
					extauthSettings = &extauth.Settings{
						ExtauthzServerRef: &ref,
						HttpService: &extauth.HttpService{
							// test that prefix with traling slash is handled correctly
							PathPrefix: "/prefix/",
						},
					}
				})

				It("should allow ext auth envoy", func() {
					Eventually(func() (int, error) {
						resp, err := http.Get(fmt.Sprintf("http://%s:%d/prefixallow", "localhost", envoyPort))
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusOK))
				})

				It("should deny ext auth envoy", func() {
					Eventually(func() (int, error) {
						resp, err := http.Get(fmt.Sprintf("http://%s:%d/prefixdeny", "localhost", envoyPort))
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusUnauthorized))
				})
			})
		})
	})
})

func getProxyCustomAuth(envoyPort uint32, upstream core.ResourceRef) *gloov1.Proxy {

	extauthCfg := &extauth.VhostExtension{
		AuthConfig: &extauth.VhostExtension_CustomAuth{}}
	return getProxyExtAuth(envoyPort, upstream, extauthCfg)
}
