package dynamodb

import (
	"fmt"
	"sync/atomic"

	"github.com/solo-io/solo-projects/test/services"

	"github.com/aws/aws-sdk-go/aws/endpoints"

	"github.com/solo-io/gloo/test/ginkgo/parallel"
)

const (
	// basePort is the starting port for the ext auth server, but it is not a special value
	basePort = uint32(14000)

	defaultAddress      = "127.0.0.1"
	imageName           = "amazon/dynamodb-local:1.22.0"
	dynamoContainerPort = 8000
	containerNameFormat = "dynamodb_%d"
)

type Factory struct {
	basePort uint32
}

func NewFactory() *Factory {
	return &Factory{
		basePort: basePort,
	}
}

func (f *Factory) NewInstance() *Instance {
	instancePort := advancePort(&f.basePort)
	containerName := fmt.Sprintf(containerNameFormat, instancePort)

	if services.RunningInDocker() {
		return &Instance{
			dockerRunArgs: []string{
				"-d",
				"--rm",
				"-p", fmt.Sprintf("%d:%d", dynamoContainerPort, dynamoContainerPort),
				"--net", services.GetContainerNetwork(),
				imageName,
			},
			containerName: containerName,
			port:          dynamoContainerPort,
			// At the moment, the DynamoDB Instance relies on the services/docker.go file
			// That file uses some name adjustments to work in CI
			// It's not ideal that we need to do this, but it is a temporary solution
			address: services.GetUpdatedContainerName(containerName),
			region:  endpoints.UsEast2RegionID,
		}
	}

	return &Instance{
		dockerRunArgs: []string{
			"-d",
			"--rm",
			"-p", fmt.Sprintf("%d:%d", instancePort, dynamoContainerPort),
			"--net", services.GetContainerNetwork(),
			imageName,
		},
		containerName: containerName,
		port:          instancePort,
		address:       defaultAddress,
		region:        endpoints.UsEast2RegionID,
	}
}

func advancePort(p *uint32) uint32 {
	return atomic.AddUint32(p, 2) + uint32(parallel.GetPortOffset())
}
