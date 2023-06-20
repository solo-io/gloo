package extauth

import (
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/ext-auth-service/pkg/server"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	extauthrunner "github.com/solo-io/solo-projects/projects/extauth/pkg/runner"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Instance struct {
	address string

	serverSettings server.Settings
}

func (i *Instance) GetServerSettings() server.Settings {
	return i.serverSettings
}

func (i *Instance) UpdateServerSettings(mutator func(*server.Settings)) {
	mutator(&i.serverSettings)
}

func (i *Instance) RunWithXds(ctx context.Context, xdsPort uint32) {
	settings := extauthrunner.Settings{
		GlooAddress: fmt.Sprintf("%s:%d", i.Address(), xdsPort),

		ExtAuthSettings: i.serverSettings,
	}

	err := extauthrunner.RunWithSettings(ctx, settings)
	if ctx.Err() == nil {
		Expect(err).NotTo(HaveOccurred())
	}
}

func (i *Instance) Address() string {
	return i.address
}

func (i *Instance) Port() uint32 {
	return uint32(i.serverSettings.ServerPort)
}

func (i *Instance) Url() string {
	return fmt.Sprintf("%s:%d", i.Address(), i.Port())
}

func (i *Instance) GetHealthCheckResponse(opts ...grpc.CallOption) (*healthpb.HealthCheckResponse, error) {
	conn, err := grpc.Dial(i.Url(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	return healthpb.NewHealthClient(conn).Check(context.Background(), &healthpb.HealthCheckRequest{
		Service: i.serverSettings.ServiceName,
	}, opts...)
}

func (i *Instance) IsHealthy() bool {
	response, err := i.GetHealthCheckResponse()
	if err != nil {
		return false
	}
	if response.Status != healthpb.HealthCheckResponse_SERVING {
		return false
	}
	return true
}

func (i *Instance) EventuallyIsHealthy() {
	EventuallyWithOffset(1, func(g Gomega) {
		g.Expect(i.IsHealthy()).To(BeTrue())
	}, "5s", ".1s").Should(Succeed())
}

func (i *Instance) GetServerUpstream() *gloov1.Upstream {
	return &gloov1.Upstream{
		Metadata: &core.Metadata{
			Name:      "extauth-server",
			Namespace: defaults.GlooSystem,
		},
		UseHttp2: &wrappers.BoolValue{Value: true},
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
