package e2e_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
	"github.com/solo-io/glue/module/example"
	"github.com/solo-io/glue/test/e2e/helpers"
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
	var vmName string
	BeforeSuite(func() {
		// if a minikube vm exists, we can skip creating and tearing down
		vmName = os.Getenv("MINIKUBE_VM")
		if vmName == "" {
			vmName = "test-" + uuid.New()
			err := helpers.StartMinikube(vmName)
			Must(err)
		}
		err := helpers.BuildContainers(vmName)
		Must(err)
		err = helpers.CreateKubeResources(vmName)
		Must(err)
	})
	AfterSuite(func() {
		if os.Getenv("MINIKUBE_VM") == "" {
			err := helpers.DeleteMinikube(vmName)
			Must(err)
		} else {
			//err := helpers.DeleteKubeResources()
			//Must(err)
		}
	})
	Describe("E2e", func() {
		Describe("updating glue config", func() {
			It("dynamically updates envoy with new routes", func() {
				randomPath := "/" + uuid.New()
				result, err := curlEnvoy(randomPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(ContainSubstring("< HTTP/1.1 404"))
				rules := []example.ExampleRule{
					newExampleRule(time.Second, randomPath, helloService, helloService, 8080),
				}
				err = updateGlueConfig(rules)
				Expect(err).NotTo(HaveOccurred())
				Eventually(func() string {
					res, err := curlEnvoy(randomPath)
					Expect(err).NotTo(HaveOccurred())
					return res
				}).Should(ContainSubstring("< HTTP/1.1 200"))
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
	return helpers.TestRunner("curl", "-v", "http://envoy:8080"+path)
}

func updateGlueConfig(rules []example.ExampleRule) error {
	templ := template.New("glue-config")
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
			Name:    upstreamName,
			Address: upstreamAddr,
			Port:    upstreamPort,
		},
	}
}
