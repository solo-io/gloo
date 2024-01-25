package make_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/solo-io/go-utils/docker"

	"github.com/onsi/gomega/types"

	"github.com/solo-io/skv2/codegen/util"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMake(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Make Suite")
}

const (
	StandardGlooImage = "quay.io/solo-io/gloo-ee:1.16.0-beta1"
	FipsGlooImage     = "quay.io/solo-io/gloo-ee-fips:1.16.0-beta1"
)

var _ = BeforeSuite(func() {
	for _, image := range []string{StandardGlooImage, FipsGlooImage} {
		_, err := docker.PullIfNotPresent(context.Background(), image, 3)
		Expect(err).NotTo(HaveOccurred(), "can pull image locally")
	}
})

var (
	moduleRoot = util.GetModuleRoot()
)

type EnvVar struct {
	Name, Value string
}

type MakeVar struct {
	Name, ExpectedValue string
}

// ExpectMakeVarsWithEnvVars expects that if you assign the provided envVars,
// that the output of each make command is equal to the provided string
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

	for _, makeVar := range expectedMakeVars {
		ExpectMakeOutputWithOffset(1, fmt.Sprintf("print-%s", makeVar.Name), Equal(makeVar.ExpectedValue))
	}
}

// ExpectMakeOutputWithOffset expects that the output of a single make target is equal to the provided matcher
// To provide flags to the target, separate them from the target name with a space:
//
//	ExpectMakeOutputWithOffset(1, "docker-push --ignore-errors", Equal("some output"))
func ExpectMakeOutputWithOffset(offset int, target string, outputMatcher types.GomegaMatcher) {
	makeArgs := append([]string{
		"--directory",
		moduleRoot,
	}, strings.Split(target, " ")...)

	cmd := exec.Command("make", makeArgs...)
	out, err := cmd.CombinedOutput()

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
