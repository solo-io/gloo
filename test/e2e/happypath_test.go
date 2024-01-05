package e2e_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	testhelpers "github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/services/envoy"
	"google.golang.org/protobuf/types/known/wrapperspb"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gatewaydefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/healthcheck"
	routerV1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/router"
	static_plugin_gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/stats"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/services"
	"github.com/solo-io/gloo/test/v1helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
)

var _ = Describe("Happy path", func() {

	var (
		ctx           context.Context
		cancel        context.CancelFunc
		envoyInstance *envoy.Instance
		tu            *v1helpers.TestUpstream
		envoyPort     uint32
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		envoyInstance = envoyFactory.NewInstance()

		tu = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())
		envoyPort = envoyInstance.HttpPort
	})

	AfterEach(func() {
		envoyInstance.Clean()
		cancel()
	})

	Describe("in memory", func() {

		Describe("Rest Eds Enabled", func() {

			var (
				args                TestArgs
				transportApiVersion envoy_config_core_v3.ApiVersion
				restEdsEnabled      bool
			)

			BeforeEach(func() {
				transportApiVersion = envoy_config_core_v3.ApiVersion_V3
				restEdsEnabled = true
				args = setupBeforeEach(ctx, restEdsEnabled, transportApiVersion, tu, envoyInstance, envoyPort)
			})

			It("should not crash", func() {
				shouldNotCrash(args)
			})
			It("should not crash with multiple methods", func() {
				shouldNotCrashMultipleMethods(args)
			})
			It("correctly configures envoy to emit virtual cluster statistics", func() {
				correctlyConfiguresEnvoyToEmitVirtualClusterStatistics(args)
			})
			It("correctly passes suppressEnvoyHeaders config", func() {
				correctlyPassesSuppressEnvoyHeadersConfig(args)
			})
			It("correctly did not pass suppressEnvoyHeaders config", func() {
				correctlyDidNotPassSuppressEnvoyHeadersConfig(args)
			})
			It("passes health check", func() {
				passesHealthCheck(args)
			})

			Context("ssl", func() {
				var (
					sslport1 uint32
					sslport2 uint32
					err      error
				)

				BeforeEach(func() {
					err, args = sslSetupBeforeEach(ctx, tu, sslport1, sslport2, args)
					Expect(err).NotTo(HaveOccurred())

					// we need to overwrite these ports so that we don't reuse the same port in subsequent tests
					sslport1 = args.Ssl.sslport1
					sslport2 = args.Ssl.sslport2
				})

				Context("simple ssl", func() {
					BeforeEach(func() {
						writeUpstream(args)
					})

					It("should work with ssl", func() {
						simpleSslShouldWorkWithSsl(args)
					})
				})

				Context("sni", func() {
					BeforeEach(func() {
						setHosts(args)
						writeUpstream(args)
					})

					It("should work with ssl", func() {
						sniShouldWorkWithSsl(args)
					})
				})
			})

			Context("sad path", func() {
				It("should error the proxy with two listeners with the same bind address", func() {
					sadPathShouldError(args)
				})
			})

		})

		Describe("Rest Eds Disabled", func() {
			var (
				args                TestArgs
				transportApiVersion envoy_config_core_v3.ApiVersion
				restEdsEnabled      bool
			)

			BeforeEach(func() {
				transportApiVersion = envoy_config_core_v3.ApiVersion_V3
				restEdsEnabled = false
				args = setupBeforeEach(ctx, restEdsEnabled, transportApiVersion, tu, envoyInstance, envoyPort)
			})

			It("should not crash", func() {
				shouldNotCrash(args)
			})
			It("should not crash with multiple methods", func() {
				shouldNotCrashMultipleMethods(args)
			})
			It("correctly configures envoy to emit virtual cluster statistics", func() {
				correctlyConfiguresEnvoyToEmitVirtualClusterStatistics(args)
			})
			It("correctly passes suppressEnvoyHeaders config", func() {
				correctlyPassesSuppressEnvoyHeadersConfig(args)
			})
			It("correctly did not pass suppressEnvoyHeaders config", func() {
				correctlyDidNotPassSuppressEnvoyHeadersConfig(args)
			})
			It("passes health check", func() {
				passesHealthCheck(args)
			})

			Context("ssl", func() {
				var (
					sslport1 uint32
					sslport2 uint32
					err      error
				)

				BeforeEach(func() {
					err, args = sslSetupBeforeEach(ctx, tu, sslport1, sslport2, args)
					Expect(err).NotTo(HaveOccurred())

					// we need to overwrite these ports so that we don't reuse the same port in subsequent tests
					sslport1 = args.Ssl.sslport1
					sslport2 = args.Ssl.sslport2
				})

				Context("simple ssl", func() {
					BeforeEach(func() {
						writeUpstream(args)
					})

					It("should work with ssl", func() {
						simpleSslShouldWorkWithSsl(args)
					})
				})

				Context("sni", func() {
					BeforeEach(func() {
						setHosts(args)
						writeUpstream(args)
					})

					It("should work with ssl", func() {
						sniShouldWorkWithSsl(args)
					})
				})
			})

			Context("sad path", func() {
				It("should error the proxy with two listeners with the same bind address", func() {
					sadPathShouldError(args)
				})
			})
		})
	})
})

type TestCase struct {
	Title               string
	RestEdsEnabled      *wrappers.BoolValue
	TransportApiVersion envoy_config_core_v3.ApiVersion
}

func setupBeforeEach(ctx context.Context, restEdsEnabled bool, transportApiVersion envoy_config_core_v3.ApiVersion, tu *v1helpers.TestUpstream, envoyInstance *envoy.Instance, envoyPort uint32) TestArgs {
	ns := defaults.GlooSystem
	ro := &services.RunOptions{
		NsToWrite: ns,
		NsToWatch: []string{"default", ns},
		WhatToRun: services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		},
		Settings: &gloov1.Settings{
			Gloo: &gloov1.GlooOptions{
				EnableRestEds: &wrappers.BoolValue{Value: restEdsEnabled},
			},
		},
	}
	testClients := services.RunGlooGatewayUdsFds(ctx, ro)
	envoyInstance.ApiVersion = transportApiVersion.String()
	err := envoyInstance.RunWithRoleAndRestXds(ns+"~"+gatewaydefaults.GatewayProxyName, testClients.GlooPort, testClients.RestXdsPort)
	Expect(err).NotTo(HaveOccurred())

	up := tu.Upstream
	_, err = testClients.UpstreamClient.Write(up, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	args := TestArgs{
		TestUpstreamReachable: func() {
			v1helpers.TestUpstreamReachableWithOffset(3, envoyPort, tu, nil)
		},
		EnvoyPort:     envoyPort,
		Up:            up,
		TestClients:   testClients,
		Ctx:           ctx,
		EnvoyInstance: envoyInstance,
	}

	return args
}

// TODO: make these before each helpers more similarly named + give more similar signatures
func sslSetupBeforeEach(ctx context.Context, tu *v1helpers.TestUpstream, sslport1, sslport2 uint32, args TestArgs) (error, TestArgs) {
	hellos := make(chan sslConn, 100)
	sslSecret := &gloov1.Secret{
		Metadata: &core.Metadata{
			Name:      "secret",
			Namespace: "default",
		},
		Kind: &gloov1.Secret_Tls{
			Tls: &gloov1.TlsSecret{
				PrivateKey: testhelpers.PrivateKey(),
				CertChain:  testhelpers.Certificate(),
				RootCa:     testhelpers.Certificate(),
			},
		},
	}

	_, err := args.TestClients.SecretClient.Write(sslSecret, clients.WriteOpts{})
	if err != nil {
		return err, TestArgs{}
	}

	copyUp := *tu.Upstream
	copyUp.Metadata.Name = copyUp.Metadata.Name + "-ssl"
	port := tu.Upstream.UpstreamType.(*gloov1.Upstream_Static).Static.Hosts[0].Port
	addr := tu.Upstream.UpstreamType.(*gloov1.Upstream_Static).Static.Hosts[0].Addr

	sslport1 = v1helpers.StartSslProxyWithHelloCB(ctx, port, func(chi *tls.ClientHelloInfo) {
		hellos <- sslConn{sni: chi.ServerName, port: sslport1}
	})
	sslport2 = v1helpers.StartSslProxyWithHelloCB(ctx, port, func(chi *tls.ClientHelloInfo) {
		hellos <- sslConn{sni: chi.ServerName, port: sslport2}
	})

	ref := sslSecret.Metadata.Ref()

	copyUp.UpstreamType = &gloov1.Upstream_Static{
		Static: &static_plugin_gloo.UpstreamSpec{
			Hosts: []*static_plugin_gloo.Host{{
				Addr: addr,
				Port: sslport1,
			}, {
				Addr: addr,
				Port: sslport2,
			}},
		},
	}
	copyUp.SslConfig = &ssl.UpstreamSslConfig{
		SslSecrets: &ssl.UpstreamSslConfig_SecretRef{
			SecretRef: ref,
		},
	}

	args.Ssl = SslArgs{
		upSsl:    &copyUp,
		sslport1: sslport1,
		sslport2: sslport2,
		hellos:   hellos,
	}

	return nil, args
}

/**************************************/
/* 		    Test Content    		  */
/**************************************/

type sslConn struct {
	sni  string
	port uint32
}
type SslArgs struct {
	upSsl    *gloov1.Upstream
	hellos   chan sslConn
	sslport1 uint32
	sslport2 uint32
}
type TestArgs struct {
	TestUpstreamReachable func()
	EnvoyPort             uint32
	Up                    *gloov1.Upstream
	TestClients           services.TestClients
	Ctx                   context.Context
	EnvoyInstance         *envoy.Instance
	Ssl                   SslArgs
}

func shouldNotCrash(args TestArgs) {
	proxy := getTrivialProxyForUpstream(defaults.GlooSystem, args.EnvoyPort, args.Up.Metadata.Ref())
	_, err := args.TestClients.ProxyClient.Write(proxy, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	args.TestUpstreamReachable()
}

func shouldNotCrashMultipleMethods(args TestArgs) {
	proxy := getTrivialProxyForUpstream(defaults.GlooSystem, args.EnvoyPort, args.Up.Metadata.Ref())
	proxy.Listeners[0].ListenerType.(*gloov1.Listener_HttpListener).HttpListener.
		VirtualHosts[0].Routes[0].Matchers = []*matchers.Matcher{
		{
			PathSpecifier: &matchers.Matcher_Prefix{Prefix: "/"},
			Methods:       []string{"GET", "POST"},
		},
	}

	_, err := args.TestClients.ProxyClient.Write(proxy, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	args.TestUpstreamReachable()
}

func correctlyConfiguresEnvoyToEmitVirtualClusterStatistics(args TestArgs) {
	proxy := getTrivialProxyForUpstream(defaults.GlooSystem, args.EnvoyPort, args.Up.Metadata.Ref())

	// Set a virtual cluster matching everything
	proxy.Listeners[0].GetHttpListener().VirtualHosts[0].Options = &gloov1.VirtualHostOptions{
		Stats: &stats.Stats{
			VirtualClusters: []*stats.VirtualCluster{{
				Name:    "test-vc",
				Pattern: ".*",
			}},
		},
	}

	_, err := args.TestClients.ProxyClient.Write(proxy, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	// This will hit the virtual host with the above virtual cluster config
	args.TestUpstreamReachable()

	response, err := http.Get(fmt.Sprintf("http://localhost:%d/stats", args.EnvoyInstance.AdminPort))
	Expect(err).NotTo(HaveOccurred())
	Expect(response).NotTo(BeNil())
	//goland:noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	Expect(err).NotTo(HaveOccurred())
	statsString := string(body)

	// Verify that stats for the above virtual cluster are present
	Expect(statsString).To(ContainSubstring("vhost.virt1.vcluster.test-vc."))
}

func correctlyPassesSuppressEnvoyHeadersConfig(args TestArgs) {
	proxy := getTrivialProxyForUpstream(defaults.GlooSystem, args.EnvoyPort, args.Up.Metadata.Ref())

	// configuring an http listener option to set suppressEnvoyHeaders to true
	//projects/gloo/api/v1/options/router/router.proto
	proxy.Listeners[0].GetHttpListener().Options = &gloov1.HttpListenerOptions{
		Router: &routerV1.Router{
			SuppressEnvoyHeaders: wrapperspb.Bool(true),
		},
	}

	_, err := args.TestClients.ProxyClient.Write(proxy, clients.WriteOpts{
		Ctx: args.Ctx,
	})
	Expect(err).NotTo(HaveOccurred())

	args.TestUpstreamReachable()

	// This will hit the virtual host with the above virtual cluster config
	response, err := http.Get(fmt.Sprintf("http://%s:%d/", "localhost", args.EnvoyInstance.HttpPort))
	Expect(err).NotTo(HaveOccurred())
	Expect(response.Header).NotTo(HaveKey("X-Envoy-Upstream-Service-Time"))

	cfg, err := args.EnvoyInstance.ConfigDump()
	Expect(err).NotTo(HaveOccurred())

	// We expect the envoy configuration to contain these properties in the configuration dump
	Expect(cfg).To(MatchRegexp("\"suppress_envoy_headers\": true"))
}

func correctlyDidNotPassSuppressEnvoyHeadersConfig(args TestArgs) {
	proxy := getTrivialProxyForUpstream(defaults.GlooSystem, args.EnvoyPort, args.Up.Metadata.Ref())

	// Set a virtual cluster listener that is blank and has no options
	proxy.Listeners[0].GetHttpListener().Options = &gloov1.HttpListenerOptions{}

	_, err := args.TestClients.ProxyClient.Write(proxy, clients.WriteOpts{
		Ctx: args.Ctx,
	})
	Expect(err).NotTo(HaveOccurred())
	args.TestUpstreamReachable()

	// This will hit the virtual host with the above virtual cluster config
	response, err := http.Get(fmt.Sprintf("http://%s:%d/", "localhost", args.EnvoyInstance.HttpPort))
	Expect(err).NotTo(HaveOccurred())
	Expect(response.Header).To(HaveKey("X-Envoy-Upstream-Service-Time"))

	cfg, err := args.EnvoyInstance.ConfigDump()
	Expect(err).NotTo(HaveOccurred())

	// We expect the envoy configuration to NOT contain these properties in the configuration dump
	Expect(cfg).To(Not(MatchRegexp("\"suppress_envoy_headers\": true")))
}

func passesHealthCheck(args TestArgs) {
	proxy := getTrivialProxyForUpstream(defaults.GlooSystem, args.EnvoyPort, args.Up.Metadata.Ref())

	// Set a virtual cluster matching everything
	proxy.Listeners[0].GetHttpListener().Options = &gloov1.HttpListenerOptions{
		HealthCheck: &healthcheck.HealthCheck{
			Path: "/healthy",
		},
	}

	proxy.Listeners[0].GetHttpListener().VirtualHosts[0].Routes[0].Action = &gloov1.Route_DirectResponseAction{
		DirectResponseAction: &gloov1.DirectResponseAction{
			Status: 400,
			Body:   "only health checks work on me. sorry!",
		},
	}
	_, err := args.TestClients.ProxyClient.Write(proxy, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	Eventually(func() error {
		res, err := http.Get(fmt.Sprintf("http://%s:%d/healthy", "localhost", args.EnvoyPort))
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
			return errors.Errorf("bad status code: %v", res.StatusCode)
		}
		return nil
	}, time.Second*10, time.Second/2).ShouldNot(HaveOccurred())

	res, err := http.Post(fmt.Sprintf("http://localhost:%v/healthcheck/fail", args.EnvoyInstance.AdminPort), "", nil)
	Expect(err).NotTo(HaveOccurred())
	Expect(res.StatusCode).To(Equal(200))

	Eventually(func() error {
		res, err := http.Get(fmt.Sprintf("http://%s:%d/healthy", "localhost", args.EnvoyPort))
		if err != nil {
			return err
		}
		if res.StatusCode != 503 {
			return errors.Errorf("bad status code: %v", res.StatusCode)
		}
		return nil
	}, time.Second*10, time.Second/2).ShouldNot(HaveOccurred())
}

func simpleSslShouldWorkWithSsl(args TestArgs) {
	proxycli := args.TestClients.ProxyClient
	proxy := getTrivialProxyForUpstream(defaults.GlooSystem, args.EnvoyPort, args.Up.Metadata.Ref())
	_, err := proxycli.Write(proxy, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	args.TestUpstreamReachable()
}

func sniShouldWorkWithSsl(args TestArgs) {
	proxycli := args.TestClients.ProxyClient
	proxy := getTrivialProxyForUpstream(defaults.GlooSystem, args.EnvoyPort, args.Up.Metadata.Ref())
	_, err := proxycli.Write(proxy, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	match1 := sslConn{sni: "solo-sni-test", port: args.Ssl.sslport1}
	match2 := sslConn{sni: "solo-sni-test2", port: args.Ssl.sslport2}

	matched1 := false
	matched2 := false

	timeout := time.After(5 * time.Second)
	for {
		args.TestUpstreamReachable()
		select {
		case <-timeout:
			Fail("timedout waiting for sni")
		case clienthello := <-args.Ssl.hellos:
			Expect(clienthello).To(SatisfyAny(Equal(match1), Equal(match2)))
			if clienthello == match1 {
				matched1 = true
			} else {
				matched2 = true
			}
		}
		if matched1 && matched2 {
			break
		}
	}
}

func sadPathShouldError(args TestArgs) {
	// create a proxy with two identical listeners to see errors come up
	proxy := getTrivialProxyForUpstream(defaults.GlooSystem, args.EnvoyPort, args.Up.Metadata.Ref())
	proxy.Listeners = append(proxy.Listeners, proxy.Listeners[0])

	// persist the proxy
	_, err := args.TestClients.ProxyClient.Write(proxy, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())

	// eventually the proxy is rejected
	testhelpers.EventuallyResourceRejected(func() (resources.InputResource, error) {
		return args.TestClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
	})
}

/**************************************/
/*          Helper Functions          */
/**************************************/

func writeUpstream(args TestArgs) {
	_, err := args.TestClients.UpstreamClient.Write(args.Ssl.upSsl, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())
}

func setHosts(args TestArgs) {
	args.Ssl.upSsl.GetStatic().GetHosts()[0].SniAddr = "solo-sni-test"
	args.Ssl.upSsl.GetStatic().GetHosts()[1].SniAddr = "solo-sni-test2"
}

func getTrivialProxyForUpstream(ns string, bindPort uint32, upstream *core.ResourceRef) *gloov1.Proxy {
	proxy := getTrivialProxy(ns, bindPort)
	proxy.Listeners[0].ListenerType.(*gloov1.Listener_HttpListener).HttpListener.
		VirtualHosts[0].Routes[0].Action.(*gloov1.Route_RouteAction).RouteAction.
		Destination.(*gloov1.RouteAction_Single).Single.DestinationType =
		&gloov1.Destination_Upstream{Upstream: upstream}
	return proxy
}

func getTrivialProxyForService(ns string, bindPort uint32, service *core.ResourceRef, svcPort uint32) *gloov1.Proxy {
	proxy := getTrivialProxy(ns, bindPort)
	proxy.Listeners[0].ListenerType.(*gloov1.Listener_HttpListener).HttpListener.
		VirtualHosts[0].Routes[0].Action.(*gloov1.Route_RouteAction).RouteAction.
		Destination.(*gloov1.RouteAction_Single).Single.DestinationType =
		&gloov1.Destination_Kube{
			Kube: &gloov1.KubernetesServiceDestination{
				Ref:  service,
				Port: svcPort,
			},
		}
	return proxy
}

func getTrivialProxy(ns string, bindPort uint32) *gloov1.Proxy {
	return &gloov1.Proxy{
		Metadata: &core.Metadata{
			Name:      gatewaydefaults.GatewayProxyName,
			Namespace: ns,
		},
		Listeners: []*gloov1.Listener{{
			Name:        "listener",
			BindAddress: "::",
			BindPort:    bindPort,
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: []*gloov1.VirtualHost{{
						Name:    "virt1",
						Domains: []string{"*"},
						Routes: []*gloov1.Route{{
							Action: &gloov1.Route_RouteAction{
								RouteAction: &gloov1.RouteAction{
									Destination: &gloov1.RouteAction_Single{
										Single: &gloov1.Destination{},
									},
								},
							},
						}},
					}},
				},
			},
		}},
	}
}

// getNonSpecialIP returns a non-special IP that Kubernetes will allow in an endpoint.
func getNonSpecialIP(instance *envoy.Instance) string {
	if instance.UseDocker {
		return instance.LocalAddr()
	}

	ifaces, err := net.Interfaces()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}
			if isNonSpecialIP(ip) {
				return ip.String()
			}
		}
	}
	Fail("no ip address available", 1)
	return ""
}

// isNonSpecialIP is adapted from ValidateNonSpecialIP in k8s.io/kubernetes/pkg/apis/core/validation/validation.go
//
// Specifically disallowed are unspecified, loopback addresses, and link-local addresses
// which tend to be used for node-centric purposes (e.g. metadata service).
func isNonSpecialIP(ip net.IP) bool {
	if ip == nil {
		return false // must be a valid IP address
	}
	if ip.IsUnspecified() {
		return false // may not be unspecified
	}
	if ip.IsLoopback() {
		return false // may not be in the loopback range (127.0.0.0/8, ::1/128)
	}
	if ip.IsLinkLocalUnicast() {
		return false // may not be in the link-local range (169.254.0.0/16, fe80::/10)
	}
	if ip.IsLinkLocalMulticast() {
		return false // may not be in the link-local multicast range (224.0.0.0/24, ff02::/10)
	}
	return true
}
