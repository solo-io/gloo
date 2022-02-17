package bootstrap_handler

import (
	"context"

	"github.com/solo-io/solo-projects/pkg/license"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/apiserverutils"
	"k8s.io/client-go/rest"
)

func NewBootstrapHandler(
	config *rest.Config,
	licensedFeatureProvider *license.LicensedFeatureProvider,
) rpc_edge_v1.BootstrapApiServer {
	return &bootstrapHandler{
		config:                  config,
		licensedFeatureProvider: licensedFeatureProvider,
	}
}

type bootstrapHandler struct {
	config                  *rest.Config
	licensedFeatureProvider *license.LicensedFeatureProvider
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

func (h *bootstrapHandler) IsGraphqlEnabled(_ context.Context, _ *rpc_edge_v1.GraphqlCheckRequest) (*rpc_edge_v1.GraphqlCheckResponse, error) {
	graphqlFeatureState := h.licensedFeatureProvider.GetStateForLicensedFeature(license.GraphQL)
	return &rpc_edge_v1.GraphqlCheckResponse{
		Enabled: graphqlFeatureState.Enabled,
	}, nil
}
