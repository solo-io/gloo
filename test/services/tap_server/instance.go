package tap_server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/tap-extension-examples/pkg/data_scrubber"
	"github.com/solo-io/tap-extension-examples/pkg/tap_server_builder"

	// TODO this alias is fixing the fact that the tap-extension-examples
	// directory needs to be fixed due to the fact that the package names in
	// the `tap-extension-examples` repository are not quite right
	tap_service "github.com/solo-io/tap-extension-examples/pkg/tap_service"
)

type InstanceConfig struct {
	// Whether the tap server should scrub credit card numbers, social security
	// numbers, etc. from the messages received
	EnableDataScrubbing bool
}

type Instance struct {
	address      string
	HostPort     int
	DataScrubber *data_scrubber.DataScrubber
	TapRequests  chan tap_service.TapRequest
	TapServer    tap_server_builder.TapServer
}

func (i *Instance) Address() string {
	return i.address
}

func (i *Instance) Port() uint32 {
	return uint32(i.HostPort)
}

func (i *Instance) Url() string {
	return fmt.Sprintf("%s:%d", i.Address(), i.Port())
}

func (i *Instance) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		i.Clean()
	}()

	ExpectWithOffset(1, i.TapRequests).NotTo(BeNil())
	i.TapServer = tap_server_builder.
		NewTapServerBuilder().
		WithDataScrubber(i.DataScrubber).
		WithTapMessageChannel(i.TapRequests).
		BuildHttp()
	listenAddress := fmt.Sprintf(":%d", i.Port())
	go i.TapServer.Run(listenAddress)
}

func (i *Instance) Clean() {
	if i == nil {
		return
	}
	i.TapServer.Stop()
}

func (i *Instance) EventuallyIsHealthy() {
	EventuallyWithOffset(1, func(g Gomega) {
		g.Expect(i.IsHealthy()).To(BeTrue())
	}, "5s", ".1s").Should(Succeed())
}

func (i *Instance) IsHealthy() bool {
	// check if i.HostPort is open
	log.Printf("checking to see if tap server is available on port %d\n", i.Port())
	conn, err := net.DialTimeout("tcp", i.Url(), time.Second*5)
	if conn != nil {
		conn.Close()
	}
	return err == nil
}

func (i *Instance) Logs() *tap_service.TapRequest {
	select {
	case request := <-i.TapRequests:
		return &request
	default:
		return nil
	}
}

func (i *Instance) GetServerUpstream() *gloov1.Upstream {
	return &gloov1.Upstream{
		Metadata: &core.Metadata{
			Namespace: defaults.GlooSystem,
			Name:      "tap-server",
		},
		UseHttp2: &wrappers.BoolValue{Value: false},
		UpstreamType: &gloov1.Upstream_Static{
			Static: &gloov1static.UpstreamSpec{
				Hosts: []*gloov1static.Host{{
					Addr: i.Address(),
					Port: i.Port(),
				}},
			},
		},
	}
}
