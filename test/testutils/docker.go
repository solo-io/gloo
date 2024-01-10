package testutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/avast/retry-go"
	"github.com/onsi/gomega/types"
	"github.com/solo-io/gloo/test/ginkgo/parallel"
	"github.com/solo-io/go-utils/docker"
	"github.com/solo-io/skv2/codegen/util"

	. "github.com/onsi/gomega"
)

var (
	moduleRoot = util.GetModuleRoot()
)

// DockerTag executes a `docker tag` command
func DockerTag(source, dest string) error {
	return docker.Command("tag", source, dest).Run()
}

// DockerPush executes a `docker push` command
func DockerPush(image string) error {
	return docker.Command("push", image).Run()
}

// CopyImageFileToLocal executes a series of docker commands to copy a file from a docker image to the local filesystem
func CopyImageFileToLocal(imageName string, pathToSource, pathToDestination string) error {
	tmpContainerName := fmt.Sprintf("tmp-container-%d", parallel.GetParallelProcessCount())

	// If running in a non "linux/amd64" environment, you need to add "--platform", "linux/amd64" after "create" or it will use the warning as the image name
	cmd := exec.Command("docker", "create", "--name", tmpContainerName, imageName)
	containerIdRaw, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	containerId := strings.TrimSpace(string(containerIdRaw))

	defer func() {
		// Cleanup the container we created
		err = docker.Command("rm", tmpContainerName).Run()
	}()

	return retry.Do(func() error {
		containerPath := fmt.Sprintf("%s:%s", containerId, pathToSource)
		copyCommand := exec.Command("docker", "cp", containerPath, pathToDestination)
		return copyCommand.Run()
	},
		// Retry a few times to account for the fact that the container may not be ready yet
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	)

}

func DockerValidPushTest(validTag string, StandardImage string) {

	// This functionality relies on permissions to push to quay.io, which is only enabled in CI
	ValidateRequirementsAndNotifyGinkgo(DefinedEnv(GcloudBuildId))

	err := DockerTag(StandardImage, validTag)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "can re-tag image locally")

	err = DockerPush(validTag)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "can push to quay.io for existing repository")
}

func DockerInvalidPushTest(invalidTag string, StandardImage string) {

	err := DockerTag(StandardImage, invalidTag)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "can re-tag image locally")

	err = DockerPush(invalidTag)
	ExpectWithOffset(1, err).To(HaveOccurred(), "can NOT push to quay.io for non-existent repository")

}

func DockerPushImagesTest(messageFmt string, dockerImages []string) {
	// The point of this test is to ensure that the docker-push target attempts to push every image that we support
	// As a result, we attempt to push images, and then assert that the output contains the expected error message
	// NOTE: At the moment these images do not exist (since they aren't built in the same pipeline as our tests).
	// If that assumption changes, we may need to re-work these tests if the log output changes
	countMatcher := WithTransform(func(output string) int {
		return strings.Count(output, messageFmt)
	}, Equal(len(dockerImages)))

	// If this test fails, it's likely because you added a new image to the list of images that we attempt to build and publish during our CI pipeline
	// If this is true, you will need to request that we configure the Quay repository for this image
	// If you do not, we will hit a failure during our release pipeline:
	// https://github.com/solo-io/solo-projects/issues/5372#issuecomment-1732184633
	ExpectMakeOutputWithOffset(1, "docker-push --ignore-errors", countMatcher)
}

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

const (
	ExpectStandardCrypto = false
	ExpectBoringCrypto   = true
)

func ValidateCrypto(imageName string, binaryPath string, binaryLocalPath string, expectFips bool) {
	pwd, err := os.Getwd()
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "can get working directory")

	err = CopyImageFileToLocal(imageName, binaryPath, binaryLocalPath)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "can copy binary from image to local filesystem")

	expectedCrypto := "standard"
	unexpectedCrypto := "boring"

	if expectFips {
		expectedCrypto = "boring"
		unexpectedCrypto = "standard"
	}

	// Expected crypto
	target := fmt.Sprintf("BINARY=%s validate-%s-crypto --ignore-errors", filepath.Join(pwd, binaryLocalPath), unexpectedCrypto)
	ExpectMakeOutputWithOffset(1, target, ContainSubstring(fmt.Sprintf("validate-%s-crypto] Error 1 (ignored)", unexpectedCrypto)))

	// Fail this one
	target = fmt.Sprintf("BINARY=%s validate-%s-crypto --ignore-errors", filepath.Join(pwd, binaryLocalPath), expectedCrypto)
	ExpectMakeOutputWithOffset(1, target, And(
		ContainSubstring("goversion -crypto"),
		Not(ContainSubstring("Error 1 (ignored)")),
	))
}

type EnvVar struct {
	Name, Value string
}

type MakeVar struct {
	Name, ExpectedValue string
}

// ExpectMakeOutputWithOffset expects that the output of a single make target is equal to the provided matcher
// To provide flags to the target, separate them from the target name with a space:
//
//	ExpectMakeOutputWithOffset(1, "docker-push --ignore-errors", Equal("some output"))
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
