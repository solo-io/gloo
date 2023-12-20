package e2e_test

import (
	"testing"

	"github.com/solo-io/gloo/test/services/envoy"
	glooe_envoy "github.com/solo-io/solo-projects/test/services/envoy"
	"github.com/solo-io/solo-projects/test/services/tap_server"

	"github.com/solo-io/solo-projects/test/e2e"

	. "github.com/onsi/ginkgo/v2"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/test/helpers"
)

var (
	envoyFactory       envoy.Factory
	tapServerFactory   *tap_server.Factory
	testContextFactory *e2e.TestContextFactory

	namespace = defaults.GlooSystem
)

var _ = BeforeSuite(func() {
	envoyFactory = glooe_envoy.NewFactory()
	tapServerFactory = tap_server.NewFactory()

	testContextFactory = &e2e.TestContextFactory{
		EnvoyFactory:     envoyFactory,
		TapServerFactory: tapServerFactory,
	}
})

var _ = AfterSuite(func() {
	envoyFactory.Clean()
})

// NOTE: Please read the README.md for these tests
func TestE2e(t *testing.T) {

	// set default port to an unprivileged port for local testing.
	// 8081 is used by validation. see here:
	// test/services/gateway.go:233
	defaults.HttpPort = 8083

	// NOTE: if tests fails and spits out logs continuously, please read the README.md for these tests.
	helpers.RegisterCommonFailHandlers()
	helpers.SetupLog()
	// RegisterFailHandler(Fail)

	RunSpecs(t, "E2e Suite")
}
