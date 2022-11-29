package installer

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/solo-io/solo-projects/test/kubeutils"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/anyvendor/pkg/modutils"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/k8s-utils/testutils/helper"
	"github.com/solo-io/solo-projects/install/helm/gloo-ee/generate"
	k8syamlutil "sigs.k8s.io/yaml"
)

const (
	helmChartName  = "gloo-ee"
	installTimeout = time.Minute * 2 // A sensible default that should provide enough time for Gloo to be installed
)

var _ Installer = new(GlooInstaller)

type GlooInstaller struct {
	config     InstallConfig
	testHelper *helper.SoloTestHelper
}

type InstallConfig struct {
	ClusterName      string
	InstallNamespace string
	HelmValuesFile   string
}

func NewGlooInstaller(config InstallConfig) (*GlooInstaller, error) {
	goModFile, err := modutils.GetCurrentModPackageFile()
	if err != nil {
		return nil, err
	}

	glooLicenseKey := kubeutils.LicenseKey()
	if glooLicenseKey == "" {
		return nil, errors.New("LicenseKey required, but not provided")
	}

	testHelper, err := helper.NewSoloTestHelper(func(defaults helper.TestConfig) helper.TestConfig {
		defaults.RootDir = filepath.Dir(goModFile)
		defaults.HelmChartName = helmChartName
		defaults.LicenseKey = kubeutils.LicenseKey()
		defaults.InstallNamespace = config.InstallNamespace
		defaults.Verbose = true
		defaults.DeployTestRunner = false
		return defaults
	})
	if err != nil {
		return nil, err
	}

	return &GlooInstaller{
		config:     config,
		testHelper: testHelper,
	}, nil

}

func (g *GlooInstaller) GetContext() (string, string) {
	return g.config.ClusterName, g.config.InstallNamespace
}

func (g *GlooInstaller) Install(ctx context.Context) error {
	log.Printf("installing gloo in [%s] cluster to namespace [%s]", g.config.ClusterName, g.config.InstallNamespace)

	if err := g.validateHelmValuesFile(); err != nil {
		return err
	}

	return g.testHelper.InstallGloo(
		ctx,
		helper.GATEWAY,
		installTimeout,
		helper.ExtraArgs("--values", g.config.HelmValuesFile))
}

func (g *GlooInstaller) validateHelmValuesFile() error {
	valueBytes, err := ioutil.ReadFile(g.config.HelmValuesFile)
	if err != nil {
		return err
	}

	// This Go type is the source of truth for the Helm docs
	var structuredHelmValues generate.HelmConfig

	// This ensures that an error will be raised if there is an unstructured helm value
	// defined but there is not the equivalent type defined in our Go struct
	//
	// When an error occurs, this means the Go type needs to be amended
	// to include the new field (which is the source of truth for our docs)
	return errors.Wrapf(
		k8syamlutil.UnmarshalStrict(valueBytes, &structuredHelmValues),
		"Helm Values File: %s", g.config.HelmValuesFile)
}

func (g *GlooInstaller) Uninstall(ctx context.Context) error {
	log.Printf("uninstalling gloo in [%s] cluster from namespace [%s]", g.config.ClusterName, g.config.InstallNamespace)

	return g.testHelper.UninstallGloo()
}
