package localhelpers

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/solo-io/gloo-testing/helpers"

	"github.com/onsi/ginkgo"
)

type FunctionDiscoveryFactory struct {
	srcpath           string
	funcDiscoveryPath string
	wasbuilt          bool
}

func NewFunctionDiscoveryFactory() (*FunctionDiscoveryFactory, error) {
	funcDiscoveryPath := os.Getenv("FUNC_DISCOVERY_BINARY")

	if funcDiscoveryPath != "" {
		return &FunctionDiscoveryFactory{
			funcDiscoveryPath: funcDiscoveryPath,
		}, nil
	}
	srcpath := filepath.Join(helpers.SoloDirectory(), "gloo-function-discovery")
	funcDiscoveryPath = filepath.Join(srcpath, "gloo-function-discovery")
	gf := &FunctionDiscoveryFactory{
		srcpath:           srcpath,
		funcDiscoveryPath: funcDiscoveryPath,
	}
	gf.build()
	return gf, nil
}

func (fdf *FunctionDiscoveryFactory) build() error {
	if fdf.srcpath == "" {
		if fdf.funcDiscoveryPath == "" {
			return errors.New("can't build funcDiscovery and none provided")
		}
		return nil
	}
	fdf.wasbuilt = true

	cmd := exec.Command("go", "build", "-v", "-i", "-gcflags", "-N -l", "-o", "gloo-function-discovery", "main.go")
	cmd.Dir = fdf.srcpath
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (fdf *FunctionDiscoveryFactory) NewFunctionDiscoveryInstance() (*FunctionDiscoveryInstance, error) {

	fdi := &FunctionDiscoveryInstance{
		funcDiscoveryPath: fdf.funcDiscoveryPath,
	}

	return fdi, nil
}

func (fdf *FunctionDiscoveryFactory) Clean() error {
	return nil
}

type FunctionDiscoveryInstance struct {
	funcDiscoveryPath string

	cmd *exec.Cmd
}

func (fdi *FunctionDiscoveryInstance) Run(datadir string) error {

	cmd := exec.Command(fdi.funcDiscoveryPath,
		"--storage.type=file",
		"--storage.refreshrate=1s",
		"--secrets.type=file",
		"--secrets.refreshrate=1s",
	)

	cmd.Dir = datadir
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	err := cmd.Start()
	if err != nil {
		return err
	}
	fdi.cmd = cmd
	return nil
}

func (fdi *FunctionDiscoveryInstance) Clean() error {
	if fdi.cmd != nil {
		fdi.cmd.Process.Kill()
		fdi.cmd.Wait()
	}
	return nil
}
