package e2e_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als"

	"github.com/fgrosse/zaptest"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
)

var _ = Describe("Leftmost x-forwarded-for address Local E2E", func() {

	var (
		ctx           context.Context
		cancel        context.CancelFunc
		testClients   services.TestClients
		envoyInstance *services.EnvoyInstance
		envoyPort     uint32
	)

	BeforeEach(func() {

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		what := services.What{
			DisableFds: true,
			DisableUds: true,
		}

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, "gloo-system", nil, nil, &gloov1.Settings{})

	})

	setupProxy := func(leftmostXffAddress bool) {
		var err error
		envoyInstance, err = envoyFactory.NewEnvoyInstance()
		Expect(err).NotTo(HaveOccurred())
		err = envoyInstance.Run(testClients.GlooPort)
		Expect(err).NotTo(HaveOccurred())
		path := "/dev/stdout"

		if !envoyInstance.UseDocker() {
			tmpfile, err := ioutil.TempFile("", "")
			Expect(err).NotTo(HaveOccurred())
			path = tmpfile.Name()
			envoyInstance.AccessLogs = path
		}

		envoyPort = defaults.HttpPort
		proxy := getProxy(envoyPort, leftmostXffAddress, path)

		_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", envoyInstance.AdminPort), nil)
		Expect(err).NotTo(HaveOccurred())
		client := &http.Client{}
		Eventually(func() (int, error) {
			response, err := client.Do(request)
			if response == nil {
				return 0, err
			}
			return response.StatusCode, err
		}, 20*time.Second, 1*time.Second).Should(Equal(200))

	}

	AfterEach(func() {
		cancel()
		if envoyInstance != nil {
			_ = envoyInstance.Clean()
		}
	})

	Context("With envoy", func() {

		It("sets leftmost xff address as downstream remote address", func() {
			setupProxy(true)

			request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/", envoyPort), nil)
			Expect(err).NotTo(HaveOccurred())
			address := "192.168.2.1"
			request.Header.Add("x-forwarded-for", address+",192.123.3.1,192.123.3.2")
			Eventually(func() int {
				response, _ := http.DefaultClient.Do(request)
				if response == nil {
					return 0
				}
				return response.StatusCode
			}, 3*time.Second, 1*time.Second).Should(Equal(200))

			// check access logs to make sure downstream remote ip is set to left most x-forwarded-for address
			err = checkAccessLogs(envoyInstance, func(logs string) bool {
				return strings.Contains(logs, fmt.Sprintf("[%s]", address))
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("sanity check regular x-forwarded-for behaviour", func() {
			setupProxy(false)

			request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/", envoyPort), nil)
			Expect(err).NotTo(HaveOccurred())
			address := "192.168.2.1"
			request.Header.Add("x-forwarded-for", "192.123.3.1,192.123.3.2,"+address)
			Eventually(func() int {
				response, _ := http.DefaultClient.Do(request)
				if response == nil {
					return 0
				}
				return response.StatusCode
			}, 3*time.Second, 1*time.Second).Should(Equal(200))

			// check access logs to make sure downstream remote ip is set to right most x-forwarded-for address (default behaviour)
			err = checkAccessLogs(envoyInstance, func(logs string) bool {
				return strings.Contains(logs, fmt.Sprintf("[%s]", address))
			})
			Expect(err).NotTo(HaveOccurred())
		})

	})

})

func getProxy(envoyPort uint32, leftmostXffAddress bool, accessLogPath string) *gloov1.Proxy {

	vhosts := []*gloov1.VirtualHost{
		{
			Name:    "virt1",
			Domains: []string{"*"},
			Routes: []*gloov1.Route{{
				Action: &gloov1.Route_DirectResponseAction{
					DirectResponseAction: &gloov1.DirectResponseAction{
						Status: uint32(200),
					},
				}}},
		},
	}

	p := &gloov1.Proxy{
		Metadata: &core.Metadata{
			Name:      "proxy",
			Namespace: "default",
		},
		Listeners: []*gloov1.Listener{{
			Name:        "listener",
			BindAddress: "0.0.0.0",
			BindPort:    envoyPort,
			Options: &gloov1.ListenerOptions{
				AccessLoggingService: &als.AccessLoggingService{
					AccessLog: []*als.AccessLog{
						{
							OutputDestination: &als.AccessLog_FileSink{
								FileSink: &als.FileSink{
									Path: accessLogPath,
									OutputFormat: &als.FileSink_StringFormat{
										StringFormat: "[%DOWNSTREAM_REMOTE_ADDRESS_WITHOUT_PORT%]",
									},
								},
							},
						},
					},
				},
			},
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: vhosts,
					Options: &gloov1.HttpListenerOptions{
						LeftmostXffAddress: &wrappers.BoolValue{
							Value: leftmostXffAddress,
						},
					},
				},
			},
		}},
	}

	return p
}

func checkAccessLogs(ei *services.EnvoyInstance, logsPresent func(logs string) bool) error {
	var (
		logs string
		err  error
	)

	if ei.UseDocker() {
		logs, err = ei.Logs()
		if err != nil {
			return err
		}
	} else {
		file, err := os.OpenFile(ei.AccessLogs, os.O_RDONLY, 0777)
		if err != nil {
			return err
		}
		var byt []byte
		byt, err = ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		logs = string(byt)
	}

	if logs == "" {
		return fmt.Errorf("logs should not be empty")
	}
	if !logsPresent(logs) {
		return fmt.Errorf("no access logs present")
	}
	return nil
}
