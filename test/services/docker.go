package services

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/gomega"

	"github.com/onsi/ginkgo/v2"
	"github.com/pkg/errors"
)

const (
	// defaultNetwork is the default docker network driver
	// https://docs.docker.com/network/drivers/bridge/
	defaultNetwork = "bridge"

	// cloudbuildNetwork is the docker network driver used by Google Cloudbuild
	// https://cloud.google.com/build/docs/build-config-file-schema#network
	cloudbuildNetwork = "cloudbuild"
)

func RunContainer(containerName string, args []string) error {
	updatedContainerName := GetUpdatedContainerName(containerName)
	runArgs := []string{"run", "--name", updatedContainerName}
	runArgs = append(runArgs, args...)
	fmt.Fprintln(ginkgo.GinkgoWriter, args)
	cmd := exec.Command("docker", runArgs...)
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "Unable to start "+containerName+" container")
	}
	return nil
}

// ContainerExistsWithName returns an empty string if the container does not exist
func ContainerExistsWithName(containerName string) string {
	updatedContainerName := GetUpdatedContainerName(containerName)
	cmd := exec.Command("docker", "ps", "-aq", "-f", "name=^/"+updatedContainerName+"$")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() [%s %s] failed with %s\n", cmd.Path, cmd.Args, err)
	}
	return string(out)
}

func ExecOnContainer(containerName string, args []string) ([]byte, error) {
	updatedContainerName := GetUpdatedContainerName(containerName)
	arguments := []string{"exec", updatedContainerName}
	arguments = append(arguments, args...)
	cmd := exec.Command("docker", arguments...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to execute command %v on [%s] container [%s]", arguments, containerName, out)
	}
	return out, nil
}

func MustKillAndRemoveContainer(containerName string) {
	err := KillAndRemoveContainer(containerName)
	Expect(err).ToNot(HaveOccurred())
	// CI host may be extremely CPU-bound as it's often building test assets in tandem with other tests,
	// as well as other CI builds running in parallel. When that happens, the tests can run much slower,
	// thus they need a longer timeout. see https://github.com/solo-io/solo-projects/issues/1701#issuecomment-620873754
	Eventually(ContainerExistsWithName(containerName), "30s", "2s").Should(BeEmpty())
}

func KillAndRemoveContainer(containerName string) error {
	updatedContainerName := GetUpdatedContainerName(containerName)
	cmd := exec.Command("docker", "rm", "-f", updatedContainerName)
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "Error stopping and removing container "+containerName)
	}
	return WaitUntilContainerRemoved(containerName)
}

// poll docker for removal of the container named containerName - block until
// successful or fail after a small number of retries
func WaitUntilContainerRemoved(containerName string) error {
	// if this function returns nil, it means the container is still running
	isContainerRemoved := func() bool {
		cmd := exec.Command("docker", "inspect", containerName)
		return cmd.Run() != nil
	}
	for i := 0; i < 5; i++ {
		if isContainerRemoved() {
			return nil
		}
		fmt.Println("Waiting for removal of container " + containerName)
		time.Sleep(1 * time.Second)
	}
	return errors.New("Unable to delete container " + containerName)
}

func RunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); os.IsNotExist(err) {
		// magic docker env file doesn't exist. not running in docker
		return false
	}
	return true
}

func GetDockerHost(containerName string) string {
	if RunningInDocker() {
		return GetUpdatedContainerName(containerName)
	} else {
		return "127.0.0.1"
	}
}

func GetContainerNetwork() string {
	network := defaultNetwork
	if RunningInDocker() {
		if !runningInCloudbuild() {
			// We error loudly here so that if/when we move off of Google Cloudbuild, we can clean up this logic
			ginkgo.Fail("Running in docker but not in cloudbuild. Could not determine docker network")
		}
		network = cloudbuildNetwork
	}
	return network
}

func runningInCloudbuild() bool {
	gcloudId := os.Getenv("GCLOUD_BUILD_ID")
	return len(gcloudId) > 0
}

func GetUpdatedContainerName(containerName string) string {
	gcloudId := os.Getenv("GCLOUD_BUILD_ID")
	if len(gcloudId) > 0 {
		// we are running in CI - let's suffix our container with gcloud build ID
		// so a concurrent build on the host doesn't try to create a container with
		// a conflicting name
		return containerName + "_" + gcloudId
	}
	return containerName
}
