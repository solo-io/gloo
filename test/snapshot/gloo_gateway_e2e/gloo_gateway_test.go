package bugs_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/test/snapshot/gloo_gateway_e2e/testcase"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/solo-io/skv2/contrib/pkg/sets"
)

const (
	outputSnapToFileEnv = "OUTPUT_SNAPSHOT_TO_FILE"
)

var _ = Describe("Snapshot Test", func() {
	for _, test := range testcase.Testcases {
		test := test // pike
		it := It
		if test.Focus {
			it = FIt
		} else if test.Skip {
			it = PIt
		}
		When(test.When, func() {
			it(test.It, func() {
				t := translator.NewDefaultTranslator()

				var err error
				var inputResources []*client.Object

				if test.TestInput.TranslatorInputYamlFile != "" {
					// TODO: load yaml
				} else {
					inputResources = test.TestInput.TranslatorInput
				}

				runTest(t, test.Assertions)

			})

		})
	}
})

func runTest(
	t translator.Translator,
	assertions []testcase.TranslatorAssertion,
) {

	out, _, _ := t.Translate()

	if os.Getenv(outputSnapToFileEnv) == "true" {
		o, err := out.MarshalJSON()
		Expect(err).NotTo(HaveOccurred())

		dir := "/tmp/gp-snapshot-outputs"
		err = os.MkdirAll(dir, os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		fileName := fmt.Sprintf("%s.json", strings.ToLower(strings.ReplaceAll(CurrentSpecReport().LeafNodeText, " ", "-")))
		tmpfile, err := os.Create(filepath.Join(dir, fileName))
		Expect(err).NotTo(HaveOccurred())

		defer tmpfile.Close()

		fmt.Fprintf(GinkgoWriter, "Writing snapshot output to %s\n", tmpfile.Name())
		_, err = io.Copy(tmpfile, bytes.NewReader(o))
		Expect(err).NotTo(HaveOccurred())
	}

	if len(errs) > 0 {
		for ws, err := range errs {
			fmt.Fprintf(GinkgoWriter, "WS: (%s), Error: %v\n", sets.Key(ws), err)
		}
		Fail("Errors occurred during translation")
	}

	for _, outs := range out.TctxOutput {
		for _, outputs := range outs {
			Expect(outputs.Err).NotTo(HaveOccurred())
			outputs.Outputs.ForEach(func(obj client.Object) {
				// Check idempotence, no output should appear in more than one output
				Expect(allOutputs.Contains(obj)).To(BeFalse(), fmt.Sprintf("%s already exists in another output", sets.TypedKey(obj)))
				allOutputs.Insert(obj)
			})
		}
	}

	for _, outputsMap := range out.GlobalOutput {
		for _, outputs := range outputsMap {
			Expect(outputs.Err).NotTo(HaveOccurred())
			outputs.Outputs.ForEach(func(obj client.Object) {
				// Check idempotence, no output should appear in more than one output
				Expect(allOutputs.Contains(obj)).To(BeFalse(), fmt.Sprintf("%s already exists in another output", sets.TypedKey(obj)))
				allOutputs.Insert(obj)
			})
		}
	}

	for _, assertion := range assertions {
		assertion(allOutputs)
	}
}

// CleanObject removes fields from the object that are not needed, which can decrease size of JSON file by ~ 50%
func cleanObject(input, filename string) {
	var (
		rootIn = map[string]interface{}{}
	)

	Expect(json.NewDecoder(strings.NewReader(input)).Decode(&rootIn)).NotTo(HaveOccurred())

	for k := range rootIn {
		nodes, ok := rootIn[k].([]any)
		if !ok {
			continue
		}

		for i, node := range nodes {
			if node, ok := node.(map[string]any); ok {
				metadata := node["metadata"].(map[string]any)

				for k := range metadata {
					switch k {
					case "managedFields", "resourceVersion", "uid", "generation":
						delete(metadata, k)
					}
				}

				node["metadata"] = metadata
			}

			nodes[i] = node
		}

		rootIn[k] = nodes
	}

	pwd, err := filepath.Abs(".")
	Expect(err).NotTo(HaveOccurred())

	filename = filepath.Join(pwd, "testcase", filename)

	outFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	enc := json.NewEncoder(outFile)
	err = enc.Encode(rootIn)
	Expect(err).NotTo(HaveOccurred())
}

func loadYamlResources(input testcase.TestInput) (error, []*client.Object) {
	// Decode YAML and append objects to the list
	for _, doc := range yaml.SplitYAML(input.TranslatorInputYaml.Filename) {
		obj, _, err := decode([]byte(doc), nil, nil)
		if err != nil {
			log.Fatalf("Error decoding YAML: %v", err)
		}
		objects = append(objects, obj)
	}
}

// unstructuredSerializer returns a serializer for unstructured objects.
func unstructuredSerializer() runtime.Serializer {
	return serializer.NewCodecFactory(scheme.Scheme).WithoutConversion()
}
