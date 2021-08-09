package bootstrap_handler

import (
	"context"

	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/apiserverutils"
	"k8s.io/client-go/rest"
)

func NewBootstrapHandler(
	config *rest.Config,
) rpc_edge_v1.BootstrapApiServer {
	return &bootstrapHandler{
		config: config,
	}
}

type bootstrapHandler struct {
	config *rest.Config
}

func (h *bootstrapHandler) IsGlooFedEnabled(ctx context.Context, _ *rpc_edge_v1.GlooFedCheckRequest) (*rpc_edge_v1.GlooFedCheckResponse, error) {
	glooFedEnabled, err := apiserverutils.IsGlooFedEnabled(ctx, h.config)
	if err != nil {
		return nil, err
	}

	return &rpc_edge_v1.GlooFedCheckResponse{
		Enabled: glooFedEnabled,
	}, nil
}
