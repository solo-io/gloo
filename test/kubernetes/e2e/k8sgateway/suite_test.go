package k8sgateway_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/route_options"
	"github.com/solo-io/gloo/test/kubernetes/testutils/cluster"
	"github.com/solo-io/gloo/test/kubernetes/testutils/gloogateway"
	"github.com/solo-io/gloo/test/kubernetes/testutils/runtime"
	"github.com/solo-io/skv2/codegen/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)

	ctx := context.Background()
	var testInstallation *e2e.TestInstallation
	r := require.New(t)

	runtimeContext := runtime.NewContext()

	// Construct the cluster.Context for this suite
	clusterContext := cluster.MustKindContext(runtimeContext.ClusterName)

	testCluster = &e2e.TestCluster{
		RuntimeContext: runtimeContext,
		ClusterContext: clusterContext,
	}

	t.Run("before", func(t *testing.T) {
		fmt.Println("before suite")

		testInstallation = testCluster.RegisterTestInstallation(
			&gloogateway.Context{
				InstallNamespace:   "k8s-gw-deployer-test",
				ValuesManifestFile: filepath.Join(util.MustGetThisDir(), "manifests", "k8s-gateway-test-helm.yaml"),
			},
		)

		err := testInstallation.InstallGlooGateway(ctx, testInstallation.Actions.Glooctl().NewTestHelperInstallAction())
		r.NoError(err)
	})

	t.Run("route_options", func(t *testing.T) {
		suite.Run(t, route_options.NewExample(ctx, testInstallation))
	})

	t.Run("after", func(t *testing.T) {
		fmt.Println("after suite")
		err := testInstallation.UninstallGlooGateway(ctx, testInstallation.Actions.Glooctl().NewTestHelperUninstallAction())
		r.NoError(err)

		testCluster.UnregisterTestInstallation(testInstallation)
	})
}
