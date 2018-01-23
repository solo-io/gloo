package e2e_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/pborman/uuid"
	"github.com/solo-io/glue/test/e2e/helpers"
	"io/ioutil"
	"os"
	"strings"
	"fmt"
)

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
		Describe("envoy resolves configured rules", func(){
			BeforeEach(func(){
				randomPath := "/"+uuid.New()
				result, err := curlEnvoy(randomPath)
				Must(err)
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
	return helpers.TestRunner("curl", "http://envoy:8080"+path)
}

func updateGlueConfig(contents string) error {
	tmpCmFile, err := ioutil.TempFile("", "glue=configmap.yml")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(tmpCmFile.Name(), []byte(contents), 0644); err != nil {
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

func newGlueConfig(path string) 