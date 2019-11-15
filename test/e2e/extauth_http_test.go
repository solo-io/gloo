package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/test/services"

	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"

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
		extauthServer   *gloov1.Upstream

		envoyInstance *services.EnvoyInstance
		testUpstream  *v1helpers.TestUpstream
		envoyPort     = uint32(8080)
		cache         memory.InMemoryResourceCache

		err error
	)

	startLocalHttpExtAuthServer := func(ctx context.Context, port uint32) {
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
		srv.Addr = fmt.Sprintf(":%d", port)
		go func() {
			<-ctx.Done()
			_ = srv.Close()
		}()
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			Expect(err).NotTo(HaveOccurred())
		}
	}

	BeforeEach(func() {
		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())

		// Get test clients
		cache = memory.NewInMemoryResourceCache()
		testClients = services.GetTestClients(cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		// Start Envoy
		envoyInstance, err = envoyFactory.NewEnvoyInstance()
		Expect(err).NotTo(HaveOccurred())

		err = envoyInstance.Run(testClients.GlooPort)
		Expect(err).NotTo(HaveOccurred())

		// Start HTTP ext auth server
		extauthPort := atomic.AddUint32(&baseExtauthPort, 1) + uint32(config.GinkgoConfig.ParallelNode*1000)
		go func(testCtx context.Context) {
			defer GinkgoRecover()
			startLocalHttpExtAuthServer(testCtx, extauthPort)
		}(ctx)
		extauthServer = &gloov1.Upstream{
			Metadata: core.Metadata{
				Name:      "extauth-server",
				Namespace: "default",
			},
			UpstreamSpec: &gloov1.UpstreamSpec{
				UpstreamType: &gloov1.UpstreamSpec_Static{
					Static: &gloov1static.UpstreamSpec{
						Hosts: []*gloov1static.Host{{
							Addr: envoyInstance.LocalAddr(),
							Port: extauthPort,
						}},
					},
				},
			},
		}
		_, err = testClients.UpstreamClient.Write(extauthServer, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		// Create a test upstream
		testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())
		_, err = testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		// Start Gloo, pointing it to the ext auth server
		ref := extauthServer.Metadata.Ref()
		extauthSettings = &extauth.Settings{
			ExtauthzServerRef: &ref,
			HttpService:       &extauth.HttpService{},
		}
	})

	JustBeforeEach(func() {
		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}
		settings := &gloov1.Settings{
			Extauth: extauthSettings,
		}
		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, nil, settings)
	})

	AfterEach(func() {
		cancel()
		if envoyInstance != nil {
			_ = envoyInstance.Clean()
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
				proxyFromStorage, err := testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				if err != nil {
					return core.Status{}, err
				}

				return proxyFromStorage.Status, nil
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
				ref := extauthServer.Metadata.Ref()
				extauthSettings = &extauth.Settings{
					ExtauthzServerRef: &ref,
					HttpService: &extauth.HttpService{
						// test that prefix with trailing slash is handled correctly
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

func getProxyCustomAuth(envoyPort uint32, upstream core.ResourceRef) *gloov1.Proxy {
	extauthCfg := &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_CustomAuth{
			CustomAuth: &extauth.CustomAuth{},
		}}
	return getProxyExtAuth(envoyPort, upstream, extauthCfg)
}
