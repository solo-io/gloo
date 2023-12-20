package tap_server

import (
	"github.com/solo-io/gloo/test/ginkgo/parallel"
	"github.com/solo-io/tap-extension-examples/pkg/data_scrubber"
	tap_service "github.com/solo-io/tap-extension-examples/pkg/tap_service"
)

const (
	basePort uint32 = 9831
)

type Factory struct {
	basePort uint32
}

func NewFactory() *Factory {
	return &Factory{
		basePort: basePort,
	}
}

func (f *Factory) NewInstance(address string, instanceConfig *InstanceConfig) *Instance {
	instance := &Instance{
		address:     address,
		HostPort:    int(advancePort(&f.basePort)),
		TapRequests: make(chan tap_service.TapRequest),
	}
	if instanceConfig.EnableDataScrubbing {
		instance.DataScrubber = &data_scrubber.DataScrubber{}
		instance.DataScrubber.Init()
	}
	return instance
}

func advancePort(p *uint32) uint32 {
	return parallel.AdvancePortSafeListen(p)
}
