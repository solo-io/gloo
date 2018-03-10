package helpers

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"os"

	"fmt"
	"strings"
	"time"

	"bytes"
	"text/template"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/log"
)

func SetupKubeForTest(namespace string) error {
	context := os.Getenv("KUBECTL_CONTEXT")
	if context == "" {
		current, err := KubectlOut("config", "current-context")
		if err != nil {
			return errors.Wrap(err, "getting currrent context")
		}
		context = strings.TrimSuffix(current, "\n")
	}
	if err := kubectl("config", "set-context", context, "--namespace="+namespace); err != nil {
		return errors.Wrap(err, "setting context")
	}
	return kubectl("create", "namespace", namespace)
}

func TeardownKube(namespace string) error {
	return kubectl("delete", "namespace", namespace)
}

func SetupKubeForE2eTest(namespace string, buildImages, push bool) error {
	if err := SetupKubeForTest(namespace); err != nil {
		return err
	}
	if buildImages {
		if err := BuildPushContainers(push); err != nil {
			return err
		}
	}
	kubeResourcesDir := filepath.Join(E2eDirectory(), "kube_resources")

	installBytes, err := exec.Command("helm", "template", HelmDirectory(),
		"--namespace", namespace,
		"-n", "test",
		"-f", filepath.Join(kubeResourcesDir, "helm-values.yaml")).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "running helm template: %v", installBytes)
	}

	err = ioutil.WriteFile(filepath.Join(kubeResourcesDir, "install.yml"), installBytes, 0644)
	if err != nil {
		return errors.Wrap(err, "writing generated install template")
	}

	tmpl, err := template.New("Test_Resources").ParseFiles(filepath.Join(kubeResourcesDir, "testing-resources.yaml.tmpl"))
	if err != nil {
		return errors.Wrap(err, "parsing template from testing-resources.yaml.tmpl")
	}

	buf := &bytes.Buffer{}
	data := &struct{ Namespace string }{Namespace: namespace}
	if err := tmpl.ExecuteTemplate(buf, "testing-resources.yaml.tmpl", data); err != nil {
		return errors.Wrap(err, "executing template")
	}

	err = ioutil.WriteFile(filepath.Join(kubeResourcesDir, "testing-resources.yml"), buf.Bytes(), 0644)
	if err != nil {
		return errors.Wrap(err, "writing generated test resources template")
	}

	// order matters here
	resources := []string{
		"install.yml",
		"testing-resources.yml",
	}
	for _, resource := range resources {
		if err := kubectl("create", "-f", filepath.Join(kubeResourcesDir, resource)); err != nil {
			return errors.Wrapf(err, "creating kube resource from "+resource)
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
		return errors.Wrap(err, "waiting for pods to start")
	}
	TestRunner("curl", "test-ingress:19000/logging?config=debug")
	return nil
}

func kubectl(args ...string) error {
	cmd := exec.Command("kubectl", args...)
	log.Printf("k command: %v", cmd.Args)
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

// waitPodsRunning waits for all pods to be running
func waitPodsRunning(podNames ...string) error {
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

// waitPodsTerminated waits for all pods to be terminated
func waitPodsTerminated(podNames ...string) error {
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
	log.Debugf("trying command %v", args)
	return KubectlOut(args...)
}

func waitPodStatus(pod, status string, finished func(output string) bool) error {
	timeout := time.Second * 20
	interval := time.Millisecond * 1000
	tick := time.Tick(interval)

	log.Debugf("waiting %v for pod %v to be %v...", timeout, pod, status)
	for {
		select {
		case <-time.After(timeout):
			return fmt.Errorf("timed out waiting for %v to be %v", pod, status)
		case <-tick:
			out, err := KubectlOut("get", "pod", "-l", "gloo="+pod)
			if err != nil {
				return fmt.Errorf("failed getting pod: %v", err)
			}
			if strings.Contains(out, "CrashLoopBackOff") {
				out, _ := KubectlOut("logs", "-l", "gloo="+pod)
				return errors.Errorf("%v in crash loop with logs %v", pod, out)
			}
			if finished(out) {
				return nil
			}
		}
	}
}

func waitNamespaceStatus(namespace, status string, finished func(output string) bool) error {
	timeout := time.Second * 20
	interval := time.Millisecond * 1000
	tick := time.Tick(interval)

	log.Debugf("waiting %v for namespace %v to be %v...", timeout, namespace, status)
	for {
		select {
		case <-time.After(timeout):
			return fmt.Errorf("timed out waiting for %v to be %v", namespace, status)
		case <-tick:
			out, err := KubectlOut("get", "namespace", namespace)
			if err != nil {
				return fmt.Errorf("failed getting pod: %v", err)
			}
			if finished(out) {
				return nil
			}
		}
	}
}
