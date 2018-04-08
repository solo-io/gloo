package localhelpers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/keybase/go-ps"
	"github.com/solo-io/gloo/test/helpers"

	"github.com/onsi/ginkgo"
)

type FunctionDiscoveryFactory struct {
	srcpath           string
	funcDiscoveryPath string
	wasbuilt          bool
}

type FunctionDiscoveryInstance struct {
	funcDiscoveryPath string
	srcpath           string
	cmd               *exec.Cmd
}

func NewFunctionDiscoveryFactory() (*FunctionDiscoveryFactory, error) {
	funcDiscoveryPath := os.Getenv("FUNC_DISCOVERY_BINARY")

	if funcDiscoveryPath != "" {
		return &FunctionDiscoveryFactory{
			funcDiscoveryPath: funcDiscoveryPath,
		}, nil
	}
	srcpath := filepath.Join(helpers.GlooSoloDirectory(), "cmd", "function-discovery")
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
		srcpath:           fdf.srcpath,
	}

	return fdi, nil
}

func (fdf *FunctionDiscoveryFactory) Clean() error {
	return nil
}

func (fdi *FunctionDiscoveryInstance) Run(datadir string) error {

	cmd := exec.Command(fdi.funcDiscoveryPath,
		"--storage.type=file",
		"--storage.refreshrate=1s",
		"--secrets.type=file",
		"--secrets.refreshrate=1s",
		"--files.type=file",
		"--files.refreshrate=1s",
	)

	cmd.Dir = datadir
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	cmd, err := fdi.run(cmd)
	if err != nil {
		return err
	}
	fdi.cmd = cmd
	return nil
}

func (fdi *FunctionDiscoveryInstance) waitForExternalProcess() error {
	for {
		procs, err := ps.Processes()
		if err != nil {
			panic(err)
		}
		for _, proc := range procs {
			str, err := proc.Path()
			if err != nil {
				continue
			}
			if strings.Contains(str, fdi.srcpath) {
				return nil
			}

		}
		time.Sleep(time.Second)
	}
}

func (fdi *FunctionDiscoveryInstance) run(c *exec.Cmd) (*exec.Cmd, error) {
	if os.Getenv("USE_DEBUGGER") == "1" {
		fmt.Println("Please run the following command in your debugger:\n")
		fmt.Printf("%v %v\n", c.Path, c.Args)
		fmt.Printf("CWD %v\n", c.Dir)
		fmt.Println("looking for processes started from", fdi.srcpath)

		return nil, fdi.waitForExternalProcess()
	}
	return c, c.Start()
}

func (fdi *FunctionDiscoveryInstance) Clean() error {
	if fdi.cmd != nil {
		fdi.cmd.Process.Kill()
		fdi.cmd.Wait()
	}
	return nil
}
