package redis

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/solo-io/solo-projects/test/testutils"

	"github.com/solo-io/gloo/test/ginkgo/parallel"
)

const (
	defaultBinaryPath = "redis-server"
	defaultBasePort   = uint32(6379)
	defaultAddress    = "127.0.0.1"
)

type Factory struct {
	binaryPath string
	basePort   uint32
}

func NewFactory() *Factory {
	// This is not currently used in our CI pipeline, but we decided to support it anyway
	// It has always existed in the codebase, and provides an easy mechanism to inject behavior
	// It also follows a pattern established where each of the binaries (Consul, Vault, Envoy) can be overridden
	binaryPath := os.Getenv(testutils.RedisBinary)
	if binaryPath == "" {
		binaryPath = defaultBinaryPath
	}

	return &Factory{
		binaryPath: binaryPath,
		basePort:   defaultBasePort,
	}
}

func (f *Factory) NewInstance() *Instance {
	instancePort := parallel.AdvancePortSafeListen(&f.basePort)

	command := exec.Command(f.binaryPath, "--port", fmt.Sprintf("%d", instancePort))

	return &Instance{
		startCmd: command,
		port:     instancePort,
		address:  defaultAddress,
	}
}
