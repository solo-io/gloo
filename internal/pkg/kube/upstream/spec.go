package upstream

import (
	"github.com/mitchellh/mapstructure"
	"github.com/solo-io/glue/pkg/api/types/v1"
)

const Kubernetes v1.UpstreamType = "kubernetes"

type Spec struct {
	ServiceName      string `json:"service_name"`
	ServiceNamespace string `json:"service_namespace"`
	ServicePortName  string `json:"service_port_name"`
}

func FromMap(m map[string]interface{}) (Spec, error) {
	var spec Spec
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &spec,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return spec, err
	}

	return spec, decoder.Decode(m)
}

func ToMap(spec Spec) map[string]interface{} {
	return map[string]interface{}{
		"service_name":      spec.ServiceName,
		"service_namespace": spec.ServiceNamespace,
		"service_port_name": spec.ServicePortName,
	}
}
