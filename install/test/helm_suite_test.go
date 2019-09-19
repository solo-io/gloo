package test

import (
	"bytes"
	"io"
	"os/exec"
	"sync"
	"testing"

	. "github.com/solo-io/go-utils/manifesttestutils"
	"github.com/solo-io/go-utils/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHelm(t *testing.T) {
	RegisterFailHandler(Fail)
	testutils.RegisterCommonFailHandlers()
	RunSpecs(t, "Helm Suite")
}

const (
	namespace = "gloo-system"
)

var (
	version      string
	testManifest TestManifest
	// use a mutex to prevent these tests from running in parallel
	makefileSerializer sync.Mutex
)

func MakeCmd(dir string, args ...string) error {
	makeCmd := exec.Command("make", args...)
	makeCmd.Dir = dir

	makeCmd.Stdout = GinkgoWriter
	makeCmd.Stderr = GinkgoWriter
	err := makeCmd.Run()

	return err
}

func ExecHelm(dir string, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	if err := MakeCmd(dir, "init-helm"); err != nil {
		return nil, nil, err
	}
	helmCmd := exec.Command("helm", args...)
	helmCmd.Dir = dir

	outBuf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	helmCmd.Stdout = outBuf
	helmCmd.Stderr = errBuf
	err := helmCmd.Run()
	return outBuf, errBuf, err
}

func ExecHelmTemplate(dir string, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	templateArgs := append([]string{"template"}, args...)
	return ExecHelm(dir, templateArgs...)
}

func WriteTestManifest(f io.Writer, baseArgs []string, customHelmArgs ...string) error {
	helmTemplateArgs := append(baseArgs, customHelmArgs...)
	stdOut, stdErr, err := ExecHelmTemplate("../..", helmTemplateArgs...)
	if err != nil {
		_, _ = io.Copy(GinkgoWriter, stdErr)
		return err
	}
	if _, err := f.Write(stdOut.Bytes()); err != nil {
		return err
	}
	return nil
}

func WriteGlooETestManifest(f io.Writer, customHelmArgs ...string) error {
	baseArgs := []string{"install/helm/gloo-ee", "--namespace", "gloo-system", "--name=glooe"}
	return WriteTestManifest(f, baseArgs, customHelmArgs...)
}

func WriteGlooOsWithRoUiTestManifest(f io.Writer, customHelmArgs ...string) error {
	baseArgs := []string{"install/helm/gloo-os-with-ui", "--namespace", "gloo-system", "--name=gloo"}
	return WriteTestManifest(f, baseArgs, customHelmArgs...)
}
