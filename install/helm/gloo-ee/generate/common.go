package generate

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	glooGenerate "github.com/solo-io/gloo/install/helm/gloo/generate"
)

func GetGlooOsVersion(filesets ...*GenerationFiles) (string, error) {
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
	return glooGoModPackageVersion()
}

func glooGoModPackageVersion() (string, error) {
	cmd := exec.Command("go", "list", "-f", "'{{ .Version }}'", "-m", "github.com/solo-io/gloo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	cleanedOutput := strings.Trim(strings.TrimSpace(string(output)), "'")
	return strings.TrimPrefix(cleanedOutput, "v"), nil
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
