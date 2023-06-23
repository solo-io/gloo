package make_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/solo-io/skv2/codegen/util"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMake(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Make Suite")
}

type EnvVar struct {
	Name, Value string
}

type MakeVar struct {
	Name, ExpectedValue string
}

func ExpectMakeVarsWithEnvVars(envVars []*EnvVar, expectedMakeVars []*MakeVar) {
	for _, envVar := range envVars {
		err := os.Setenv(envVar.Name, envVar.Value)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
	}
	defer func() {
		for _, envVar := range envVars {
			err := os.Unsetenv(envVar.Name)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
		}
	}()

	rootDirectory := util.GetModuleRoot()
	for _, makeVar := range expectedMakeVars {
		cmd := exec.Command("make", "--directory", rootDirectory, fmt.Sprintf("print-%s", makeVar.Name))
		out, err := cmd.CombinedOutput()
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		ExpectWithOffset(1, getRelevantOutput(string(out))).To(Equal(makeVar.ExpectedValue))
	}
}

func getRelevantOutput(output string) string {
	// We trim lines that are produced in our CI pipeline
	// These are not present locally, so the trim is a no-op
	relevantOutput := strings.TrimSpace(output)
	relevantOutput = strings.TrimPrefix(relevantOutput, "make[1]: Entering directory '/workspace/solo-projects'")
	relevantOutput = strings.TrimSuffix(relevantOutput, "make[1]: Leaving directory '/workspace/solo-projects'")
	return strings.TrimSpace(relevantOutput)
}
