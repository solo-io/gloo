package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/log"
)

// minikube.go provides helper methods for running tests on minikube

const (
	// glooo labels
	testrunner        = "testrunner"
	helloservice      = "helloservice"
	helloservice2     = "helloservice-2"
	envoy             = "ingress"
	gloo              = "control-plane"
	ingress           = "ingress-controller"
	k8sd              = "k8s-service-discovery"
	funcitonDiscovery = "function-discovery"
	upstreamForEvents = "upstream-for-events"
	eventEmitter      = "event-emitter"
)

// ErrMinikubeNotInstalled indicates minikube binary is not found
var ErrMinikubeNotInstalled = fmt.Errorf("minikube not found in path")

type MinikubeInstance struct {
	vmName             string
	ephemeral          bool
	deployGloo         bool
	ephemeralNamespace string
}

func NewMinikube(deployGloo bool, ephemeralNamespace ...string) *MinikubeInstance {
	var ephemeral bool
	vmName := os.Getenv("MINIKUBE_VM")
	if vmName == "" {
		ephemeral = true
		vmName = "test-minikube-" + RandString(6)
	}
	var namespace string
	if len(ephemeralNamespace) > 0 {
		namespace = ephemeralNamespace[0]
	}
	return &MinikubeInstance{
		vmName:             vmName,
		ephemeral:          ephemeral,
		deployGloo:         deployGloo,
		ephemeralNamespace: namespace,
	}
}

func (mkb *MinikubeInstance) Addr() (string, error) {
	out, err := exec.Command("minikube", "ip", "-p", mkb.vmName).CombinedOutput()
	return "https://" + strings.TrimSuffix(string(out), "\n") + ":8443", err
}

func (mkb *MinikubeInstance) Setup() error {
	if mkb.ephemeral {
		if err := mkb.startMinikube(); err != nil {
			return err
		}
	}
	if mkb.deployGloo {
		if err := mkb.buildContainers(); err != nil {
			return err
		}
		if err := mkb.createE2eResources(); err != nil {
			return err
		}
	}
	if mkb.ephemeralNamespace != "" {
		if err := kubectl("create", "namespace", mkb.ephemeralNamespace); err != nil {
			return err
		}
	}
	return nil
}

func (mkb *MinikubeInstance) Teardown() error {
	if mkb.ephemeral {
		return mkb.deleteMinikube()
	}
	if mkb.deployGloo {
		if err := mkb.deleteE2eResources(); err != nil {
			return err
		}
	}
	if mkb.ephemeralNamespace != "" {
		if err := kubectl("delete", "namespace", mkb.ephemeralNamespace); err != nil {
			return err
		}
	}
	return nil
}

// startMinikube starts a minikube vm with the given name
func (mkb *MinikubeInstance) startMinikube() error {
	log.Debugf("starting minikube %v", mkb.vmName)
	return minikube("start", "-p", mkb.vmName)
}

// deleteMinikube deletes the given minikube vm
func (mkb *MinikubeInstance) deleteMinikube() error {
	log.Debugf("deleting minikube %v", mkb.vmName)
	return minikube("delete", "-p", mkb.vmName)
}

// setMinikubeDockerEnv sets the docker env for the current process
func (mkb *MinikubeInstance) setMinikubeDockerEnv() error {
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
			log.Debugf("unexpected line was not parsed: %v", line)
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

// buildContainers builds all docker containers needed for test
func (mkb *MinikubeInstance) buildContainers() error {
	if os.Getenv("SKIP_BUILD") == "1" {
		return nil
	}
	if err := mkb.setMinikubeDockerEnv(); err != nil {
		return err
	}
	for _, path := range []string{
		filepath.Join(SoloDirectory(), "gloo"),
		filepath.Join(SoloDirectory(), "gloo-ingress"),
		filepath.Join(SoloDirectory(), "gloo-k8s-service-discovery"),
		filepath.Join(E2eDirectory(), "containers", "helloservice"),
		filepath.Join(E2eDirectory(), "containers", "testrunner"),
		filepath.Join(E2eDirectory(), "containers", "event-emitter"),
		filepath.Join(E2eDirectory(), "containers", "upstream-for-events"),
	} {
		log.Debugf("TEST: building container %v", filepath.Base(path))
		cmd := exec.Command("make", "docker")
		cmd.Dir = path
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		cmd = exec.Command("make", "clean")
		cmd.Dir = path
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
	return nil
}

// createE2eResources creates all the kube resources contained in kube_resources dir
func (mkb *MinikubeInstance) createE2eResources() error {

	kubeResourcesDir := filepath.Join(E2eDirectory(), "kube_resources")

	namespace := mkb.ephemeralNamespace
	if namespace == "" {
		namespace = "gloo-system"
	}

	installBytes, err := exec.Command("helm", "template", HelmDirectory(), "--namespace", namespace, "-n", "test-gloo").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "running helm template: %v", installBytes)
	}

	err = ioutil.WriteFile(filepath.Join(kubeResourcesDir, "install.yml"), installBytes, 0644)
	if err != nil {
		return errors.Wrap(err, "writing generated install template")
	}

	if err := kubectl("config", "set-context", mkb.vmName, "--namespace="+namespace); err != nil {
		return err
	}
	// order matters here
	resources := []string{
		"install.yml",

		"helloservice-deployment.yml",
		"helloservice-service.yml",

		"helloservice-2-deployment.yml",
		"helloservice-2-service.yml",

		"event-emitter-deployment.yml",
		"event-upstream-deployment.yml",
		"event-upstream-service.yml",

		"test-runner-pod.yml",
	}
	for _, resource := range resources {
		if err := kubectl("create", "-f", filepath.Join(kubeResourcesDir, resource)); err != nil {
			return err
		}
	}
	if err := waitPodsRunning(testrunner,
		helloservice,
		helloservice2,
		envoy,
		gloo,
		ingress,
		k8sd,
		funcitonDiscovery,
		upstreamForEvents,
		eventEmitter); err != nil {
		return err
	}
	TestRunner("curl", "test-ingress:19000/logging?config=debug")
	return nil
}

// deleteE2eResources deletes all the kube resources contained in kube_resources dir
func (mkb *MinikubeInstance) deleteE2eResources() error {
	if err := kubectl("config", "set-context", mkb.vmName, "--namespace=gloo-system"); err != nil {
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
	return waitPodsTerminated(testrunner,
		helloservice,
		helloservice2,
		envoy,
		gloo,
		ingress,
		funcitonDiscovery,
		k8sd,
		upstreamForEvents,
		eventEmitter)
}

// DeleteContext deletes the context from the kubeconfig
func (mkb *MinikubeInstance) DeleteContext() error {
	return kubectl("config", "delete-context", mkb.vmName)
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
