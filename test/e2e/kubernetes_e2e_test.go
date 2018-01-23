package e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
	"github.com/solo-io/glue/test/e2e/helpers"
	"io/ioutil"
	"os"
	"strings"
	"fmt"
	"github.com/solo-io/glue/module/example"
	"text/template"
	"bytes"
	"time"
)

const glueConfigTmpl = `
{{range .}}
- example_rule:
  timeout: {{.Timeout}}
  match:
    prefix: {{.Match.Prefix}}
  upstream:
    name: {{.Upstream.Name}}
    address: {{.Upstream.Address}}
    port: {{.Upstream.Port}}
{{end}}
`

const helloService = "helloservice"

var _ = Describe("Kubernetes Deployment", func() {
	vmName := "test-" + uuid.New()
	BeforeSuite(func() {
		err := helpers.StartMinikube(vmName)
		Must(err)
		err = helpers.BuildContainers(vmName)
		Must(err)
		err = helpers.CreateKubeResources(vmName)
		Must(err)
	})
	AfterSuite(func() {
		err := helpers.DeleteMinikube(vmName)
		Must(err)
	})
	Describe("E2e", func() {
		Describe("updating glue config", func(){
			It("dynamically updates envoy with new routes", func(){
				randomPath := "/"+uuid.New()
				result, err := curlEnvoy(randomPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(ContainSubstring("< HTTP/1.1 404"))
				rules := []example.ExampleRule{
					newExampleRule(time.Second, randomPath, helloService, helloService, 8080),
				}
				err = updateGlueConfig(rules)
				Expect(err).NotTo(HaveOccurred())
				result, err = curlEnvoy(randomPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(ContainSubstring("< HTTP/1.1 200"))
			})

		})
	})
})

func Must(err error) {
	if err != nil {
		Fail(err.Error())
	}
}

func curlEnvoy(path string) (string, error) {
	return helpers.TestRunner("curl", "http://envoy:8080"+path, "-v")
}

func updateGlueConfig(rules []example.ExampleRule) error {templ := template.New("glue-config")
	t, err := templ.Parse(glueConfigTmpl)
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, rules); err != nil {
		return err
	}
	tmpCmFile, err := ioutil.TempFile("", "glue=configmap.yml")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(tmpCmFile.Name(), buf.Bytes(), 0644); err != nil {
		return err
	}
	defer os.Remove(tmpCmFile.Name())
	out, err := helpers.KubectlOut("apply", "-f", tmpCmFile.Name())
	if err != nil {
		return err
	}
	if !strings.Contains(out, "updated") {
		return fmt.Errorf("expected 'updated' in kubectl output: %v", out)
	}
	return nil
}

func newExampleRule(timeout time.Duration, path, upstreamName, upstreamAddr string, upstreamPort int) example.ExampleRule {
	return example.ExampleRule{
		Timeout: timeout,
		Match: example.Match{
			Prefix: path,
		},
		Upstream: example.Upstream{
			Name: upstreamName,
			Address: upstreamAddr,
			Port: upstreamPort,
		},
	}
}