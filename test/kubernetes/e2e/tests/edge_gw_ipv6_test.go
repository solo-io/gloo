package tests_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/solo-io/gloo/pkg/utils/envutils"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	. "github.com/solo-io/gloo/test/kubernetes/e2e/tests"
	"github.com/solo-io/gloo/test/kubernetes/testutils/gloogateway"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/go-utils/log"
)

// TestGlooGatewayEdgeIpv6Gateway is the function which executes a series of tests against a given installation where
// the k8s Gateway controller is disabled
func TestGlooGatewayEdgeIpv6Gateway(t *testing.T) {
	if os.Getenv("ENABLE_IPV6_ONLY") != "true" {
		log.Warnf("ENABLE_IPV6_ONLY not set to true, skipping")
		return
	}
	ctx := context.Background()
	installNs, nsEnvPredefined := envutils.LookupOrDefault(testutils.InstallNamespace, "gloo-gateway-edge-ipv6-test")
	testInstallation := e2e.CreateTestInstallation(
		t,
		&gloogateway.Context{
			InstallNamespace:          installNs,
			ProfileValuesManifestFile: e2e.EdgeGatewayProfilePath,
			ValuesManifestFile:        e2e.ManifestPath("edge-gateway-ipv6-test-helm.yaml"),
		},
	)

	testHelper := e2e.MustTestHelper(ctx, testInstallation)
	// Set the env to the install namespace if it is not already set
	if !nsEnvPredefined {
		os.Setenv(testutils.InstallNamespace, installNs)
	}

	// We register the cleanup function _before_ we actually perform the installation.
	// This allows us to uninstall Gloo Gateway, in case the original installation only completed partially
	t.Cleanup(func() {
		if !nsEnvPredefined {
			os.Unsetenv(testutils.InstallNamespace)
		}
		if t.Failed() {
			testInstallation.PreFailHandler(ctx)
		}

		testInstallation.UninstallGlooGatewayWithTestHelper(ctx, testHelper)
	})

	// Install Gloo Gateway with only Gloo Edge Gateway APIs enabled
	testInstallation.InstallGlooGatewayWithTestHelper(ctx, testHelper, 5*time.Minute)

	EdgeGwIpv6SuiteRunner().Run(ctx, t, testInstallation)
}
