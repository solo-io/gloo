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
)

// TestK8sGatewayIstioRevision is the function which executes a series of tests against a given installation with
// k8s gateway enabled and Istio installed with revisions
func TestK8sGatewayIstioRevision(t *testing.T) {
	ctx := context.Background()
	installNs, nsEnvPredefined := envutils.LookupOrDefault(testutils.InstallNamespace, "istio-rev-k8s-gw-test")
	testInstallation := e2e.CreateTestInstallation(
		t,
		&gloogateway.Context{
			InstallNamespace:          installNs,
			ProfileValuesManifestFile: e2e.KubernetesGatewayProfilePath,
			ValuesManifestFile:        e2e.ManifestPath("istio-revision-k8s-gateway-helm.yaml"),
		},
	)

	testHelper := e2e.MustTestHelper(ctx, testInstallation)
	err := testInstallation.AddIstioctl(ctx)
	if err != nil {
		t.Fatalf("failed to get istioctl: %v", err)
	}

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

			// Generate istioctl bug report
			testInstallation.CreateIstioBugReport(ctx)
		}

		testInstallation.UninstallGlooGatewayWithTestHelper(ctx, testHelper)

		// Uninstall Istio
		err = testInstallation.UninstallIstio()
		if err != nil {
			t.Fatalf("failed to uninstall: %v\n", err)
		}
	})

	// Install Istio before Gloo Gateway to make sure istiod is present before istio-proxy
	err = testInstallation.InstallRevisionedIstio(ctx, "test-revision", "minimal")
	if err != nil {
		t.Fatalf("failed to install: %v", err)
	}

	// Install Gloo Gateway
	testInstallation.InstallGlooGatewayWithTestHelper(ctx, testHelper, 5*time.Minute)

	RevisionIstioK8sGatewaySuiteRunner().Run(ctx, t, testInstallation)
}
