package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	glooGenerate "github.com/solo-io/gloo/install/helm/gloo/generate"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-projects/install/helm/gloo-ee/generate"
)

const (
	valuesTemplate       = "install/helm/gloo-ee/values-template.yaml"
	valuesOutput         = "install/helm/gloo-ee/values.yaml"
	distributionOutput   = "install/distribution/values.yaml"
	chartTemplate        = "install/helm/gloo-ee/Chart-template.yaml"
	chartOutput          = "install/helm/gloo-ee/Chart.yaml"
	requirementsTemplate = "install/helm/gloo-ee/requirements-template.yaml"
	requirementsOutput   = "install/helm/gloo-ee/requirements.yaml"

	gopkgToml    = "Gopkg.toml"
	depsToml     = "Deps.toml"
	constraint   = "constraint"
	glooPkg      = "github.com/solo-io/gloo"
	glooiPkg     = "github.com/solo-io/gloo-i"
	nameConst    = "name"
	versionConst = "version"

	always       = "Always"
	ifNotPresent = "IfNotPresent"
)

var (
	osGlooVersion string
)

func main() {
	var version, repoPrefixOverride = "", ""
	var err error
	if len(os.Args) < 2 {
		panic("Must provide version as argument")
	} else {
		version = os.Args[1]

		if len(os.Args) == 3 {
			repoPrefixOverride = os.Args[2]
		}
	}

	osGlooVersion, err = GetVersionFromToml(gopkgToml, glooPkg)
	if err != nil {
		log.Fatalf("failed to determine open source Gloo version. Cause: %v", err)
	}
	log.Printf("Open source gloo version is: %v", osGlooVersion)

	log.Printf("Generating helm files.")
	if err := generateValuesYamls(version, repoPrefixOverride); err != nil {
		log.Fatalf("generating values.yaml failed!: %v", err)
	}
	if err := generateChartYaml(version); err != nil {
		log.Fatalf("generating Chart.yaml failed!: %v", err)
	}
	if err := generateRequirementsYaml(); err != nil {
		log.Fatalf("unable to parse Gopkg.toml for proper gloo version: %v", err)
	}
}

func readConfig() (generate.HelmConfig, error) {
	var config generate.HelmConfig
	if err := readYaml(valuesTemplate, &config); err != nil {
		return config, err
	}
	return config, nil
}

func readYaml(path string, obj interface{}) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "failed reading server config file: %s", path)
	}

	if err := yaml.Unmarshal(bytes, obj); err != nil {
		return errors.Wrap(err, "failed parsing configuration file")
	}

	return nil
}

func writeYaml(obj interface{}, path string) error {
	bytes, err := yaml.Marshal(obj)
	if err != nil {
		return errors.Wrapf(err, "failed marshaling config struct")
	}

	err = ioutil.WriteFile(path, bytes, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failing writing config file")
	}
	return nil
}

func generateValuesYaml(version, pullPolicy, outputFile, repositoryPrefix string) error {
	config, err := readConfig()
	if err != nil {
		return err
	}

	config.Gloo.Gloo.Deployment.Image.Tag = version
	for _, v := range config.Gloo.GatewayProxies {
		v.PodTemplate.Image.Tag = version
	}
	if config.Gloo.IngressProxy != nil {
		config.Gloo.IngressProxy.Deployment.Image.Tag = version
	}
	// Use open source gloo version for discovery and gateway
	config.Gloo.Discovery.Deployment.Image.Tag = osGlooVersion
	config.Gloo.Gateway.Deployment.Image.Tag = osGlooVersion
	config.RateLimit.Deployment.Image.Tag = version
	config.Observability.Deployment.Image.Tag = version
	config.Gloo.ApiServer.Deployment.Server.Image.Tag = version
	config.Gloo.ApiServer.Deployment.Envoy.Image.Tag = version
	config.Gloo.ApiServer.Deployment.Ui.Image.Tag = version

	config.Gloo.Gloo.Deployment.Image.PullPolicy = pullPolicy
	for _, v := range config.Gloo.GatewayProxies {
		v.PodTemplate.Image.PullPolicy = pullPolicy
	}
	if config.Gloo.IngressProxy != nil {
		config.Gloo.IngressProxy.Deployment.Image.PullPolicy = pullPolicy
	}
	config.Gloo.Discovery.Deployment.Image.PullPolicy = pullPolicy
	config.Gloo.Gateway.Deployment.Image.PullPolicy = pullPolicy
	config.RateLimit.Deployment.Image.PullPolicy = pullPolicy
	config.Observability.Deployment.Image.PullPolicy = pullPolicy
	config.Redis.Deployment.Image.PullPolicy = pullPolicy
	config.Gloo.ApiServer.Deployment.Ui.Image.PullPolicy = pullPolicy
	config.Gloo.ApiServer.Deployment.Server.Image.PullPolicy = pullPolicy
	config.Gloo.ApiServer.Deployment.Envoy.Image.PullPolicy = pullPolicy

	if err = updateExtensionsImageVersionAndPullPolicy(config, version, pullPolicy); err != nil {
		return err
	}

	if version == "dev" {
		config.Gloo.Gloo.Deployment.Image.PullPolicy = always
		config.Gloo.Discovery.Deployment.Image.PullPolicy = always
		config.Gloo.Gateway.Deployment.Image.PullPolicy = always
		config.Gloo.Gateway.ConversionJob.Image.PullPolicy = always
		for _, v := range config.Gloo.GatewayProxies {
			v.PodTemplate.Image.PullPolicy = always
		}
	}

	if repositoryPrefix != "" {
		config.Global.Image.Registry = repositoryPrefix
	}

	return writeYaml(&config, outputFile)
}

func updateExtensionsImageVersionAndPullPolicy(config generate.HelmConfig, version, pullPolicy string) (err error) {
	bytes, err := yaml.Marshal(config.Global.Extensions)
	if err != nil {
		return err
	}
	var glooEeExtensions generate.GlooEeExtensions
	err = yaml.Unmarshal(bytes, &glooEeExtensions)
	if err != nil {
		return err
	}
	glooEeExtensions.ExtAuth.Deployment.Image.Tag = version
	glooEeExtensions.ExtAuth.Deployment.Image.PullPolicy = pullPolicy
	config.Global.Extensions = glooEeExtensions
	return nil
}

func generateValuesYamls(version, repositoryPrefix string) error {
	// Generate values for standard manifest
	standardPullPolicy := ifNotPresent
	if version == "dev" {
		standardPullPolicy = always
	}
	if err := generateValuesYaml(version, standardPullPolicy, valuesOutput, repositoryPrefix); err != nil {
		return err
	}

	// Generate values for distribution
	if err := generateValuesYaml(version, ifNotPresent, distributionOutput, ""); err != nil {
		return err
	}
	return nil
}

func generateChartYaml(version string) error {
	var chart glooGenerate.Chart
	if err := readYaml(chartTemplate, &chart); err != nil {
		return err
	}

	chart.Version = version

	return writeYaml(&chart, chartOutput)
}

func generateRequirementsYaml() error {
	var dl generate.DependencyList
	if err := readYaml(requirementsTemplate, &dl); err != nil {
		return err
	}
	for i, v := range dl.Dependencies {
		if v.Name == "gloo" {
			dl.Dependencies[i].Version = osGlooVersion
		}
	}
	return writeYaml(dl, requirementsOutput)
}

func GetVersionFromToml(filename, pkg string) (string, error) {
	config, err := toml.LoadFile(filename)
	if err != nil {
		return "", err
	}

	tomlTree := config.Get(constraint)
	var (
		tree    []*toml.Tree
		version string
	)

	switch typedTree := tomlTree.(type) {
	case []*toml.Tree:
		tree = typedTree
	default:
		return "", fmt.Errorf("unable to parse toml tree")
	}

	for _, v := range tree {
		if v.Get(nameConst) == pkg && v.Get(versionConst) != "" {
			version = v.Get(versionConst).(string)
		}
	}

	return version, nil
}
