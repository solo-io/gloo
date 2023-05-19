package generate

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	glooGenerate "github.com/solo-io/gloo/install/helm/gloo/generate"
	"github.com/solo-io/k8s-utils/installutils/helmchart"
)

const (
	glooFedDependencyName = "gloo-fed"
	glooOsDependencyName  = "gloo"
	glooOsModuleName      = "github.com/solo-io/gloo"
)

func GetGlooOsVersion(args *GenerationArguments, filesets ...*GenerationFiles) (string, error) {
	if args.GlooRepoOverride != "" {
		// in the event we provide a glooRepoOverride, we want to use an identical version of gloo and gloo-ee
		return args.Version, nil
	}
	return getDependencyVersion(glooOsDependencyName, glooOsModuleName, filesets...)
}

func getDependencyVersion(dependencyName, moduleName string, filesets ...*GenerationFiles) (string, error) {
	var dl DependencyList
	for _, fs := range filesets {
		if err := readYaml(fs.RequirementsTemplate, &dl); err != nil {
			return "", err
		}
		for _, v := range dl.Dependencies {
			if v.Name == dependencyName && v.Version != "" {
				return v.Version, nil
			}
		}
	}
	return goModPackageVersion(moduleName)
}

func goModPackageVersion(moduleName string) (string, error) {
	cmd := exec.Command("go", "list", "-f", "'{{ .Version }}'", "-m", moduleName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	cleanedOutput := strings.Trim(strings.TrimSpace(string(output)), "'")
	return strings.TrimPrefix(cleanedOutput, "v"), nil
}

func readYaml(path string, obj interface{}) error {
	bytes, err := os.ReadFile(path)
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

	err = os.WriteFile(path, bytes, os.ModePerm)
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

func generateRequirementsYaml(requirementsTemplatePath, outputPath, osGlooVersion string, glooFedVersion string, glooFedRepoOverride string, glooRepositoryOverride string) error {
	var dl DependencyList
	if err := readYaml(requirementsTemplatePath, &dl); err != nil {
		return err
	}
	for i, v := range dl.Dependencies {
		if v.Name == glooOsDependencyName && glooRepositoryOverride != "" {
			fmt.Print("overriding gloo repository\n")
			dl.Dependencies[i].Repository = glooRepositoryOverride
		}

		if v.Name == glooOsDependencyName && v.Version == "" {
			dl.Dependencies[i].Version = osGlooVersion
		}
		if v.Name == glooFedDependencyName {
			if v.Version == "" {
				dl.Dependencies[i].Version = glooFedVersion
			}
			if glooFedRepoOverride != "" {
				dl.Dependencies[i].Repository = glooFedRepoOverride
			}
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

func writeDocs(docs helmchart.HelmValues, path string) error {
	err := os.WriteFile(path, []byte(docs.ToMarkdown()), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed writing helm values file")
	}
	return nil
}
