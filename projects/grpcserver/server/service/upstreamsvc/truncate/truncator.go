package truncate

import (
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options"
)

//go:generate mockgen -destination mocks/mock_truncator.go github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/truncate UpstreamTruncator

// UpstreamTruncators truncate fields on upstreams
type UpstreamTruncator interface {
	// Truncates fields on the provided upstream
	Truncate(upstream *gloov1.Upstream)
}

type Truncator struct{}

func NewUpstreamTruncator() Truncator {
	return Truncator{}
}

// Truncate the descriptors field on gRPC Kube upstreams, as they can be upwards of 300KB
func (t Truncator) Truncate(upstream *gloov1.Upstream) {
	switch upstream.GetUpstreamType().(type) {
	case *gloov1.Upstream_Kube:
		switch upstream.GetKube().GetServiceSpec().GetPluginType().(type) {
		case *options.ServiceSpec_Grpc:
			if upstream.GetKube().GetServiceSpec().GetGrpc() != nil {
				upstream.GetKube().GetServiceSpec().GetGrpc().Descriptors = nil
			}
		}
	}
}
