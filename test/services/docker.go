package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	. "github.com/onsi/gomega"

	"github.com/onsi/ginkgo"
	"github.com/pkg/errors"
)

var (
	dockerDefaultNetwork = "bridge" // if unspecified, docker containers are created on the default bridge network
)

func RunContainer(containerName string, args []string) error {
	fmt.Fprintln(ginkgo.GinkgoWriter, args)
	cmd := exec.Command("docker", args...)
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "Unable to start "+containerName+" container")
	}
	return nil
}

func ContainerExistsWithName(containerName string) bool {
	cmd := exec.Command("docker", "ps", "-aq", "-f", "name=^/"+containerName+"$")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() [%s %s] failed with %s\n", cmd.Path, cmd.Args, err)
	}
	return string(out) != ""
}

func MustStopContainer(containerName string) {
	err := StopContainer(containerName)
	Expect(err).ToNot(HaveOccurred())
	Eventually(ContainerExistsWithName(containerName), "10s", "1s").Should(BeFalse())
}

func StopContainer(containerName string) error {
	cmd := exec.Command("docker", "kill", containerName)
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "Error stopping container "+containerName)
	}
	return nil
}

func RunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); os.IsNotExist(err) {
		// magic docker env file doesn't exist. not running in docker
		return false
	}
	return true
}

func GetContainerNetwork() string {
	network := dockerDefaultNetwork
	if RunningInDocker() {
		// assume in CI
		network = "cloudbuild"
	}
	return network
}

// If docker containers are running on the same host and their own docker cli is configured on the same
// docker daemon (e.g., google cloudbuild), we can determine their IPs (even if they're on the default bridge network)

// This function is unused for now -- obsolete because docker containers started on the same network
// (except the default bridge network) are addressable by hostname (i.e., their container name)
func GetSiblingDockerIp(containerName, containerNetwork string) string {
	inspect := "docker inspect " + containerName + " -f \"{{json .NetworkSettings.Networks." + containerNetwork + ".IPAddress }}\""
	cmd := exec.Command("bash", "-c", inspect)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() [%s] failed with %s\n", inspect, err)
	}
	siblingIp := strings.TrimSuffix(strings.TrimSpace(string(out)), "\n")
	siblingIp = strings.ReplaceAll(siblingIp, "\"", "")
	fmt.Printf("determined sibling docker ip: %s for container %s\n", siblingIp, containerName)
	return siblingIp
}
