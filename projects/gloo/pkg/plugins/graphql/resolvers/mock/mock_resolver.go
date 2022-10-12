package mock

import (
	"encoding/json"

	"github.com/rotisserie/eris"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
)

var (
	StaticResolverTypedExtensionConfigName = "io.solo.graphql.resolver.static"
)

func TranslateMockResolver(r *v1beta1.MockResolver) (*v3.TypedExtensionConfig, error) {
	envoyConfig := &v2.StaticResolver{
		Response: nil,
	}
	switch response := r.Response.(type) {
	case *v1beta1.MockResolver_SyncResponse:
		{
			out, err := json.Marshal(response.SyncResponse)
			if err != nil {
				return nil, eris.Wrapf(err, "unable to marshal async response %s to JSON", response.SyncResponse.String())
			}
			envoyConfig.Response = &v2.StaticResolver_SyncResponse{
				SyncResponse: string(out),
			}
		}
	case *v1beta1.MockResolver_AsyncResponse_:
		{
			out, err := json.Marshal(response.AsyncResponse.GetResponse())
			if err != nil {
				return nil, eris.Wrapf(err, "unable to marshal async response %s to JSON", response.AsyncResponse.GetResponse().String())
			}
			envoyConfig.Response = &v2.StaticResolver_AsyncResponse_{
				AsyncResponse: &v2.StaticResolver_AsyncResponse{
					Response: string(out),
					DelayMs:  uint32(response.AsyncResponse.GetDelay().AsDuration().Milliseconds()),
				},
			}
		}
	case *v1beta1.MockResolver_ErrorResponse:
		{
			envoyConfig.Response = &v2.StaticResolver_ErrorResponse{
				ErrorResponse: response.ErrorResponse,
			}
		}
	default:
		{
			return nil, eris.Errorf("unknown mock resolver responser type %T", r.Response)
		}
	}

	marshalledEnvoyConfig, err := utils.MessageToAny(envoyConfig)
	if err != nil {
		return nil, eris.Wrapf(err, "unable to marshal envoyConfig")
	}
	return &v3.TypedExtensionConfig{
		Name:        StaticResolverTypedExtensionConfigName,
		TypedConfig: marshalledEnvoyConfig,
	}, nil
}
