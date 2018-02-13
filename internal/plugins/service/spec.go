package service

import (
	"github.com/pkg/errors"
	"github.com/solo-io/glue/internal/pkg/mapstruct"
	"github.com/solo-io/glue/pkg/api/types/v1"
)

type UpstreamSpec struct {
	Hosts []Host `json:"hosts"`
}

type Host struct {
	IP   string `json:"ip"`
	Port uint32 `json:"port"`
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
		"hosts": spec.Hosts,
	}
}

func (s *UpstreamSpec) validateUpstream() error {
	if len(s.Hosts) == 0 {
		return errors.New("most provide at least 1 host")
	}
	for _, host := range s.Hosts {
		if host.IP == "" {
			return errors.New("ip cannot be empty")
		}
		if host.Port == 0 {
			return errors.New("port cannot be empty")
		}
	}
	return nil
}
