package extauth_test

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	v12 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-projects/test/e2e"

	"github.com/solo-io/gloo/test/ginkgo/parallel"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	baseExtauthPort = uint32(27000)
)

var _ = Describe("External http", func() {

	var (
		testContext   *e2e.TestContext
		extauthServer *gloov1.Upstream
	)

	startLocalHttpExtAuthServer := func(ctx context.Context, port uint32) {
		srvMux := http.NewServeMux()
		fail := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}
		ok := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}
		checkHeaders := func(w http.ResponseWriter, r *http.Request) {
			if len(r.Header.Get("allowed")) > 0 {
				w.WriteHeader(http.StatusOK)
				return
			}
			if len(r.Header.Get("pattern")) > 0 {
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		}
		srvMux.HandleFunc("/deny", fail)
		srvMux.HandleFunc("/allow", ok)
		srvMux.HandleFunc("/prefix/prefixdeny", fail)
		srvMux.HandleFunc("/prefix/prefixallow", ok)
		srvMux.HandleFunc("/headers", checkHeaders)

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
		testContext = testContextFactory.NewTestContext()
		testContext.BeforeEach()

		// Start HTTP ext auth server
		extauthPort := atomic.AddUint32(&baseExtauthPort, 1) + uint32(parallel.GetPortOffset())
		go func(testCtx context.Context) {
			defer GinkgoRecover()
			startLocalHttpExtAuthServer(testCtx, extauthPort)
		}(testContext.Ctx())
		extauthServer = &gloov1.Upstream{
			Metadata: &core.Metadata{
				Name:      "extauth-server",
				Namespace: "default",
			},
			UpstreamType: &gloov1.Upstream_Static{
				Static: &gloov1static.UpstreamSpec{
					Hosts: []*gloov1static.Host{{
						Addr: testContext.EnvoyInstance().LocalAddr(),
						Port: extauthPort,
					}},
				},
			},
		}
		testContext.ResourcesToCreate().Upstreams = append(testContext.ResourcesToCreate().Upstreams, extauthServer)

		// Point Gloo settings to the ext auth server
		testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
			settings.Extauth = &extauth.Settings{
				ExtauthzServerRef: extauthServer.Metadata.Ref(),
				ServiceType: &extauth.Settings_HttpService{
					HttpService: &extauth.HttpService{},
				},
			}
		})
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	AfterEach(func() {
		testContext.AfterEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	Context("custom sanity tests", func() {

		JustBeforeEach(func() {
			testContext.PatchDefaultVirtualService(func(vs *v12.VirtualService) *v12.VirtualService {
				vsBuilder := helpers.BuilderFromVirtualService(vs)
				vsBuilder.WithRouteOptions(e2e.DefaultRouteName, &gloov1.RouteOptions{
					Extauth: &extauth.ExtAuthExtension{
						Spec: &extauth.ExtAuthExtension_CustomAuth{
							CustomAuth: &extauth.CustomAuth{},
						}},
				}).WithRouteActionToUpstreamRef(e2e.DefaultRouteName, testContext.TestUpstream().Upstream.GetMetadata().Ref())
				return vsBuilder.Build()
			})
		})

		// The "/deny" path on the custom ext auth server should always return 401 - logic in the `startLocalHttpExtAuthServer` function.
		It("should deny ext auth envoy", func() {
			httpBuilder := testContext.GetHttpRequestBuilder().WithPath("deny")
			Eventually(func(g Gomega) {
				g.Expect(testutils.DefaultHttpClient.Do(httpBuilder.Build())).To(HaveHTTPStatus(http.StatusUnauthorized))
			}, "5s", "0.5s").Should(Succeed())
		})

		// The "/allow" path on the custom ext auth server should always return 200 - logic in the `startLocalHttpExtAuthServer` function.
		It("should allow ext auth envoy", func() {
			httpBuilder := testContext.GetHttpRequestBuilder().WithPath("allow")
			Eventually(func(g Gomega) {
				g.Expect(testutils.DefaultHttpClient.Do(httpBuilder.Build())).To(HaveHTTPStatus(http.StatusOK))
			}, "5s", "0.5s").Should(Succeed())
		})

		Context("with allowed headers", func() {
			BeforeEach(func() {
				testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
					settings.Extauth.ServiceType.(*extauth.Settings_HttpService).HttpService.Request = &extauth.HttpService_Request{
						AllowedHeaders:      []string{"allowed"},
						AllowedHeadersRegex: []string{"pa[ter]+n"},
					}
				})
			})

			It("should allow ext auth with exact header match present", func() {
				httpBuilder := testContext.GetHttpRequestBuilder().WithPath("headers").WithHeader("allowed", "foo")
				Eventually(func(g Gomega) {
					g.Expect(testutils.DefaultHttpClient.Do(httpBuilder.Build())).To(HaveHTTPStatus(http.StatusOK))
				}, "5s", "0.5s").Should(Succeed())
			})

			It("should allow ext auth with regex header match present", func() {
				httpBuilder := testContext.GetHttpRequestBuilder().WithPath("headers").WithHeader("pattern", "foo")
				Eventually(func(g Gomega) {
					g.Expect(testutils.DefaultHttpClient.Do(httpBuilder.Build())).To(HaveHTTPStatus(http.StatusOK))
				}, "5s", "0.5s").Should(Succeed())
			})

			It("should deny ext auth without header match present", func() {
				httpBuilder := testContext.GetHttpRequestBuilder().WithPath("headers").WithHeader("denied", "foo")
				Eventually(func(g Gomega) {
					g.Expect(testutils.DefaultHttpClient.Do(httpBuilder.Build())).To(HaveHTTPStatus(http.StatusUnauthorized))
				}, "5s", "0.5s").Should(Succeed())
			})
		})

		Context("with path prefix", func() {

			BeforeEach(func() {
				testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
					settings.Extauth.ServiceType.(*extauth.Settings_HttpService).HttpService.PathPrefix = "/prefix/"
				})
			})

			It("should allow ext auth envoy", func() {
				httpBuilder := testContext.GetHttpRequestBuilder().WithPath("prefixallow")
				Eventually(func(g Gomega) {
					g.Expect(testutils.DefaultHttpClient.Do(httpBuilder.Build())).To(HaveHTTPStatus(http.StatusOK))
				}, "5s", "0.5s").Should(Succeed())
			})

			It("should deny ext auth envoy", func() {
				httpBuilder := testContext.GetHttpRequestBuilder().WithPath("prefixdeny")
				Eventually(func(g Gomega) {
					g.Expect(testutils.DefaultHttpClient.Do(httpBuilder.Build())).To(HaveHTTPStatus(http.StatusUnauthorized))
				}, "5s", "0.5s").Should(Succeed())
			})
		})

	})
})
