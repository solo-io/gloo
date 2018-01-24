package helpers

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pborman/uuid"
)

// minikube.go provides helper methods for running tests on minikube

const (
	testrunner   = "testrunner"
	helloservice = "helloservice"
	envoy        = "envoy"
	glue         = "glue"
)

// ErrMinikubeNotInstalled indicates minikube binary is not found
var ErrMinikubeNotInstalled = fmt.Errorf("minikube not found in path")

type MinikubeInstance struct {
	vmName    string
	ephemeral bool
}

func NewMinikube() *MinikubeInstance {
	var ephemeral bool
	vmName := os.Getenv("MINIKUBE_VM")
	if vmName == "" {
		ephemeral = true
		vmName = "test-minikube" + uuid.New()
	}
	return &MinikubeInstance{
		vmName:    vmName,
		ephemeral: ephemeral,
	}
}

func (mkb *MinikubeInstance) Setup() error {
	if mkb.ephemeral {
		if err := mkb.StartMinikube(); err != nil {
			return err
		}
	}
	if err := mkb.BuildContainers(); err != nil {
		return err
	}
	return mkb.CreateKubeResources()
}

func (mkb *MinikubeInstance) Teardown() error {
	if mkb.ephemeral {
		return mkb.DeleteMinikube()
	}
	return mkb.DeleteKubeResources()
}

// StartMinikube starts a minikube vm with the given name
func (mkb *MinikubeInstance) StartMinikube() error {
	log.Printf("starting minikube %v", mkb.vmName)
	return minikube("start", "-p", mkb.vmName)
}

// DeleteMinikube deletes the given minikube vm
func (mkb *MinikubeInstance) DeleteMinikube() error {
	log.Printf("deleting minikube %v", mkb.vmName)
	return minikube("delete", "-p", mkb.vmName)
}

// SetMinikubeDockerEnv sets the docker env for the current process
func (mkb *MinikubeInstance) SetMinikubeDockerEnv() error {
	bashEnv, err := minikubeOutput("docker-env", "-p", mkb.vmName)
	if err != nil {
		return err
	}
	lines := strings.Split(bashEnv, "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "export ") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			log.Printf("unexpected line was not parsed: %v", line)
			continue
		}
		key := parts[0]
		val := strings.TrimPrefix(parts[1], "\"")
		val = strings.TrimSuffix(val, "\"")
		if err := os.Setenv(key, val); err != nil {
			return err
		}
	}
	return nil
}

// BuildContainers builds all docker containers needed for test
func (mkb *MinikubeInstance) BuildContainers() error {
	if err := mkb.SetMinikubeDockerEnv(); err != nil {
		return err
	}
	containerDir := filepath.Join(E2eDirectory(), "containers")
	return filepath.Walk(containerDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() || path == containerDir {
			return nil
		}
		log.Printf("TEST: building container %v", filepath.Base(path))
		cmd := exec.Command(filepath.Join(path, "build.sh"))
		cmd.Dir = path
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		return cmd.Run()
	})
}

// CreateKubeResources creates all the kube resources contained in kube_resources dir
func (mkb *MinikubeInstance) CreateKubeResources() error {
	kubeResourcesDir := filepath.Join(E2eDirectory(), "kube_resources")
	if err := kubectl("config", "set-context", mkb.vmName, "--namespace=glue-system"); err != nil {
		return err
	}
	// order matters here
	resources := []string{
		"namespace.yml",

		"glue-configmap.yml",
		"glue-deployment.yml",
		"glue-service.yml",

		"envoy-configmap.yml",
		"envoy-deployment.yml",
		"envoy-service.yml",

		"helloservice-deployment.yml",
		"helloservice-service.yml",

		"test-runner-pod.yml",
	}
	for _, resource := range resources {
		if err := kubectl("create", "-f", filepath.Join(kubeResourcesDir, resource)); err != nil {
			return err
		}
	}
	return WaitPodsRunning(testrunner, helloservice, envoy, glue)
}

func kubectl(args ...string) error {
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func KubectlOut(args ...string) (string, error) {
	cmd := exec.Command("kubectl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%s (%v)", out, err)
	}
	return string(out), err
}

// DeleteKubeResources deletes all the kube resources contained in kube_resources dir
func (mkb *MinikubeInstance) DeleteKubeResources() error {
	if err := kubectl("config", "set-context", mkb.vmName, "--namespace=glue-system"); err != nil {
		return err
	}
	kubeResourcesDir := filepath.Join(E2eDirectory(), "kube_resources")
	if err := kubectl("delete", "-f", filepath.Join(kubeResourcesDir, "namespace.yml")); err != nil {
		return err
	}
	// test runner pod is slow to tear down sometimes, just force it
	if err := kubectl("delete", "-f", filepath.Join(kubeResourcesDir, "test-runner-pod.yml"), "--force"); err != nil {
		return err
	}
	return WaitPodsTerminated(testrunner, helloservice, envoy, glue)
}

// DeleteContext deletes the context from the kubeconfig
func (mkb *MinikubeInstance) DeleteContext() error {
	return kubectl("config", "delete-context", mkb.vmName)
}

// WaitPodsRunning waits for all pods to be running
func WaitPodsRunning(podNames ...string) error {
	for _, pod := range podNames {
		finished := func(output string) bool {
			return strings.Contains(output, "Running")
		}
		if err := waitPodStatus(pod, "Running", finished); err != nil {
			return err
		}
	}
	return nil
}

// WaitPodsTerminated waits for all pods to be terminated
func WaitPodsTerminated(podNames ...string) error {
	for _, pod := range podNames {
		finished := func(output string) bool {
			return !strings.Contains(output, pod)
		}
		if err := waitPodStatus(pod, "terminated", finished); err != nil {
			return err
		}
	}
	return nil
}

// TestRunner executes a command inside the TestRunner container
func TestRunner(args ...string) (string, error) {
	args = append([]string{"exec", "-i", testrunner, "--"}, args...)
	return KubectlOut(args...)
}

func waitPodStatus(pod, status string, finished func(output string) bool) error {
	timeout := time.Second * 20
	interval := time.Millisecond * 1000
	tick := time.Tick(interval)

	log.Printf("waiting %v for pod %v to be %v...", timeout, pod, status)
	for {
		select {
		case <-time.After(timeout):
			return fmt.Errorf("timed out waiting for %v to be %v", pod, status)
		case <-tick:
			out, err := KubectlOut("get", "pod", "-l", "app="+pod)
			if err != nil {
				return fmt.Errorf("failed getting pod: %v", err)
			}
			if finished(out) {
				return nil
			}
		}
	}
}

func minikube(args ...string) error {
	cmd := exec.Command("minikube", args...)
	// useful for logging during tests
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == exec.ErrNotFound {
		return ErrMinikubeNotInstalled
	}
	return err
}

func minikubeOutput(args ...string) (string, error) {
	cmd := exec.Command("minikube", args...)
	out, err := cmd.Output()
	if err == exec.ErrNotFound {
		return "", ErrMinikubeNotInstalled
	}
	return string(out), err
}
