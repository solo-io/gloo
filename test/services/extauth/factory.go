package extauth

import (
	"sync/atomic"

	. "github.com/onsi/gomega"
	"github.com/solo-io/ext-auth-service/pkg/service"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"

	extauthserver "github.com/solo-io/ext-auth-service/pkg/server"
	"github.com/solo-io/gloo/test/ginkgo/parallel"
)

const (
	// basePort is the starting port for the ext auth server
	// This was the previous static port used in tests, but it is not a special value
	basePort = uint32(9100)
)

type Factory struct {
	basePort uint32
}

func NewFactory() *Factory {
	return &Factory{
		basePort: basePort,
	}
}

func (f *Factory) NewInstance(address string) *Instance {
	serverSettings, err := extauthserver.NewSettings()
	Expect(err).NotTo(HaveOccurred(), "failed to create extauth server settings")

	serverSettings.DebugPort = 0
	serverSettings.ServerPort = int(advancePort(&f.basePort))
	serverSettings.SigningKey = "hello"
	serverSettings.UserIdHeader = extauth.DefaultAuthHeader
	serverSettings.HealthCheckHttpPath = "/healthcheck"
	serverSettings.HealthCheckHttpPort = int(advancePort(&f.basePort))
	// The number of seconds that the server will remain alive, but actively failing health checks
	// After this time elapses, the server will exit
	// This is useful for testing that the server will exit when it fails health checks
	// This is useful in production to ensure that the server will handle in flight requests before exiting
	serverSettings.HealthCheckFailTimeout = 0 // seconds
	// These settings are required for the server to add the userID to the dynamic metadata
	serverSettings.MetadataSettings = service.DynamicMetadataSettings{
		Enabled:   true,
		UserIdKey: "authUserId",
	}
	serverSettings.LogSettings = extauthserver.LogSettings{
		LoggerName: "ext-auth-service",
		DebugMode:  "true",
	}

	return &Instance{
		address:        address,
		serverSettings: serverSettings,
	}
}

func advancePort(p *uint32) uint32 {
	return atomic.AddUint32(p, 1) + uint32(parallel.GetPortOffset())
}
