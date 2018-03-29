package nats

import (
	"github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/internal/detector"
	"github.com/solo-io/gloo-plugins/nats-streaming"
	"github.com/solo-io/gloo/pkg/log"
)

const (
	// TODO: add more cluster ids, or extend with config options
	defaultClusterID = "test-cluster"
)

type natsDetector struct {
	clusterID string
}

func NewNatsDetector(clusterID string) detector.Interface {
	if clusterID == "" {
		clusterID = defaultClusterID
	}
	return &natsDetector{
		clusterID: defaultClusterID,
	}
}

// if it detects the upstream is a known functional type, give us the
// service info and annotations to mark it with
func (d *natsDetector) DetectFunctionalService(addr string) (*v1.ServiceInfo, map[string]string, error) {
	// try to connect to the addr as though it's a NATS cluster
	log.Debugf("trying to connect to nats cluster at nats://%v with cluster id  %s", addr, d.clusterID)
	c, err := stan.Connect(d.clusterID, "gloo-function-discovery", stan.NatsURL("nats://"+addr))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to connect to nats-streaming cluster")
	}
	defer c.Close()

	svcInfo := &v1.ServiceInfo{
		Type: natsstreaming.ServiceTypeNatsStreaming,
		Properties: natsstreaming.EncodeServiceProperties(natsstreaming.ServiceProperties{
			ClusterID: d.clusterID,
		}),
	}

	return svcInfo, nil, nil
}
