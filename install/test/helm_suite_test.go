package test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	soloprojectsinstall "github.com/solo-io/solo-projects/pkg/install"

	"github.com/solo-io/gloo/pkg/cliutil/helm"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/ghodss/yaml"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"helm.sh/helm/v3/pkg/action"
	k8syamlutil "sigs.k8s.io/yaml"

	. "github.com/solo-io/go-utils/manifesttestutils"
	"github.com/solo-io/go-utils/testutils"
	"helm.sh/helm/v3/pkg/chart/loader"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	helmValues, err := buildHelmValues(chartPath, values)
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
	f, err := ioutil.TempFile("", "*.yaml")
	Expect(err).NotTo(HaveOccurred(), "Should be able to write a temp file for the helm unit test manifest")
	defer os.Remove(f.Name())

	_, err = f.Write([]byte(rel.Manifest))
	Expect(err).NotTo(HaveOccurred(), "Should be able to write the release manifest to the temp file for the helm unit tests")

	// also need to add in the hooks, which are not included in the release manifest
	// be sure to skip the resources that we duplicate because of Helm hook weirdness (see the comment on helm.GetNonCleanupHooks)
	nonCleanupHooks, err := helm.GetNonCleanupHooks(rel.Hooks)
	Expect(err).NotTo(HaveOccurred(), "Should be able to get the non-cleanup hooks in the helm unit test setup")

	for _, hook := range nonCleanupHooks {
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

func readValuesFile(filePath string) (map[string]interface{}, error) {
	mapFromFile := map[string]interface{}{}

	bytes, err := ioutil.ReadFile(filePath)
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
