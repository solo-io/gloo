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
	"github.com/solo-io/glue/pkg/log"
	"github.com/solo-io/glue/test/e2e/helpers"
)

const glueConfigTmpl = `
apiVersion: v1
data:
  glue.yml: |
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
kind: ConfigMap
metadata:
  name: glue-config
  namespace: glue-system
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
			err := helpers.DeleteKubeResources()
			Must(err)
		}
	})
	Describe("E2e", func() {
		Describe("updating glue config", func() {
			It("responds 503 for a route with misconfigured upstream", func() {
				curlEventuallyShouldRespond("/broken", "< HTTP/1.1 503", time.Minute*3)
			})
			Context("update glue with new rules", func() {
				randomPath := "/" + uuid.New()
				It("responds 404 before update", func() {
					curlEventuallyShouldRespond(randomPath, "< HTTP/1.1 404")
				})
				It("responds 200 after update", func() {
					rules := []example.ExampleRule{
						newExampleRule(time.Second, randomPath, helloService, helloService, 8080),
					}
					err := updateGlueConfig(rules)
					Expect(err).NotTo(HaveOccurred())
					curlEventuallyShouldRespond(randomPath, "< HTTP/1.1 200", time.Minute*30)
				})
			})

		})
	})
})

func curlEventuallyShouldRespond(path, substr string, timeout ...time.Duration) {
	t := time.Second * 20
	if len(timeout) > 0 {
		t = timeout[0]
	}
	// for some useful-ish output
	tick := time.Tick(t / 8)
	Eventually(func() string {
		res, err := curlEnvoy(path)
		if err != nil {
			res = err.Error()
		}
		select {
		default:
			break
		case <-tick:
			log.Printf("curl output: %v", res)
		}
		return res
	}, t).Should(ContainSubstring(substr))
}

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
	if !strings.Contains(out, `configmap "glue-config" configured`) {
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
