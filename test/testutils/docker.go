package testutils

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/avast/retry-go/v4"
	"github.com/solo-io/gloo/test/ginkgo/parallel"
	"github.com/solo-io/go-utils/docker"
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

	dockerArgs := []string{"create"}
	// If running in a non "linux/amd64" environment, you need to add "--platform", "linux/amd64" after "create" or it will use the warning as the image name
	if runtime.GOARCH != "amd64" {
		dockerArgs = append(dockerArgs, "--platform", "linux/amd64")
	}
	dockerArgs = append(dockerArgs, "--name", tmpContainerName, imageName)

	cmd := exec.Command("docker", dockerArgs...)

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
