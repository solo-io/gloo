package helpers

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/solo-io/gloo/pkg/log"
)

var ImageTag = "testing-" + RandString(4)

// builds and pushes all docker containers needed for test
func BuildPushContainers(push bool) error {
	if os.Getenv("SKIP_BUILD") == "1" {
		return nil
	}
	for _, path := range []string{
		filepath.Join(SoloDirectory(), "gloo"),
		filepath.Join(E2eDirectory(), "containers", "helloservice"),
		filepath.Join(E2eDirectory(), "containers", "testrunner"),
		filepath.Join(E2eDirectory(), "containers", "event-emitter"),
		filepath.Join(E2eDirectory(), "containers", "upstream-for-events"),
		filepath.Join(E2eDirectory(), "containers", "grpc-test-service"),
	} {
		os.Setenv("IMAGE_TAG", ImageTag)
		dockerUser := os.Getenv("DOCKER_USER")
		if dockerUser == "" {
			dockerUser = "soloio"
		}
		fullImage := dockerUser + "/" + filepath.Base(path) + ":" + ImageTag
		log.Debugf("TEST: building fullImage %v", fullImage)
		cmd := exec.Command("make", "docker")
		cmd.Dir = path
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		if push {
			cmd = exec.Command("docker", "push", fullImage)
			cmd.Stdout = os.Stderr
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
		}
		cmd = exec.Command("make", "clean")
		cmd.Dir = path
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
	return nil
}
