package consul

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

type endpointController struct {
}

func (p *plugin) WatchEndpoints(writeNamespace string, upstreamsToTrack v1.UpstreamList, opts clients.WatchOpts) (<-chan v1.EndpointList, <-chan error, error) {
	return nil, nil, errors.Errorf("not implemented")
}
