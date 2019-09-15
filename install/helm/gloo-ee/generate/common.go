package generate

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	glooGenerate "github.com/solo-io/gloo/install/helm/gloo/generate"

	"github.com/pelletier/go-toml"
)

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

func generateRequirementsYaml(requirementsTemplatePath, outputPath, osGLooVersion string) error {
	var dl DependencyList
	if err := readYaml(requirementsTemplatePath, &dl); err != nil {
		return err
	}
	for i, v := range dl.Dependencies {
		if v.Name == "gloo" {
			dl.Dependencies[i].Version = osGLooVersion
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
