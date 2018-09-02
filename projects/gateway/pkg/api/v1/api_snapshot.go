package v1

import (
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type ApiSnapshot struct {
	Gateways        GatewaysByNamespace
	VirtualServices VirtualServicesByNamespace
}

func (s ApiSnapshot) Clone() ApiSnapshot {
	return ApiSnapshot{
		Gateways:        s.Gateways.Clone(),
		VirtualServices: s.VirtualServices.Clone(),
	}
}

func (s ApiSnapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, gateway := range snapshotForHashing.Gateways.List() {
		resources.UpdateMetadata(gateway, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		gateway.SetStatus(core.Status{})
	}
	for _, virtualService := range snapshotForHashing.VirtualServices.List() {
		resources.UpdateMetadata(virtualService, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		virtualService.SetStatus(core.Status{})
	}
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}
