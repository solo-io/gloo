package bootstrap_handler

import (
	"context"

	gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	"github.com/solo-io/solo-projects/pkg/license"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/apiserverutils"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func NewBootstrapHandler(
	mgr manager.Manager,
	licensedFeatureProvider *license.LicensedFeatureProvider,
) rpc_edge_v1.BootstrapApiServer {
	return &bootstrapHandler{
		mgr:                     mgr,
		licensedFeatureProvider: licensedFeatureProvider,
	}
}

type bootstrapHandler struct {
	mgr                     manager.Manager
	licensedFeatureProvider *license.LicensedFeatureProvider
}

func (h *bootstrapHandler) IsGlooFedEnabled(ctx context.Context, _ *rpc_edge_v1.GlooFedCheckRequest) (*rpc_edge_v1.GlooFedCheckResponse, error) {
	glooFedEnabled, err := apiserverutils.IsGlooFedEnabled(ctx, h.mgr.GetConfig())
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

func (h *bootstrapHandler) GetConsoleOptions(ctx context.Context, _ *rpc_edge_v1.GetConsoleOptionsRequest) (*rpc_edge_v1.GetConsoleOptionsResponse, error) {
	glooClientset := gloo_v1.NewClientset(h.mgr.GetClient())
	consoleOptions, err := apiserverutils.GetConsoleOptions(ctx, glooClientset.Settings())
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.GetConsoleOptionsResponse{
		Options: &rpc_edge_v1.ConsoleOptions{
			ReadOnly:           consoleOptions.GetReadOnly(),
			ApiExplorerEnabled: consoleOptions.GetApiExplorerEnabled(),
		},
	}, nil
}
