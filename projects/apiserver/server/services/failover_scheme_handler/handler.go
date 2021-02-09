package failover_scheme_handler

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	rpc_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.rpc/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/apiserverutils"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewFailoverSchemeHandler(
	failoverSchemeClient fedv1.FailoverSchemeClient,
) rpc_v1.FailoverSchemeApiServer {
	return &failoverSchemeHandler{
		failoverSchemeClient: failoverSchemeClient,
	}
}

type failoverSchemeHandler struct {
	failoverSchemeClient fedv1.FailoverSchemeClient
}

func (k *failoverSchemeHandler) GetFailoverScheme(ctx context.Context, request *rpc_v1.GetFailoverSchemeRequest) (*rpc_v1.GetFailoverSchemeResponse, error) {
	failovers, err := k.failoverSchemeClient.ListFailoverScheme(ctx)
	if err != nil {
		return nil, err
	}
	for _, failover := range failovers.Items {
		if failover.Spec.GetPrimary().Equal(request.GetUpstreamRef()) {
			return &rpc_v1.GetFailoverSchemeResponse{
				FailoverScheme: BuildRpcFailoverScheme(failover),
			}, nil
		}
	}
	return &rpc_v1.GetFailoverSchemeResponse{
		FailoverScheme: &rpc_v1.FailoverScheme{},
	}, nil
}

func (k *failoverSchemeHandler) GetFailoverSchemeYaml(ctx context.Context, request *rpc_v1.GetFailoverSchemeYamlRequest) (*rpc_v1.GetFailoverSchemeYamlResponse, error) {
	failoverScheme, err := k.failoverSchemeClient.GetFailoverScheme(ctx, client.ObjectKey{
		Namespace: request.GetFailoverSchemeRef().GetNamespace(),
		Name:      request.GetFailoverSchemeRef().GetName(),
	})
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get failover scheme")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	content, err := yaml.Marshal(failoverScheme)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to marshal kube resource into yaml")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &rpc_v1.GetFailoverSchemeYamlResponse{
		YamlData: &rpc_v1.ResourceYaml{
			Yaml: string(content),
		},
	}, nil
}

func BuildRpcFailoverScheme(failoverScheme fedv1.FailoverScheme) *rpc_v1.FailoverScheme {
	return &rpc_v1.FailoverScheme{
		Metadata: apiserverutils.ToMetadata(failoverScheme.ObjectMeta),
		Spec:     &failoverScheme.Spec,
		Status:   &failoverScheme.Status,
	}
}
