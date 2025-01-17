package istio_automtls

import (
	"github.com/solo-io/gloo/projects/gloo/constants"
	"google.golang.org/protobuf/types/known/structpb"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
)

const EnvoyTransportSocketMatch = "envoy.transport_socket_match"

// AddIstioAutomtlsMetadata adds metadata used by the transport_socket_match
// to select the mTLS transport socket. The Envoy metadata label is added
// based on the presence of the Istio workload label "security.istio.io/tlsMode=istio".
func AddIstioAutomtlsMetadata(
	metadata *envoy_config_core_v3.Metadata,
	workloadLabels map[string]string,
	enableAutoMtls bool,
) *envoy_config_core_v3.Metadata {
	if enableAutoMtls {
		// Valid label values are 'istio', 'disabled'
		// https://github.com/istio/api/blob/5b3f065ee1c2802fb4bc6010ac847c181caa6cc3/label/labels.gen.go#L285
		if value, ok := workloadLabels[constants.IstioTlsModeLabel]; ok && value == constants.IstioMutualTLSModeLabel {
			metadata.GetFilterMetadata()[EnvoyTransportSocketMatch] = &structpb.Struct{
				Fields: map[string]*structpb.Value{
					constants.TLSModeLabelShortname: {
						Kind: &structpb.Value_StringValue{
							StringValue: constants.IstioMutualTLSModeLabel,
						},
					},
				},
			}
		}
	}
	return metadata
}
