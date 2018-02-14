package kubernetes

import (
	"github.com/pkg/errors"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/mapstruct"
)

type UpstreamSpec struct {
	ServiceName      string `json:"service_name"`
	ServiceNamespace string `json:"service_namespace"`
	ServicePort      string `json:"service_port"`
}

func DecodeUpstreamSpec(generic v1.UpstreamSpec) (UpstreamSpec, error) {
	var s UpstreamSpec
	if err := mapstruct.Decode(generic, &s); err != nil {
		return s, err
	}
	return s, s.validateUpstream()
}

func EncodeUpstreamSpec(spec UpstreamSpec) v1.UpstreamSpec {
	return v1.UpstreamSpec{
		"service_name":      spec.ServiceName,
		"service_namespace": spec.ServiceNamespace,
		"service_port":      spec.ServicePort,
	}
}

func (s *UpstreamSpec) validateUpstream() error {
	if s.ServiceName == "" {
		return errors.New("service name must be set")
	}
	if s.ServiceNamespace == "" {
		return errors.New("service namespace must be set")
	}
	return nil
}
