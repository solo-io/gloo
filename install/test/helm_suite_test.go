package test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/testutils"
	. "github.com/solo-io/k8s-utils/manifesttestutils"
	"github.com/solo-io/solo-projects/install/helm/gloo-ee/generate"
	fedgenerate "github.com/solo-io/solo-projects/install/helm/gloo-fed/generate"
	soloprojectsinstall "github.com/solo-io/solo-projects/pkg/install"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	k8syamlutil "sigs.k8s.io/yaml"
)

const (
	topLevelDir = "../.."
	chartsDir   = "../helm/"
	namespace   = defaults.GlooSystem
)

func TestHelm(t *testing.T) {
	RegisterFailHandler(Fail)
	testutils.RegisterCommonFailHandlers()
	RunSpecs(t, "Helm Suite")
}

var _ = BeforeSuite(func() {
	err := MakeCmd(topLevelDir, "init-helm")
	Expect(err).NotTo(HaveOccurred(), "Should be able to `make init-helm`")
})

func MakeCmd(dir string, args ...string) error {
	makeCmd := exec.Command("make", args...)
	makeCmd.Dir = dir

	makeCmd.Stdout = GinkgoWriter
	makeCmd.Stderr = GinkgoWriter
	err := makeCmd.Run()

	return err
}

type helmValues struct {
	valuesFile string
	valuesArgs []string // each entry should look like `path.to.helm.field=value`
}

func BuildTestManifest(chartName string, namespace string, values helmValues) (TestManifest, error) {
	chartPath := path.Join(chartsDir, chartName)
	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	values.valuesArgs = append(values.valuesArgs, "license_key=\"placeholder-license-key\"")
	helmValues, err := buildHelmValues(chartPath, values)
	if err != nil {
		return nil, err
	}

	err = validateHelmValues(chartName, helmValues)
	if err != nil {
		return nil, err
	}

	client, err := buildRenderer(namespace)
	if err != nil {
		return nil, err
	}

	rel, err := client.Run(chartRequested, helmValues)
	if err != nil {
		return nil, err
	}

	// the test manifest utils can only read from a file, ugh
	f, err := os.CreateTemp("", "*.yaml")
	Expect(err).NotTo(HaveOccurred(), "Should be able to write a temp file for the helm unit test manifest")
	defer os.Remove(f.Name())

	_, err = f.Write([]byte(rel.Manifest))
	Expect(err).NotTo(HaveOccurred(), "Should be able to write the release manifest to the temp file for the helm unit tests")

	for _, hook := range rel.Hooks {
		manifest := hook.Manifest
		_, err = f.Write([]byte("\n---\n" + manifest))
		Expect(err).NotTo(HaveOccurred(), "Should be able to write the hook manifest to the temp file for the helm unit tests")
	}

	return NewTestManifest(f.Name()), nil
}

// each entry in valuesArgs should look like `path.to.helm.field=value`
func buildHelmValues(chartPath string, values helmValues) (map[string]interface{}, error) {
	// read the chart's base values file first
	finalValues, err := readValuesFile(path.Join(chartPath, "values.yaml"))
	if err != nil {
		return nil, err
	}

	for _, v := range values.valuesArgs {
		err := strvals.ParseInto(v, finalValues)
		if err != nil {
			return nil, err
		}
	}

	if values.valuesFile != "" {
		// these lines ripped out of Helm internals
		// https://github.com/helm/helm/blob/release-3.0/pkg/cli/values/options.go
		mapFromFile, err := readValuesFile(values.valuesFile)
		if err != nil {
			return nil, err
		}

		// Merge with the previous map
		finalValues = mergeMaps(finalValues, mapFromFile)
	}

	return finalValues, nil
}

// validateHelmValues ensures that the unstructured helm values that are provided
// to a chart match the Go type used to generate the Helm documentation
// Returns nil if all the provided values are all included in the Go struct
// Returns an error if a provided value is not included in the Go struct.
func validateHelmValues(chartName string, unstructuredHelmValues map[string]interface{}) error {
	if chartName == soloprojectsinstall.GlooEnterpriseChartName {
		// This Go type is the source of truth for the Helm docs
		var structuredHelmValues generate.HelmConfig
		unstructuredHelmValueBytes, err := json.Marshal(unstructuredHelmValues)
		if err != nil {
			return err
		}
		// This ensures that an error will be raised if there is an unstructured helm value
		// defined but there is not the equivalent type defined in our Go struct
		//
		// When an error occurs, this means the Go type needs to be amended
		// to include the new field (which is the source of truth for our docs)
		return errors.Wrapf(
			k8syamlutil.UnmarshalStrict(unstructuredHelmValueBytes, &structuredHelmValues),
			"Gloo EE Helm Values: %s", unstructuredHelmValues)
	} else if chartName == soloprojectsinstall.GlooFedChartName {
		var structuredHelmValues fedgenerate.HelmConfig
		unstructuredHelmValueBytes, err := json.Marshal(unstructuredHelmValues)
		if err != nil {
			return err
		}
		return errors.Wrapf(
			k8syamlutil.UnmarshalStrict(unstructuredHelmValueBytes, &structuredHelmValues),
			"Gloo Fed Helm Values: %s", unstructuredHelmValues)
	} else {
		return errors.New(fmt.Sprintf("unsupported chart name: %s", chartName))
	}
}

func readValuesFile(filePath string) (map[string]interface{}, error) {
	mapFromFile := map[string]interface{}{}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// NOTE: This is not the default golang yaml.Unmarshal, because that implementation
	// does not unmarshal into a map[string]interface{}; it unmarshals the file into a map[interface{}]interface{}
	// https://github.com/go-yaml/yaml/issues/139
	if err := k8syamlutil.Unmarshal(bytes, &mapFromFile); err != nil {
		return nil, err
	}

	return mapFromFile, nil
}

func buildRenderer(namespace string) (*action.Install, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	noOpDebugLog := func(format string, v ...interface{}) {}

	if err := actionConfig.Init(
		settings.RESTClientGetter(),
		defaults.GlooSystem,
		os.Getenv("HELM_DRIVER"),
		noOpDebugLog,
	); err != nil {
		return nil, err
	}

	renderer := action.NewInstall(actionConfig)
	renderer.DryRun = true
	renderer.Namespace = namespace
	renderer.ReleaseName = soloprojectsinstall.GlooEnterpriseReleaseName
	renderer.Namespace = "gloo-system"
	renderer.ClientOnly = true

	return renderer, nil
}

// stolen from Helm internals
// https://github.com/helm/helm/blob/release-3.0/pkg/cli/values/options.go#L88
func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

func makeUnstructured(yam string) *unstructured.Unstructured {
	jsn, err := yaml.YAMLToJSON([]byte(yam))
	Expect(err).NotTo(HaveOccurred())
	runtimeObj, err := runtime.Decode(unstructured.UnstructuredJSONScheme, jsn)
	Expect(err).NotTo(HaveOccurred())
	return runtimeObj.(*unstructured.Unstructured)
}
