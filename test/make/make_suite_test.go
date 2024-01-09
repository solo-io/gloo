package make_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/docker"
)

const (
	StandardGlooImage = "quay.io/solo-io/gloo:1.16.0-beta1"
	StandardSdsImage  = "quay.io/solo-io/sds:1.16.0-beta1"
	FipsSdsImage      = "quay.io/solo-io/sds-fips:1.17.0-beta2-9037"
)

var _ = BeforeSuite(func() {
	for _, image := range []string{StandardGlooImage, StandardSdsImage, FipsSdsImage} {
		_, err := docker.PullIfNotPresent(context.Background(), image, 3)
		Expect(err).NotTo(HaveOccurred(), "can pull image locally")
	}
})

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
