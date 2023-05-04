package e2e_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/fgrosse/zaptest"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/solo-io/ext-auth-service/pkg/server"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/ginkgo/parallel"
	"github.com/solo-io/gloo/test/helpers"
	gloohelpers "github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/rate-limiter/pkg/cache/redis"
	ratelimitserver "github.com/solo-io/rate-limiter/pkg/server"
	rlv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	extauthrunner "github.com/solo-io/solo-projects/projects/extauth/pkg/runner"
	"github.com/solo-io/solo-projects/test/services"
	ratelimitservice "github.com/solo-io/solo-projects/test/services/ratelimit"
	"github.com/solo-io/solo-projects/test/v1helpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

const (
	rateLimitAddr       = "127.0.0.1"
	redisAddr           = "127.0.0.1"
	numRateLimitConfigs = 1000
)

var _ = Describe("XDS interfaces", func() {

	var (
		ctx              context.Context
		cancel           context.CancelFunc
		err              error
		testClients      services.TestClients
		settings         extauthrunner.Settings
		glooSettings     *gloov1.Settings
		cache            memory.InMemoryResourceCache
		extAuthServer    *gloov1.Upstream
		baseExtauthPort  = uint32(27000)
		envoyPort        = uint32(8080)
		redisPort        = uint32(6380)
		redisSession     *gexec.Session
		testUpstream     *v1helpers.TestUpstream
		envoyInstance    *services.EnvoyInstance
		rlServerSettings ratelimitserver.Settings
	)
	// Set up envoy and test clients which is used for both rate limit and extauth
	BeforeEach(func() {
		// These tests create a lot of resources and take some extra time to run so we want to run them as nightlies
		// TODO: run in memory tests nightly with this enabled
		if _, found := os.LookupEnv("RUN_XDS_SCALE_TESTS"); !found {
			Skip("Skipping tests with large xds messages")
		}
		ctx, cancel = context.WithCancel(context.Background())
		cache = memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		envoyInstance, err = envoyFactory.NewEnvoyInstance()
		Expect(err).NotTo(HaveOccurred())

		err = envoyInstance.Run(testClients.GlooPort)
		Expect(err).NotTo(HaveOccurred())
		testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())
		var opts clients.WriteOpts
		up := testUpstream.Upstream
		_, err := testClients.UpstreamClient.Write(up, opts)
		Expect(err).NotTo(HaveOccurred())
	})
	JustBeforeEach(func() {
		// The tests start gloo in later BeforeEach blocks so the US won't be accepted until here
		gloohelpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
			return testClients.UpstreamClient.Read(testUpstream.Upstream.Metadata.Namespace, testUpstream.Upstream.Metadata.Name, clients.ReadOpts{})
		})
	})
	AfterEach(func() {
		cancel()
		if envoyInstance != nil {
			_ = envoyInstance.Clean()
		}
	})
	// Show that we can send very large configuration to Ext-Auth via xds
	// Create a very large API key (which is read and sent over xds) and show that requests still work as expected.
	Context("Extauth XDS", func() {
		waitForHealthyExtauthService := func() {
			// First make sure gloo has found the extauth upstream
			gloohelpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.UpstreamClient.Read(extAuthServer.Metadata.Namespace, extAuthServer.Metadata.Name, clients.ReadOpts{})
			})
			extAuthHealthServerAddr := "localhost:" + strconv.Itoa(settings.ExtAuthSettings.ServerPort)
			conn, err := grpc.Dial(extAuthHealthServerAddr, grpc.WithInsecure())
			Expect(err).ToNot(HaveOccurred())

			// make sure that extauth is up and serving
			healthClient := grpc_health_v1.NewHealthClient(conn)
			var header metadata.MD
			Eventually(func() (grpc_health_v1.HealthCheckResponse_ServingStatus, error) {
				resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
					Service: settings.ExtAuthSettings.ServiceName,
				}, grpc.Header(&header))
				return resp.GetStatus(), err
			}, "10s", ".1s").Should(Equal(grpc_health_v1.HealthCheckResponse_SERVING))
		}
		BeforeEach(func() {
			extAuthPort := atomic.AddUint32(&baseExtauthPort, 1) + uint32(parallel.GetPortOffset())
			extAuthHealthPort := atomic.AddUint32(&baseExtauthPort, 1) + uint32(parallel.GetPortOffset())

			logger := zaptest.LoggerWriter(GinkgoWriter)
			contextutils.SetFallbackLogger(logger.Sugar())

			extauthAddr := "localhost"
			if runtime.GOOS == "darwin" {
				extauthAddr = "host.docker.internal"
			}

			extAuthServer = &gloov1.Upstream{
				Metadata: &core.Metadata{
					Name:      fmt.Sprintf("extauth-server-%d", extAuthPort),
					Namespace: "default",
				},
				UseHttp2: &wrappers.BoolValue{Value: true},
				UpstreamType: &gloov1.Upstream_Static{
					Static: &gloov1static.UpstreamSpec{
						Hosts: []*gloov1static.Host{{
							Addr: extauthAddr,
							Port: extAuthPort,
						}},
					},
				},
			}
			logger.Info(fmt.Sprintf("Expect to see extauth as upstream named %s", extAuthServer.Metadata.Name))
			_, err := testClients.UpstreamClient.Write(extAuthServer, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			ref := extAuthServer.Metadata.Ref()
			extauthSettings := &extauth.Settings{
				ExtauthzServerRef: ref,
				RequestBody: &extauth.BufferSettings{
					MaxRequestBytes:     0,
					AllowPartialMessage: false,
					PackAsBytes:         false,
				},
			}

			s := extauthrunner.NewSettings()
			settings = extauthrunner.Settings{
				GlooAddress: fmt.Sprintf("localhost:%d", testClients.GlooPort),
				ExtAuthSettings: server.Settings{
					DebugPort:                   0,
					ServerPort:                  int(extAuthPort),
					SigningKey:                  "hello",
					UserIdHeader:                "X-User-Id",
					HealthCheckFailTimeout:      1,
					HealthCheckHttpPort:         int(extAuthHealthPort),
					HealthCheckHttpPath:         s.ExtAuthSettings.HealthCheckHttpPath,
					HealthLivenessCheckHttpPath: s.ExtAuthSettings.HealthLivenessCheckHttpPath,
					LogSettings: server.LogSettings{
						// Note(yuval-k): Disable debug logs as they are noisy. If you are writing new
						// tests, uncomment this while developing to increase verbosity. I couldn't find
						// a good way to wire this to GinkgoWriter
						// DebugMode:  "1",
						LoggerName: "extauth-service-test",
					},
				},
			}
			glooSettings = &gloov1.Settings{Extauth: extauthSettings}

			what := services.What{
				DisableGateway: true,
				DisableUds:     true,
				DisableFds:     true,
			}
			services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{Ctx: ctx, Cache: cache, LocalGlooPort: int32(testClients.GlooPort), What: what, Namespace: defaults.GlooSystem, Settings: glooSettings})
			go func(testCtx context.Context) {
				defer GinkgoRecover()
				err := extauthrunner.RunWithSettings(testCtx, settings)
				if testCtx.Err() == nil {
					Expect(err).NotTo(HaveOccurred())
				}
			}(ctx)

			waitForHealthyExtauthService()

			logger.Info("Creating secrets and authconfig")
			_, err = testClients.AuthConfigClient.Write(&extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      getApiKeyExtAuthExtension().GetConfigRef().Name,
					Namespace: getApiKeyExtAuthExtension().GetConfigRef().Namespace,
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_ApiKeyAuth{
						ApiKeyAuth: getApiKeyAuthConfigLarge(),
					},
				}},
			}, clients.WriteOpts{Ctx: ctx})
			Expect(err).NotTo(HaveOccurred())

			apiKey1 := &extauth.ApiKey{
				ApiKey: "secretApiKey1",
			}

			secret1 := &gloov1.Secret{
				Metadata: &core.Metadata{
					Name:      "secret1",
					Namespace: "default",
				},
				Kind: &gloov1.Secret_ApiKey{
					ApiKey: apiKey1,
				},
			}

			apiKey2 := &extauth.ApiKey{
				ApiKey: "secretApiKey2",
			}

			secret2 := &gloov1.Secret{
				Metadata: &core.Metadata{
					Name:      "secret2",
					Namespace: "default",
					Labels:    map[string]string{"team": "infrastructure"},
				},
				Kind: &gloov1.Secret_ApiKey{
					ApiKey: apiKey2,
				},
			}
			// Make a key that's very large
			unreasonableKey := strings.Repeat("very long string", 500000)
			apiKey3 := &extauth.ApiKey{
				ApiKey: unreasonableKey,
			}

			secret3 := &gloov1.Secret{
				Metadata: &core.Metadata{
					Name:      "secret3",
					Namespace: "default",
					Labels:    map[string]string{"team": "infrastructure"},
				},
				Kind: &gloov1.Secret_ApiKey{
					ApiKey: apiKey3,
				},
			}
			_, err = testClients.SecretClient.Write(secret1, clients.WriteOpts{})
			Expect(err).ToNot(HaveOccurred())

			_, err = testClients.SecretClient.Write(secret2, clients.WriteOpts{})
			Expect(err).ToNot(HaveOccurred())

			_, err = testClients.SecretClient.Write(secret3, clients.WriteOpts{})
			Expect(err).ToNot(HaveOccurred())
			proxy := getProxyExtAuthApiKeyAuth(envoyPort, testUpstream.Upstream.Metadata.Ref())

			_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
			})

			logger.Info("Resources persisted")
		})

		AfterEach(func() {
			cancel()
		})

		// Run all the API key sanity tests
		It("should deny ext auth envoy without apikey", func() {
			Eventually(func() (int, error) {
				resp, err := http.Get(fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort))
				if err != nil {
					return 0, err
				}
				defer resp.Body.Close()
				_, _ = io.ReadAll(resp.Body)
				return resp.StatusCode, nil
			}, "20s", "1s").Should(Equal(http.StatusUnauthorized))
		})

		It("should deny ext auth envoy with incorrect apikey", func() {
			Eventually(func() (int, error) {
				req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
				req.Header.Add("api-key", "badApiKey")
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return 0, err
				}
				defer resp.Body.Close()
				_, _ = io.ReadAll(resp.Body)
				return resp.StatusCode, nil
			}, "20s", "1s").Should(Equal(http.StatusUnauthorized))
		})

		It("should accept ext auth envoy with correct apikey -- secret ref match", func() {
			Eventually(func() (int, error) {
				req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
				req.Header.Add("api-key", "secretApiKey1")
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return 0, err
				}
				defer resp.Body.Close()
				_, _ = io.ReadAll(resp.Body)
				return resp.StatusCode, nil
			}, "20s", "1s").Should(Equal(http.StatusOK))
		})

		It("should accept ext auth envoy with correct apikey -- label match", func() {
			Eventually(func() (int, error) {
				req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
				req.Header.Add("api-key", "secretApiKey2")
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return 0, err
				}
				defer resp.Body.Close()
				_, _ = io.ReadAll(resp.Body)
				return resp.StatusCode, nil
			}, "20s", "1s").Should(Equal(http.StatusOK))
		})
	})
	Context("Rate limit xds", func() {

		BeforeEach(func() {
			glooSettings = &gloov1.Settings{}
			glooSettings.Ratelimit = &ratelimit.ServiceSettings{
				SetDescriptors: getManySetDescriptors(numRateLimitConfigs),
			}
			rlServerSettings = ratelimitserver.NewSettings()
			rlServerSettings.HealthFailTimeout = 2 // seconds
			rlServerSettings.RateLimitPort = int(atomic.AddUint32(&baseRateLimitPort, 1) + uint32(parallel.GetPortOffset()))
			rlServerSettings.ReadyPort = int(atomic.AddUint32(&baseRateLimitPort, 1) + uint32(parallel.GetPortOffset()))
			rlServerSettings.RedisSettings = redis.Settings{}

			rlServerSettings.RedisSettings = redis.NewSettings()
			rlServerSettings.RedisSettings.Url = fmt.Sprintf("%s:%d", redisAddr, redisPort)
			rlServerSettings.RedisSettings.SocketType = "tcp"
			rlServerSettings.RedisSettings.Clustered = false

			command := exec.Command(getRedisPath(), "--port", fmt.Sprintf("%d", redisPort))
			redisSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			// give redis a chance to start
			Eventually(redisSession.Out, "5s").Should(gbytes.Say("Ready to accept connections"))

			// add the rl service as a static upstream
			rlserver := &gloov1.Upstream{
				Metadata: &core.Metadata{
					Name:      "rl-server",
					Namespace: "default",
				},
				UseHttp2: &wrappers.BoolValue{Value: true},
				UpstreamType: &gloov1.Upstream_Static{
					Static: &gloov1static.UpstreamSpec{
						Hosts: []*gloov1static.Host{{
							Addr: envoyInstance.LocalAddr(),
							Port: uint32(rlServerSettings.RateLimitPort),
						}},
					},
				},
			}

			_, err = testClients.UpstreamClient.Write(rlserver, clients.WriteOpts{})
			Expect(err).ToNot(HaveOccurred())
			ref := rlserver.Metadata.Ref()
			rlSettings := &ratelimit.Settings{
				RatelimitServerRef: ref,
				DenyOnFail:         true, // ensures ConsistentlyNotRateLimited() calls will not pass unless server is healthy
			}

			glooSettings.RatelimitServer = rlSettings

			what := services.What{
				DisableGateway: true,
				DisableUds:     true,
				DisableFds:     true,
			}
			services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{Ctx: ctx, Cache: cache, LocalGlooPort: int32(testClients.GlooPort), What: what, Namespace: defaults.GlooSystem, Settings: glooSettings})

		})

		AfterEach(func() {
			redisSession.Terminate().Wait("5s")
			cancel()
		})
		It("should honor rate limit rules with a subset of the SetActions", func() {

			proxy := proxyBuilderWithManyHosts(newRateLimitingProxyBuilder(envoyPort, testUpstream.Upstream.Metadata.Ref())).build()

			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			isServerHealthy := ratelimitservice.RunRateLimitServer(ctx, rateLimitAddr, testClients.GlooPort, rlServerSettings)
			Eventually(isServerHealthy, "10s").Should(BeTrue())
			gloohelpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
			}, "90s", "3s")
			EventuallyRateLimited("host1", envoyPort)
		})

	})

})

func getLongActionList(actionValue string) []*rlv1alpha1.Action {
	// add extra actions to increase the size of the rate limit snapshot and demonstrate that ratelimiting works with a
	// subset of the descriptors
	numActions := 2
	actions := make([]*rlv1alpha1.Action, numActions)
	for i := 0; i < numActions; i++ {
		actions[i] = &rlv1alpha1.Action{
			ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
				GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: fmt.Sprintf("%s%d", actionValue, i)},
			},
		}
	}
	return actions

}
func proxyBuilderWithManyHosts(builder *rateLimitingProxyBuilder) *rateLimitingProxyBuilder {
	for i := 0; i < numRateLimitConfigs; i++ {
		actions := []*rlv1alpha1.Action{{
			ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
				GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: fmt.Sprintf("foo%d", i)},
			}},
			{ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
				GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: fmt.Sprintf("bar%d", i)},
			}},
			{ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
				GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: fmt.Sprintf("baz%d", i)},
			}},
		}
		actions = append(actions, getLongActionList(fmt.Sprintf("uniqueActionName%d", i))...)
		rateLimits := []*rlv1alpha1.RateLimitActions{{
			SetActions: actions,
		}}
		builder = builder.withVirtualHost(fmt.Sprintf("host%d", i), virtualHostConfig{rateLimitConfig: rateLimits})
	}
	return builder
}
func getManySetDescriptors(totalSets int) []*rlv1alpha1.SetDescriptor {

	setDescriptors := make([]*rlv1alpha1.SetDescriptor, totalSets)
	for i := 0; i < totalSets; i++ {
		descriptors := []*rlv1alpha1.SimpleDescriptor{
			{
				Key:   "generic_key",
				Value: fmt.Sprintf("foo%d", i),
			},
			{
				Key:   "generic_key",
				Value: fmt.Sprintf("bar%d", i),
			},
		}
		setDescriptors[i] = &rlv1alpha1.SetDescriptor{

			SimpleDescriptors: descriptors,
			RateLimit: &rlv1alpha1.RateLimit{
				Unit:            rlv1alpha1.RateLimit_MINUTE,
				RequestsPerUnit: 2,
			},
		}
	}
	return setDescriptors
}
func getApiKeyAuthConfigLarge() *extauth.ApiKeyAuth {
	return &extauth.ApiKeyAuth{
		ApiKeySecretRefs: []*core.ResourceRef{
			{
				Namespace: "default",
				Name:      "secret1",
			},
			{
				Namespace: "default",
				Name:      "secret3",
			},
		},
		LabelSelector: map[string]string{"team": "infrastructure"},
	}
}
