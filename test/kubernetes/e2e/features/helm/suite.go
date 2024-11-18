package helm

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/solo-io/gloo/test/kubernetes/testutils/helper"
	"github.com/solo-io/skv2/codegen/util"
	"github.com/solo-io/solo-kit/pkg/code-generator/schemagen"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is the entire Suite of tests for the Helm Tests
type testingSuite struct {
	*base.BaseTestingSuite
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		base.NewBaseTestingSuite(ctx, testInst, e2e.MustTestHelper(ctx, testInst), base.SimpleTestCase{}, helmTestCases),
	}
}

func (s *testingSuite) TestProductionRecommendations() {
	envoyDeployment := s.GetKubectlOutput("-n", s.TestInstallation.Metadata.InstallNamespace, "get", "deployment", "gateway-proxy", "-o", "yaml")
	s.Contains(envoyDeployment, "readinessProbe:")
	s.Contains(envoyDeployment, "/envoy-hc")
	s.Contains(envoyDeployment, "readyReplicas: 1")
}

func (s *testingSuite) TestChangedConfigMapTriggersRollout() {
	expectConfigDumpToContain := func(str string) {
		adminCli, shutdown, err := admincli.NewPortForwardedClient(s.Ctx, "deployment/gateway-proxy", s.TestHelper.InstallNamespace)
		s.NoError(err)
		defer shutdown()

		var b bytes.Buffer
		dump := io.Writer(&b)
		err = adminCli.ConfigDumpCmd(s.Ctx, nil).WithStdout(dump).Run().Cause()
		s.NoError(err)

		s.Contains(b.String(), str)
	}

	getChecksum := func() string {
		return s.GetKubectlOutput("-n", s.TestInstallation.Metadata.InstallNamespace, "get", "deployment", "gateway-proxy", "-o", "jsonpath='{.spec.template.metadata.annotations.checksum/gateway-proxy-envoy-config}'")
	}

	// The default value is 250000
	expectConfigDumpToContain(`"global_downstream_max_connections": 250000`)
	oldChecksum := getChecksum()

	// A change in the config map should trigger a new deployment anyway
	s.UpgradeWithCustomValuesFile(configMapChangeSetup)

	// We upgrade Gloo with a new value of `globalDownstreamMaxConnections` on envoy
	// This should cause the checkup annotation on the deployment to change and therefore
	// the deployment should be updated with the new value
	expectConfigDumpToContain(`"global_downstream_max_connections": 12345`)
	newChecksum := getChecksum()
	s.NotEqual(oldChecksum, newChecksum)
}

func (s *testingSuite) TestApplyCRDs() {
	var crdsByFileName = map[string]v1.CustomResourceDefinition{}
	crdDir := filepath.Join(util.GetModuleRoot(), "install", "helm", s.TestHelper.HelmChartName, "crds")

	err := filepath.Walk(crdDir, func(crdFile string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || info.Name() == "README.md" {
			return nil
		}

		// Parse the file, and extract the CRD- will fail for any non-yaml files or files not containing a CRD
		crd, err := schemagen.GetCRDFromFile(crdFile)
		if err != nil {
			return eris.Wrap(err, "error getting CRD from "+crdFile)
		}
		crdsByFileName[crdFile] = crd

		// continue traversing
		return nil
	})
	s.NoError(err)

	for crdFile, crd := range crdsByFileName {
		// Apply the CRD
		err := s.TestHelper.ApplyFile(s.Ctx, crdFile)
		s.NoError(err)

		// Ensure the CRD is eventually accepted
		out, _, err := s.TestHelper.Execute(s.Ctx, "get", "crd", crd.GetName())
		s.NoError(err)
		s.Contains(out, crd.GetName())

		// Ensure the CRD has the gloo-gateway category
		out, _, err = s.TestHelper.Execute(s.Ctx, "get", "crd", crd.GetName(), "-o", "json")
		s.NoError(err)

		var crdJson v1.CustomResourceDefinition
		s.NoError(json.Unmarshal([]byte(out), &crdJson))
		s.Contains(crdJson.Spec.Names.Categories, CommonCRDCategory)

		// Ensure the CRD has the solo-io category iff it's an enterprise CRD
		s.Equal(
			slices.Contains(enterpriseCRDs, crd.GetName()),
			slices.Contains(crdJson.Spec.Names.Categories, enterpriseCRDCategory),
		)
	}
}

func (s *testingSuite) UpgradeWithCustomValuesFile(valuesFile string) {
	_, err := s.TestHelper.UpgradeGloo(s.Ctx, 600*time.Second, helper.WithExtraArgs([]string{
		// Do not reuse the existing values as we need to install the new chart with the new version of the images
		"--values", valuesFile,
	}...))
	s.TestInstallation.Assertions.Require.NoError(err)
}
