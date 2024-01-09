package make_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/solo-io/skv2/codegen/util"
)

const (
	StandardGlooImage = "quay.io/solo-io/gloo:1.16.0-beta1"
	SdsGlooImage      = "quay.io/solo-io/sds:1.16.0-beta1-fips"
)

func TestMake(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Make Suite")
}

var (
	moduleRoot = util.GetModuleRoot()
)

type EnvVar struct {
	Name, Value string
}

type MakeVar struct {
	Name, ExpectedValue string
}

func ExpectMakeVarsWithEnvVars(envVars []*EnvVar, makeVars []*MakeVar) {
	for _, envVar := range envVars {
		err := os.Setenv(envVar.Name, envVar.Value)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
	}
	for _, makeVar := range makeVars {
		cmd := exec.Command("make", "-C", "../..", fmt.Sprintf("print-%s", makeVar.Name))
		out, err := cmd.CombinedOutput()
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		output := strings.TrimSpace(string(out))

		ExpectWithOffset(1, output).To(ContainSubstring(makeVar.ExpectedValue))
	}
	for _, envVar := range envVars {
		err := os.Unsetenv(envVar.Name)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
	}
}

func ExpectMakeOutputWithOffset(offset int, target string, outputMatcher types.GomegaMatcher) {
	makeArgs := append([]string{
		"--directory",
		moduleRoot,
	}, strings.Split(target, " ")...)

	cmd := exec.Command("make", makeArgs...)
	out, err := cmd.CombinedOutput()

	fmt.Println(string(out))
	ExpectWithOffset(offset+1, err).NotTo(HaveOccurred(), "make command should succeed")
	ExpectWithOffset(offset+1, out).To(WithTransform(getRelevantOutput, outputMatcher), "make command should produce expected output")
}

func getRelevantOutput(rawOutput []byte) string {
	// We trim lines that are produced in our CI pipeline
	// These are not present locally, so the trim is a no-op
	relevantOutput := strings.TrimSpace(string(rawOutput))
	relevantOutput = strings.TrimPrefix(relevantOutput, "make[1]: Entering directory '/workspace/solo-projects'")
	relevantOutput = strings.TrimSuffix(relevantOutput, "make[1]: Leaving directory '/workspace/solo-projects'")
	return strings.TrimSpace(relevantOutput)
}
