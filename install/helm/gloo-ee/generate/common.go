package generate

import (
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	glooGenerate "github.com/solo-io/gloo/install/helm/gloo/generate"
	"github.com/solo-io/go-utils/versionutils"
	"github.com/solo-io/solo-projects/pkg/version"
)

func GetGlooOsVersion(pathToGopkgTomlDir string, filesets ...*GenerationFiles) (string, error) {
	var dl DependencyList
	for _, fs := range filesets {
		if err := readYaml(fs.RequirementsTemplate, &dl); err != nil {
			return "", err
		}
		for _, v := range dl.Dependencies {
			if v.Name == "gloo" && v.Version != "" {
				return v.Version, nil
			}
		}
	}
	return getGlooOsVersionFromToml(pathToGopkgTomlDir)
}

func getGlooOsVersionFromToml(pathToGopkgTomlDir string) (string, error) {
	tomlTree, err := versionutils.ParseFullTomlFromDir(pathToGopkgTomlDir)
	if err != nil {
		return "", err
	}

	glooVersion, err := versionutils.GetDependencyVersionInfo(version.GlooPkg, tomlTree)
	if err != nil {
		return "", err
	}

	return glooVersion.Version, nil
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

func readConfig(path string) (HelmConfig, error) {
	var config HelmConfig
	if err := readYaml(path, &config); err != nil {
		return config, err
	}
	return config, nil
}

func generateRequirementsYaml(requirementsTemplatePath, outputPath, osGlooVersion string) error {
	var dl DependencyList
	if err := readYaml(requirementsTemplatePath, &dl); err != nil {
		return err
	}
	for i, v := range dl.Dependencies {
		if v.Name == "gloo" && v.Version == "" {
			dl.Dependencies[i].Version = osGlooVersion
		}
	}
	return writeYaml(dl, outputPath)
}

func (gc *GenerationConfig) generateChartYaml(chartTemplate, chartOutput, chartVersion string) error {
	var chart glooGenerate.Chart
	if err := readYaml(chartTemplate, &chart); err != nil {
		return err
	}

	chart.Version = chartVersion

	return writeYaml(&chart, chartOutput)
}
